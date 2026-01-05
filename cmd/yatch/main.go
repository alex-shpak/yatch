package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alex-shpak/yatch/lib/yatch"
)

func main() {
	var (
		in      = flag.String("in", "", "input YAML file (default: stdin)")
		out     = flag.String("out", "", "output YAML file (default: stdout)")
		comment = flag.Bool("comment", false, "update a comment instead of the node value")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s [options] <jsonpath> <value>

Patch YAML files using JSONPath expressions.

Options:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  # Update from stdin to stdout
  %[1]s '$.image.tag' 'v1.2.3' < input.yaml > output.yaml

  # Update files in-place
  %[1]s -in values.yaml -out values.yaml '$.image.tag' '2.0'

  # Update a comment
  %[1]s --comment '$.service.port' 'Exposed port' < values.yaml
`, os.Args[0])
	}

	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := exec(*in, *out, *comment, flag.Arg(0), flag.Arg(1)); err != nil {
		slog.Error("Error", "err", err)
		os.Exit(1)
	}
}

func exec(in, out string, comment bool, jsonpath, value string) error {
	reader, err := reader(in)
	if err != nil {
		return fmt.Errorf("error opening input: %w", err)
	}
	defer reader.Close()

	yamlfile, err := yatch.NewFile(reader)
	if err != nil {
		return fmt.Errorf("error parsing input content: %w", err)
	}

	if comment {
		err = yamlfile.PatchComment(jsonpath, value)
	} else {
		err = yamlfile.Patch(jsonpath, value)
	}

	if err != nil {
		return fmt.Errorf("error while patching: %w", err)
	}

	writer, err := writer(out)
	if err != nil {
		return fmt.Errorf("error opening output: %w", err)
	}
	defer writer.Close()

	if _, err = writer.Write(yamlfile.Content()); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}

func reader(in string) (io.ReadCloser, error) {
	if in == "" {
		return os.Stdin, nil
	}
	return os.Open(in)
}

func writer(out string) (io.WriteCloser, error) {
	if out == "" {
		return os.Stdout, nil
	}
	return os.Create(out)
}
