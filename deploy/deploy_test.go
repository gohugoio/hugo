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

package deploy

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/afero"
	"gocloud.dev/blob"
)

func TestFindDiffs(t *testing.T) {
	hash1 := []byte("hash 1")
	hash2 := []byte("hash 2")
	makeLocal := func(path string, size int64, hash []byte) *localFile {
		return &localFile{Path: path, UploadSize: size, md5: hash}
	}
	makeRemote := func(path string, size int64, hash []byte) *blob.ListObject {
		return &blob.ListObject{Key: path, Size: size, MD5: hash}
	}

	tests := []struct {
		Description string
		Local       []*localFile
		Remote      []*blob.ListObject
		Force       bool
		WantUpdates []*fileToUpload
		WantDeletes []string
	}{
		{
			Description: "empty -> no diffs",
		},
		{
			Description: "local == remote -> no diffs",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 3, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
				makeRemote("ccc", 3, hash2),
			},
		},
		{
			Description: "local == remote with force flag true -> diffs",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 3, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
				makeRemote("ccc", 3, hash2),
			},
			Force: true,
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonForce},
				{makeLocal("bbb", 2, nil), reasonForce},
				{makeLocal("ccc", 3, nil), reasonForce},
			},
		},
		{
			Description: "local == remote with route.Force true -> diffs",
			Local: []*localFile{
				{Path: "aaa", UploadSize: 1, matcher: &matcher{Force: true}, md5: hash1},
				makeLocal("bbb", 2, hash1),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonForce},
			},
		},
		{
			Description: "extra local file -> upload",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("bbb", 2, nil), reasonNotFound},
			},
		},
		{
			Description: "extra remote file -> delete",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, hash1),
				makeRemote("bbb", 2, hash2),
			},
			WantDeletes: []string{"bbb"},
		},
		{
			Description: "diffs in size or md5 -> upload",
			Local: []*localFile{
				makeLocal("aaa", 1, hash1),
				makeLocal("bbb", 2, hash1),
				makeLocal("ccc", 1, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("aaa", 1, nil),
				makeRemote("bbb", 1, hash1),
				makeRemote("ccc", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("aaa", 1, nil), reasonMD5Missing},
				{makeLocal("bbb", 2, nil), reasonSize},
				{makeLocal("ccc", 1, nil), reasonMD5Differs},
			},
		},
		{
			Description: "mix of updates and deletes",
			Local: []*localFile{
				makeLocal("same", 1, hash1),
				makeLocal("updated", 2, hash1),
				makeLocal("updated2", 1, hash2),
				makeLocal("new", 1, hash1),
				makeLocal("new2", 2, hash2),
			},
			Remote: []*blob.ListObject{
				makeRemote("same", 1, hash1),
				makeRemote("updated", 1, hash1),
				makeRemote("updated2", 1, hash1),
				makeRemote("stale", 1, hash1),
				makeRemote("stale2", 1, hash1),
			},
			WantUpdates: []*fileToUpload{
				{makeLocal("new", 1, nil), reasonNotFound},
				{makeLocal("new2", 2, nil), reasonNotFound},
				{makeLocal("updated", 2, nil), reasonSize},
				{makeLocal("updated2", 1, nil), reasonMD5Differs},
			},
			WantDeletes: []string{"stale", "stale2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			local := map[string]*localFile{}
			for _, l := range tc.Local {
				local[l.Path] = l
			}
			remote := map[string]*blob.ListObject{}
			for _, r := range tc.Remote {
				remote[r.Key] = r
			}
			gotUpdates, gotDeletes := findDiffs(local, remote, tc.Force)
			sort.Slice(gotUpdates, func(i, j int) bool { return gotUpdates[i].Local.Path < gotUpdates[j].Local.Path })
			sort.Slice(gotDeletes, func(i, j int) bool { return gotDeletes[i] < gotDeletes[j] })
			if diff := cmp.Diff(gotUpdates, tc.WantUpdates, cmpopts.IgnoreUnexported(localFile{})); diff != "" {
				t.Errorf("updates differ:\n%s", diff)

			}
			if diff := cmp.Diff(gotDeletes, tc.WantDeletes); diff != "" {
				t.Errorf("deletes differ:\n%s", diff)
			}
		})
	}
}

