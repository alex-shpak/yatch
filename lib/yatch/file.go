package yatch

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// YAMLFile represents a parsed YAML file and its content.
type YAMLFile struct {
	content []byte
	node    *yaml.Node
}

// NewFile
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

func (file *YAMLFile) Find(jsonpath string) (node []*yaml.Node, err error) {
	path, err := yamlpath.NewPath(jsonpath)
	if err != nil {
		return
	}

	return path.Find(file.node)
}

func (file *YAMLFile) Patch(jsonpath, value string) error {
	nodes, err := file.Find(jsonpath)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.Kind != yaml.ScalarNode || node.Style == yaml.LiteralStyle || node.Style == yaml.FoldedStyle {
			return errors.ErrUnsupported
		}

		if err := file.Replace(node, value); err != nil {
			return err
		}
	}

	return nil
}

func (file *YAMLFile) Replace(node *yaml.Node, value string) error {
	buffer := bytes.Buffer{}
	rewriter := NewRewriter(
		bytes.NewReader(file.content),
		&buffer,
	)

	decorated := node.Style == yaml.DoubleQuotedStyle || node.Style == yaml.SingleQuotedStyle

	if node.Style == yaml.DoubleQuotedStyle {
		value = fmt.Sprintf("\"%s\"", value)
	} else if node.Style == yaml.SingleQuotedStyle {
		value = fmt.Sprintf("'%s'", value)
	}

	// Copy content before yaml node
	rewriter.CopyLines(node.Line)
	rewriter.CopyBytes(node.Column)

	// Discard current value, including quotes
	rewriter.Discard(len(node.Value))
	if decorated {
		rewriter.Discard(2)
	}

	// Write new value
	rewriter.Write([]byte(value))

	// Copy the rest of the content
	rewriter.Copy()

	file.setContent(buffer.Bytes())
	return nil
}

func (file *YAMLFile) setContent(content []byte) error {
	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return err
	}

	file.content = content
	file.node = &node
	return nil
}
