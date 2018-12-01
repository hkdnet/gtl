package gtl

import (
	"errors"
	"strings"
)

// Lexer is a lexer for typed_lang
type Lexer struct {
	source string
	cur    int

	hasNext bool
}

var keywordMap map[string]TokenType

// Token is a token of typed_lang
type Token struct {
	TokenType TokenType
	Text      string
}

var (
	// ErrUnknownToken is an error for lexer, which means the source is not a valid typed_lang
	ErrUnknownToken = errors.New("Unknown token")
)

func init() {
	keywordMap = make(map[string]TokenType)
	keywordMap["true"] = KeywordTrue
	keywordMap["false"] = KeywordFalse
	keywordMap["if"] = KeywordIf
	keywordMap["then"] = KeywordThen
	keywordMap["else"] = KeywordElse
	keywordMap["iszero"] = KeywordIsZero
}

// NewLexer returns a new lexer from source string
func NewLexer(source string) *Lexer {
	return &Lexer{source, 0, true}
}

// HasNext returns whether this lexer has more tokens or not
func (l *Lexer) HasNext() bool {
	return l.hasNext
}

// NextToken returns a next token, and increments its cursor
func (l *Lexer) NextToken() (*Token, error) {
	beg := l.cur
	if beg == len(l.source) {
		l.hasNext = false
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
		for kw, tt := range keywordMap {
			if size := len(kw); len(l.source) >= idx+size && l.source[idx:idx+size] == kw {
				l.cur += size
				return &Token{tt, kw}, nil
			}
		}
		mode = Word
		for ; idx < len(l.source); idx++ {
			c := l.source[idx : idx+1]
			if isWhitespace(c) {
				break
			}
			if c == "(" || c == ")" || c == "." {
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
	case c == "0":
		mode = Number
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
