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
  - 'item11'
  - item222
object:
  # comments
  key0: true
  key1: 21
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

	t.Run("Patch List", func(t *testing.T) {
		if err := file.Patch("$.list[*]", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "7ffca437780e24301be4b48f278da04f" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
			return
		}
	})

	t.Run("Patch Object", func(t *testing.T) {
		if err := file.Patch("$.object.key2", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "b25063e19e6a931dc5d67af595d2c054" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
			return
		}

	})

	t.Run("Patch Json", func(t *testing.T) {
		if err := file.Patch("$.json.key", "patched"); err != nil {
			t.Error(err)
			return
		}

		t.Log(string(file.Content()))
		if checksum := checksum(file.Content()); checksum != "466703bb0d0b985fc9cbc94c06b980d0" {
			t.Errorf("checksum mismatch, expected: %s", checksum)
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
}

func checksum(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}
