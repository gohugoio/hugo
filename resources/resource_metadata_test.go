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

package resources

import (
	"testing"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"

	qt "github.com/frankban/quicktest"
)

func TestAssignMetadata(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})

	var foo1, foo2, foo3, logo1, logo2, logo3 resource.Resource
	var resources resource.Resources

	for _, this := range []struct {
		metaData   []map[string]interface{}
		assertFunc func(err error)
	}{
		{[]map[string]interface{}{
			{
				"title": "My Resource",
				"name":  "My Name",
				"src":   "*",
			},
		}, func(err error) {
			c.Assert(logo1.Title(), qt.Equals, "My Resource")
			c.Assert(logo1.Name(), qt.Equals, "My Name")
			c.Assert(foo2.Name(), qt.Equals, "My Name")

		}},
		{[]map[string]interface{}{
			{
				"title": "My Logo",
				"src":   "*loGo*",
			},
			{
				"title": "My Resource",
				"name":  "My Name",
				"src":   "*",
			},
		}, func(err error) {
			c.Assert(logo1.Title(), qt.Equals, "My Logo")
			c.Assert(logo2.Title(), qt.Equals, "My Logo")
			c.Assert(logo1.Name(), qt.Equals, "My Name")
			c.Assert(foo2.Name(), qt.Equals, "My Name")
			c.Assert(foo3.Name(), qt.Equals, "My Name")
			c.Assert(foo3.Title(), qt.Equals, "My Resource")

		}},
		{[]map[string]interface{}{
			{
				"title": "My Logo",
				"src":   "*loGo*",
				"params": map[string]interface{}{
					"Param1": true,
					"icon":   "logo",
				},
			},
			{
				"title": "My Resource",
				"src":   "*",
				"params": map[string]interface{}{
					"Param2": true,
					"icon":   "resource",
				},
			},
		}, func(err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(logo1.Title(), qt.Equals, "My Logo")
			c.Assert(foo3.Title(), qt.Equals, "My Resource")
			_, p1 := logo2.Params()["param1"]
			_, p2 := foo2.Params()["param2"]
			_, p1_2 := foo2.Params()["param1"]
			_, p2_2 := logo2.Params()["param2"]

			icon1 := logo2.Params()["icon"]
			icon2 := foo2.Params()["icon"]

			c.Assert(p1, qt.Equals, true)
			c.Assert(p2, qt.Equals, true)

			// Check merge
			c.Assert(p2_2, qt.Equals, true)
			c.Assert(p1_2, qt.Equals, false)

			c.Assert(icon1, qt.Equals, "logo")
			c.Assert(icon2, qt.Equals, "resource")

		}},
		{[]map[string]interface{}{
			{
				"name": "Logo Name #:counter",
				"src":  "*logo*",
			},
			{
				"title": "Resource #:counter",
				"name":  "Name #:counter",
				"src":   "*",
			},
		}, func(err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(logo2.Title(), qt.Equals, "Resource #2")
			c.Assert(logo2.Name(), qt.Equals, "Logo Name #1")
			c.Assert(logo1.Title(), qt.Equals, "Resource #4")
			c.Assert(logo1.Name(), qt.Equals, "Logo Name #2")
			c.Assert(foo2.Title(), qt.Equals, "Resource #1")
			c.Assert(foo1.Title(), qt.Equals, "Resource #3")
			c.Assert(foo1.Name(), qt.Equals, "Name #2")
			c.Assert(foo3.Title(), qt.Equals, "Resource #5")

			c.Assert(resources.GetMatch("logo name #1*"), qt.Equals, logo2)

		}},
		{[]map[string]interface{}{
			{
				"title": "Third Logo #:counter",
				"src":   "logo3.png",
			},
			{
				"title": "Other Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(logo3.Title(), qt.Equals, "Third Logo #1")
			c.Assert(logo3.Name(), qt.Equals, "Name #3")
			c.Assert(logo2.Title(), qt.Equals, "Other Logo #1")
			c.Assert(logo2.Name(), qt.Equals, "Name #1")
			c.Assert(logo1.Title(), qt.Equals, "Other Logo #2")
			c.Assert(logo1.Name(), qt.Equals, "Name #2")

		}},
		{[]map[string]interface{}{
			{
				"title": "Third Logo",
				"src":   "logo3.png",
			},
			{
				"title": "Other Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(logo3.Title(), qt.Equals, "Third Logo")
			c.Assert(logo3.Name(), qt.Equals, "Name #3")
			c.Assert(logo2.Title(), qt.Equals, "Other Logo #1")
			c.Assert(logo2.Name(), qt.Equals, "Name #1")
			c.Assert(logo1.Title(), qt.Equals, "Other Logo #2")
			c.Assert(logo1.Name(), qt.Equals, "Name #2")

		}},
		{[]map[string]interface{}{
			{
				"name": "third-logo",
				"src":  "logo3.png",
			},
			{
				"title": "Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			c.Assert(err, qt.IsNil)
			c.Assert(logo3.Title(), qt.Equals, "Logo #3")
			c.Assert(logo3.Name(), qt.Equals, "third-logo")
			c.Assert(logo2.Title(), qt.Equals, "Logo #1")
			c.Assert(logo2.Name(), qt.Equals, "Name #1")
			c.Assert(logo1.Title(), qt.Equals, "Logo #2")
			c.Assert(logo1.Name(), qt.Equals, "Name #2")

		}},
		{[]map[string]interface{}{
			{
				"title": "Third Logo #:counter",
			},
		}, func(err error) {
			// Missing src
			c.Assert(err, qt.Not(qt.IsNil))

		}},
		{[]map[string]interface{}{
			{
				"title": "Title",
				"src":   "[]",
			},
		}, func(err error) {
			// Invalid pattern
			c.Assert(err, qt.Not(qt.IsNil))

		}},
	} {

		foo2 = spec.newGenericResource(nil, nil, nil, "/b/foo2.css", "foo2.css", media.CSSType)
		logo2 = spec.newGenericResource(nil, nil, nil, "/b/Logo2.png", "Logo2.png", pngType)
		foo1 = spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType)
		logo1 = spec.newGenericResource(nil, nil, nil, "/a/logo1.png", "logo1.png", pngType)
		foo3 = spec.newGenericResource(nil, nil, nil, "/b/foo3.css", "foo3.css", media.CSSType)
		logo3 = spec.newGenericResource(nil, nil, nil, "/b/logo3.png", "logo3.png", pngType)

		resources = resource.Resources{
			foo2,
			logo2,
			foo1,
			logo1,
			foo3,
			logo3,
		}

		this.assertFunc(AssignMetadata(this.metaData, resources...))
	}

}
