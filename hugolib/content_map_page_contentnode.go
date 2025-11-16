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
	"iter"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugolib/sitesmatrix"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/resources/resource"
)

var _ contentNode = (*contentNodeSeq)(nil)

var (
	_ contentNode                      = (*resourceSource)(nil)
	_ contentNodeContentWeightProvider = (*pageState)(nil)
	_ contentNodeForSites              = (*pageState)(nil)
	_ contentNodePage                  = (*contentNodes)(nil)
	_ contentNodePage                  = (*pageMetaSource)(nil)
	_ contentNodePage                  = (*pageState)(nil)
	_ contentNodeSourceEntryIDProvider = (*pageState)(nil)
	_ contentNodeSourceEntryIDProvider = (*pageMetaSource)(nil)
	_ contentNodeSourceEntryIDProvider = (*resourceSource)(nil)
	_ contentNodeLookupContentNode     = (*contentNodesMap)(nil)
	_ contentNodeLookupContentNode     = (contentNodes)(nil)
	_ contentNodeSingle                = (*pageMetaSource)(nil)
	_ contentNodeSingle                = (*pageState)(nil)
	_ contentNodeSingle                = (*resourceSource)(nil)
)

var (
	_ contentNode    = (*contentNodesMap)(nil)
	_ contentNodeMap = (*contentNodesMap)(nil)
)

func (n contentNodesMap) ForEeachInAllDimensions(f func(contentNode) bool) {
	for _, nn := range n {
		if !f(nn) {
			return
		}
	}
}

func (n contentNodesMap) Path() string {
	return n.sample().Path()
}

type (
	contentNode interface {
		Path() string
		contentNodeForEach
	}

	contentNodeBuildStateResetter interface {
		resetBuildState()
	}

	contentNodeSeq2 iter.Seq2[sitesmatrix.Vector, contentNode]

	contentNodeForSite interface {
		contentNode
		siteVector() sitesmatrix.Vector
	}

	contentNodeForEach interface {
		// forEeachContentNode iterates over all content nodes.
		// calling f with the site vector and the content node.
		// If f returns false, the iteration stops.
		// It returns false if the iteration was stopped early.
		forEeachContentNode(f func(sitesmatrix.Vector, contentNode) bool) bool
	}

	contentNodeForSites interface {
		contentNode
		siteVectors() sitesmatrix.VectorIterator
	}

	contentNodePage interface {
		contentNode
		nodeCategoryPage() // Marker interface.
	}

	contentNodeSingle interface {
		contentNode
		nodeCategorySingle() // Marker interface.
	}

	contentNodeLookupContentNode interface {
		contentNode
		lookupContentNode(v sitesmatrix.Vector) contentNode
	}

	contentNodeLookupContentNodes interface {
		lookupContentNodes(v sitesmatrix.Vector, fallback bool) iter.Seq[contentNodeForSite]
	}

	contentNodeSampleProvider interface {
		// sample is used to get a sample contentNode from a collection.
		sample() contentNode
	}

	contentNodeCascadeProvider interface {
		getCascade() *page.PageMatcherParamsConfigs
	}

	contentNodeContentWeightProvider interface {
		contentWeight() int
	}

	contentNodeIsEmptyProvider interface {
		isEmpty() bool
	}

	contentNodeSourceEntryIDProvider interface {
		nodeSourceEntryID() any
	}

	contentNodeMap interface {
		contentNode
		contentNodeLookupContentNode
	}
)

type contentNodeSeq iter.Seq[contentNode]

func (n contentNodeSeq) Path() string {
	return n.sample().Path()
}

func (n contentNodeSeq) isEmpty() bool {
	for range n {
		return false
	}
	return true
}

func (n contentNodeSeq) forEeachContentNode(f func(sitesmatrix.Vector, contentNode) bool) bool {
	for nn := range n {
		if !nn.forEeachContentNode(f) {
			return false
		}
	}
	return true
}

func (n contentNodeSeq) sample() contentNode {
	for nn := range n {
		return nn
	}
	return nil
}

func (n contentNodeSeq2) forEeachContentNode(f func(sitesmatrix.Vector, contentNode) bool) bool {
	for _, nn := range n {
		if !nn.forEeachContentNode(f) {
			return false
		}
	}
	return true
}

type contentNodes []contentNode

func (n contentNodes) Path() string {
	return n.sample().Path()
}

func (m contentNodes) isEmpty() bool {
	return len(m) == 0
}

func (n contentNodes) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for _, nn := range n {
		if !nn.forEeachContentNode(f) {
			return false
		}
	}
	return true
}

func (n contentNodes) lookupContentNode(v sitesmatrix.Vector) contentNode {
	for _, nn := range n {
		if vv := nn.(contentNodeLookupContentNode).lookupContentNode(v); vv != nil {
			return vv
		}
	}
	return nil
}

func (m *contentNodes) nodeCategoryPage() {
	// Marker method.
}

func (n contentNodes) sample() contentNode {
	if len(n) == 0 {
		panic("pageMetaSourcesSlice is empty")
	}
	return n[0]
}

type contentNodesMap map[sitesmatrix.Vector]contentNode

var cnh helperContentNode

type helperContentNode struct{}

