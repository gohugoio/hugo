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
	"html/template"
	"strings"
	texttemplate "text/template"
	"text/template/parse"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// decl keeps track of the variable mappings, i.e. $mysite => .Site etc.
type decl map[string]string

const (
	paramsIdentifier = "Params"
)

// Containers that may contain Params that we will not touch.
var reservedContainers = map[string]bool{
	// Aka .Site.Data.Params which must stay case sensitive.
	"Data": true,
}

type templateContext struct {
	decl     decl
	visited  map[string]bool
	lookupFn func(name string) *parse.Tree

	// The last error encountered.
	err error

	// Only needed for shortcodes
	isShortcode bool

	// Set when we're done checking for config header.
	configChecked bool

	// Contains some info about the template
	tpl.Info
}

func (c templateContext) getIfNotVisited(name string) *parse.Tree {
	if c.visited[name] {
		return nil
	}
	c.visited[name] = true
	return c.lookupFn(name)
}

func newTemplateContext(lookupFn func(name string) *parse.Tree) *templateContext {
	return &templateContext{
		Info:     tpl.Info{Config: tpl.DefaultConfig},
		lookupFn: lookupFn,
		decl:     make(map[string]string),
		visited:  make(map[string]bool)}

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

func applyTemplateTransformersToHMLTTemplate(isShortcode bool, templ *template.Template) (tpl.Info, error) {
	return applyTemplateTransformers(isShortcode, templ.Tree, createParseTreeLookup(templ))
}

func applyTemplateTransformersToTextTemplate(isShortcode bool, templ *texttemplate.Template) (tpl.Info, error) {
	return applyTemplateTransformers(isShortcode, templ.Tree,
		func(nn string) *parse.Tree {
			tt := templ.Lookup(nn)
			if tt != nil {
				return tt.Tree
			}
			return nil
		})
}

func applyTemplateTransformers(isShortcode bool, templ *parse.Tree, lookupFn func(name string) *parse.Tree) (tpl.Info, error) {
	if templ == nil {
		return tpl.Info{}, errors.New("expected template, but none provided")
	}

	c := newTemplateContext(lookupFn)
	c.isShortcode = isShortcode

	err := c.applyTransformations(templ.Root)

	return c.Info, err
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
func (c *templateContext) applyTransformations(n parse.Node) error {
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
		if len(x.Decl) == 1 && len(x.Cmds) == 1 {
			// maps $site => .Site etc.
			c.decl[x.Decl[0].Ident[0]] = x.Cmds[0].String()
		}

		for _, cmd := range x.Cmds {
			c.applyTransformations(cmd)
		}

	case *parse.CommandNode:
		c.collectInner(x)

		for _, elem := range x.Args {
			switch an := elem.(type) {
			case *parse.FieldNode:
				c.updateIdentsIfNeeded(an.Ident)
			case *parse.VariableNode:
				c.updateIdentsIfNeeded(an.Ident)
			case *parse.PipeNode:
				c.applyTransformations(an)
			case *parse.ChainNode:
				// site.Params...
				if len(an.Field) > 1 && an.Field[0] == paramsIdentifier {
					c.updateIdentsIfNeeded(an.Field)
				}
			}
		}
	}

	return c.err
}

func (c *templateContext) applyTransformationsToNodes(nodes ...parse.Node) {
	for _, node := range nodes {
		c.applyTransformations(node)
	}
}

func (c *templateContext) updateIdentsIfNeeded(idents []string) {
	index := c.decl.indexOfReplacementStart(idents)

	if index == -1 {
		return
	}

	for i := index; i < len(idents); i++ {
		idents[i] = strings.ToLower(idents[i])
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
	if !c.isShortcode {
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
		m, err := cast.ToStringMapE(s.Text)
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
	if !c.isShortcode {
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

// indexOfReplacementStart will return the index of where to start doing replacement,
// -1 if none needed.
func (d decl) indexOfReplacementStart(idents []string) int {

	l := len(idents)

	if l == 0 {
		return -1
	}

	if l == 1 {
		first := idents[0]
		if first == "" || first == paramsIdentifier || first[0] == '$' {
			// This can not be a Params.x
			return -1
		}
	}

	var lookFurther bool
	var needsVarExpansion bool
	for _, ident := range idents {
		if ident[0] == '$' {
			lookFurther = true
			needsVarExpansion = true
			break
		} else if ident == paramsIdentifier {
			lookFurther = true
			break
		}
	}

	if !lookFurther {
		return -1
	}

	var resolvedIdents []string

	if !needsVarExpansion {
		resolvedIdents = idents
	} else {
		var ok bool
		resolvedIdents, ok = d.resolveVariables(idents)
		if !ok {
			return -1
		}
	}

	var paramFound bool
	for i, ident := range resolvedIdents {
		if ident == paramsIdentifier {
			if i > 0 {
				container := resolvedIdents[i-1]
				if reservedContainers[container] {
					// .Data.Params.someKey
					return -1
				}
			}

			paramFound = true
			break
		}
	}

	if !paramFound {
		return -1
	}

	var paramSeen bool
	idx := -1
	for i, ident := range idents {
		if ident == "" || ident[0] == '$' {
			continue
		}

		if ident == paramsIdentifier {
			paramSeen = true
			idx = -1

		} else {
			if paramSeen {
				return i
			}
			if idx == -1 {
				idx = i
			}
		}
	}
	return idx

}

func (d decl) resolveVariables(idents []string) ([]string, bool) {
	var (
		replacements []string
		replaced     []string
	)

	// An Ident can start out as one of
	// [Params] [$blue] [$colors.Blue]
	// We need to resolve the variables, so
	// $blue => [Params Colors Blue]
	// etc.
	replacements = []string{idents[0]}

	// Loop until there are no more $vars to resolve.
	for i := 0; i < len(replacements); i++ {

		if i > 20 {
			// bail out
			return nil, false
		}

		potentialVar := replacements[i]

		if potentialVar == "$" {
			continue
		}

		if potentialVar == "" || potentialVar[0] != '$' {
			// leave it as is
			replaced = append(replaced, strings.Split(potentialVar, ".")...)
			continue
		}

		replacement, ok := d[potentialVar]

		if !ok {
			// Temporary range vars. We do not care about those.
			return nil, false
		}

		if !d.isKeyword(replacement) {
			// This can not be .Site.Params etc.
			return nil, false
		}

		replacement = strings.TrimPrefix(replacement, ".")

		if replacement == "" {
			continue
		}

		if replacement[0] == '$' {
			// Needs further expansion
			replacements = append(replacements, strings.Split(replacement, ".")...)
		} else {
			replaced = append(replaced, strings.Split(replacement, ".")...)
		}
	}

	return append(replaced, idents[1:]...), true

}

func (d decl) isKeyword(s string) bool {
	return !strings.ContainsAny(s, " -\"")
}
