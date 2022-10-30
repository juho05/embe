package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

type Define struct {
	Name    Token
	Start   Position
	End     Position
	Content []Token
}

type Defines struct {
	defines map[string][]Define
}

func NewDefines() *Defines {
	return &Defines{
		defines: make(map[string][]Define),
	}
}

func (d *Defines) GetDefines(at Position) []Define {
	defines := make([]Define, 0, 10)
	for _, defs := range d.defines {
		for _, def := range defs {
			if def.IsInScope(at) {
				defines = append(defines, def)
			}
		}
	}
	return defines
}

func (d *Defines) GetDefine(name string, at Position) (Define, bool) {
	if d, ok := d.defines[name]; ok {
		for _, def := range d {
			if def.IsInScope(at) {
				return def, true
			}
		}
	}
	return Define{}, false
}

func (d *Defines) addDefine(token Token, pos Position, content []Token) {
	d.undefine(token)
	if _, ok := d.defines[token.Lexeme]; !ok {
		d.defines[token.Lexeme] = make([]Define, 0, 1)
	}
	d.defines[token.Lexeme] = append(d.defines[token.Lexeme], Define{
		Name:    token,
		Start:   pos,
		Content: content,
	})
}

func (d *Defines) copy() *Defines {
	if d == nil {
		return nil
	}

	defines := make(map[string][]Define, len(d.defines))
	for k, v := range d.defines {
		ds := make([]Define, len(v))
		copy(ds, v)
		defines[k] = ds
	}
	return &Defines{
		defines: defines,
	}
}

func (d *Defines) undefine(token Token) {
	if def, ok := d.defines[token.Lexeme]; ok {
		lastDef := def[len(def)-1]
		if lastDef.IsInScope(token.Pos) {
			lastDef.End = Position{
				Line:   token.Pos.Line,
				Column: 0,
			}
		}
		def[len(def)-1] = lastDef
	}
}

func (d *Define) String() string {
	str := fmt.Sprintf("#define %s ", d.Name.Lexeme)
	for _, t := range d.Content {
		str += t.Lexeme
	}
	return strings.TrimSpace(str)
}

func (d *Define) IsInScope(at Position) bool {
	return at.Path == d.Start.Path && at.Line >= d.Start.Line && (at.Line != d.Start.Line || at.Column >= d.Start.Column) && (d.End == (Position{}) || (at.Line <= d.End.Line && (at.Line != d.End.Line || at.Column <= d.End.Column)))
}

type preprocessor struct {
	tokens  []Token
	defines *Defines
	index   int
	errors  []error
	files   map[string][][]rune
	path    string
	stack   []string
	open    func(name string) (io.ReadCloser, error)
}

func Preprocess(tokens []Token, absPath string, open func(name string) (io.ReadCloser, error), stack []string, defines *Defines) ([]Token, map[string][][]rune, *Defines, []string, []error) {
	eof := tokens[len(tokens)-1]

	if stack == nil {
		stack = []string{absPath}
	}

	defines = defines.copy()

	if defines == nil {
		defines = NewDefines()
	}

	for _, ds := range defines.defines {
		for _, d := range ds {
			pos := d.Name.Pos
			pos.Line = 0
			pos.Path = absPath
			defines.addDefine(d.Name, pos, d.Content)
		}
	}

	if open == nil {
		open = func(name string) (io.ReadCloser, error) {
			return os.Open(name)
		}
	}

	p := &preprocessor{
		tokens:  tokens,
		defines: defines,
		errors:  make([]error, 0),
		files:   make(map[string][][]rune),
		stack:   stack,
		path:    absPath,
		open:    open,
	}
	p.preprocess()

	// make sure EOF was not removed
	if p.tokens[len(p.tokens)-1].Type != TkEOF {
		p.tokens = append(p.tokens, eof)
	}

	return p.tokens, p.files, p.defines, p.stack, p.errors
}

func (p *preprocessor) preprocess() {
	for p.index < len(p.tokens) {
		if p.tokens[p.index].Type == TkPreprocessor {
			err := p.directive(p.tokens[p.index].Lexeme)
			if err != nil {
				p.errors = append(p.errors, err)
				return
			}
		} else {
			p.index++
		}
	}
	p.replace()
}

