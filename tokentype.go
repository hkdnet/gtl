package gtl

// TokenType is an enum for token, which represents what a token is.
type TokenType uint8

const (
	// EOF is an end of file
	EOF TokenType = iota
	// Word is an Word, which may be a variable name, function name, or keyword such as if, etc.
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
	// Number is "0"
	Number
	// KeywordTrue is "true"
	KeywordTrue
	// KeywordFalse is "false"
	KeywordFalse
	// KeywordIf is "if"
	KeywordIf
	// KeywordThen is "then"
	KeywordThen
	// KeywordElse is "else"
	KeywordElse
	// KeywordIsZero is "iszero"
	KeywordIsZero
)
