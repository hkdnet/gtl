package gtl

import "testing"

func (t *Token) isEqual(other *Token) bool {
	if t == nil {
		return other == nil
	}
	return t.TokenType == other.TokenType &&
		t.Text == other.Text
}

func TestNextToken(t *testing.T) {
	type testcase struct {
		src      string
		want     *Token
		curAfter int
	}
	testcases := []testcase{
		{"", &Token{EOF, ""}, 0},
		{"a", &Token{Word, "a"}, 1},
		{"ab", &Token{Word, "ab"}, 2},
		{"a b", &Token{Word, "a"}, 1},
		{"a)", &Token{Word, "a"}, 1},
		{"a.", &Token{Word, "a"}, 1},
		{"(", &Token{LParen, "("}, 1},
		{"(a", &Token{LParen, "("}, 1},
		{")", &Token{RParen, ")"}, 1},
		{" )", &Token{RParen, ")"}, 2},
		{"{", &Token{LBlace, "{"}, 1},
		{"{a", &Token{LBlace, "{"}, 1},
		{"}", &Token{RBlace, "}"}, 1},
		{" }", &Token{RBlace, "}"}, 2},
		{"->", &Token{Arrow, "->"}, 2},
		{".", &Token{Dot, "."}, 1},
		{"0", &Token{Number, "0"}, 1},
		{"true", &Token{Keyword, "true"}, 4},
		{"false", &Token{Keyword, "false"}, 5},
		{"if", &Token{Keyword, "if"}, 2},
		{"then", &Token{Keyword, "then"}, 4},
		{"else", &Token{Keyword, "else"}, 4},
	}
	for i, v := range testcases {
		l := NewLexer(v.src)
		got, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		if !got.isEqual(v.want) {
			t.Errorf("case %d: want %#v but got %#v\n", i, v.want, got)
		}
		if l.cur != v.curAfter {
			t.Errorf("case %d: want cursor to be %d but %d\n", i, v.curAfter, l.cur)
		}
	}
	// unknown token
	{
		l := NewLexer("❗")
		_, err := l.NextToken()
		if err == nil {
			t.Error("next token should return with unknown token")
		} else if err != ErrUnknownToken {
			t.Errorf("err should be ErrUnknownToken but got %v", err)
		}
	}
}
