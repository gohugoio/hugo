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

package hugolib

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestBranchMap(t *testing.T) {
	c := qt.New(t)

	m := newBranchMap(nil)

	walkAndGetOne := func(c *qt.C, m *branchMap, s string) contentNodeProvider {
		var result contentNodeProvider
		h := func(np contentNodeProvider) bool {
			if np.Key() != s {
				return false
			}
			result = np
			return true
		}

		q := branchMapQuery{
			Deep: true,
			Branch: branchMapQueryCallBacks{
				Key:      newBranchMapQueryKey("", true),
				Page:     h,
				Resource: h,
			},
			Leaf: branchMapQueryCallBacks{
				Page:     h,
				Resource: h,
			},
		}

		c.Assert(m.Walk(q), qt.IsNil)
		c.Assert(result, qt.Not(qt.IsNil))

		return result
	}

	c.Run("Node methods", func(c *qt.C) {
		m := newBranchMap(nil)
		bn, ln := &contentNode{key: "/my/section"}, &contentNode{key: "/my/section/mypage"}
		m.InsertBranch(&contentNode{key: "/my"}) // We need a root section.
		b := m.InsertBranch(bn)
		b.InsertPage(ln.key, ln)

		branch := walkAndGetOne(c, m, "/my/section").(contentNodeInfoProvider)
		page := walkAndGetOne(c, m, "/my/section/mypage").(contentNodeInfoProvider)
		c.Assert(branch.SectionsEntries(), qt.DeepEquals, []string{"my", "section"})
		c.Assert(page.SectionsEntries(), qt.DeepEquals, []string{"my", "section"})
	})

	c.Run("Tree relation", func(c *qt.C) {
		for _, test := range []struct {
			name   string
			s1     string
			s2     string
			expect int
		}{
			{"Sibling", "/blog/sub1", "/blog/sub2", 1},
			{"Root child", "", "/blog", 0},
			{"Child", "/blog/sub1", "/blog/sub1/sub2", 0},
			{"New root", "/blog/sub1", "/docs/sub2", -1},
		} {
			c.Run(test.name, func(c *qt.C) {
				c.Assert(m.TreeRelation(test.s1, test.s2), qt.Equals, test.expect)
			})
		}
	})

	home, blog, blog_sub, blog_sub2, docs, docs_sub := &contentNode{}, &contentNode{key: "/blog"}, &contentNode{key: "/blog/sub"}, &contentNode{key: "/blog/sub2"}, &contentNode{key: "/docs"}, &contentNode{key: "/docs/sub"}
	docs_sub2, docs_sub2_sub := &contentNode{key: "/docs/sub2"}, &contentNode{key: "/docs/sub2/sub"}

	article1, article2 := &contentNode{}, &contentNode{}

	image1, image2, image3 := &contentNode{}, &contentNode{}, &contentNode{}
	json1, json2, json3 := &contentNode{}, &contentNode{}, &contentNode{}
	xml1, xml2 := &contentNode{}, &contentNode{}

	c.Assert(m.InsertBranch(home), qt.Not(qt.IsNil))
	c.Assert(m.InsertBranch(docs), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/docs/data1.json", json1), qt.IsNil)
	c.Assert(m.InsertBranch(docs_sub), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/docs/sub/data2.json", json2), qt.IsNil)
	c.Assert(m.InsertBranch(docs_sub2), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/docs/sub2/data1.xml", xml1), qt.IsNil)
	c.Assert(m.InsertBranch(docs_sub2_sub), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/docs/sub2/sub/data2.xml", xml2), qt.IsNil)
	c.Assert(m.InsertBranch(blog), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/blog/logo.png", image3), qt.IsNil)
	c.Assert(m.InsertBranch(blog_sub), qt.Not(qt.IsNil))
	c.Assert(m.InsertBranch(blog_sub2), qt.Not(qt.IsNil))
	c.Assert(m.InsertResource("/blog/sub2/data3.json", json3), qt.IsNil)

	blogSection := m.Get("/blog")
	c.Assert(blogSection.n, qt.Equals, blog)

	_, section := m.LongestPrefix("/blog/asdfadf")
	c.Assert(section, qt.Equals, blogSection)

	blogSection.InsertPage("/blog/my-article", article1)
	blogSection.InsertPage("/blog/my-article2", article2)
	c.Assert(blogSection.InsertResource("/blog/my-article/sunset.jpg", image1), qt.IsNil)
	c.Assert(blogSection.InsertResource("/blog/my-article2/sunrise.jpg", image2), qt.IsNil)

	type querySpec struct {
		key              string
		isBranchKey      bool
		isPrefix         bool
		noRecurse        bool
		doBranch         bool
		doBranchResource bool
		doPage           bool
		doPageResource   bool
	}

	type queryResult struct {
		query  branchMapQuery
		result []string
	}

	newQuery := func(spec querySpec) *queryResult {
		qr := &queryResult{}

		addResult := func(typ, key string) {
			qr.result = append(qr.result, fmt.Sprintf("%s:%s", typ, key))
		}

		var (
			handleSection        func(np contentNodeProvider) bool
			handlePage           func(np contentNodeProvider) bool
			handleLeafResource   func(np contentNodeProvider) bool
			handleBranchResource func(np contentNodeProvider) bool

			keyBranch branchMapQueryKey
			keyLeaf   branchMapQueryKey
		)

		if spec.isBranchKey {
			keyBranch = newBranchMapQueryKey(spec.key, spec.isPrefix)
		} else {
			keyLeaf = newBranchMapQueryKey(spec.key, spec.isPrefix)
		}

		if spec.doBranch {
			handleSection = func(np contentNodeProvider) bool {
				addResult("section", np.Key())
				return false
			}
		}

		if spec.doPage {
			handlePage = func(np contentNodeProvider) bool {
				addResult("page", np.Key())
				return false
			}
		}

		if spec.doPageResource {
			handleLeafResource = func(np contentNodeProvider) bool {
				addResult("resource", np.Key())
				return false
			}
		}

		if spec.doBranchResource {
			handleBranchResource = func(np contentNodeProvider) bool {
				addResult("resource-branch", np.Key())
				return false
			}
		}

		qr.query = branchMapQuery{
			NoRecurse: spec.noRecurse,
			Branch: branchMapQueryCallBacks{
				Key:      keyBranch,
				Page:     handleSection,
				Resource: handleBranchResource,
			},
			Leaf: branchMapQueryCallBacks{
				Key:      keyLeaf,
				Page:     handlePage,
				Resource: handleLeafResource,
			},
		}

		return qr
	}

	for _, test := range []struct {
		name   string
		spec   querySpec
		expect []string
	}{
		{
			"Branch",
			querySpec{key: "/blog", isBranchKey: true, doBranch: true},
			[]string{"section:/blog"},
		},
		{
			"Branch pages",
			querySpec{key: "/blog", isBranchKey: true, doPage: true},
			[]string{"page:/blog/my-article", "page:/blog/my-article2"},
		},
		{
			"Branch resources",
			querySpec{key: "/docs/", isPrefix: true, isBranchKey: true, doBranchResource: true},
			[]string{"resource-branch:/docs/sub/data2.json", "resource-branch:/docs/sub2/data1.xml", "resource-branch:/docs/sub2/sub/data2.xml"},
		},
		{
			"Branch section and resources",
			querySpec{key: "/docs/", isPrefix: true, isBranchKey: true, doBranch: true, doBranchResource: true},
			[]string{"section:/docs/sub", "resource-branch:/docs/sub/data2.json", "section:/docs/sub2", "resource-branch:/docs/sub2/data1.xml", "section:/docs/sub2/sub", "resource-branch:/docs/sub2/sub/data2.xml"},
		},
		{
			"Branch section and page resources",
			querySpec{key: "/blog", isPrefix: false, isBranchKey: true, doBranchResource: true, doPageResource: true},
			[]string{"resource-branch:/blog/logo.png", "resource:/blog/my-article/sunset.jpg", "resource:/blog/my-article2/sunrise.jpg"},
		},
		{
			"Branch section and pages",
			querySpec{key: "/blog", isBranchKey: true, doBranch: true, doPage: true},
			[]string{"section:/blog", "page:/blog/my-article", "page:/blog/my-article2"},
		},
		{
			"Branch pages and resources",
			querySpec{key: "/blog", isBranchKey: true, doPage: true, doPageResource: true},
			[]string{"page:/blog/my-article", "resource:/blog/my-article/sunset.jpg", "page:/blog/my-article2", "resource:/blog/my-article2/sunrise.jpg"},
		},
		{
			"Leaf page",
			querySpec{key: "/blog/my-article", isBranchKey: false, doPage: true},
			[]string{"page:/blog/my-article"},
		},
		{
			"Leaf page and resources",
			querySpec{key: "/blog/my-article", isBranchKey: false, doPage: true, doPageResource: true},
			[]string{"page:/blog/my-article", "resource:/blog/my-article/sunset.jpg"},
		},
		{
			"Root sections",
			querySpec{key: "/", isBranchKey: true, isPrefix: true, doBranch: true, noRecurse: true},
			[]string{"section:/blog", "section:/docs"},
		},
		{
			"All sections",
			querySpec{key: "", isBranchKey: true, isPrefix: true, doBranch: true},
			[]string{"section:", "section:/blog", "section:/blog/sub", "section:/blog/sub2", "section:/docs", "section:/docs/sub", "section:/docs/sub2", "section:/docs/sub2/sub"},
		},
	} {
		c.Run(test.name, func(c *qt.C) {
			qr := newQuery(test.spec)
			c.Assert(m.Walk(qr.query), qt.IsNil)
			c.Assert(qr.result, qt.DeepEquals, test.expect)
		})
	}
}
