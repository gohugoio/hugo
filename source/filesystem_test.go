package source

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestEmptySourceFilesystem(t *testing.T) {
	src := new(Filesystem)
	if len(src.Files()) != 0 {
		t.Errorf("new filesystem should contain 0 files.")
	}
}

type TestPath struct {
	filename string
	logical  string
	content  string
	section  string
	dir      string
}

func TestAddFile(t *testing.T) {
	tests := platformPaths
	for _, test := range tests {
		base := platformBase
		srcDefault := new(Filesystem)
		srcWithBase := &Filesystem{
			Base: base,
		}

		for _, src := range []*Filesystem{srcDefault, srcWithBase} {

			p := test.filename
			if !filepath.IsAbs(test.filename) {
				p = filepath.Join(src.Base, test.filename)
			}

			if err := src.add(p, bytes.NewReader([]byte(test.content))); err != nil {
				if err.Error() == "source: missing base directory" {
					continue
				}
				t.Fatalf("%s add returned an error: %s", p, err)
			}

			if len(src.Files()) != 1 {
				t.Fatalf("%s Files() should return 1 file", p)
			}

			f := src.Files()[0]
			if f.LogicalName() != test.logical {
				t.Errorf("Filename (Base: %q) expected: %q, got: %q", src.Base, test.logical, f.LogicalName())
			}

			b := new(bytes.Buffer)
			b.ReadFrom(f.Contents)
			if b.String() != test.content {
				t.Errorf("File (Base: %q) contents should be %q, got: %q", src.Base, test.content, b.String())
			}

			if f.Section() != test.section {
				t.Errorf("File section (Base: %q) expected: %q, got: %q", src.Base, test.section, f.Section())
			}

			if f.Dir() != test.dir {
				t.Errorf("Dir path (Base: %q) expected: %q, got: %q", src.Base, test.dir, f.Dir())
			}
		}
	}
}