func (p *preprocessor) directive(directive string) error {
	startIndex := p.index
	switch directive {
	case "#include":
		if p.peek().Type == TkLiteral && p.peek().DataType == DTString {
			p.index++
			err := p.include(p.index-1, p.tokens[p.index].Literal.(string))
			if err != nil {
				return err
			}
		} else {
			p.errors = append(p.errors, p.newError("Expected file name after #include."))
		}
	case "#define":
		if p.peek().Type == TkIdentifier {
			p.index++
			nameIndex := p.index
			for p.peek().Type != TkNewLine && p.peek().Type != TkEOF {
				p.index++
			}
			replace := make([]Token, p.index-nameIndex)
			copy(replace, p.tokens[nameIndex+1:p.index+1])
			p.defines.addDefine(p.tokens[nameIndex], p.tokens[nameIndex].Pos, replace)
			p.index++
		} else {
			p.errors = append(p.errors, p.newError("Expected name after #define."))
		}
	case "#undef":
		if p.peek().Type == TkIdentifier {
			p.index++
			p.defines.undefine(p.tokens[p.index])
			p.index++
		} else {
			p.errors = append(p.errors, p.newError("Expected name after #undef."))
		}
	case "#ifdef", "#ifndef":
		if p.peek().Type == TkIdentifier {
			p.index++
			if _, ok := p.defines.GetDefine(p.tokens[p.index].Lexeme, p.tokens[p.index].Pos); (p.tokens[p.index-1].Lexeme == "#ifdef") != ok {
				for (p.peek().Type != TkPreprocessor || p.peek().Lexeme != "#endif") && p.peek().Type != TkEOF {
					p.index++
				}
				p.index++
			}
			if p.peek().Type == TkNewLine {
				p.index++
			}
		} else {
			p.errors = append(p.errors, p.newError("Expected name after #ifdef."))
		}
	case "#endif":
		if p.peek().Type == TkNewLine && p.index > 0 && p.tokens[p.index-1].Type == TkNewLine {
			p.index++
		}
	default:
		p.errors = append(p.errors, p.newErrorAt("Unknown preprocessor directive.", p.tokens[p.index]))
	}
	if p.index >= len(p.tokens) {
		p.index = len(p.tokens) - 1
	}
	p.tokens = slices.Delete(p.tokens, startIndex, p.index+1)
	p.index = startIndex
	return nil
}

func (p *preprocessor) include(keywordIndex int, path string) error {
	if filepath.Ext(path) != ".mb" {
		path += ".mb"
	}

	path = filepath.Join(filepath.Dir(p.path), path)
	if runtime.GOOS == "windows" {
		path = strings.ToLower(path)
	}

	file, err := p.open(path)
	if err != nil {
		return p.newErrorAt(fmt.Sprintf("Unable to open file `%s`: %s", path, err), p.tokens[keywordIndex+1])
	}
	defer file.Close()
	tokens, lines, errs := Scan(file, path)
	p.files[path] = lines
	if len(errs) > 0 {
		p.errors = append(p.errors, errs[:len(errs)-1]...)
		return errs[len(errs)-1]
	}

	if slices.Contains(p.stack, path) {
		files := make([]string, len(p.stack))
		index := 0
		for i, f := range p.stack {
			files[i] = filepath.Base(f)
			if f == path {
				index = i
			}
		}
		return p.newErrorAt(fmt.Sprintf("Include cycle detected: %s -> %s", strings.Join(files[index:], " -> "), filepath.Base(path)), p.tokens[keywordIndex+1])
	}

	p.stack = append(p.stack, path)
	tokens, files, defines, stack, errs := Preprocess(tokens, path, p.open, p.stack, p.defines)
	p.stack = stack[:len(stack)-1]
	for k, v := range files {
		p.files[k] = v
	}
	defines = defines.copy()
	for _, ds := range defines.defines {
		for _, d := range ds {
			pos := d.Name.Pos
			pos.Line += p.tokens[keywordIndex].Pos.Line
			pos.Path = p.path
			p.defines.addDefine(d.Name, pos, d.Content)
		}
	}
	if len(errs) > 0 {
		p.errors = append(p.errors, errs[:len(errs)-1]...)
		return errs[len(errs)-1]
	}

	tokens = tokens[:len(tokens)-1]

	newTokens := make([]Token, len(p.tokens)+len(tokens))
	copy(newTokens, p.tokens[:keywordIndex])
	copy(newTokens[keywordIndex:], tokens)
	if keywordIndex < len(p.tokens)-3 {
		copy(newTokens[keywordIndex+len(tokens):], p.tokens[keywordIndex+3:])
	}
	i := len(newTokens) - 1
	for ; newTokens[i].Type == 0 && newTokens[i].Lexeme == ""; i-- {
	}
	p.tokens = newTokens[:i+1]
	p.index -= 2
	return nil
}

func (p *preprocessor) replace() {
	for i := 0; i < len(p.tokens); i++ {
		token := p.tokens[i]
		if token.Type != TkIdentifier {
			continue
		}
		if d, ok := p.defines.GetDefine(token.Lexeme, token.Pos); ok {
			if len(d.Content) == 0 {
				size := 1
				if i < len(p.tokens)-1 && p.tokens[i+1].Type == TkNewLine {
					size = 2
				}
				slices.Delete(p.tokens, i, i+size)
				i--
				continue
			}
			newTokens := make([]Token, len(p.tokens)+len(d.Content)-1)
			copy(newTokens, p.tokens[:i])
			copy(newTokens[i:], d.Content)
			if i < len(p.tokens)-1 {
				copy(newTokens[i+len(d.Content):], p.tokens[i+1:])
			}
			for j := i; j < len(newTokens) && j < i+len(d.Content); j++ {
				newTokens[j].Indent = token.Indent
				newTokens[j].Pos = token.Pos
				newTokens[j].EndPos = token.EndPos
			}
			p.tokens = newTokens
			i--
		}
	}
}

func (p *preprocessor) peek() Token {
	if p.index+1 >= len(p.tokens) {
		return Token{
			Type: TkEOF,
		}
	}
	return p.tokens[p.index+1]
}

func (p *preprocessor) newError(message string) error {
	return p.newErrorAt(message, p.peek())
}

func (p *preprocessor) newErrorAt(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
	}
}
