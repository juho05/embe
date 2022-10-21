package parser

import (
	"fmt"
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
	defines map[string][]*Define
}

func NewDefines() *Defines {
	return &Defines{
		defines: make(map[string][]*Define),
	}
}

func (d *Defines) GetDefines(at Position) []*Define {
	defines := make([]*Define, 0, 10)
	for _, defs := range d.defines {
		for _, def := range defs {
			if def.IsInScope(at) {
				defines = append(defines, def)
			}
		}
	}
	return defines
}

func (d *Defines) GetDefine(name string, at Position) (*Define, bool) {
	if d, ok := d.defines[name]; ok {
		for _, def := range d {
			if def.IsInScope(at) {
				return def, true
			}
		}
	}
	return nil, false
}

func (d *Defines) addDefine(token Token, content []Token) {
	d.undefine(token)
	if _, ok := d.defines[token.Lexeme]; !ok {
		d.defines[token.Lexeme] = make([]*Define, 0, 1)
	}
	d.defines[token.Lexeme] = append(d.defines[token.Lexeme], &Define{
		Name:    token,
		Start:   token.Pos,
		Content: content,
	})
}

func (d *Defines) undefine(token Token) {
	if def, ok := d.defines[token.Lexeme]; ok {
		lastDef := def[len(def)-1]
		if lastDef.IsInScope(token.Pos) {
			lastDef.End = Position{
				Line:   token.Pos.Line - 1,
				Column: 0,
			}
		}
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
	return at.Line >= d.Start.Line && (at.Line != d.Start.Line || at.Column >= d.Start.Column) && (d.End == (Position{}) || (at.Line <= d.End.Line && (at.Line != d.End.Line || at.Column <= d.End.Column)))
}

type preprocessor struct {
	tokens  []Token
	defines *Defines
	index   int
	errors  []error
}

func Preprocess(tokens []Token) ([]Token, *Defines, []error) {
	p := &preprocessor{
		tokens:  tokens,
		defines: NewDefines(),
		errors:  make([]error, 0),
	}
	p.preprocess()
	return p.tokens, p.defines, p.errors
}

func (p *preprocessor) preprocess() {
	for p.index < len(p.tokens) {
		if p.tokens[p.index].Type == TkPreprocessor {
			p.directive(p.tokens[p.index].Lexeme)
		} else {
			p.index++
		}
	}
	p.replace()
}

func (p *preprocessor) directive(directive string) {
	startIndex := p.index
	switch directive {
	case "#define":
		if p.peek().Type == TkIdentifier {
			p.index++
			nameIndex := p.index
			for p.peek().Type != TkNewLine && p.peek().Type != TkEOF {
				p.index++
			}
			replace := make([]Token, p.index-nameIndex)
			copy(replace, p.tokens[nameIndex+1:p.index+1])
			p.defines.addDefine(p.tokens[nameIndex], replace)
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
			newTokens := make([]Token, len(p.tokens)+len(d.Content))
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
