// Copyright 2016 The Hugo Authors. All rights reserved.
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

package tplimpl

import (
	"errors"
	"fmt"
	"strings"

	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
)

type templateType int

const (
	templateUndefined templateType = iota
	templateShortcode
	templatePartial
)

type templateContext struct {
	visited          map[string]bool
	templateNotFound map[string]bool
	deferNodes       map[string]*parse.ListNode
	lookupFn         func(name string) *templateState

	// The last error encountered.
	err error

	// Set when we're done checking for config header.
	configChecked bool

	t *templateState

	// Store away the return node in partials.
	returnNode *parse.CommandNode
}

func (c templateContext) getIfNotVisited(name string) *templateState {
	if c.visited[name] {
		return nil
	}
	c.visited[name] = true
	templ := c.lookupFn(name)
	if templ == nil {
		// This may be a inline template defined outside of this file
		// and not yet parsed. Unusual, but it happens.
		// Store the name to try again later.
		c.templateNotFound[name] = true
	}

	return templ
}

func newTemplateContext(
	t *templateState,
	lookupFn func(name string) *templateState,
) *templateContext {
	return &templateContext{
		t:                t,
		lookupFn:         lookupFn,
		visited:          make(map[string]bool),
		templateNotFound: make(map[string]bool),
		deferNodes:       make(map[string]*parse.ListNode),
	}
}

func applyTemplateTransformers(
	t *templateState,
	lookupFn func(name string) *templateState,
) (*templateContext, error) {
	if t == nil {
		return nil, errors.New("expected template, but none provided")
	}

	c := newTemplateContext(t, lookupFn)
	tree := getParseTree(t.Template)

	_, err := c.applyTransformations(tree.Root)

	if err == nil && c.returnNode != nil {
		// This is a partial with a return statement.
		c.t.parseInfo.HasReturn = true
		tree.Root = c.wrapInPartialReturnWrapper(tree.Root)
	}

	return c, err
}

func getParseTree(templ tpl.Template) *parse.Tree {
	templ = unwrap(templ)
	if text, ok := templ.(*texttemplate.Template); ok {
		return text.Tree
	}
	return templ.(*htmltemplate.Template).Tree
}

const (
	// We parse this template and modify the nodes in order to assign
	// the return value of a partial to a contextWrapper via Set. We use
	// "range" over a one-element slice so we can shift dot to the
	// partial's argument, Arg, while allowing Arg to be falsy.
	partialReturnWrapperTempl = `{{ $_hugo_dot := $ }}{{ $ := .Arg }}{{ range (slice .Arg) }}{{ $_hugo_dot.Set ("PLACEHOLDER") }}{{ end }}`

	doDeferTempl = `{{ doDefer ("PLACEHOLDER1") ("PLACEHOLDER2") }}`
)

var (
	partialReturnWrapper *parse.ListNode
	doDefer              *parse.ListNode
)

func init() {
	templ, err := texttemplate.New("").Parse(partialReturnWrapperTempl)
	if err != nil {
		panic(err)
	}
	partialReturnWrapper = templ.Tree.Root

	templ, err = texttemplate.New("").Funcs(texttemplate.FuncMap{"doDefer": func(string, string) string { return "" }}).Parse(doDeferTempl)
	if err != nil {
		panic(err)
	}
	doDefer = templ.Tree.Root
}

// wrapInPartialReturnWrapper copies and modifies the parsed nodes of a
// predefined partial return wrapper to insert those of a user-defined partial.
func (c *templateContext) wrapInPartialReturnWrapper(n *parse.ListNode) *parse.ListNode {
	wrapper := partialReturnWrapper.CopyList()
	rangeNode := wrapper.Nodes[2].(*parse.RangeNode)
	retn := rangeNode.List.Nodes[0]
	setCmd := retn.(*parse.ActionNode).Pipe.Cmds[0]
	setPipe := setCmd.Args[1].(*parse.PipeNode)
	// Replace PLACEHOLDER with the real return value.
	// Note that this is a PipeNode, so it will be wrapped in parens.
	setPipe.Cmds = []*parse.CommandNode{c.returnNode}
	rangeNode.List.Nodes = append(n.Nodes, retn)

	return wrapper
}

// applyTransformations do 2 things:
// 1) Parses partial return statement.
// 2) Tracks template (partial) dependencies and some other info.
func (c *templateContext) applyTransformations(n parse.Node) (bool, error) {
	switch x := n.(type) {
	case *parse.ListNode:
		if x != nil {
			c.applyTransformationsToNodes(x.Nodes...)
		}
	case *parse.ActionNode:
		c.applyTransformationsToNodes(x.Pipe)
	case *parse.IfNode:
		c.applyTransformationsToNodes(x.Pipe, x.List, x.ElseList)
	case *parse.WithNode:
		c.handleDefer(x)
		c.applyTransformationsToNodes(x.Pipe, x.List, x.ElseList)
	case *parse.RangeNode:
		c.applyTransformationsToNodes(x.Pipe, x.List, x.ElseList)
	case *parse.TemplateNode:
		subTempl := c.getIfNotVisited(x.Name)
		if subTempl != nil {
			c.applyTransformationsToNodes(getParseTree(subTempl.Template).Root)
		}
	case *parse.PipeNode:
		c.collectConfig(x)
		for i, cmd := range x.Cmds {
			keep, _ := c.applyTransformations(cmd)
			if !keep {
				x.Cmds = append(x.Cmds[:i], x.Cmds[i+1:]...)
			}
		}

	case *parse.CommandNode:
		c.collectInner(x)
		keep := c.collectReturnNode(x)

		for _, elem := range x.Args {
			switch an := elem.(type) {
			case *parse.PipeNode:
				c.applyTransformations(an)
			}
		}
		return keep, c.err
	}

	return true, c.err
}

