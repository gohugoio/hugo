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

package paths

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestBaseURL(t *testing.T) {
	c := qt.New(t)
	b, err := newBaseURLFromString("http://example.com")
	c.Assert(err, qt.IsNil)
	c.Assert(b.String(), qt.Equals, "http://example.com")

	p, err := b.WithProtocol("webcal://")
	c.Assert(err, qt.IsNil)
	c.Assert(p, qt.Equals, "webcal://example.com")

	p, err = b.WithProtocol("webcal")
	c.Assert(err, qt.IsNil)
	c.Assert(p, qt.Equals, "webcal://example.com")

	_, err = b.WithProtocol("mailto:")
	c.Assert(err, qt.Not(qt.IsNil))

	b, err = newBaseURLFromString("mailto:hugo@rules.com")
	c.Assert(err, qt.IsNil)
	c.Assert(b.String(), qt.Equals, "mailto:hugo@rules.com")

	// These are pretty constructed
	p, err = b.WithProtocol("webcal")
	c.Assert(err, qt.IsNil)
	c.Assert(p, qt.Equals, "webcal:hugo@rules.com")

	p, err = b.WithProtocol("webcal://")
	c.Assert(err, qt.IsNil)
	c.Assert(p, qt.Equals, "webcal://hugo@rules.com")

	// Test with "non-URLs". Some people will try to use these as a way to get
	// relative URLs working etc.
	b, err = newBaseURLFromString("/")
	c.Assert(err, qt.IsNil)
	c.Assert(b.String(), qt.Equals, "/")

	b, err = newBaseURLFromString("")
	c.Assert(err, qt.IsNil)
	c.Assert(b.String(), qt.Equals, "")

	// BaseURL with sub path
	b, err = newBaseURLFromString("http://example.com/sub")
	c.Assert(err, qt.IsNil)
	c.Assert(b.String(), qt.Equals, "http://example.com/sub")
	c.Assert(b.HostURL(), qt.Equals, "http://example.com")
}
