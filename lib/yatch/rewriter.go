package yatch

import (
	"bufio"
	"io"
)

type (
	Rewriter interface {
		CopyLines(n int) error
		CopyBytes(n int) error
		CopyAll() error
		Discard(n int) error
		Write(write []byte) error
		WriteString(write string) error
	}

	rewriter struct {
		reader *bufio.Reader
		writer io.Writer
	}
)

func NewRewriter(reader io.Reader, writer io.Writer) Rewriter {
	return &rewriter{
		reader: bufio.NewReader(reader),
		writer: writer,
	}
}

func (rw *rewriter) CopyAll() error {
	if _, err := rw.reader.WriteTo(rw.writer); err != nil {
		return err
	}
	return nil
}

func (rw *rewriter) CopyLines(n int) error {
	for i := 0; i < n-1; i++ {
		bytes, err := rw.reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		_, err = rw.writer.Write(bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rw *rewriter) CopyBytes(n int) error {
	prefix := make([]byte, n-1)
	if _, err := rw.reader.Read(prefix); err != nil {
		return err
	}
	if _, err := rw.writer.Write(prefix); err != nil {
		return err
	}

	return nil
}

func (rw *rewriter) Discard(n int) error {
	if _, err := rw.reader.Discard(n); err != nil {
		return err
	}
	return nil
}

func (rw *rewriter) Write(bytes []byte) error {
	if _, err := rw.writer.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (rw *rewriter) WriteString(write string) error {
	return rw.Write([]byte(write))
}
