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

type weightedContentNode struct {
	n      contentNode
	weight int
	term   *pageWithOrdinal
}

var _ contentNode = (*contentNodeSeq)(nil)

type (
	contentNode interface {
		identity.IdentityProvider
		identity.ForEeachIdentityProvider
		Path() string
		PathInfo() *paths.Path
		isContentNodeBranch() bool
		contentNodeBuildStateResetter
		resource.StaleMarker

		contentNodeForEach
	}

	contentNodeBuildStateResetter interface {
		resetBuildState()
	}

	contentNodeSeq  iter.Seq[contentNode]
	contentNodeSeq2 iter.Seq2[sitesmatrix.Vector, contentNode]

	contentNodeForSite interface {
		contentNode
		siteVector() sitesmatrix.Vector
	}

	contentNodeForEach interface {
		// forEeachContentNode iterates over all content nodes.
		// It returns false if the iteration was stopped early.
		forEeachContentNode(func(sitesmatrix.Vector, contentNode) bool) bool
	}

	contentNodeForSites interface {
		contentNode
		siteVectors() sitesmatrix.VectorIterator
	}

	contentNodePage interface {
		contentNode
		nodeCategoryPage() // Marker interface.
	}

	contentNodeLookupContentNode interface {
		contentNode
		lookupContentNode(v sitesmatrix.Vector) contentNode
	}

	contentNodeLookupContentNodes interface {
		lookupContentNodes(v sitesmatrix.Vector, fallback bool) iter.Seq[contentNodeForSite]
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

	contentNodeMap interface {
		contentNode
		contentNodeLookupContentNode
	}

	helperContentNode struct{}
)

var (
	_ contentNode                      = (*resourceSource)(nil)
	_ contentNodeContentWeightProvider = (*pageState)(nil)
	_ contentNodeForSites              = (*pageState)(nil)
	_ contentNodePage                  = (*contentNodeSlice)(nil)
	_ contentNodePage                  = (*pageMetaSource)(nil)
	_ contentNodePage                  = (*pageState)(nil)
)

var contentNodeHelper helperContentNode

func (h helperContentNode) isPageNode(n contentNode) bool {
	n = h.one(n)
	switch n.(type) {
	case contentNodePage:
		return true
	default:
		return false
	}
}

func (helperContentNode) one(n contentNode) contentNode {
	var nn contentNode
	n.forEeachContentNode(func(_ sitesmatrix.Vector, n contentNode) bool {
		nn = n
		return false
	})
	return nn
}

var (
	_ contentNode    = (*contentNodes[contentNodePage])(nil)
	_ contentNodeMap = (*contentNodes[contentNodePage])(nil)
)

func (r contentNodeSeq) isEmpty() bool {
	for range r {
		return false
	}
	return true
}

func contentNodeToContentNodesPage(n contentNode) (contentNodes[contentNodePage], bool) {
	switch v := n.(type) {
	case contentNodes[contentNodePage]:
		return v, false
	case *pageState:
		return contentNodes[contentNodePage]{v.s.siteVector: v}, true
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

type contentNodes[V contentNode] map[sitesmatrix.Vector]V

func (n contentNodes[V]) isEmpty() bool {
	return len(n) == 0
}

func (n contentNodes[V]) lookupContentNode(v sitesmatrix.Vector) contentNode {
	if vv, ok := n[v]; ok {
		return vv
	}
	return nil
}

func (n contentNodes[V]) siteVectors() sitesmatrix.VectorIterator {
	return sitesmatrix.VectorIteratorFunc(func(yield func(v sitesmatrix.Vector) bool) bool {
		for k := range n {
			if !yield(k) {
				return false
			}
		}
		return true
	})
}

func (n contentNodes[V]) one() contentNode {
	for _, nn := range n {
		return nn
	}
	return nil
}

func (ps contentNodes[V]) forEeachContentNode(f func(v sitesmatrix.Vector, n contentNode) bool) bool {
	for vec, nn := range ps {
		if !f(vec, nn) {
			return false
		}
	}
	return true
}

func (n contentNodes[V]) ForEeachInAllDimensions(f func(contentNode) bool) {
	for _, nn := range n {
		if f(nn) {
			return
		}
	}
}

func (n contentNodes[V]) Path() string {
	return n.one().Path()
}

func (n contentNodes[V]) PathInfo() *paths.Path {
	return n.one().PathInfo()
}

func (n contentNodes[V]) isContentNodeBranch() bool {
	return n.one().isContentNodeBranch()
}

func (n contentNodes[V]) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
}

func (n contentNodes[V]) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, nn := range n {
		if nn.ForEeachIdentity(f) {
			return true
		}
	}
	return false
}

func (n contentNodes[V]) resetBuildState() {
	for _, nn := range n {
		nn.resetBuildState()
	}
}

func (n contentNodes[V]) MarkStale() {
	for _, nn := range n {
		resource.MarkStale(nn)
	}
}

func (n contentNodeSeq) ForEeach(f func(n contentNode) bool) bool {
	for nn := range n {
		if !f(nn) {
			return false
		}
	}
	return true
}

func (n contentNodeSeq) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for nn := range n {
		if nn.ForEeachIdentity(f) {
			return true
		}
	}
	return false
}

func (n contentNodeSeq) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
}

func (n contentNodeSeq) MarkStale() {
	for nn := range n {
		resource.MarkStale(nn)
	}
}

func (n contentNodeSeq) resetBuildState() {
	for nn := range n {
		nn.resetBuildState()
	}
}

func (n contentNodeSeq) Path() string {
	return n.one().Path()
}

func (n contentNodeSeq) PathInfo() *paths.Path {
	return n.one().PathInfo()
}

func (n contentNodeSeq) isContentNodeBranch() bool {
	return n.one().isContentNodeBranch()
}

func (n contentNodeSeq) forEeachContentNode(f func(sitesmatrix.Vector, contentNode) bool) bool {
	for nn := range n {
		if !nn.forEeachContentNode(f) {
			return false
		}
	}
	return true
}

func (n contentNodeSeq) one() contentNode {
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

func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
