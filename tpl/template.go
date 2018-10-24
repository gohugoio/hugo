// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/afero"

	"html/template"
	texttemplate "text/template"
	"text/template/parse"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/metrics"
	"github.com/pkg/errors"
)

var (
	_ TemplateExecutor = (*TemplateAdapter)(nil)
)

// TemplateHandler manages the collection of templates.
type TemplateHandler interface {
	TemplateFinder
	AddTemplate(name, tpl string) error
	AddLateTemplate(name, tpl string) error
	LoadTemplates(prefix string) error

	NewTextTemplate() TemplateParseFinder

	MarkReady()
	RebuildClone()
}

// TemplateFinder finds templates.
type TemplateFinder interface {
	Lookup(name string) (Template, bool)
}

// Template is the common interface between text/template and html/template.
type Template interface {
	Execute(wr io.Writer, data interface{}) error
	Name() string
}

// TemplateParser is used to parse ad-hoc templates, e.g. in the Resource chain.
type TemplateParser interface {
	Parse(name, tpl string) (Template, error)
}

// TemplateParseFinder provides both parsing and finding.
type TemplateParseFinder interface {
	TemplateParser
	TemplateFinder
}

// TemplateExecutor adds some extras to Template.
type TemplateExecutor interface {
	Template
	ExecuteToString(data interface{}) (string, error)
	Tree() string
}

// TemplateDebugger prints some debug info to stdoud.
type TemplateDebugger interface {
	Debug()
}

// TemplateAdapter implements the TemplateExecutor interface.
type TemplateAdapter struct {
	Template
	Metrics metrics.Provider

	// The filesystem where the templates are stored.
	Fs afero.Fs

	// Maps to base template if relevant.
	NameBaseTemplateName map[string]string
}

var baseOfRe = regexp.MustCompile("template: (.*?):")

func extractBaseOf(err string) string {
	m := baseOfRe.FindStringSubmatch(err)
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

// Execute executes the current template. The actual execution is performed
// by the embedded text or html template, but we add an implementation here so
// we can add a timer for some metrics.
func (t *TemplateAdapter) Execute(w io.Writer, data interface{}) (execErr error) {
	defer func() {
		// Panics in templates are a little bit too common (nil pointers etc.)
		// See https://github.com/gohugoio/hugo/issues/5327
		if r := recover(); r != nil {
			execErr = t.addFileContext(t.Name(), fmt.Errorf(`panic in Execute: %s. See "https://github.com/gohugoio/hugo/issues/5327" for the reason why we cannot provide a better error message for this.`, r))
		}
	}()

	if t.Metrics != nil {
		defer t.Metrics.MeasureSince(t.Name(), time.Now())
	}

	execErr = t.Template.Execute(w, data)
	if execErr != nil {
		execErr = t.addFileContext(t.Name(), execErr)
	}

	return
}

// The identifiers may be truncated in the log, e.g.
// "executing "main" at <$scaled.SRelPermalin...>: can't evaluate field SRelPermalink in type *resource.Image"
var identifiersRe = regexp.MustCompile("at \\<(.*?)(\\.{3})?\\>:")

func (t *TemplateAdapter) extractIdentifiers(line string) []string {
	m := identifiersRe.FindAllStringSubmatch(line, -1)
	identifiers := make([]string, len(m))
	for i := 0; i < len(m); i++ {
		identifiers[i] = m[i][1]
	}
	return identifiers
}

func (t *TemplateAdapter) addFileContext(name string, inerr error) error {
	if strings.HasPrefix(t.Name(), "_internal") {
		return inerr
	}

	f, realFilename, err := t.fileAndFilename(t.Name())
	if err != nil {
		return inerr

	}
	defer f.Close()

	master, hasMaster := t.NameBaseTemplateName[name]

	ferr := errors.Wrap(inerr, "execute of template failed")

	// Since this can be a composite of multiple template files (single.html + baseof.html etc.)
	// we potentially need to look in both -- and cannot rely on line number alone.
	lineMatcher := func(m herrors.LineMatcher) bool {
		if m.FileError.LineNumber() != m.LineNumber {
			return false
		}
		if !hasMaster {
			return true
		}

		identifiers := t.extractIdentifiers(m.FileError.Error())

		for _, id := range identifiers {
			if strings.Contains(m.Line, id) {
				return true
			}
		}
		return false
	}

	fe, ok := herrors.WithFileContext(ferr, realFilename, f, lineMatcher)
	if ok || !hasMaster {
		return fe
	}

	// Try the base template if relevant
	f, realFilename, err = t.fileAndFilename(master)
	if err != nil {
		return err
	}
	defer f.Close()

	fe, ok = herrors.WithFileContext(ferr, realFilename, f, lineMatcher)

	if !ok {
		// Return the most specific.
		return ferr

	}
	return fe

}

func (t *TemplateAdapter) fileAndFilename(name string) (afero.File, string, error) {
	fs := t.Fs
	filename := filepath.FromSlash(name)

	fi, err := fs.Stat(filename)
	if err != nil {
		return nil, "", err
	}
	f, err := fs.Open(filename)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to open template file %q:", filename)
	}

	return f, fi.(hugofs.RealFilenameInfo).RealFilename(), nil
}

// ExecuteToString executes the current template and returns the result as a
// string.
func (t *TemplateAdapter) ExecuteToString(data interface{}) (string, error) {
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	if err := t.Execute(b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Tree returns the template Parse tree as a string.
// Note: this isn't safe for parallel execution on the same template
// vs Lookup and Execute.
func (t *TemplateAdapter) Tree() string {
	var tree *parse.Tree
	switch tt := t.Template.(type) {
	case *template.Template:
		tree = tt.Tree
	case *texttemplate.Template:
		tree = tt.Tree
	default:
		panic("Unknown template")
	}

	if tree == nil || tree.Root == nil {
		return ""
	}
	s := tree.Root.String()

	return s
}

// TemplateFuncsGetter allows to get a map of functions.
type TemplateFuncsGetter interface {
	GetFuncs() map[string]interface{}
}

// TemplateTestMocker adds a way to override some template funcs during tests.
// The interface is named so it's not used in regular application code.
type TemplateTestMocker interface {
	SetFuncs(funcMap map[string]interface{})
}
