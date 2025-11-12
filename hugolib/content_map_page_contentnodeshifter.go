// Copyright 2025 The Hugo Authors. All rights reserved.
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

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/resources/resource"
)

type contentNodeShifter struct {
	conf config.AllProvider // Used for logging/debugging.
}

func (s *contentNodeShifter) Delete(n contentNode, vec sitesmatrix.Vector) (contentNode, bool, bool) {
	switch v := n.(type) {
	case contentNodesMap:
		deleted, wasDeleted := v[vec]
		if wasDeleted {
			delete(v, vec)
			resource.MarkStale(deleted)
		}
		return deleted, wasDeleted, len(v) == 0
	case contentNodeForSite:
		if v.siteVector() != vec {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	default:
		v = v.(contentNodeSingle) // Ensure single node.
		resource.MarkStale(v)
		return v, true, true

	}
}

func (s *contentNodeShifter) DeleteFunc(v contentNode, f func(n contentNode) bool) bool {
	switch ss := v.(type) {
	case contentNodeSingle:
		if f(ss) {
			resource.MarkStale(ss)
			return true
		}
		return false
	case contentNodes:
		for i, n := range ss {
			if f(n) {
				resource.MarkStale(n)
				ss = append(ss[:i], ss[i+1:]...)
			}
		}
		return len(ss) == 0
	case contentNodesMap:
		for k, n := range ss {
			if f(n) {
				resource.MarkStale(n)
				delete(ss, k)
			}
		}
		return len(ss) == 0
	default:
		panic(fmt.Sprintf("DeleteFunc: unknown type %T", v))
	}
}

func (s *contentNodeShifter) ForEeachInAllDimensions(n contentNode, f func(contentNode) bool) {
	if n == nil {
		return
	}
	if v, ok := n.(interface {
		// Implemented by all the list nodes.
		ForEeachInAllDimensions(f func(contentNode) bool)
	}); ok {
		v.ForEeachInAllDimensions(f)
		return
	}
	f(n)
}

func (s *contentNodeShifter) ForEeachInDimension(n contentNode, vec sitesmatrix.Vector, d int, f func(contentNode) bool) {
LOOP1:
	for vec2, v := range contentNodeToSeq2(n) {
		for i, v := range vec2 {
			if i != d && v != vec[i] {
				continue LOOP1
			}
		}
		if !f(v) {
			return
		}
	}
}

func (s *contentNodeShifter) Insert(old, new contentNode) (contentNode, contentNode, bool) {
	new = new.(contentNodeSingle) // Ensure single node.

	switch vv := old.(type) {
	case contentNodeSingle:
		return contentNodes{vv, new}, old, false
	case contentNodes:
		s := make(contentNodes, 0, len(vv)+1)
		s = append(s, vv...)
		s = append(s, new)
		return s, old, false
	case contentNodesMap:
		switch new := new.(type) {
		case contentNodeForSite:
			oldp := vv[new.siteVector()]
			updated := oldp != new
			if updated {
				resource.MarkStale(oldp)
			}
			vv[new.siteVector()] = new
			return vv, oldp, updated
		default:
			s := make(contentNodes, 0, len(vv)+1)
			for _, v := range vv {
				s = append(s, v)
			}
			s = append(s, new)
			return s, vv, false
		}
	default:
		panic(fmt.Sprintf("Insert: unknown type %T", old))
	}
}

func (s *contentNodeShifter) Shift(n contentNode, siteVector sitesmatrix.Vector, fallback bool) (contentNode, bool) {
	switch v := n.(type) {
	case contentNodeLookupContentNode:
		if vv := v.lookupContentNode(siteVector); vv != nil {
			return vv, true
		}
	default:
		panic(fmt.Sprintf("Shift: unknown type %T for %q", n, n.Path()))
	}

	if !fallback {
		// Done
		return nil, false
	}

	if vvv := cnh.findContentNodeForSiteVector(siteVector, fallback, contentNodeToSeq(n)); vvv != nil {
		return vvv, true
	}

	return nil, false
}
