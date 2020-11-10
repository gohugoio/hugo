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
	"io"
	"os"
	"path"
	"strings"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/resources/resource"

	radix "github.com/armon/go-radix"
	"github.com/pkg/errors"
)

var noTaxonomiesFilter = func(s string, n *contentNode) bool {
	return n != nil && n.IsView()
}

func newBranchMap(createBranchNode func(elem ...string) *contentNode) *branchMap {
	return &branchMap{
		branches:         newNodeTree("branches"),
		createBranchNode: createBranchNode,
	}
}

func newBranchMapQueryKey(value string, isPrefix bool) branchMapQueryKey {
	return branchMapQueryKey{Value: value, isPrefix: isPrefix, isSet: true}
}

func newContentBranchNode(n *contentNode) *contentBranchNode {
	return &contentBranchNode{
		n:             n,
		resources:     &contentBranchNodeTree{nodes: newNodeTree("resources")},
		pages:         &contentBranchNodeTree{nodes: newNodeTree("pages")},
		pageResources: &contentBranchNodeTree{nodes: newNodeTree("pageResources")},
		refs:          make(map[interface{}]ordinalWeight),
	}
}

func newNodeTree(name string) nodeTree {
	tree := &defaultNodeTree{nodeTree: radix.New()}
	return tree
	// return &nodeTreeUpdateTracer{name: name, nodeTree: tree}
}

type branchMap struct {
	// branches stores *contentBranchNode
	branches nodeTree

	createBranchNode func(elem ...string) *contentNode
}

func (m *branchMap) GetBranchOrLeaf(key string) *contentNode {
	s, branch := m.LongestPrefix(key)
	if branch == nil {
		return nil
	}

	if key == s {
		// A branch node.
		return branch.n
	}
	n, found := branch.pages.nodes.Get(key)
	if !found {
		return nil
	}

	return n.(*contentNode)
}

func (m *branchMap) GetNode(key string) *contentNode {
	n, _ := m.GetNodeAndTree(key)
	return n
}

func (m *branchMap) GetNodeAndTree(key string) (*contentNode, nodeTree) {
	s, branch := m.LongestPrefix(key)
	if branch == nil {
		return nil, nil
	}

	if key == s {
		// It's a branch node (e.g. a section).
		return branch.n, m.branches
	}

	findFirst := func(s string, trees ...*contentBranchNodeTree) (*contentNode, nodeTree) {
		for _, tree := range trees {
			if v, found := tree.nodes.Get(s); found {
				return v.(*contentNode), tree.nodes
			}
		}
		return nil, nil
	}

	return findFirst(key, branch.pages, branch.pageResources, branch.resources)
}

// GetFirstSection walks up the tree from s and returns the first
// section below root.
func (m *branchMap) GetFirstSection(s string) (string, *contentNode) {
	for {
		k, v, found := m.branches.LongestPrefix(s)

		if !found {
			return "", nil
		}

		// /blog
		if strings.Count(k, "/") <= 1 {
			return k, v.(*contentBranchNode).n
		}

		s = path.Dir(s)

	}
}

// InsertBranch inserts or updates a branch.
func (m *branchMap) InsertBranch(n *contentNode) *contentBranchNode {
	_, b := m.InsertRootAndBranch(n)
	return b
}

func (m *branchMap) InsertResource(key string, n *contentNode) error {
	if err := validateSectionMapKey(key); err != nil {
		return err
	}

	_, v, found := m.branches.LongestPrefix(key)
	if !found {
		return errors.Errorf("no section found for resource %q", key)
	}

	v.(*contentBranchNode).resources.nodes.Insert(key, n)

	return nil
}

