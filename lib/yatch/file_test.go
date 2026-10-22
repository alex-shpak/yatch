package yatch_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/alex-shpak/yatch/lib/yatch"
)

const inputFile = `
key: value # comment
list:
  - "item0" # comment
  - 'item11'
  - item222
object:
  # comments
  key0: true
  key1: 21
  key2:
    value222 # comment
json: { "key": "value" }
multiline:
  line0: |
    This is a multiline string
    that spans multiple lines
    This is a multiline string
    that spans multiple lines

  line1: >
    This is a multiline string
    that spans multiple lines
  line2: normal string
`

func TestFile(t *testing.T) {
	file, err := yatch.NewFile(bytes.NewReader([]byte(inputFile)))
	if err != nil {
		t.Error(err)
		return
	}
	doc := file.FirstDocument()

	t.Run("Patch List", func(t *testing.T) {
		if err := doc.Patch("$.list[*]", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "f407ed29f9287e02d5527b8f38019348" {
			t.Errorf("checksum mismatch, got: %s", checksum)
			return
		}
	})

	t.Run("Patch Object", func(t *testing.T) {
		if err := doc.Patch("$.object.key2", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "123785f69652a177f40145cbed577270" {
			t.Errorf("checksum mismatch, got: %s", checksum)
			return
		}

	})

	t.Run("Patch Json", func(t *testing.T) {
		if err := doc.Patch("$.json.key", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "bcf753a12943203f6e69aee11e22f622" {
			t.Errorf("checksum mismatch, got: %s", checksum)
			return
		}
	})

	t.Run("Patch Map", func(t *testing.T) {
		err := doc.Patch("$.object", "patched")
		if err == nil {
			t.Error("Must return and error for maps")
			return
		}
	})

	t.Run("Patch Literal Multiline", func(t *testing.T) {
		err := doc.Patch("$.multiline.line0", "patched")
		if err == nil {
			t.Error("Must return and error for literal value")
			return
		}
	})

	t.Run("Patch Folded Multiline", func(t *testing.T) {
		err := doc.Patch("$.multiline.line1", "patched")
		if err == nil {
			t.Error("Must return and error for folded value")
			return
		}
	})

	t.Run("Patch Comment Multiline", func(t *testing.T) {
		if err := doc.PatchComment("$.key", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "c537445eba30451ffaf1f709af210660" {
			t.Errorf("checksum mismatch, got: %s", checksum)
			return
		}
	})
}

func TestDocuments(t *testing.T) {
	file, err := yatch.NewFile(bytes.NewReader([]byte("kind: A\nimage: a\n---\nkind: B\nimage: \"b\"\nonly: x\n")))
	if err != nil {
		t.Fatal(err)
	}
	docs := file.Documents()

	t.Run("Patch Document", func(t *testing.T) {
		if err := docs[1].Patch("$.image", "patched"); err != nil {
			t.Error(err)
			return
		}

		if want := "kind: A\nimage: a\n---\nkind: B\nimage: \"patched\"\nonly: x\n"; string(file.Content()) != want {
			t.Errorf("unexpected content:\n%s", file.Content())
		}
	})

	t.Run("Find Document", func(t *testing.T) {
		if node, err := docs[1].Find("$.image"); err != nil || node.GetToken().Value != "patched" {
			t.Errorf("expected 'patched', got: %v, %v", node, err)
		}
	})

	t.Run("Not Found", func(t *testing.T) {
		if _, err := docs[0].Find("$.only"); !errors.Is(err, yatch.ErrNodeNotFound) {
			t.Errorf("expected ErrNodeNotFound, got: %v", err)
		}
	})

	t.Run("Multiline Value", func(t *testing.T) {
		before := string(file.Content())
		if err := docs[0].Patch("$.image", "x\n---\nd: w"); !errors.Is(err, yatch.ErrValueMultiline) {
			t.Errorf("expected ErrValueMultiline, got: %v", err)
		}

		if string(file.Content()) != before {
			t.Errorf("content changed:\n%s", file.Content())
		}
	})

	t.Run("First Document Skips Directives", func(t *testing.T) {
		file, err := yatch.NewFile(bytes.NewReader([]byte("%YAML 1.2\n---\na: 1\n")))
		if err != nil {
			t.Error(err)
			return
		}

		if err := file.FirstDocument().Patch("$.a", "2"); err != nil || string(file.Content()) != "%YAML 1.2\n---\na: 2\n" {
			t.Errorf("unexpected result: %v\n%s", err, file.Content())
		}
	})
}

func checksum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
