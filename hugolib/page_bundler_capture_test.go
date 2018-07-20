// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/gohugoio/hugo/common/loggers"

	jww "github.com/spf13/jwalterweatherman"

	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/source"
	"github.com/stretchr/testify/require"
)

type storeFilenames struct {
	sync.Mutex
	filenames []string
	copyNames []string
	dirKeys   []string
}

func (s *storeFilenames) handleSingles(fis ...*fileInfo) {
	s.Lock()
	defer s.Unlock()
	for _, fi := range fis {
		s.filenames = append(s.filenames, filepath.ToSlash(fi.Filename()))
	}
}

func (s *storeFilenames) handleBundles(d *bundleDirs) {
	s.Lock()
	defer s.Unlock()
	var keys []string
	for _, b := range d.bundles {
		res := make([]string, len(b.resources))
		i := 0
		for _, r := range b.resources {
			res[i] = path.Join(r.Lang(), filepath.ToSlash(r.Filename()))
			i++
		}
		sort.Strings(res)
		keys = append(keys, path.Join("__bundle", b.fi.Lang(), filepath.ToSlash(b.fi.Filename()), "resources", strings.Join(res, "|")))
	}
	s.dirKeys = append(s.dirKeys, keys...)
}

func (s *storeFilenames) handleCopyFiles(files ...pathLangFile) {
	s.Lock()
	defer s.Unlock()
	for _, file := range files {
		s.copyNames = append(s.copyNames, filepath.ToSlash(file.Filename()))
	}
}

func (s *storeFilenames) sortedStr() string {
	s.Lock()
	defer s.Unlock()
	sort.Strings(s.filenames)
	sort.Strings(s.dirKeys)
	sort.Strings(s.copyNames)
	return "\nF:\n" + strings.Join(s.filenames, "\n") + "\nD:\n" + strings.Join(s.dirKeys, "\n") +
		"\nC:\n" + strings.Join(s.copyNames, "\n") + "\n"
}

func TestPageBundlerCaptureSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" && os.Getenv("CI") == "" {
		t.Skip("Skip TestPageBundlerCaptureSymlinks as os.Symlink needs administrator rights on Windows")
	}

	assert := require.New(t)
	ps, clean, workDir := newTestBundleSymbolicSources(t)
	sourceSpec := source.NewSourceSpec(ps, ps.BaseFs.Content.Fs)
	defer clean()

	fileStore := &storeFilenames{}
	logger := loggers.NewErrorLogger()
	c := newCapturer(logger, sourceSpec, fileStore, nil)

	assert.NoError(c.capture())

	// Symlink back to content skipped to prevent infinite recursion.
	assert.Equal(uint64(3), logger.LogCountForLevelsGreaterThanorEqualTo(jww.LevelWarn))

	expected := `
F:
/base/a/page_s.md
/base/a/regular.md
/base/symbolic1/s1.md
/base/symbolic1/s2.md
/base/symbolic3/circus/a/page_s.md
/base/symbolic3/circus/a/regular.md
D:
__bundle/en/base/symbolic2/a1/index.md/resources/en/base/symbolic2/a1/logo.png|en/base/symbolic2/a1/page.md
C:
/base/symbolic3/s1.png
/base/symbolic3/s2.png
`

	got := strings.Replace(fileStore.sortedStr(), filepath.ToSlash(workDir), "", -1)
	got = strings.Replace(got, "//", "/", -1)

	if expected != got {
		diff := helpers.DiffStringSlices(strings.Fields(expected), strings.Fields(got))
		t.Log(got)
		t.Fatalf("Failed:\n%s", diff)
	}
}

func TestPageBundlerCaptureBasic(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	fs, cfg := newTestBundleSources(t)
	assert.NoError(loadDefaultSettingsFor(cfg))
	assert.NoError(loadLanguageSettings(cfg, nil))
	ps, err := helpers.NewPathSpec(fs, cfg)
	assert.NoError(err)

	sourceSpec := source.NewSourceSpec(ps, ps.BaseFs.Content.Fs)

	fileStore := &storeFilenames{}

	c := newCapturer(loggers.NewErrorLogger(), sourceSpec, fileStore, nil)

	assert.NoError(c.capture())

	expected := `
F:
/work/base/_1.md
/work/base/a/1.md
/work/base/a/2.md
/work/base/assets/pages/mypage.md
D:
__bundle/en/work/base/_index.md/resources/en/work/base/_1.png
__bundle/en/work/base/a/b/index.md/resources/en/work/base/a/b/ab1.md
__bundle/en/work/base/b/my-bundle/index.md/resources/en/work/base/b/my-bundle/1.md|en/work/base/b/my-bundle/2.md|en/work/base/b/my-bundle/c/logo.png|en/work/base/b/my-bundle/custom-mime.bep|en/work/base/b/my-bundle/sunset1.jpg|en/work/base/b/my-bundle/sunset2.jpg
__bundle/en/work/base/c/bundle/index.md/resources/en/work/base/c/bundle/logo-은행.png
__bundle/en/work/base/root/index.md/resources/en/work/base/root/1.md|en/work/base/root/c/logo.png
C:
/work/base/assets/pic1.png
/work/base/assets/pic2.png
/work/base/images/hugo-logo.png
`

	got := fileStore.sortedStr()

	if expected != got {
		diff := helpers.DiffStringSlices(strings.Fields(expected), strings.Fields(got))
		t.Log(got)
		t.Fatalf("Failed:\n%s", diff)
	}
}