// InsertBranchAndRoot inserts or updates a branch.
// The return values are the branch's root or nil and then the branch itself.
func (m *branchMap) InsertRootAndBranch(n *contentNode) (root *contentBranchNode, branch *contentBranchNode) {
	mustValidateSectionMapKey(n.key)
	if v, found := m.branches.Get(n.key); found {
		// Update existing.
		branch = v.(*contentBranchNode)
		branch.n = n
		return
	}

	if strings.Count(n.key, "/") > 1 {
		// Make sure we have a root section.
		s, v, found := m.branches.LongestPrefix(n.key)
		if !found || s == "" {
			nkey := n.KeyParts()
			rkey := nkey[:1]
			root = newContentBranchNode(m.createBranchNode(rkey...))
			m.branches.Insert(root.n.key, root)
		} else {
			root = v.(*contentBranchNode)
		}
	}

	if branch == nil {
		branch = newContentBranchNode(n)
		m.branches.Insert(n.key, branch)
	}

	return
}

// GetLeaf gets the leaf node identified with s, nil if not found.
func (m *branchMap) GetLeaf(s string) *contentNode {
	_, branch := m.LongestPrefix(s)
	if branch != nil {
		n, found := branch.pages.nodes.Get(s)
		if found {
			return n.(*contentNode)
		}
	}
	// Not  found.
	return nil
}

// LongestPrefix returns the branch with the longest prefix match of s.
func (m *branchMap) LongestPrefix(s string) (string, *contentBranchNode) {
	k, v, found := m.branches.LongestPrefix(s)
	if !found {
		return "", nil
	}
	return k, v.(*contentBranchNode)
}

// Returns
// 0 if s2 is a descendant of s1
// 1 if s2 is a sibling of s1
// else -1
func (m *branchMap) TreeRelation(s1, s2 string) int {
	if s1 == "" && s2 != "" {
		return 0
	}

	if strings.HasPrefix(s1, s2) {
		return 0
	}

	for {
		s2 = s2[:strings.LastIndex(s2, "/")]
		if s2 == "" {
			break
		}

		if s1 == s2 {
			return 0
		}

		if strings.HasPrefix(s1, s2) {
			return 1
		}
	}

	return -1
}

