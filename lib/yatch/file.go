package yatch

import (
	"bytes"
	"io"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

// File represents a parsed YAML file and its content.
type File struct {
	content []byte
	file    *ast.File
}

// NewFile creates a new YAMLFile from a reader.
func NewFile(reader io.Reader) (*File, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	file := &File{}
	err = file.setContent(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (file *File) Content() []byte {
	return file.content
}

func (file *File) Replace(token *Token, value string) error {
	if strings.ContainsAny(value, "\r\n") { // a line break could add documents or keys around the token
		return ErrValueMultiline
	}

	content := bytes.Buffer{}
	rewriter := NewRewriter(
		bytes.NewReader(file.content),
		&content,
	)

	discard := token.Len()
	value = token.Render(value)

	if err := rewriter.CopyLines(token.Position.Line); err != nil { // Copy lines before yaml node
		return err
	}

	if err := rewriter.CopyBytes(token.Position.Column); err != nil { // Copy bytes before yaml node
		return err
	}

	if err := rewriter.Discard(discard); err != nil { // Discard original yaml node
		return err
	}

	if err := rewriter.WriteString(value); err != nil { // Write new value
		return err
	}

	if err := rewriter.CopyAll(); err != nil { // Copy the rest of the content
		return err
	}

	return file.setContent(content.Bytes())
}

func (file *File) setContent(content []byte) error {
	parsed, err := parser.ParseBytes(content, 0) // Verify that file is still valid
	if err != nil {
		return err
	}

	file.content = content
	file.file = parsed
	return nil
}
