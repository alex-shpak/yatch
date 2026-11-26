package yatch

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/token"
)

type Token token.Token

func (t Token) Len() int {
	return len(strings.TrimSpace(t.Origin))
}

func (t Token) Render(value string) string {
	switch t.Type {
	case token.DoubleQuoteType:
		return strconv.Quote(value)
	case token.CommentType:
		return fmt.Sprintf("%s %s", string(token.CommentCharacter), value)
	case token.StringType, token.SingleQuoteType:
		bytes, _ := yaml.MarshalWithOptions(value, yaml.UseSingleQuote(true))
		return strings.TrimSpace(string(bytes))
	default:
		return value
	}
}
