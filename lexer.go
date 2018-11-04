package gtl

import (
	"errors"
	"strings"
)

type Lexer struct {
	source string
	cur    int
}
type TokenType uint8

const (
	EOF = iota
	IDENTIFIER
	LPAREN
	RPAREN
	ARROW
	DOT
)

type Token struct {
	tokenType TokenType
	text      string
}

var ErrUnknownToken = errors.New("Unknown token")

func (t *Token) IsEqual(other *Token) bool {
	if t == nil {
		return other == nil
	}
	return t.tokenType == other.tokenType &&
		t.text == other.text
}

func NewLexer(source string) *Lexer {
	return &Lexer{source, 0}
}

var NotFound error = errors.New("no new token")

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
		mode = IDENTIFIER
		for ; idx < len(l.source); idx++ {
			if l.source[idx:idx+1] == " " {
				break
			}
		}
		l.cur = idx
		return &Token{mode, l.source[beg:idx]}, nil
	case c == "(":
		mode = LPAREN
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == ")":
		mode = RPAREN
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == ".":
		mode = DOT
		l.cur++
		return &Token{mode, l.source[beg : beg+1]}, nil
	case c == "-":
		if l.source[idx+1:idx+2] == ">" {
			l.cur += 2
			return &Token{ARROW, l.source[beg : beg+2]}, nil
		}
	}

	return nil, ErrUnknownToken
}

func isWhitespace(s string) bool {
	return s == "\t" || s == "\n" || s == "\r" || s == "\f" || s == " "
}
