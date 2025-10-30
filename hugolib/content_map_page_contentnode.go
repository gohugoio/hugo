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

var contentNodeHelper helperContentNode

var (
	_ contentNode    = (*contentNodesMap)(nil)
	_ contentNodeMap = (*contentNodesMap)(nil)
)

func (n contentNodesMap) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, nn := range n {
		if nn.ForEeachIdentity(f) {
			return true
		}
	}
	return false
}

func (n contentNodesMap) ForEeachInAllDimensions(f func(contentNode) bool) {
	for _, nn := range n {
		if f(nn) {
			return
		}
	}
}

func (n contentNodesMap) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
}

func (n contentNodesMap) MarkStale() {
	for _, nn := range n {
		resource.MarkStale(nn)
	}
}

func (n contentNodesMap) Path() string {
	return n.one().Path()
}

func (n contentNodesMap) PathInfo() *paths.Path {
	return n.one().PathInfo()
}

type (
	contentNode interface {
		contentNodeBuildStateResetter

		contentNodeForEach
		identity.ForEeachIdentityProvider
		identity.IdentityProvider
		resource.StaleMarker

		Path() string
		PathInfo() *paths.Path
		isContentNodeBranch() bool
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

	helperContentNode struct{}
)

type contentNodeSeq iter.Seq[contentNode]

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

func (n contentNodeSeq) Path() string {
	return n.one().Path()
}

func (n contentNodeSeq) PathInfo() *paths.Path {
	return n.one().PathInfo()
}

func (n contentNodeSeq) isContentNodeBranch() bool {
	return n.one().isContentNodeBranch()
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

func (n contentNodeSeq) one() contentNode {
	for nn := range n {
		return nn
	}
	return nil
}

func (n contentNodeSeq) resetBuildState() {
	for nn := range n {
		nn.resetBuildState()
	}
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

func (n contentNodes) ForEeachIdentity(f func(identity.Identity) bool) bool {
	for _, nn := range n {
		if nn.ForEeachIdentity(f) {
			return true
		}
	}
	return false
}

func (n contentNodes) GetIdentity() identity.Identity {
	return n.one().GetIdentity()
}

func (n contentNodes) MarkStale() {
	for _, nn := range n {
		nn.MarkStale()
	}
}

func (n contentNodes) Path() string {
	return n.one().Path()
}

func (n contentNodes) PathInfo() *paths.Path {
	return n.one().PathInfo()
}

func (n contentNodes) isContentNodeBranch() bool {
	return false
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

func (n contentNodes) one() contentNode {
	if len(n) == 0 {
		panic("pageMetaSourcesSlice is empty")
	}
	return n[0]
}

func (n contentNodes) resetBuildState() {
	// Nothing to do for now.
}

type contentNodesMap map[sitesmatrix.Vector]contentNode

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

type weightedContentNode struct {
	n      contentNode
	weight int
	term   *pageWithOrdinal
}

func (n contentNodesMap) isContentNodeBranch() bool {
	return n.one().isContentNodeBranch()
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

func (n contentNodesMap) one() contentNode {
	for _, nn := range n {
		return nn
	}
	return nil
}

func (n contentNodesMap) resetBuildState() {
	for _, nn := range n {
		nn.resetBuildState()
	}
}

func (n contentNodesMap) siteVectors() sitesmatrix.VectorIterator {
	return sitesmatrix.VectorIteratorFunc(func(yield func(v sitesmatrix.Vector) bool) bool {
		for k := range n {
			if !yield(k) {
				return false
			}
		}
		return true
	})
}
