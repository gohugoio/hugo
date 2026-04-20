// Copyright 2019 The Hugo Authors. All rights reserved.
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

package rst

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config/security"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

func TestGetRstArgs(t *testing.T) {
	c := qt.New(t)

	c.Assert(getRstArgs("/ignored/path/rst2html.py", false, true), qt.DeepEquals, []string{
		"--leave-comments",
		"--initial-header-level=2",
		rst2ShortSyntaxHighlightArg,
	})

	c.Assert(getRstArgs("/ignored/path/rst2html.py", false, false), qt.DeepEquals, []string{
		"--leave-comments",
		"--initial-header-level=2",
	})

	c.Assert(getRstArgs(`C:\Python\Scripts\rst2html.py`, true, true), qt.DeepEquals, []string{
		`C:\Python\Scripts\rst2html.py`,
		"--leave-comments",
		"--initial-header-level=2",
		rst2ShortSyntaxHighlightArg,
	})
}

func TestGetRstContentUsesShortSyntaxHighlight(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("scripted rst2html test not portable to Windows")
	}

	c := qt.New(t)
	workDir := t.TempDir()
	rstBinary := filepath.Join(workDir, "rst2html")

	err := os.WriteFile(rstBinary, []byte(`#!/bin/sh
printf '<html>\n<body>\n%s\n</body>\n</html>\n' "$*"
`), 0o755)
	c.Assert(err, qt.IsNil)

	t.Setenv("PATH", workDir)

	sc := security.DefaultConfig
	sc.Exec.Allow = security.MustNewWhitelist("^rst2html$")

	conv := &rstConverter{
		ctx: converter.DocumentContext{DocumentName: "test.rst"},
		cfg: converter.ProviderConfig{
			Logger: loggers.NewDefault(),
			Exec:   hexec.New(sc, "", loggers.NewDefault()),
		},
	}

	b, err := conv.getRstContent([]byte("ignored"), conv.ctx)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, strings.Join(append(append([]string(nil), rst2BaseArgs...), rst2ShortSyntaxHighlightArg), " "))
}

func TestGetRstContentFallbackWithoutShortSyntaxHighlight(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("scripted rst2html test not portable to Windows")
	}

	c := qt.New(t)
	workDir := t.TempDir()
	rstBinary := filepath.Join(workDir, "rst2html")

	err := os.WriteFile(rstBinary, []byte(`#!/bin/sh
case "$*" in
  *"--syntax-highlight=short --help"*)
    printf 'Usage: rst2html\n  --syntax-highlight=short\n'
    exit 0
    ;;
  *"--syntax-highlight=short"*)
    exit 0
    ;;
esac
printf '<html>\n<body>\n%s\n</body>\n</html>\n' "$*"
`), 0o755)
	c.Assert(err, qt.IsNil)

	t.Setenv("PATH", workDir)

	sc := security.DefaultConfig
	sc.Exec.Allow = security.MustNewWhitelist("^rst2html$")

	conv := &rstConverter{
		ctx: converter.DocumentContext{DocumentName: "test.rst"},
		cfg: converter.ProviderConfig{
			Logger: loggers.NewDefault(),
			Exec:   hexec.New(sc, "", loggers.NewDefault()),
		},
	}

	b, err := conv.getRstContent([]byte("ignored"), conv.ctx)
	c.Assert(err, qt.IsNil)
	c.Assert(string(b), qt.Equals, strings.Join(rst2BaseArgs, " "))
}

func TestGetRstContentUnsupportedShortSyntaxHighlight(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("scripted rst2html test not portable to Windows")
	}

	c := qt.New(t)
	workDir := t.TempDir()
	rstBinary := filepath.Join(workDir, "rst2html")

	err := os.WriteFile(rstBinary, []byte(`#!/bin/sh
case "$*" in
  *"--syntax-highlight=short"*)
    printf 'rst2html: error: unrecognized arguments: --syntax-highlight=short\n' >&2
    exit 2
    ;;
esac
printf '<html>\n<body>\n%s\n</body>\n</html>\n' "$*"
`), 0o755)
	c.Assert(err, qt.IsNil)

	t.Setenv("PATH", workDir)

	sc := security.DefaultConfig
	sc.Exec.Allow = security.MustNewWhitelist("^rst2html$")

	conv := &rstConverter{
		ctx: converter.DocumentContext{DocumentName: "test.rst"},
		cfg: converter.ProviderConfig{
			Logger: loggers.NewDefault(),
			Exec:   hexec.New(sc, "", loggers.NewDefault()),
		},
	}

	_, err = conv.getRstContent([]byte("ignored"), conv.ctx)
	c.Assert(err, qt.ErrorMatches, ".*does not support --syntax-highlight=short; please upgrade Docutils \\(0\\.13\\+ required\\).*")
}

func TestConvert(t *testing.T) {
	if !Supports() {
		t.Skip("rst not installed")
	}
	c := qt.New(t)
	sc := security.DefaultConfig
	sc.Exec.Allow = security.MustNewWhitelist("rst", "python")

	p, err := Provider.New(
		converter.ProviderConfig{
			Logger: loggers.NewDefault(),
			Exec:   hexec.New(sc, "", loggers.NewDefault()),
		})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<div class=\"document\">\n\n\n<p>testContent</p>\n</div>")
}
