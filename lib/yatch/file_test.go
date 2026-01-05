package yatch_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
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

	t.Run("Patch List", func(t *testing.T) {
		if err := file.Patch("$.list[*]", "patched"); err != nil {
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
		if err := file.Patch("$.object.key2", "patched"); err != nil {
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
		if err := file.Patch("$.json.key", "patched"); err != nil {
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
		err := file.Patch("$.object", "patched")
		if err == nil {
			t.Error("Must return and error for maps")
			return
		}
	})

	t.Run("Patch Literal Multiline", func(t *testing.T) {
		err := file.Patch("$.multiline.line0", "patched")
		if err == nil {
			t.Error("Must return and error for literal value")
			return
		}
	})

	t.Run("Patch Folded Multiline", func(t *testing.T) {
		err := file.Patch("$.multiline.line1", "patched")
		if err == nil {
			t.Error("Must return and error for folded value")
			return
		}
	})

	t.Run("Patch Comment Multiline", func(t *testing.T) {
		if err := file.PatchComment("$.key", "patched"); err != nil {
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

func checksum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
