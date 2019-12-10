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
	template "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"

	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type templateType int

const (
	templateUndefined templateType = iota
	templateShortcode
	templatePartial
)

type templateContext struct {
	visited  map[string]bool
	notFound map[string]bool
	lookupFn func(name string) *parse.Tree

	// The last error encountered.
	err error

	typ templateType

	// Set when we're done checking for config header.
	configChecked bool

	// Contains some info about the template
	tpl.Info

	// Store away the return node in partials.
	returnNode *parse.CommandNode
}

func (c templateContext) getIfNotVisited(name string) *parse.Tree {
	if c.visited[name] {
		return nil
	}
	c.visited[name] = true
	templ := c.lookupFn(name)
	if templ == nil {
		// This may be a inline template defined outside of this file
		// and not yet parsed. Unusual, but it happens.
		// Store the name to try again later.
		c.notFound[name] = true
	}

	return templ
}

func newTemplateContext(lookupFn func(name string) *parse.Tree) *templateContext {
	return &templateContext{
		Info:     tpl.Info{Config: tpl.DefaultConfig},
		lookupFn: lookupFn,
		visited:  make(map[string]bool),
		notFound: make(map[string]bool)}
}

func createParseTreeLookup(templ *template.Template) func(nn string) *parse.Tree {
	return func(nn string) *parse.Tree {
		tt := templ.Lookup(nn)
		if tt != nil {
			return tt.Tree
		}
		return nil
	}
}

func applyTemplateTransformersToHMLTTemplate(typ templateType, templ *template.Template) (*templateContext, error) {
	return applyTemplateTransformers(typ, templ.Tree, createParseTreeLookup(templ))
}

func applyTemplateTransformersToTextTemplate(typ templateType, templ *texttemplate.Template) (*templateContext, error) {
	return applyTemplateTransformers(typ, templ.Tree,
		func(nn string) *parse.Tree {
			tt := templ.Lookup(nn)
			if tt != nil {
				return tt.Tree
			}
			return nil
		})
}

func applyTemplateTransformers(typ templateType, templ *parse.Tree, lookupFn func(name string) *parse.Tree) (*templateContext, error) {
	if templ == nil {
		return nil, errors.New("expected template, but none provided")
	}

	c := newTemplateContext(lookupFn)
	c.typ = typ

	_, err := c.applyTransformations(templ.Root)

	if err == nil && c.returnNode != nil {
		// This is a partial with a return statement.
		c.Info.HasReturn = true
		templ.Root = c.wrapInPartialReturnWrapper(templ.Root)
	}

	return c, err
}

const (
	partialReturnWrapperTempl = `{{ $_hugo_dot := $ }}{{ $ := .Arg }}{{ with .Arg }}{{ $_hugo_dot.Set ("PLACEHOLDER") }}{{ end }}`
)

var partialReturnWrapper *parse.ListNode

func init() {
	templ, err := texttemplate.New("").Parse(partialReturnWrapperTempl)
	if err != nil {
		panic(err)
	}
	partialReturnWrapper = templ.Tree.Root
}

func (c *templateContext) wrapInPartialReturnWrapper(n *parse.ListNode) *parse.ListNode {
	wrapper := partialReturnWrapper.CopyList()
	withNode := wrapper.Nodes[2].(*parse.WithNode)
	retn := withNode.List.Nodes[0]
	setCmd := retn.(*parse.ActionNode).Pipe.Cmds[0]
	setPipe := setCmd.Args[1].(*parse.PipeNode)
	// Replace PLACEHOLDER with the real return value.
	// Note that this is a PipeNode, so it will be wrapped in parens.
	setPipe.Cmds = []*parse.CommandNode{c.returnNode}
	withNode.List.Nodes = append(n.Nodes, retn)

	return wrapper

}

// The truth logic in Go's template package is broken for certain values
// for the if and with keywords. This works around that problem by wrapping
// the node passed to if/with in a getif conditional.
// getif works slightly different than the Go built-in in that it also
// considers any IsZero methods on the values (as in time.Time).
// See https://github.com/gohugoio/hugo/issues/5738
func (c *templateContext) wrapWithGetIf(p *parse.PipeNode) {
	if len(p.Cmds) == 0 {
		return
	}

	// getif will return an empty string if not evaluated as truthful,
	// which is when we need the value in the with clause.
	firstArg := parse.NewIdentifier("getif")
	secondArg := p.CopyPipe()
	newCmd := p.Cmds[0].Copy().(*parse.CommandNode)

	// secondArg is a PipeNode and will behave as it was wrapped in parens, e.g:
	// {{ getif (len .Params | eq 2) }}
	newCmd.Args = []parse.Node{firstArg, secondArg}

	p.Cmds = []*parse.CommandNode{newCmd}

}

// applyTransformations do 3 things:
// 1) Make all .Params.CamelCase and similar into lowercase.
// 2) Wraps every with and if pipe in getif
// 3) Collects some information about the template content.
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
		c.wrapWithGetIf(x.Pipe)
	case *parse.WithNode:
		c.applyTransformationsToNodes(x.Pipe, x.List, x.ElseList)
		c.wrapWithGetIf(x.Pipe)
	case *parse.RangeNode:
		c.applyTransformationsToNodes(x.Pipe, x.List, x.ElseList)
	case *parse.TemplateNode:
		subTempl := c.getIfNotVisited(x.Name)
		if subTempl != nil {
			c.applyTransformationsToNodes(subTempl.Root)
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
//    {{ $_hugo_config:= `{ "version": 1 }` }}
func (c *templateContext) collectConfig(n *parse.PipeNode) {
	if c.typ != templateShortcode {
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
		errMsg := "failed to decode $_hugo_config in template"
		m, err := maps.ToStringMapE(s.Text)
		if err != nil {
			c.err = errors.Wrap(err, errMsg)
			return
		}
		if err := mapstructure.WeakDecode(m, &c.Info.Config); err != nil {
			c.err = errors.Wrap(err, errMsg)
		}
	}

}

// collectInner determines if the given CommandNode represents a
// shortcode call to its .Inner.
func (c *templateContext) collectInner(n *parse.CommandNode) {
	if c.typ != templateShortcode {
		return
	}
	if c.Info.IsInner || len(n.Args) == 0 {
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

		if c.hasIdent(idents, "Inner") {
			c.Info.IsInner = true
			break
		}
	}

}

func (c *templateContext) collectReturnNode(n *parse.CommandNode) bool {
	if c.typ != templatePartial || c.returnNode != nil {
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
