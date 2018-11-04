package gtl

import (
	"errors"
	"strings"
)

// Lexer is a lexer for typed_lang
type Lexer struct {
	source string
	cur    int
}

// TokenType is an enum for token, which represents what a token is.
type TokenType uint8

const (
	// EOF is an end of file
	EOF = iota
	// Word is an Word, which may be a variable name or function name, etc.
	Word
	// LParen is "("
	LParen
	// RParen is ")"
	RParen
	// LBlace is "{"
	LBlace
	// RBlace is "}"
	RBlace
	// Arrow is "->"
	Arrow
	// Dot is "."
	Dot
)

// Token is a token of typed_lang
type Token struct {
	tokenType TokenType
	text      string
}

var (
	// ErrUnknownToken is an error for lexer, which means the source is not a valid typed_lang
	ErrUnknownToken = errors.New("Unknown token")
)

// NewLexer returns a new lexer from source string
func NewLexer(source string) *Lexer {
	return &Lexer{source, 0}
}

// NextToken returns a next token, and increments its cursor
func (l *Lexer) NextToken() (*Token, error) {
	beg := l.cur
	if beg == len(l.source) {
		return &Token{EOF, ""}, nil
	}
	if beg > len(l.source) {
		return nil, errors.New("NextToken is called after EOF")
	}

	idx := l.cur
	var mode TokenType
	c := l.source[idx : idx+1]
	switch {
	case isWhitespace(c):
		l.cur++
		return l.NextToken()
	case strings.Contains("abcdefghijklmnopqrstuvwxyz", c):
		mode = Word
		for ; idx < len(l.source); idx++ {
			if l.source[idx:idx+1] == " " {
				break
			}
		}
		l.cur = idx
		return &Token{mode, l.source[beg:idx]}, nil
	case c == "(":
		mode = LParen
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == ")":
		mode = RParen
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == "{":
		mode = LBlace
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == "}":
		mode = RBlace
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == ".":
		mode = Dot
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == "-":
		if l.source[idx+1:idx+2] == ">" {
			l.cur += 2
			return &Token{Arrow, l.source[beg : beg+2]}, nil
		}
	}

	return nil, ErrUnknownToken
}

func isWhitespace(s string) bool {
	return s == "\t" || s == "\n" || s == "\r" || s == "\f" || s == " "
}
