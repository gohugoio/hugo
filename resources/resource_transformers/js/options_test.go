// Copyright 2020 The Hugo Authors. All rights reserved.
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

package js

import (
	"testing"

	"github.com/gohugoio/hugo/media"

	"github.com/evanw/esbuild/pkg/api"

	qt "github.com/frankban/quicktest"
)

// This test is added to test/warn against breaking the "stability" of the
// cache key. It's sometimes needed to break this, but should be avoided if possible.
func TestOptionKey(t *testing.T) {
	c := qt.New(t)

	opts := map[string]interface{}{
		"TargetPath": "foo",
		"Target":     "es2018",
	}

	key := (&buildTransformation{optsm: opts}).Key()

	c.Assert(key.Value(), qt.Equals, "jsbuild_7891849149754191852")
}

func TestToBuildOptions(t *testing.T) {
	c := qt.New(t)

	opts, err := toBuildOptions(Options{mediaType: media.JavascriptType})

	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle: true,
		Target: api.ESNext,
		Format: api.FormatIIFE,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})

	opts, err = toBuildOptions(Options{
		Target:    "es2018",
		Format:    "cjs",
		Minify:    true,
		mediaType: media.JavascriptType,
		AvoidTDZ:  true,
	})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		AvoidTDZ:          true,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})

	opts, err = toBuildOptions(Options{
		Target: "es2018", Format: "cjs", Minify: true, mediaType: media.JavascriptType,
		SourceMap: "inline"})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Sourcemap:         api.SourceMapInline,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})

	opts, err = toBuildOptions(Options{
		Target: "es2018", Format: "cjs", Minify: true, mediaType: media.JavascriptType,
		SourceMap: "inline"})
	c.Assert(err, qt.IsNil)
	c.Assert(opts, qt.DeepEquals, api.BuildOptions{
		Bundle:            true,
		Target:            api.ES2018,
		Format:            api.FormatCommonJS,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		Sourcemap:         api.SourceMapInline,
		Stdin: &api.StdinOptions{
			Loader: api.LoaderJS,
		},
	})

}
