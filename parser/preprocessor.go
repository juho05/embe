package parser

import (
	"golang.org/x/exp/slices"
)

type preprocessor struct {
	tokens  []Token
	defines map[string][]Token
	index   int
	errors  []error
}

func Preprocess(tokens []Token) ([]Token, []error) {
	p := &preprocessor{
		tokens:  tokens,
		defines: make(map[string][]Token),
		errors:  make([]error, 0),
	}
	p.preprocess()
	return p.tokens, p.errors
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
			p.defines[p.tokens[nameIndex].Lexeme] = replace
			p.index++
		} else {
			p.errors = append(p.errors, p.newError("Expected name after #define."))
		}
	case "#ifdef", "#ifndef":
		if p.peek().Type == TkIdentifier {
			p.index++
			if _, ok := p.defines[p.tokens[p.index].Lexeme]; (p.tokens[p.index-1].Lexeme == "#ifdef") != ok {
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
		if d, ok := p.defines[token.Lexeme]; ok {
			if len(d) == 0 {
				size := 1
				if i < len(p.tokens)-1 && p.tokens[i+1].Type == TkNewLine {
					size = 2
				}
				slices.Delete(p.tokens, i, i+size)
				i--
				continue
			}
			for i := range d {
				d[i].Indent = token.Indent
				d[i].Pos = token.Pos
				d[i].EndPos = token.EndPos
			}
			newTokens := make([]Token, len(p.tokens)+len(d))
			copy(newTokens, p.tokens[:i])
			copy(newTokens[i:], d)
			if i < len(p.tokens)-1 {
				copy(newTokens[i+len(d):], p.tokens[i+1:])
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
