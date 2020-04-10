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

package config

import (
	"errors"
	"testing"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/types"

	qt "github.com/frankban/quicktest"

	"github.com/spf13/viper"
)

func TestBuild(t *testing.T) {
	c := qt.New(t)

	v := viper.New()
	v.Set("build", map[string]interface{}{
		"useResourceCacheWhen": "always",
	})

	b := DecodeBuild(v)

	c.Assert(b.UseResourceCacheWhen, qt.Equals, "always")

	v.Set("build", map[string]interface{}{
		"useResourceCacheWhen": "foo",
	})

	b = DecodeBuild(v)

	c.Assert(b.UseResourceCacheWhen, qt.Equals, "fallback")

	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.Equals, true)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.Equals, false)

	b.UseResourceCacheWhen = "always"
	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.Equals, true)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.Equals, true)
	c.Assert(b.UseResourceCache(nil), qt.Equals, true)

	b.UseResourceCacheWhen = "never"
	c.Assert(b.UseResourceCache(herrors.ErrFeatureNotAvailable), qt.Equals, false)
	c.Assert(b.UseResourceCache(errors.New("err")), qt.Equals, false)
	c.Assert(b.UseResourceCache(nil), qt.Equals, false)

}

func TestServer(t *testing.T) {
	c := qt.New(t)

	cfg, err := FromConfigString(`[[server.headers]]
for = "/*.jpg"

[server.headers.values]
X-Frame-Options = "DENY"
X-XSS-Protection = "1; mode=block"
X-Content-Type-Options = "nosniff"
`, "toml")

	c.Assert(err, qt.IsNil)

	s := DecodeServer(cfg)

	c.Assert(s.Match("/foo.jpg"), qt.DeepEquals, []types.KeyValueStr{
		{Key: "X-Content-Type-Options", Value: "nosniff"},
		{Key: "X-Frame-Options", Value: "DENY"},
		{Key: "X-XSS-Protection", Value: "1; mode=block"}})

}
