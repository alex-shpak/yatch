package yatch

import "errors"

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrNodeNotFound        = errors.New("node not found")
	ErrNodeTypeUnsupported = errors.New("node type unsupported")
	ErrValueMultiline      = errors.New("value must be a single line")
)