func (c *templateContext) handleDefer(withNode *parse.WithNode) {
	if len(withNode.Pipe.Cmds) != 1 {
		return
	}
	cmd := withNode.Pipe.Cmds[0]
	if len(cmd.Args) != 1 {
		return
	}
	idArg := cmd.Args[0]

	p, ok := idArg.(*parse.PipeNode)
	if !ok {
		return
	}

	if len(p.Cmds) != 1 {
		return
	}

	cmd = p.Cmds[0]

	if len(cmd.Args) != 2 {
		return
	}

	idArg = cmd.Args[0]

	id, ok := idArg.(*parse.ChainNode)
	if !ok || len(id.Field) != 1 || id.Field[0] != "Defer" {
		return
	}
	if id2, ok := id.Node.(*parse.IdentifierNode); !ok || id2.Ident != "templates" {
		return
	}

	deferArg := cmd.Args[1]
	cmd.Args = []parse.Node{idArg}

	l := doDefer.CopyList()
	n := l.Nodes[0].(*parse.ActionNode)

	inner := withNode.List.CopyList()
	s := inner.String()
	if strings.Contains(s, "resources.PostProcess") {
		c.err = errors.New("resources.PostProcess cannot be used in a deferred template")
		return
	}
	innerHash := hashing.XxHashFromStringHexEncoded(s)
	deferredID := tpl.HugoDeferredTemplatePrefix + innerHash

	c.deferNodes[deferredID] = inner
	withNode.List = l

	n.Pipe.Cmds[0].Args[1].(*parse.PipeNode).Cmds[0].Args[0].(*parse.StringNode).Text = deferredID
	n.Pipe.Cmds[0].Args[2] = deferArg
}

func (c *templateContext) applyTransformationsToNodes(nodes ...parse.Node) {
	for _, node := range nodes {
		c.applyTransformations(node)
	}
}

func (c *templateContext) hasIdent(idents []string, ident string) bool {
	for _, id := range idents {
		if id == ident {
			return true
		}
	}
	return false
}

// collectConfig collects and parses any leading template config variable declaration.
// This will be the first PipeNode in the template, and will be a variable declaration
// on the form:
//
//	{{ $_hugo_config:= `{ "version": 1 }` }}
func (c *templateContext) collectConfig(n *parse.PipeNode) {
	if c.t.typ != templateShortcode {
		return
	}
	if c.configChecked {
		return
	}
	c.configChecked = true

	if len(n.Decl) != 1 || len(n.Cmds) != 1 {
		// This cannot be a config declaration
		return
	}

	v := n.Decl[0]

	if len(v.Ident) == 0 || v.Ident[0] != "$_hugo_config" {
		return
	}

	cmd := n.Cmds[0]

	if len(cmd.Args) == 0 {
		return
	}

	if s, ok := cmd.Args[0].(*parse.StringNode); ok {
		errMsg := "failed to decode $_hugo_config in template: %w"
		m, err := maps.ToStringMapE(s.Text)
		if err != nil {
			c.err = fmt.Errorf(errMsg, err)
			return
		}
		if err := mapstructure.WeakDecode(m, &c.t.parseInfo.Config); err != nil {
			c.err = fmt.Errorf(errMsg, err)
		}
	}
}

// collectInner determines if the given CommandNode represents a
// shortcode call to its .Inner.
func (c *templateContext) collectInner(n *parse.CommandNode) {
	if c.t.typ != templateShortcode {
		return
	}
	if c.t.parseInfo.IsInner || len(n.Args) == 0 {
		return
	}

	for _, arg := range n.Args {
		var idents []string
		switch nt := arg.(type) {
		case *parse.FieldNode:
			idents = nt.Ident
		case *parse.VariableNode:
			idents = nt.Ident
		}

		if c.hasIdent(idents, "Inner") || c.hasIdent(idents, "InnerDeindent") {
			c.t.parseInfo.IsInner = true
			break
		}
	}
}

func (c *templateContext) collectReturnNode(n *parse.CommandNode) bool {
	if c.t.typ != templatePartial || c.returnNode != nil {
		return true
	}

	if len(n.Args) < 2 {
		return true
	}

	ident, ok := n.Args[0].(*parse.IdentifierNode)
	if !ok || ident.Ident != "return" {
		return true
	}

	c.returnNode = n
	// Remove the "return" identifiers
	c.returnNode.Args = c.returnNode.Args[1:]

	return false
}

func findTemplateIn(name string, in tpl.Template) (tpl.Template, bool) {
	in = unwrap(in)
	if text, ok := in.(*texttemplate.Template); ok {
		if templ := text.Lookup(name); templ != nil {
			return templ, true
		}
		return nil, false
	}
	if templ := in.(*htmltemplate.Template).Lookup(name); templ != nil {
		return templ, true
	}
	return nil, false
}
