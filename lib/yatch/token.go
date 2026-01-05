package yatch

import (
	"fmt"

	"github.com/goccy/go-yaml/token"
)

type Token token.Token

func (t Token) Len() int {
	switch t.Type {
	case token.DoubleQuoteType, token.SingleQuoteType:
		return len([]byte(t.Value)) + 2
	case token.CommentType:
		return len([]byte(t.Value)) + 1
	default:
		return len([]byte(t.Value))
	}
}

func (t Token) Render(value string) string {
	switch t.Type {
	case token.DoubleQuoteType:
		return fmt.Sprintf("\"%s\"", value)
	case token.SingleQuoteType:
		return fmt.Sprintf("'%s'", value)
	case token.CommentType:
		return fmt.Sprintf("%s %s", string(token.CommentCharacter), value)
	default:
		return value
	}
}