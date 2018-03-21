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

package hugofs

import (
	"path/filepath"

	"strings"

	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCompositeLanguagFsTest(t *testing.T) {
	assert := require.New(t)

	languages := map[string]bool{
		"sv": true,
		"en": true,
		"nn": true,
	}
	msv := afero.NewMemMapFs()
	baseSv := "/content/sv"
	lfssv := NewLanguageFs("sv", languages, afero.NewBasePathFs(msv, baseSv))
	mnn := afero.NewMemMapFs()
	baseNn := "/content/nn"
	lfsnn := NewLanguageFs("nn", languages, afero.NewBasePathFs(mnn, baseNn))
	men := afero.NewMemMapFs()
	baseEn := "/content/en"
	lfsen := NewLanguageFs("en", languages, afero.NewBasePathFs(men, baseEn))

	// The order will be sv, en, nn
	composite := NewLanguageCompositeFs(lfsnn, lfsen)
	composite = NewLanguageCompositeFs(composite, lfssv)

	afero.WriteFile(msv, filepath.Join(baseSv, "f1.txt"), []byte("some sv"), 0755)
	afero.WriteFile(mnn, filepath.Join(baseNn, "f1.txt"), []byte("some nn"), 0755)
	afero.WriteFile(men, filepath.Join(baseEn, "f1.txt"), []byte("some en"), 0755)

	// Swedish is the top layer.
	assertLangFile(t, composite, "f1.txt", "sv")

	afero.WriteFile(msv, filepath.Join(baseSv, "f2.en.txt"), []byte("some sv"), 0755)
	afero.WriteFile(mnn, filepath.Join(baseNn, "f2.en.txt"), []byte("some nn"), 0755)
	afero.WriteFile(men, filepath.Join(baseEn, "f2.en.txt"), []byte("some en"), 0755)

	// English is in the middle, but the most specific language match wins.
	//assertLangFile(t, composite, "f2.en.txt", "en")

	// Fetch some specific language versions
	assertLangFile(t, composite, filepath.Join(baseNn, "f2.en.txt"), "nn")
	assertLangFile(t, composite, filepath.Join(baseEn, "f2.en.txt"), "en")
	assertLangFile(t, composite, filepath.Join(baseSv, "f2.en.txt"), "sv")

	// Read the root
	f, err := composite.Open("/")
	assert.NoError(err)
	defer f.Close()
	files, err := f.Readdir(-1)
	assert.Equal(4, len(files))
	expected := map[string]bool{
		filepath.FromSlash("/content/en/f1.txt"):    true,
		filepath.FromSlash("/content/nn/f1.txt"):    true,
		filepath.FromSlash("/content/sv/f1.txt"):    true,
		filepath.FromSlash("/content/en/f2.en.txt"): true,
	}
	got := make(map[string]bool)

	for _, fi := range files {
		fil, ok := fi.(*LanguageFileInfo)
		assert.True(ok)
		got[fil.Filename()] = true
	}
	assert.Equal(expected, got)
}

func assertLangFile(t testing.TB, fs afero.Fs, filename, match string) {
	f, err := fs.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	b, err := afero.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)
	if !strings.Contains(s, match) {
		t.Fatalf("got %q expected it to contain %q", s, match)

	}
}
