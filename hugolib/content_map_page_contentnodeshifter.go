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

	"github.com/gohugoio/hugo/common/hdebug"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/resources/resource"
)

type contentNodeShifter struct {
	conf config.AllProvider // Used for logging/debugging.
}

func (s *contentNodeShifter) Delete(n contentNode, vec sitesmatrix.Vector) (contentNode, bool, bool) {
	switch v := n.(type) {
	case contentNodes[contentNode]: // TODO1
		deleted, wasDeleted := v[vec]
		if wasDeleted {
			delete(v, vec)
			resource.MarkStale(deleted)
		}
		return deleted, wasDeleted, len(v) == 0
	case *resourceSource:
		if v.sv != vec {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	case contentNodes[contentNodePage]:
		deleted, wasDeleted := v[vec]
		if wasDeleted {
			delete(v, vec)
			resource.MarkStale(deleted)
		}
		return deleted, wasDeleted, len(v) == 0
	case *pageState:
		// TODO1 revise this entire file vs this.
		if !v.s.siteVector.HasVector(vec) {
			return nil, false, false
		}
		resource.MarkStale(v)
		return v, true, true
	case *pageMetaSource:
		resource.MarkStale(v)
		return v, true, true
	default:
		panic(fmt.Sprintf("Delete: unknown type %T", n))
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

func (s *contentNodeShifter) ForEeachInDimension(n contentNode, dims sitesmatrix.Vector, d int, f func(contentNode) bool) {
LOOP1:
	for dims2, v := range contentNodeToSeq2(n) {
		for i, v := range dims2 {
			if i != d && v != dims[i] {
				continue LOOP1
			}
		}
		if f(v) {
			return
		}
	}
}

func (s *contentNodeShifter) Insert(old, new contentNode) (contentNode, contentNode, bool) {
	switch vv := old.(type) {
	case *pageMetaSource:
		return contentNodeSlice{vv, new}, old, false
	case contentNodeSlice:
		newp, ok := new.(*pageMetaSource)
		if !ok {
			panic(fmt.Sprintf("Insert: unknown type %T", new))
		}
		return append(vv, newp), old, false
	case *pageState:
		switch new := new.(type) {
		case *pageMetaSource:
			return contentNodeSlice{vv, new}, old, false
		default:
			panic(fmt.Sprintf("Insert: unknown type %T", new))
		}
	case contentNodes[contentNodePage]:
		switch new := new.(type) {
		case *pageState:
			oldp := vv[new.s.siteVector]
			updated := oldp != new
			if updated {
				resource.MarkStale(oldp)
			}
			hdebug.AssertNotNil(new)
			vv[new.s.siteVector] = new
			return vv, oldp, updated
		case *pageMetaSource:
			s := make(contentNodeSlice, 0, len(vv)+1)
			for _, v := range vv {
				s = append(s, v)
			}
			s = append(s, new)
			return s, vv, false
		default:
			panic(fmt.Sprintf("Insert: unknown type %T", new))
		}
	case contentNodes[contentNode]:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", newp))
		}
		oldp := vv[newp.sv]
		updated := oldp != newp
		if updated {
			resource.MarkStale(oldp)
		}
		vv[newp.sv] = newp
		return vv, oldp, updated
	case *resourceSource:
		newp, ok := new.(*resourceSource)
		if !ok {
			panic(fmt.Sprintf("unknown type %T", new))
		}
		rs := resourceSourcesSlice{
			vv,
			newp,
		}
		return rs, vv, false
	case resourceSourcesSlice:
		newp := new.(*resourceSource)
		return append(vv, newp), old, false

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
	}

	if !fallback {
		// Done
		return nil, false
	}

	if vvv := contentNodeHelper.findContentNodeForSiteVector(siteVector, fallback, contentNodeToSeq(n)); vvv != nil {
		return vvv, true
	}

	return nil, false
}
