package tpl

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Some tests for Issue #1178 -- Ace
func TestAceTemplates(t *testing.T) {

	for i, this := range []struct {
		basePath     string
		innerPath    string
		baseContent  string
		innerContent string
		expect       string
		expectErr    int
	}{
		{"", filepath.FromSlash("_default/single.ace"), "", "{{ . }}", "DATA", 0},
		{filepath.FromSlash("_default/baseof.ace"), filepath.FromSlash("_default/single.ace"),
			`= content main
  h2 This is a content named "main" of an inner template. {{ . }}`,
			`= doctype html
html lang=en
  head
    meta charset=utf-8
    title Base and Inner Template
  body
    h1 This is a base template {{ . }}
    = yield main`, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Base and Inner Template</title></head><body><h1>This is a base template DATA</h1></body></html>`, 0},
	} {

		for _, root := range []string{"", os.TempDir()} {

			templ := New()

			basePath := this.basePath
			innerPath := this.innerPath

			if basePath != "" && root != "" {
				basePath = filepath.Join(root, basePath)
			}

			if innerPath != "" && root != "" {
				innerPath = filepath.Join(root, innerPath)
			}

			d := "DATA"

			err := templ.AddAceTemplate("mytemplate.ace", basePath, innerPath,
				[]byte(this.baseContent), []byte(this.innerContent))

			if err != nil && this.expectErr == 0 {
				t.Errorf("Test %d with root '%s' errored: %s", i, root, err)
			} else if err == nil && this.expectErr == 1 {
				t.Errorf("#1 Test %d with root '%s' should have errored", i, root)
			}

			var buff bytes.Buffer
			err = templ.ExecuteTemplate(&buff, "mytemplate.html", d)

			if err != nil && this.expectErr == 0 {
				t.Errorf("Test %d with root '%s' errored: %s", i, root, err)
			} else if err == nil && this.expectErr == 2 {
				t.Errorf("#2 Test with root '%s' %d should have errored", root, i)
			} else {
				result := buff.String()
				if result != this.expect {
					t.Errorf("Test %d  with root '%s' got\n%s\nexpected\n%s", i, root, result, this.expect)
				}
			}
		}
	}

}

// Test for bugs discovered by https://github.com/dvyukov/go-fuzz
func TestTplGoFuzzReports(t *testing.T) {

	// The following test case(s) also fail
	// See https://github.com/golang/go/issues/10634
	//{"{{ seq 433937734937734969526500969526500 }}", 2}}

	for i, this := range []struct {
		data      string
		expectErr int
	}{
		// Issue #1089
		//{"{{apply .C \"first\" }}", 2},
		// Issue #1090
		{"{{ slicestr \"000000\" 10}}", 2},
		// Issue #1091
		//{"{{apply .C \"first\" 0 0 0}}", 2},
		{"{{seq 3e80}}", 2},
		// Issue #1095
		{"{{apply .C \"urlize\" " +
			"\".\"}}", 2}} {
		templ := New()

		d := &Data{
			A: 42,
			B: "foo",
			C: []int{1, 2, 3},
			D: map[int]string{1: "foo", 2: "bar"},
			E: Data1{42, "foo"},
			F: []string{"a", "b", "c"},
			G: []string{"a", "b", "c", "d", "e"},
			H: "a,b,c,d,e,f",
		}

		err := templ.AddTemplate("fuzz", this.data)

		if err != nil && this.expectErr == 0 {
			t.Fatalf("Test %d errored: %s", i, err)
		} else if err == nil && this.expectErr == 1 {
			t.Fatalf("#1 Test %d should have errored", i)
		}

		err = templ.ExecuteTemplate(ioutil.Discard, "fuzz", d)

		if err != nil && this.expectErr == 0 {
			t.Fatalf("Test %d errored: %s", i, err)
		} else if err == nil && this.expectErr == 2 {
			t.Fatalf("#2 Test %d should have errored", i)
		}
	}
}

type Data struct {
	A int
	B string
	C []int
	D map[int]string
	E Data1
	F []string
	G []string
	H string
}

type Data1 struct {
	A int
	B string
}

func (Data1) Q() string {
	return "foo"
}

func (Data1) W() (string, error) {
	return "foo", nil
}

func (Data1) E() (string, error) {
	return "foo", errors.New("Data.E error")
}

func (Data1) R(v int) (string, error) {
	return "foo", nil
}

func (Data1) T(s string) (string, error) {
	return s, nil
}