// Walk walks m filtering the nodes with q.
func (m *branchMap) Walk(q branchMapQuery) error {
	if q.Branch.Key.IsZero() == q.Leaf.Key.IsZero() {
		return errors.New("must set at most one Key")
	}

	if q.Leaf.Key.IsPrefix() {
		return errors.New("prefix search is currently only implemented starting for branch keys")
	}

	if q.Exclude != nil {
		// Apply global node filters.
		applyFilterPage := func(c contentTreeNodeCallbackNew) contentTreeNodeCallbackNew {
			if c == nil {
				return nil
			}
			return func(n contentNodeProvider) bool {
				if q.Exclude(n.Key(), n.GetNode()) {
					// Skip this node, but continue walk.
					return false
				}
				return c(n)
			}
		}

		applyFilterResource := func(c contentTreeNodeCallbackNew) contentTreeNodeCallbackNew {
			if c == nil {
				return nil
			}
			return func(n contentNodeProvider) bool {
				if q.Exclude(n.Key(), n.GetNode()) {
					// Skip this node, but continue walk.
					return false
				}
				return c(n)
			}
		}

		q.Branch.Page = applyFilterPage(q.Branch.Page)
		q.Branch.Resource = applyFilterResource(q.Branch.Resource)
		q.Leaf.Page = applyFilterPage(q.Leaf.Page)
		q.Leaf.Resource = applyFilterResource(q.Leaf.Resource)

	}

	if q.BranchExclude != nil {
		cb := q.Branch.Page
		q.Branch.Page = func(n contentNodeProvider) bool {
			if q.BranchExclude(n.Key(), n.GetNode()) {
				return true
			}
			return cb(n)
		}
	}

	type depthType int

	const (
		depthAll depthType = iota
		depthBranch
		depthLeaf
	)

	newNodeProviderResource := func(s string, n, owner *contentNode, b *contentBranchNode) contentNodeProvider {
		var np contentNodeProvider
		if !q.Deep {
			np = n
		} else {
			var nInfo contentNodeInfoProvider = &contentNodeInfo{
				branch:     b,
				isResource: true,
			}

			np = struct {
				types.Identifier
				contentNodeInfoProvider
				contentGetNodeProvider
				contentGetContainerNodeProvider
				contentGetBranchProvider
			}{
				n,
				nInfo,
				n,
				owner,
				b,
			}
		}

		return np
	}

	handleBranchPage := func(depth depthType, s string, v interface{}) bool {
		bn := v.(*contentBranchNode)

		if depth <= depthBranch {

			if q.Branch.Page != nil && q.Branch.Page(m.newNodeProviderPage(s, bn.n, nil, bn, q.Deep)) {
				return false
			}

			if q.Branch.Resource != nil {
				bn.resources.nodes.Walk(func(s string, v interface{}) bool {
					n := v.(*contentNode)
					return q.Branch.Resource(newNodeProviderResource(s, n, bn.n, bn))
				})
			}
		}

		if q.OnlyBranches || depth == depthBranch {
			return false
		}

		if q.Leaf.Page != nil || q.Leaf.Resource != nil {
			bn.pages.nodes.Walk(func(s string, v interface{}) bool {
				n := v.(*contentNode)
				if q.Leaf.Page != nil && q.Leaf.Page(m.newNodeProviderPage(s, n, bn, bn, q.Deep)) {
					return true
				}
				if q.Leaf.Resource != nil {
					// Interleave the Page's resources.
					bn.pageResources.nodes.WalkPrefix(s+"/", func(s string, v interface{}) bool {
						return q.Leaf.Resource(newNodeProviderResource(s, v.(*contentNode), n, bn))
					})
				}
				return false
			})
		}

		return false
	}

	if !q.Branch.Key.IsZero() {
		// Filter by section.
		if q.Branch.Key.IsPrefix() {
			if q.Branch.Key.Value != "" && q.Leaf.Page != nil {
				// Need to include the leaf pages of the owning branch.
				s := q.Branch.Key.Value[:len(q.Branch.Key.Value)-1]
				owner := m.Get(s)
				if owner != nil {
					if handleBranchPage(depthLeaf, s, owner) {
						// Done.
						return nil
					}
				}
			}

			var level int
			if q.NoRecurse {
				level = strings.Count(q.Branch.Key.Value, "/")
			}
			m.branches.WalkPrefix(
				q.Branch.Key.Value, func(s string, v interface{}) bool {
					if q.NoRecurse && strings.Count(s, "/") > level {
						return false
					}

					depth := depthAll
					if q.NoRecurse {
						depth = depthBranch
					}

					return handleBranchPage(depth, s, v)
				},
			)

			// Done.
			return nil
		}

		// Exact match.
		section := m.Get(q.Branch.Key.Value)
		if section != nil {
			if handleBranchPage(depthAll, q.Branch.Key.Value, section) {
				return nil
			}
		}
		// Done.
		return nil
	}

	if q.OnlyBranches || q.Leaf.Key.IsZero() || !q.Leaf.HasCallback() {
		// Done.
		return nil
	}

	_, section := m.LongestPrefix(q.Leaf.Key.Value)
	if section == nil {
		return nil
	}

	// Exact match.
	v, found := section.pages.nodes.Get(q.Leaf.Key.Value)
	if !found {
		return nil
	}
	if q.Leaf.Page != nil && q.Leaf.Page(m.newNodeProviderPage(q.Leaf.Key.Value, v.(*contentNode), section, section, q.Deep)) {
		return nil
	}

	if q.Leaf.Resource != nil {
		section.pageResources.nodes.WalkPrefix(q.Leaf.Key.Value+"/", func(s string, v interface{}) bool {
			return q.Leaf.Resource(newNodeProviderResource(s, v.(*contentNode), section.n, section))
		})
	}

	return nil
}