func (h helperContentNode) findContentNodeForSiteVector(q sitesmatrix.Vector, fallback bool, candidates contentNodeSeq) contentNodeForSite {
	var (
		best         contentNodeForSite = nil
		bestDistance int
	)

	for n := range candidates {
		// The order of candidates is unstable, so we need to compare the matches to
		// get stable output. This compare will also make sure that we pick
		// language, version and role according to their individual sort order:
		// Closer is better, and matches above are better than matches below.
		if m := n.(contentNodeLookupContentNodes).lookupContentNodes(q, fallback); m != nil {
			for nn := range m {
				vec := nn.siteVector()
				if q == vec {
					// Exact match.
					return nn
				}

				distance := q.Distance(vec)

				if best == nil {
					best = nn
					bestDistance = distance
				} else {
					distanceAbs := absint(distance)
					bestDistanceAbs := absint(bestDistance)
					if distanceAbs < bestDistanceAbs {
						// Closer is better.
						best = nn
						bestDistance = distance
					} else if distanceAbs == bestDistanceAbs && distance > 0 {
						// Positive distance is better than negative.
						best = nn
						bestDistance = distance
					}
				}
			}
		}
	}

	return best
}

func (h helperContentNode) markStale(n contentNode) {
	n.forEeachContentNode(
		func(_ sitesmatrix.Vector, nn contentNode) bool {
			resource.MarkStale(nn)
			return true
		},
	)
}

func (h helperContentNode) resetBuildState(n contentNode) {
	n.forEeachContentNode(
		func(_ sitesmatrix.Vector, nn contentNode) bool {
			if nnn, ok := nn.(contentNodeBuildStateResetter); ok {
				nnn.resetBuildState()
			}
			return true
		},
	)
}

func (h helperContentNode) isBranchNode(n contentNode) bool {
	switch nn := n.(type) {
	case *pageMetaSource:
		return nn.pathInfo.IsBranchBundle()
	case *pageState:
		return nn.IsNode()
	case contentNodeSampleProvider:
		return h.isBranchNode(nn.sample())
	default:
		return false
	}
}

func (h helperContentNode) isPageNode(n contentNode) bool {
	n, _ = h.one(n)
	switch n.(type) {
	case contentNodePage:
		return true
	default:
		return false
	}
}

func (helperContentNode) one(n contentNode) (nn contentNode, hasMore bool) {
	n.forEeachContentNode(func(_ sitesmatrix.Vector, n contentNode) bool {
		if nn == nil {
			nn = n
			return true
		} else {
			hasMore = true
			return false
		}
	})
	return
}

func (h helperContentNode) GetIdentity(n contentNode) identity.Identity {
	switch nn := n.(type) {
	case identity.IdentityProvider:
		return nn.GetIdentity()
	case contentNodeSampleProvider:
		return h.GetIdentity(nn.sample())
	default:
		return identity.Anonymous
	}
}

func (h helperContentNode) PathInfo(n contentNode) *paths.Path {
	switch nn := n.(type) {
	case interface{ PathInfo() *paths.Path }:
		return nn.PathInfo()
	case contentNodeSampleProvider:
		return h.PathInfo(nn.sample())
	default:
		return nil
	}
}

func (h helperContentNode) toForEachIdentityProvider(n contentNode) identity.ForEeachIdentityProvider {
	return identity.ForEeachIdentityProviderFunc(
		func(cb func(id identity.Identity) bool) bool {
			return n.forEeachContentNode(
				func(vec sitesmatrix.Vector, nn contentNode) bool {
					return cb(h.GetIdentity(n))
				},
			)
		})
}

type weightedContentNode struct {
	n      contentNode
	weight int
	term   *pageWithOrdinal
}

func (n contentNodesMap) isEmpty() bool {
	return len(n) == 0
}

func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func contentNodeToContentNodesPage(n contentNode) (contentNodesMap, bool) {
	switch v := n.(type) {
	case contentNodesMap:
		return v, false
	case *pageState:
		return contentNodesMap{v.s.siteVector: v}, true
	default:
		panic(fmt.Sprintf("contentNodeToContentNodesPage: unexpected type %T", n))
	}
}

func contentNodeToSeq(n contentNodeForEach) contentNodeSeq {
	if nn, ok := n.(contentNodeSeq); ok {
		return nn
	}
	return func(yield func(contentNode) bool) {
		n.forEeachContentNode(func(_ sitesmatrix.Vector, nn contentNode) bool {
			return yield(nn)
		})
	}
}

func contentNodeToSeq2(n contentNodeForEach) contentNodeSeq2 {
	if nn, ok := n.(contentNodeSeq2); ok {
		return nn
	}
	return func(yield func(sitesmatrix.Vector, contentNode) bool) {
		n.forEeachContentNode(func(vec sitesmatrix.Vector, nn contentNode) bool {
			return yield(vec, nn)
		})
	}
}

func (ps contentNodesMap) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for vec, nn := range ps {
		if !f(vec, nn) {
			return false
		}
	}
	return true
}

func (n contentNodesMap) lookupContentNode(v sitesmatrix.Vector) contentNode {
	if vv, ok := n[v]; ok {
		return vv
	}
	return nil
}

func (n contentNodesMap) sample() contentNode {
	for _, nn := range n {
		return nn
	}
	return nil
}

func (n contentNodesMap) siteVectors() sitesmatrix.VectorIterator {
	return n
}

func (n contentNodesMap) ForEachVector(yield func(v sitesmatrix.Vector) bool) bool {
	for v := range n {
		if !yield(v) {
			return false
		}
	}
	return true
}

func (n contentNodesMap) LenVectors() int {
	return len(n)
}

func (n contentNodesMap) VectorSample() sitesmatrix.Vector {
	for v := range n {
		return v
	}
	panic("no vectors")
}