func TestPageBundlerCaptureMultilingual(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	fs, cfg := newTestBundleSourcesMultilingual(t)
	assert.NoError(loadDefaultSettingsFor(cfg))
	assert.NoError(loadLanguageSettings(cfg, nil))

	ps, err := helpers.NewPathSpec(fs, cfg)
	assert.NoError(err)

	sourceSpec := source.NewSourceSpec(ps, ps.BaseFs.Content.Fs)
	fileStore := &storeFilenames{}
	c := newCapturer(loggers.NewErrorLogger(), sourceSpec, fileStore, nil)

	assert.NoError(c.capture())

	expected := `
F:
/work/base/1s/mypage.md
/work/base/1s/mypage.nn.md
/work/base/bb/_1.md
/work/base/bb/_1.nn.md
/work/base/bb/en.md
/work/base/bc/page.md
/work/base/bc/page.nn.md
/work/base/be/_index.md
/work/base/be/page.md
/work/base/be/page.nn.md
D:
__bundle/en/work/base/bb/_index.md/resources/en/work/base/bb/a.png|en/work/base/bb/b.png|nn/work/base/bb/c.nn.png
__bundle/en/work/base/bc/_index.md/resources/en/work/base/bc/logo-bc.png
__bundle/en/work/base/bd/index.md/resources/en/work/base/bd/page.md
__bundle/en/work/base/bf/my-bf-bundle/index.md/resources/en/work/base/bf/my-bf-bundle/page.md
__bundle/en/work/base/lb/index.md/resources/en/work/base/lb/1.md|en/work/base/lb/2.md|en/work/base/lb/c/d/deep.png|en/work/base/lb/c/logo.png|en/work/base/lb/c/one.png|en/work/base/lb/c/page.md
__bundle/nn/work/base/bb/_index.nn.md/resources/en/work/base/bb/a.png|nn/work/base/bb/b.nn.png|nn/work/base/bb/c.nn.png
__bundle/nn/work/base/bd/index.md/resources/nn/work/base/bd/page.nn.md
__bundle/nn/work/base/bf/my-bf-bundle/index.nn.md/resources
__bundle/nn/work/base/lb/index.nn.md/resources/en/work/base/lb/c/d/deep.png|en/work/base/lb/c/one.png|nn/work/base/lb/2.nn.md|nn/work/base/lb/c/logo.nn.png
C:
/work/base/1s/mylogo.png
/work/base/bb/b/d.nn.png
`

	got := fileStore.sortedStr()

	if expected != got {
		diff := helpers.DiffStringSlices(strings.Fields(expected), strings.Fields(got))
		t.Log(got)
		t.Fatalf("Failed:\n%s", strings.Join(diff, "\n"))
	}

}

type noOpFileStore int

func (noOpFileStore) handleSingles(fis ...*fileInfo)        {}
func (noOpFileStore) handleBundles(b *bundleDirs)           {}
func (noOpFileStore) handleCopyFiles(files ...pathLangFile) {}

func BenchmarkPageBundlerCapture(b *testing.B) {
	capturers := make([]*capturer, b.N)

	for i := 0; i < b.N; i++ {
		cfg, fs := newTestCfg()
		ps, _ := helpers.NewPathSpec(fs, cfg)
		sourceSpec := source.NewSourceSpec(ps, fs.Source)

		base := fmt.Sprintf("base%d", i)
		for j := 1; j <= 5; j++ {
			js := fmt.Sprintf("j%d", j)
			writeSource(b, fs, filepath.Join(base, js, "index.md"), "content")
			writeSource(b, fs, filepath.Join(base, js, "logo1.png"), "content")
			writeSource(b, fs, filepath.Join(base, js, "sub", "logo2.png"), "content")
			writeSource(b, fs, filepath.Join(base, js, "section", "_index.md"), "content")
			writeSource(b, fs, filepath.Join(base, js, "section", "logo.png"), "content")
			writeSource(b, fs, filepath.Join(base, js, "section", "sub", "logo.png"), "content")

			for k := 1; k <= 5; k++ {
				ks := fmt.Sprintf("k%d", k)
				writeSource(b, fs, filepath.Join(base, js, ks, "logo1.png"), "content")
				writeSource(b, fs, filepath.Join(base, js, "section", ks, "logo.png"), "content")
			}
		}

		for i := 1; i <= 5; i++ {
			writeSource(b, fs, filepath.Join(base, "assetsonly", fmt.Sprintf("image%d.png", i)), "image")
		}

		for i := 1; i <= 5; i++ {
			writeSource(b, fs, filepath.Join(base, "contentonly", fmt.Sprintf("c%d.md", i)), "content")
		}

		capturers[i] = newCapturer(loggers.NewErrorLogger(), sourceSpec, new(noOpFileStore), nil, base)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := capturers[i].capture()
		if err != nil {
			b.Fatal(err)
		}
	}
}
