package yatch

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
)

// YAMLFile represents a parsed YAML file and its content.
type YAMLFile struct {
	content []byte
	file    *ast.File
}

// NewFile creates a new YAMLFile from a reader.
func NewFile(reader io.Reader) (*YAMLFile, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var file YAMLFile
	return &file, file.setContent(content)
}

func (file *YAMLFile) Content() []byte {
	return file.content
}

func (file *YAMLFile) Find(jsonpath string) ([]ast.Node, error) {
	path, err := yaml.PathString(jsonpath)
	if err != nil {
		return nil, err
	}

	node, err := path.FilterFile(file.file)
	if err != nil {
		// Treat "node not found" as an empty result, not an error
		if yaml.IsNotFoundNodeError(err) {
			return []ast.Node{}, nil
		}
		return nil, err
	}

	if node == nil {
		return nil, errors.New("node not found")
	}
	
	switch node := node.(type) {
		case *ast.SequenceNode:
			return node.Values, nil
		default:
			return []ast.Node{node}, nil
	}
}

func (file *YAMLFile) Patch(jsonpath, value string) error {
	nodes, err := file.Find(jsonpath)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		switch node.(type) {
		case *ast.StringNode, *ast.IntegerNode, *ast.FloatNode, *ast.BoolNode, *ast.NullNode:
		default:
			return errors.ErrUnsupported
		}

		tkn := node.GetToken()
		switch tkn.Type {
		case token.MappingKeyType, token.SequenceStartType, token.MappingStartType:
			return errors.ErrUnsupported
		}

		if err := file.Replace(node, value); err != nil {
			return err
		}
	}

	return nil
}

func (file *YAMLFile) Replace(node ast.Node, value string) error {
	content := bytes.Buffer{}
	rewriter := NewRewriter(
		bytes.NewReader(file.content),
		&content,
	)

	tkn := node.GetToken()
	switch tkn.Type {
	case token.DoubleQuoteType:
		value = fmt.Sprintf("\"%s\"", value)
	case token.SingleQuoteType:
		value = fmt.Sprintf("'%s'", value)
	}

	// Copy content before yaml node
	rewriter.CopyLines(tkn.Position.Line)
	rewriter.CopyBytes(tkn.Position.Column)

	discardLen := len(tkn.Value)
	switch tkn.Type {
	case token.SingleQuoteType, token.DoubleQuoteType:
		discardLen += 2
	}

	rewriter.Discard(discardLen)
	rewriter.Write([]byte(value)) // Write new value
	rewriter.Copy() // Copy the rest of the content

	file.setContent(content.Bytes())
	return nil
}

func (file *YAMLFile) setContent(content []byte) error {
	parsed, err := parser.ParseBytes(content, 0) // Verify that file is still valid
	if err != nil {
		return err
	}

	file.content = content
	file.file = parsed
	return nil
}
