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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/common/herrors"

	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client context for ESBuild.
type Client struct {
	rs  *resources.Spec
	sfs *filesystems.SourceFilesystem
}

// New creates a new client context.
func New(fs *filesystems.SourceFilesystem, rs *resources.Spec) *Client {
	return &Client{
		rs:  rs,
		sfs: fs,
	}
}

type buildTransformation struct {
	optsm map[string]interface{}
	c     *Client
}

func (t *buildTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("jsbuild", t.optsm)
}

func (t *buildTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	ctx.OutMediaType = media.JavascriptType

	opts, err := decodeOptions(t.optsm)
	if err != nil {
		return err
	}

	if opts.TargetPath != "" {
		ctx.OutPath = opts.TargetPath
	} else {
		ctx.ReplaceOutPathExtension(".js")
	}

	src, err := ioutil.ReadAll(ctx.From)
	if err != nil {
		return err
	}

	sdir, _ := path.Split(ctx.SourcePath)
	opts.sourcefile = ctx.SourcePath
	opts.resolveDir = t.c.sfs.RealFilename(sdir)
	opts.workDir = t.c.rs.WorkingDir
	opts.contents = string(src)
	opts.mediaType = ctx.InMediaType

	buildOptions, err := toBuildOptions(opts)
	if err != nil {
		return err
	}

	buildOptions.Plugins, err = createBuildPlugins(t.c, opts)
	if err != nil {
		return err
	}

	result := api.Build(buildOptions)

	if len(result.Errors) > 0 {
		first := result.Errors[0]
		loc := first.Location
		path := loc.File

		var err error
		var f afero.File
		var filename string

		if !strings.HasPrefix(path, "..") {
			// Try first in the assets fs
			var fi os.FileInfo
			fi, err = t.c.rs.BaseFs.Assets.Fs.Stat(path)
			if err == nil {
				m := fi.(hugofs.FileMetaInfo).Meta()
				filename = m.Filename()
				f, err = m.Open()
			}
		}

		if f == nil {
			path = filepath.Join(t.c.rs.WorkingDir, path)
			filename = path
			f, err = t.c.rs.Fs.Os.Open(path)
		}

		if err == nil {
			fe := herrors.NewFileError("js", 0, loc.Line, loc.Column, errors.New(first.Text))
			err, _ := herrors.WithFileContext(fe, filename, f, herrors.SimpleLineMatcher)
			f.Close()
			return err
		}

		return fmt.Errorf("%s", result.Errors[0].Text)
	}

	ctx.To.Write(result.OutputFiles[0].Contents)
	return nil
}

// Process process esbuild transform
func (c *Client) Process(res resources.ResourceTransformer, opts map[string]interface{}) (resource.Resource, error) {
	return res.Transform(
		&buildTransformation{c: c, optsm: opts},
	)
}
