package yatch

import "errors"

var (
	ErrTokenNotFound       = errors.New("token not found")
	ErrNodeNotFound        = errors.New("node not found")
	ErrNodeTypeUnsupported = errors.New("node type unsupported")
)
