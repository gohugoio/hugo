// Copyright 2024 The Hugo Authors. All rights reserved.
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

package resources_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources"
)

func TestNewResource(t *testing.T) {
	c := qt.New(t)

	spec := newTestResourceSpec(specDescriptor{c: c})

	open := hugio.NewOpenReadSeekCloser(hugio.NewReadSeekerNoOpCloserFromString("content"))

	rd := resources.ResourceSourceDescriptor{
		OpenReadSeekCloser:   open,
		TargetPath:           "a/b.txt",
		BasePathRelPermalink: "c/d",
		BasePathTargetPath:   "e/f",
		GroupIdentity:        identity.Anonymous,
	}

	r, err := spec.NewResource(rd)
	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.Not(qt.IsNil))
	c.Assert(r.RelPermalink(), qt.Equals, "/c/d/a/b.txt")

	info := resources.GetTestInfoForResource(r)
	c.Assert(info.Paths.TargetLink(), qt.Equals, "/c/d/a/b.txt")
	c.Assert(info.Paths.TargetPath(), qt.Equals, "/e/f/a/b.txt")
}