func TestNewLocalFile(t *testing.T) {
	const (
		content = "hello world!"
	)
	contentBytes := []byte(content)
	contentLen := int64(len(contentBytes))
	contentMD5 := md5.Sum(contentBytes)
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(contentBytes); err != nil {
		t.Fatal(err)
	}
	gz.Close()
	gzBytes := buf.Bytes()
	gzLen := int64(len(gzBytes))
	gzMD5 := md5.Sum(gzBytes)

	tests := []struct {
		Description         string
		Path                string
		Matcher             *matcher
		WantContent         []byte
		WantSize            int64
		WantMD5             []byte
		WantContentType     string // empty string is always OK, since content type detection is OS-specific
		WantCacheControl    string
		WantContentEncoding string
	}{
		{
			Description: "file with no suffix",
			Path:        "foo",
			WantContent: contentBytes,
			WantSize:    contentLen,
			WantMD5:     contentMD5[:],
		},
		{
			Description: "file with .txt suffix",
			Path:        "foo.txt",
			WantContent: contentBytes,
			WantSize:    contentLen,
			WantMD5:     contentMD5[:],
		},
		{
			Description:      "CacheControl from matcher",
			Path:             "foo.txt",
			Matcher:          &matcher{CacheControl: "max-age=630720000"},
			WantContent:      contentBytes,
			WantSize:         contentLen,
			WantMD5:          contentMD5[:],
			WantCacheControl: "max-age=630720000",
		},
		{
			Description:         "ContentEncoding from matcher",
			Path:                "foo.txt",
			Matcher:             &matcher{ContentEncoding: "foobar"},
			WantContent:         contentBytes,
			WantSize:            contentLen,
			WantMD5:             contentMD5[:],
			WantContentEncoding: "foobar",
		},
		{
			Description:     "ContentType from matcher",
			Path:            "foo.txt",
			Matcher:         &matcher{ContentType: "foo/bar"},
			WantContent:     contentBytes,
			WantSize:        contentLen,
			WantMD5:         contentMD5[:],
			WantContentType: "foo/bar",
		},
		{
			Description:         "gzipped content",
			Path:                "foo.txt",
			Matcher:             &matcher{Gzip: true},
			WantContent:         gzBytes,
			WantSize:            gzLen,
			WantMD5:             gzMD5[:],
			WantContentEncoding: "gzip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.Description, func(t *testing.T) {
			fs := new(afero.MemMapFs)
			if err := afero.WriteFile(fs, tc.Path, []byte(content), os.ModePerm); err != nil {
				t.Fatal(err)
			}
			lf, err := newLocalFile(fs, tc.Path, tc.Matcher)
			if err != nil {
				t.Fatal(err)
			}
			if got := lf.UploadSize; got != tc.WantSize {
				t.Errorf("got size %d want %d", got, tc.WantSize)
			}
			if got := lf.MD5(); !bytes.Equal(got, tc.WantMD5) {
				t.Errorf("got MD5 %x want %x", got, tc.WantMD5)
			}
			if got := lf.CacheControl(); got != tc.WantCacheControl {
				t.Errorf("got CacheControl %q want %q", got, tc.WantCacheControl)
			}
			if got := lf.ContentEncoding(); got != tc.WantContentEncoding {
				t.Errorf("got ContentEncoding %q want %q", got, tc.WantContentEncoding)
			}
			if tc.WantContentType != "" {
				if got := lf.ContentType(); got != tc.WantContentType {
					t.Errorf("got ContentType %q want %q", got, tc.WantContentType)
				}
			}
			// Verify the content reader last to ensure the
			// previous operations don't interfere with it.
			gotContent, err := ioutil.ReadAll(lf.UploadContentReader)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(gotContent, tc.WantContent) {
				t.Errorf("got content %q want %q", string(gotContent), string(tc.WantContent))
			}
		})
	}
}