// WalkBranches invokes cb for all branch nodes.
func (m *branchMap) WalkBranches(cb func(s string, n *contentBranchNode) bool) {
	m.branches.Walk(func(s string, v interface{}) bool {
		return cb(s, v.(*contentBranchNode))
	})
}

// WalkBranches invokes cb for all branch nodes matching the given prefix.
func (m *branchMap) WalkBranchesPrefix(prefix string, cb func(s string, n *contentBranchNode) bool) {
	m.branches.WalkPrefix(prefix, func(s string, v interface{}) bool {
		return cb(s, v.(*contentBranchNode))
	})
}

func (m *branchMap) WalkPagesAllPrefixSection(
	prefix string,
	branchExclude, exclude contentTreeNodeCallback,
	callback contentTreeNodeCallbackNew) error {
	q := branchMapQuery{
		BranchExclude: branchExclude,
		Exclude:       exclude,
		Branch: branchMapQueryCallBacks{
			Key:  newBranchMapQueryKey(prefix, true),
			Page: callback,
		},
		Leaf: branchMapQueryCallBacks{
			Page: callback,
		},
	}
	return m.Walk(q)
}

func (m *branchMap) WalkPagesLeafsPrefixSection(
	prefix string,
	branchExclude, exclude contentTreeNodeCallback,
	callback contentTreeNodeCallbackNew) error {
	q := branchMapQuery{
		BranchExclude: branchExclude,
		Exclude:       exclude,
		Branch: branchMapQueryCallBacks{
			Key:  newBranchMapQueryKey(prefix, true),
			Page: nil,
		},
		Leaf: branchMapQueryCallBacks{
			Page: callback,
		},
	}
	return m.Walk(q)
}

func (m *branchMap) WalkPagesPrefixSectionNoRecurse(
	prefix string,
	branchExclude, exclude contentTreeNodeCallback,
	callback contentTreeNodeCallbackNew) error {
	q := branchMapQuery{
		NoRecurse:     true,
		BranchExclude: branchExclude,
		Exclude:       exclude,
		Branch: branchMapQueryCallBacks{
			Key:  newBranchMapQueryKey(prefix, true),
			Page: callback,
		},
		Leaf: branchMapQueryCallBacks{
			Page: callback,
		},
	}
	return m.Walk(q)
}

func (m *branchMap) Get(key string) *contentBranchNode {
	v, found := m.branches.Get(key)
	if !found {
		return nil
	}
	return v.(*contentBranchNode)
}

func (m *branchMap) Has(key string) bool {
	_, found := m.branches.Get(key)
	return found
}

func (m *branchMap) newNodeProviderPage(s string, n *contentNode, owner, branch *contentBranchNode, deep bool) contentNodeProvider {
	if !deep {
		return n
	}

	var np contentNodeProvider
	if owner == nil {
		if s != "" {
			_, owner = m.LongestPrefix(path.Dir(s))
		}
	}

	var ownerNode *contentNode
	if owner != nil {
		ownerNode = owner.n
	}

	var nInfo contentNodeInfoProvider = &contentNodeInfo{
		branch:   branch,
		isBranch: owner != branch,
	}

	np = struct {
		types.Identifier
		contentNodeInfoProvider
		contentGetNodeProvider
		contentGetContainerBranchProvider
		contentGetContainerNodeProvider
		contentGetBranchProvider
	}{
		n,
		nInfo,
		n,
		owner,
		ownerNode,
		branch,
	}

	return np
}

