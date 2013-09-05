package source

import (
	"bytes"
	"testing"
)

func TestEmptySourceFilesystem(t *testing.T) {
	src := new(Filesystem)
	if len(src.Files()) != 0 {
		t.Errorf("new filesystem should contain 0 files.")
	}
}

func TestAddFile(t *testing.T) {
	src := new(Filesystem)
	src.add("foobar", bytes.NewReader([]byte("aaa")))
	if len(src.Files()) != 1 {
		t.Errorf("Files() should return 1 file")
	}

	f := src.Files()[0]
	if f.Name != "foobar" {
		t.Errorf("File name should be 'foobar', got: %s", f.Name)
	}

	b := new(bytes.Buffer)
	b.ReadFrom(f.Contents)
	if b.String() != "aaa" {
		t.Errorf("File contents should be 'aaa', got: %s", b.String())
	}
}
