package gtl

import "testing"

func (t *Token) isEqual(other *Token) bool {
	if t == nil {
		return other == nil
	}
	return t.tokenType == other.tokenType &&
		t.text == other.text
}

func TestNextToken(t *testing.T) {
	type testcase struct {
		src      string
		want     *Token
		curAfter int
	}
	testcases := []testcase{
		{"", &Token{EOF, ""}, 0},
		{"a", &Token{IDENTIFIER, "a"}, 1},
		{"ab", &Token{IDENTIFIER, "ab"}, 2},
		{"a b", &Token{IDENTIFIER, "a"}, 1},
		{"(", &Token{LPAREN, "("}, 1},
		{"(a", &Token{LPAREN, "("}, 1},
		{")", &Token{RPAREN, ")"}, 1},
		{" )", &Token{RPAREN, ")"}, 2},
		{"->", &Token{ARROW, "->"}, 2},
		{".", &Token{DOT, "."}, 1},
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