func (m *branchMap) debug(prefix string, w io.Writer) {
	fmt.Fprintf(w, "[%s] Start:\n", prefix)
	m.WalkBranches(func(s string, n *contentBranchNode) bool {
		var notes []string
		sectionType := "Section"
		if n.n.IsView() {
			sectionType = "View"
		}
		if n.n.p != nil {
			sectionType = n.n.p.Kind()
		}
		if n.n.p == nil {
			notes = append(notes, "MISSING_PAGE")
		}
		fmt.Fprintf(w, "[%s] %s: %q %v\n", prefix, sectionType, s, notes)
		n.pages.Walk(func(s string, n *contentNode) bool {
			fmt.Fprintf(w, "\t[%s] Page: %q\n", prefix, s)
			return false
		})
		n.pageResources.Walk(func(s string, n *contentNode) bool {
			fmt.Fprintf(w, "\t[%s] Branch Resource: %q\n", prefix, s)
			return false
		})
		n.pageResources.Walk(func(s string, n *contentNode) bool {
			fmt.Fprintf(w, "\t[%s] Leaf Resource: %q\n", prefix, s)
			return false
		})
		return false
	})
}

func (m *branchMap) debugDefault() {
	m.debug("", os.Stdout)
}

type branchMapQuery struct {
	// Restrict query to one level.
	NoRecurse bool
	// Deep/full callback objects.
	Deep bool
	// Do not navigate down to the leaf nodes.
	OnlyBranches bool
	// Global node filter. Return true to skip.
	Exclude contentTreeNodeCallback
	// Branch node filter. Return true to skip.
	BranchExclude contentTreeNodeCallback
	// Handle branch (sections and taxonomies) nodes.
	Branch branchMapQueryCallBacks
	// Handle leaf nodes (pages)
	Leaf branchMapQueryCallBacks
}

type branchMapQueryCallBacks struct {
	Key      branchMapQueryKey
	Page     contentTreeNodeCallbackNew
	Resource contentTreeNodeCallbackNew
}

func (q branchMapQueryCallBacks) HasCallback() bool {
	return q.Page != nil || q.Resource != nil
}

type branchMapQueryKey struct {
	Value string

	isSet    bool
	isPrefix bool
}

func (q branchMapQueryKey) Eq(key string) bool {
	if q.IsZero() || q.isPrefix {
		return false
	}
	return q.Value == key
}

func (q branchMapQueryKey) IsPrefix() bool {
	return !q.IsZero() && q.isPrefix
}

func (q branchMapQueryKey) IsZero() bool {
	return !q.isSet
}

type contentBranchNode struct {
	n             *contentNode
	resources     *contentBranchNodeTree
	pages         *contentBranchNodeTree
	pageResources *contentBranchNodeTree

	refs map[interface{}]ordinalWeight

	// Some default metadata if not provided in front matter.
	defaultTitle string
}

func (b *contentBranchNode) GetBranch() *contentBranchNode {
	return b
}

func (b *contentBranchNode) GetContainerBranch() *contentBranchNode {
	return b
}

func (b *contentBranchNode) InsertPage(key string, n *contentNode) {
	mustValidateSectionMapKey(key)
	b.pages.nodes.Insert(key, n)
}

func (b *contentBranchNode) InsertResource(key string, n *contentNode) error {
	mustValidateSectionMapKey(key)

	if _, _, found := b.pages.nodes.LongestPrefix(key); !found {
		return errors.Errorf("no page found for resource %q", key)
	}

	b.pageResources.nodes.Insert(key, n)

	return nil
}

func (m *contentBranchNode) newResource(n *contentNode, owner *pageState) (resource.Resource, error) {
	if owner == nil {
		panic("owner is nil")
	}

	fim := n.traits.(hugofs.FileInfoProvider).FileInfo()
	// TODO(bep) consolidate with multihost logic + clean up
	outputFormats := owner.m.outputFormats()
	seen := make(map[string]bool)
	var targetBasePaths []string

	// Make sure bundled resources are published to all of the output formats'
	// sub paths.
	for _, f := range outputFormats {
		p := f.Path
		if seen[p] {
			continue
		}
		seen[p] = true
		targetBasePaths = append(targetBasePaths, p)

	}

	meta := fim.Meta()
	r := func() (hugio.ReadSeekCloser, error) {
		return meta.Open()
	}

	target := strings.TrimPrefix(meta.Path, owner.File().Dir())

	return owner.s.ResourceSpec.New(
		resources.ResourceSourceDescriptor{
			TargetPaths:        owner.getTargetPaths,
			OpenReadSeekCloser: r,
			FileInfo:           fim,
			RelTargetFilename:  target,
			TargetBasePaths:    targetBasePaths,
			LazyPublish:        !owner.m.buildConfig.PublishResources,
			GroupIdentity:      n.GetIdentity(),
			DependencyManager:  n.GetDependencyManager(),
		})
}

