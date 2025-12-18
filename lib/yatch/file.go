package yatch

import (
	"bytes"
	"errors"
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

	file := &YAMLFile{}
	err = file.setContent(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (file *YAMLFile) Content() []byte {
	return file.content
}

func (file *YAMLFile) Find(jsonpath string) (ast.Node, error) {
	path, err := yaml.PathString(jsonpath)
	if err != nil {
		return nil, err
	}

	node, err := path.FilterFile(file.file)
	if err != nil {
		return nil, err
	}
	
	switch node := node.(type) {
		case *ast.SequenceNode:
			if len(node.Values) == 0 {
				return nil, yaml.ErrNotFoundNode
			}
			return node.Values[0], nil
		default:
			return node, nil
	}
}

func (file *YAMLFile) Patch(jsonpath, value string) error {
	node, err := file.Find(jsonpath)
	if err != nil {
		return err
	}

	switch node.(type) {
	case *ast.StringNode, *ast.IntegerNode, *ast.FloatNode, *ast.BoolNode, *ast.NullNode:
	default:
		return errors.ErrUnsupported
	}

	switch (node.GetToken()).Type {
	case token.MappingKeyType, token.SequenceStartType, token.MappingStartType:
		return errors.ErrUnsupported
	}

	if err := file.Replace(node, value); err != nil {
		return err
	}
	
	return nil
}

func (file *YAMLFile) Replace(node ast.Node, value string) error {
	content := bytes.Buffer{}
	rewriter := NewRewriter(
		bytes.NewReader(file.content),
		&content,
	)

	tkn := (*Token)(node.GetToken())
	
	value = tkn.Render(value)
	discard := tkn.Len()

	rewriter.CopyLines(tkn.Position.Line) // Copy lines before yaml node
	rewriter.CopyBytes(tkn.Position.Column) // Copy bytes before yaml node
	rewriter.Discard(discard) // Discard original yaml node
	rewriter.WriteString(value) // Write new value
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
