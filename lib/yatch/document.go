package yatch

import (
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/token"
)

// Document is a single document of a YAMLFile. It holds an index rather than a node because every patch re-parses the file.
type Document struct {
	file  *File
	index int
}

// Documents returns every document in the file, in the order the goccy-yaml parser reports them.
func (file *File) Documents() []*Document {
	docs := make([]*Document, len(file.file.Docs))
	for index := range docs {
		docs[index] = &Document{file: file, index: index}
	}

	return docs
}

// FirstDocument returns the first document, skipping %YAML and %TAG directives that goccy-yaml parses as separate documents.
func (file *File) FirstDocument() *Document {
	index := slices.IndexFunc(file.file.Docs, func(doc *ast.DocumentNode) bool {
		return doc.Body == nil || doc.Body.Type() != ast.DirectiveType
	})

	return &Document{file: file, index: max(index, 0)}
}

func (doc *Document) Find(jsonpath string) (ast.Node, error) {
	path, err := yaml.PathString(jsonpath)
	if err != nil {
		return nil, err
	}

	node, err := path.FilterNode(doc.file.file.Docs[doc.index].Body)
	if err != nil {
		return nil, err
	}

	switch node := node.(type) {
	case nil:
		return nil, ErrNodeNotFound
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

func (doc *Document) Patch(jsonpath, value string) error {
	node, err := doc.Find(jsonpath)
	if err != nil {
		return err
	}

	token := (*Token)(node.GetToken())
	return doc.file.Replace(token, value)
}

func (doc *Document) PatchComment(jsonpath, value string) error {
	node, err := doc.Find(jsonpath)
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
	return doc.file.Replace((*Token)(next), value)
}
