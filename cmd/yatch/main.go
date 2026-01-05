package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alex-shpak/yatch/lib/yatch"
)

var (
	in, out *string
	comment *bool
)

func main() {
	in = flag.String("in", "", "input file")
	out = flag.String("out", "", "output file")
	comment = flag.Bool("comment", false, "update a comment attached to the node")

	flag.Parse()

	if len(flag.Args()) != 2 {
		slog.Error("Usage: yatch [-in <input file>] [-out <output file>] [--comment] <jsonpath> <value>")
		os.Exit(1)
	}

	if err := exec(); err != nil {
		slog.Error("Fatal", "err", err)
		os.Exit(1)
	}
}

func exec() error {
	jsonpath, value := flag.Arg(0), flag.Arg(1)

	reader, err := reader()
	if err != nil {
		return fmt.Errorf("error opening input, %w", err)
	}

	writer, err := writer()
	if err != nil {
		return fmt.Errorf("error opening output, %w", err)
	}

	yamlfile, err := yatch.NewFile(reader)
	if err != nil {
		return fmt.Errorf("error parsing input content, %w", err)
	}

	if *comment {
		err = yamlfile.PatchComment(jsonpath, value)
	} else {
		err = yamlfile.Patch(jsonpath, value)
	}

	if err != nil {
		return fmt.Errorf("error patching file, %w", err)
	}

	_, err = writer.Write(yamlfile.Content())
	if err != nil {
		return fmt.Errorf("error writing output, %w", err)
	}

	return nil
}

func reader() (reader io.ReadCloser, err error) {
	if *in == "" {
		return os.Stdin, nil
	}

	return os.Open(*in)
}

func writer() (writer io.WriteCloser, err error) {
	if *out == "" {
		return os.Stdout, nil
	}

	return os.Create(*out)
}
