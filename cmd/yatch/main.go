package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alex-shpak/yatch/lib/yatch"
)

func init() {
	flag.String("in", "", "input file")
	flag.String("out", "", "output file")

	flag.Parse()

	if len(flag.Args()) != 2 {
		slog.Error("Usage: yatch [-in <input file>] [-out <output file>] <jsonpath> <value>")
		os.Exit(1)
	}
}

func main() {
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

	err = yamlfile.Patch(jsonpath, value)
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
	inFlag := flag.Lookup("in")
	if inFlag.Value.String() == "" {
		return os.Stdin, nil
	}

	return os.Open(inFlag.Value.String())
}

func writer() (writer io.WriteCloser, err error) {
	outFlag := flag.Lookup("out")
	if outFlag.Value.String() == "" {
		return os.Stdout, nil
	}

	return os.Create(outFlag.Value.String())
}
