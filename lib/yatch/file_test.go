package yatch_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/alex-shpak/yatch/lib/yatch"
)

const inputFile = `
key: value
list:
  - "item0"
  - item11
  - item222
object:
  # comments
  key0: value0
  key1: value11
  key2:
    value222
    # comments
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

	t.Run("Find", func(t *testing.T) {
		nodes, err := file.Find("list[*]")
		if err != nil {
			t.Error(err)
			return
		}

		if len(nodes) != 3 {
			t.Errorf("expected 3 nodes in 'list', got %d", len(nodes))
			return
		}
	})

	t.Run("Patch List", func(t *testing.T) {
		if err := file.Patch("list[*]", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "de52c1b69ffd3ac8eeecb9870b301fdf" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
			return
		}
	})

	t.Run("Patch Object", func(t *testing.T) {
		if err := file.Patch("object.key2", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "027febd0ee840a0e6189ad97e137548a" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
			return
		}

	})

	t.Run("Patch Json", func(t *testing.T) {
		if err := file.Patch("json.value", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "027febd0ee840a0e6189ad97e137548a" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
			return
		}

	})

	t.Run("Patch Slice", func(t *testing.T) {
		err := file.Patch("list", "patched")
		if err == nil {
			t.Error("Must return and error for lists")
			return
		}
	})

	t.Run("Patch Map", func(t *testing.T) {
		err := file.Patch("object", "patched")
		if err == nil {
			t.Error("Must return and error for maps")
			return
		}
	})

	t.Run("Patch Literal Multiline", func(t *testing.T) {
		err := file.Patch("multiline.line0", "patched")
		if err == nil {
			t.Error("Must return and error for literal value")
			return
		}
	})

	t.Run("Patch Folded Multiline", func(t *testing.T) {
		err := file.Patch("multiline.line1", "patched")
		if err == nil {
			t.Error("Must return and error for folded value")
			return
		}
	})
}

func checksum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
