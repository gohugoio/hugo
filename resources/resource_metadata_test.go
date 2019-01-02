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

	"github.com/stretchr/testify/require"
)

func TestAssignMetadata(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

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
			assert.Equal("My Resource", logo1.Title())
			assert.Equal("My Name", logo1.Name())
			assert.Equal("My Name", foo2.Name())

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
			assert.Equal("My Logo", logo1.Title())
			assert.Equal("My Logo", logo2.Title())
			assert.Equal("My Name", logo1.Name())
			assert.Equal("My Name", foo2.Name())
			assert.Equal("My Name", foo3.Name())
			assert.Equal("My Resource", foo3.Title())

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
			assert.NoError(err)
			assert.Equal("My Logo", logo1.Title())
			assert.Equal("My Resource", foo3.Title())
			_, p1 := logo2.Params()["param1"]
			_, p2 := foo2.Params()["param2"]
			_, p1_2 := foo2.Params()["param1"]
			_, p2_2 := logo2.Params()["param2"]

			icon1, _ := logo2.Params()["icon"]
			icon2, _ := foo2.Params()["icon"]

			assert.True(p1)
			assert.True(p2)

			// Check merge
			assert.True(p2_2)
			assert.False(p1_2)

			assert.Equal("logo", icon1)
			assert.Equal("resource", icon2)

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
			assert.NoError(err)
			assert.Equal("Resource #2", logo2.Title())
			assert.Equal("Logo Name #1", logo2.Name())
			assert.Equal("Resource #4", logo1.Title())
			assert.Equal("Logo Name #2", logo1.Name())
			assert.Equal("Resource #1", foo2.Title())
			assert.Equal("Resource #3", foo1.Title())
			assert.Equal("Name #2", foo1.Name())
			assert.Equal("Resource #5", foo3.Title())

			assert.Equal(logo2, resources.GetMatch("logo name #1*"))

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
			assert.NoError(err)
			assert.Equal("Third Logo #1", logo3.Title())
			assert.Equal("Name #3", logo3.Name())
			assert.Equal("Other Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Other Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

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
			assert.NoError(err)
			assert.Equal("Third Logo", logo3.Title())
			assert.Equal("Name #3", logo3.Name())
			assert.Equal("Other Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Other Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

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
			assert.NoError(err)
			assert.Equal("Logo #3", logo3.Title())
			assert.Equal("third-logo", logo3.Name())
			assert.Equal("Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

		}},
		{[]map[string]interface{}{
			{
				"title": "Third Logo #:counter",
			},
		}, func(err error) {
			// Missing src
			assert.Error(err)

		}},
		{[]map[string]interface{}{
			{
				"title": "Title",
				"src":   "[]",
			},
		}, func(err error) {
			// Invalid pattern
			assert.Error(err)

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
