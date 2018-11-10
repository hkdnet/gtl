all: nodetype_string.go tokentype_string.go

nodetype_string.go: nodetype.go
	stringer -type=NodeType

tokentype_string.go: lexer.go
	stringer -type=TokenType

deps:
	go get -u golang.org/x/tools/cmd/stringer

.PHONY: deps
