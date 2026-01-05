package yatch

import (
	"bytes"
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
	case *ast.StringNode, *ast.IntegerNode, *ast.FloatNode, *ast.BoolNode, *ast.NullNode:
		return node, nil
	case *ast.SequenceNode:
		if len(node.Values) == 0 {
			return nil, ErrNodeNotFound
		}
		return node.Values[0], nil
	default:
		return node, ErrNodeTypeUnsupported
	}
}

func (file *YAMLFile) Patch(jsonpath, value string) error {
	node, err := file.Find(jsonpath)
	if err != nil {
		return err
	}

	token := (*Token)(node.GetToken())
	return file.Replace(token, value)
}

func (file *YAMLFile) PatchComment(jsonpath, value string) error {
	node, err := file.Find(jsonpath)
	if err != nil {
		return err
	}

	switch node.GetToken().NextType() {
	case token.CommentType:
		// continue
	default:
		return ErrTokenNotFound
	}

	next := node.GetToken().Next
	return file.Replace((*Token)(next), value)
}

func (file *YAMLFile) Replace(token *Token, value string) error {
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

func (file *YAMLFile) setContent(content []byte) error {
	parsed, err := parser.ParseBytes(content, 0) // Verify that file is still valid
	if err != nil {
		return err
	}

	file.content = content
	file.file = parsed
	return nil
}
