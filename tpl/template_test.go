// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/hugofs"
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

func isAtLeastGo16() bool {
	version := runtime.Version()
	return strings.Contains(version, "1.6") || strings.Contains(version, "1.7")
}

func TestAddTemplateFileWithMaster(t *testing.T) {

	if !isAtLeastGo16() {
		t.Skip("This test only runs on Go >= 1.6")
	}

	for i, this := range []struct {
		masterTplContent  string
		overlayTplContent string
		writeSkipper      int
		expect            interface{}
	}{
		{`A{{block "main" .}}C{{end}}C`, `{{define "main"}}B{{end}}`, 0, "ABC"},
		{`A{{block "main" .}}C{{end}}C{{block "sub" .}}D{{end}}E`, `{{define "main"}}B{{end}}`, 0, "ABCDE"},
		{`A{{block "main" .}}C{{end}}C{{block "sub" .}}D{{end}}E`, `{{define "main"}}B{{end}}{{define "sub"}}Z{{end}}`, 0, "ABCZE"},
		{`tpl`, `tpl`, 1, false},
		{`tpl`, `tpl`, 2, false},
		{`{{.0.E}}`, `tpl`, 0, false},
		{`tpl`, `{{.0.E}}`, 0, false},
	} {

		hugofs.SourceFs = afero.NewMemMapFs()
		templ := New()
		overlayTplName := "ot"
		masterTplName := "mt"
		finalTplName := "tp"

		if this.writeSkipper != 1 {
			afero.WriteFile(hugofs.SourceFs, masterTplName, []byte(this.masterTplContent), 0644)
		}
		if this.writeSkipper != 2 {
			afero.WriteFile(hugofs.SourceFs, overlayTplName, []byte(this.overlayTplContent), 0644)
		}

		err := templ.AddTemplateFileWithMaster(finalTplName, overlayTplName, masterTplName)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] AddTemplateFileWithMaster didn't return an expected error", i)
			}
		} else {

			if err != nil {
				t.Errorf("[%d] AddTemplateFileWithMaster failed: %s", i, err)
				continue
			}

			resultTpl := templ.Lookup(finalTplName)

			if resultTpl == nil {
				t.Errorf("[%d] AddTemplateFileWithMaster: Result template not found", i)
				continue
			}

			var b bytes.Buffer
			err := resultTpl.Execute(&b, nil)

			if err != nil {
				t.Errorf("[%d] AddTemplateFileWithMaster execute failed: %s", i, err)
				continue
			}
			resultContent := b.String()

			if resultContent != this.expect {
				t.Errorf("[%d] AddTemplateFileWithMaster got \n%s but expected \n%v", i, resultContent, this.expect)
			}
		}

	}

}

// A Go stdlib test for linux/arm. Will remove later.
// See #1771
func TestBigIntegerFunc(t *testing.T) {
	var func1 = func(v int64) error {
		return nil
	}
	var funcs = map[string]interface{}{
		"A": func1,
	}

	tpl, err := template.New("foo").Funcs(funcs).Parse("{{ A 3e80 }}")
	if err != nil {
		t.Fatal("Parse failed:", err)
	}
	err = tpl.Execute(ioutil.Discard, "foo")

	if err == nil {
		t.Fatal("Execute should have failed")
	}

	t.Log("Got expected error:", err)

}

// A Go stdlib test for linux/arm. Will remove later.
// See #1771
type BI struct {
}

func (b BI) A(v int64) error {
	return nil
}
func TestBigIntegerMethod(t *testing.T) {

	data := &BI{}

	tpl, err := template.New("foo2").Parse("{{ .A 3e80 }}")
	if err != nil {
		t.Fatal("Parse failed:", err)
	}
	err = tpl.ExecuteTemplate(ioutil.Discard, "foo2", data)

	if err == nil {
		t.Fatal("Execute should have failed")
	}

	t.Log("Got expected error:", err)

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