type contentBranchNodeTree struct {
	nodes nodeTree
}

func (t contentBranchNodeTree) Walk(cb ...contentTreeNodeCallback) {
	cbs := newcontentTreeNodeCallbackChain(cb...)
	t.nodes.Walk(func(s string, v interface{}) bool {
		return cbs(s, v.(*contentNode))
	})
}

func (t contentBranchNodeTree) WalkPrefix(prefix string, cb ...contentTreeNodeCallback) {
	cbs := newcontentTreeNodeCallbackChain(cb...)
	t.nodes.WalkPrefix(prefix, func(s string, v interface{}) bool {
		return cbs(s, v.(*contentNode))
	})
}

func (t contentBranchNodeTree) Has(s string) bool {
	_, b := t.nodes.Get(s)
	return b
}

type defaultNodeTree struct {
	nodeTree
}

func (t *defaultNodeTree) Delete(s string) (interface{}, bool) {
	return t.nodeTree.Delete(s)
}

func (t *defaultNodeTree) DeletePrefix(s string) int {
	return t.nodeTree.DeletePrefix(s)
}

func (t *defaultNodeTree) Insert(s string, v interface{}) (interface{}, bool) {
	switch n := v.(type) {
	case *contentNode:
		n.key = s
	case *contentBranchNode:
		n.n.key = s
	}
	return t.nodeTree.Insert(s, v)
}

// Below some utils used for debugging.

// nodeTree defines the operations we use in radix.Tree.
type nodeTree interface {
	Delete(s string) (interface{}, bool)
	DeletePrefix(s string) int

	// Update ops.
	Insert(s string, v interface{}) (interface{}, bool)
	Len() int

	LongestPrefix(s string) (string, interface{}, bool)
	// Read ops
	Walk(fn radix.WalkFn)
	WalkPrefix(prefix string, fn radix.WalkFn)
	Get(s string) (interface{}, bool)
}

type nodeTreeUpdateTracer struct {
	name string
	nodeTree
}

func (t *nodeTreeUpdateTracer) Delete(s string) (interface{}, bool) {
	fmt.Printf("[%s]\t[Delete] %q\n", t.name, s)
	return t.nodeTree.Delete(s)
}

func (t *nodeTreeUpdateTracer) DeletePrefix(s string) int {
	n := t.nodeTree.DeletePrefix(s)
	fmt.Printf("[%s]\t[DeletePrefix] %q => %d\n", t.name, s, n)
	return n
}

func (t *nodeTreeUpdateTracer) Insert(s string, v interface{}) (interface{}, bool) {
	var typeInfo string
	switch n := v.(type) {
	case *contentNode:
		typeInfo = "n"
	case *contentBranchNode:
		typeInfo = fmt.Sprintf("b:isView:%t", n.n.IsView())
	}
	fmt.Printf("[%s]\t[Insert] %q %s\n", t.name, s, typeInfo)
	return t.nodeTree.Insert(s, v)
}

func mustValidateSectionMapKey(key string) {
	if err := validateSectionMapKey(key); err != nil {
		panic(err)
	}
}

func validateSectionMapKey(key string) error {
	if key == "" {
		// Home page.
		return nil
	}

	if len(key) < 2 {
		return errors.Errorf("too short key: %q", key)
	}

	if key[0] != '/' {
		return errors.Errorf("key must start with '/': %q", key)
	}

	if key[len(key)-1] == '/' {
		return errors.Errorf("key must not end with '/': %q", key)
	}

	return nil
}
