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
	"fmt"
	"html/template"
	"regexp"
	"strings"
	texttemplate "text/template"
	"text/template/parse"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// decl keeps track of the variable mappings, i.e. $mysite => .Site etc.
type decl map[string]string
type decln map[string]*parse.CommandNode

const (
	paramsIdentifier = "Params"
)

// Containers that may contain Params that we will not touch.
var reservedContainers = map[string]bool{
	// Aka .Site.Data.Params which must stay case sensitive.
	"Data": true,
}

type templateType int

const (
	templateUndefined templateType = iota
	templateShortcode
	templatePartial
)

type templateContext struct {
	decl             decl
	decln            decln // TODO1 remove
	visited          map[string]bool
	templateNotFound map[string]bool
	identityNotFound map[string]bool
	lookupFn         func(name string) *templateInfoTree

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

func (c templateContext) getIfNotVisited(name string) *templateInfoTree {
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
	} else {
		if templ.info.Manager != nil {
			c.Info.Add(templ.info)
		}
	}
	return templ
}

func newTemplateContext(info tpl.Info, lookupFn func(name string) *templateInfoTree) *templateContext {
	if info.Manager == nil {
		panic("identity manager not set")
	}
	return &templateContext{
		Info:             info,
		lookupFn:         lookupFn,
		decl:             make(map[string]string),
		decln:            make(map[string]*parse.CommandNode),
		visited:          make(map[string]bool),
		templateNotFound: make(map[string]bool),
		identityNotFound: make(map[string]bool),
	}
}

func createParseTreeLookup(templ *template.Template) func(nn string) *templateInfoTree {
	return createParseTreeLookupFor(templ, func(name string) tpl.Info { return newTemplateInfo(name) })

}

func createParseTreeLookupFor(templ *template.Template, infoFn func(name string) tpl.Info) func(nn string) *templateInfoTree {
	return func(nn string) *templateInfoTree {
		tt := templ.Lookup(nn)
		if tt != nil {
			return &templateInfoTree{
				tree: tt.Tree,
				info: infoFn(nn),
			}
		}
		return nil
	}
}

func (t *templateHandler) createParseTreeLookup(templ *template.Template) func(nn string) *templateInfoTree {
	return createParseTreeLookupFor(templ, func(name string) tpl.Info { return t.templateInfo[name] })
}

func (t *templateHandler) applyTemplateTransformersToHMLTTemplate(typ templateType, templ *template.Template) (*templateContext, error) {
	ti := &templateInfoTree{
		tree: templ.Tree,
		info: t.getOrCreateTemplateInfo(templ.Name()),
	}
	return applyTemplateTransformers(typ, ti, t.createParseTreeLookup(templ))
}

func (t *templateHandler) applyTemplateTransformersToTextTemplate(typ templateType, templ *texttemplate.Template) (*templateContext, error) {
	ti := &templateInfoTree{
		tree: templ.Tree,
		info: t.getOrCreateTemplateInfo(templ.Name()),
	}

	return applyTemplateTransformers(typ, ti,
		func(nn string) *templateInfoTree {
			tt := templ.Lookup(nn)
			if tt != nil {
				return &templateInfoTree{
					tree: tt.Tree,
					info: t.getOrCreateTemplateInfo(nn),
				}
			}
			return nil
		})
}

type templateInfoTree struct {
	info tpl.Info
	tree *parse.Tree
}

func applyTemplateTransformers(
	typ templateType,
	templ *templateInfoTree,
	lookupFn func(name string) *templateInfoTree) (*templateContext, error) {

	if templ == nil {
		return nil, errors.New("expected template, but none provided")
	}

	c := newTemplateContext(templ.info, lookupFn)
	c.typ = typ

	_, err := c.applyTransformations(templ.tree.Root)

	if err == nil && c.returnNode != nil {
		// This is a partial with a return statement.
		c.Info.HasReturn = true
		templ.tree.Root = c.wrapInPartialReturnWrapper(templ.tree.Root)
	}

	return c, err
}

const (
	partialReturnWrapperTempl = `{{ $_hugo_dot := $ }}{{ $ := .Arg }}{{ with .Arg }}{{ $_hugo_dot.Set ("PLACEHOLDER") }}{{ end }}`
	dotContextWrapperTempl    = `{{ invokeDot . "NAME" "ARGS" }}`
)

var (
	partialReturnWrapper *parse.ListNode
	dotContextWrapper    *parse.CommandNode
)

func init() {
	tt := texttemplate.New("").Funcs(texttemplate.FuncMap{
		"invokeDot": func() interface{} { return "foo" },
	})

	templ, err := tt.Parse(partialReturnWrapperTempl)
	if err != nil {
		panic(err)
	}
	partialReturnWrapper = templ.Tree.Root

	templ, err = tt.Parse(dotContextWrapperTempl)
	if err != nil {
		panic(err)
	}
	action := templ.Tree.Root.Nodes[0].(*parse.ActionNode)
	dotContextWrapper = action.Pipe.Cmds[0]
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

var (
	ignoreFuncsRe  = regexp.MustCompile("invokeDot|html")
	goBuiltInFuncs = regexp.MustCompile("len")
)

// TODO1 somehow inject (wrap dot?) a receiver object that can be
// set in .Execute. To avoid global funcs.
func (c *templateContext) wrapDot(d bool, cmd *parse.CommandNode) {
	var dotNode parse.Node
	doDebug := d || strings.Contains(cmd.String(), "blue")
	var fields string

	firstWord := cmd.Args[0]

	switch a := firstWord.(type) {
	case *parse.FieldNode:
		fields = a.String()
		//return s.evalFieldNode(dot, n, cmd.Args, final)
	case *parse.ChainNode:
		if pipe, ok := a.Node.(*parse.PipeNode); ok {
			for _, cmd := range pipe.Cmds {
				c.wrapDot(doDebug, cmd)
			}

		}
		return // TODO1
	//	fields = a.String()
	//return s.evalChainNode(dot, n, cmd.Args, final)
	case *parse.IdentifierNode:
		// Must be a function.
		if ignoreFuncsRe.MatchString(a.Ident) {
			return
		}
		if goBuiltInFuncs.MatchString(a.Ident) {
			fmt.Println(a.Ident, "==>", cmd.Args[1:])
			return
		}
		fields = a.Ident
	case *parse.PipeNode:
		for _, cmd := range a.Cmds {
			c.wrapDot(doDebug, cmd)
		}
		//s.notAFunction(cmd.Args, final)
		//return s.evalPipeline(dot, n)
		return
	case *parse.VariableNode:
		// $x.Field has $x as the first ident, Field as the second.
		fields = "." + strings.Join(a.Ident[1:], ".")
		a.Ident = a.Ident[:1]
		dotNode = a

	default:
		//fmt.Printf("UNKNOWN: %T\n", firstWord)
		return
	}

	wrapper := dotContextWrapper.Copy().(*parse.CommandNode)

	sn := wrapper.Args[2].(*parse.StringNode)
	if dotNode != nil {
		wrapper.Args[1] = dotNode
	}

	sn.Quoted = "\"" + fields + "\""
	sn.Text = fields

	args := wrapper.Args[:3]
	if len(cmd.Args) > 1 {
		args = append(args, cmd.Args[1:]...)
	}

	cmd.Args = args

	return

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

var partialRe = regexp.MustCompile(`^partial(Cached)?$|^partials\.Include(Cached)?$`)

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
			c.applyTransformationsToNodes(subTempl.tree.Root)
		}
	case *parse.PipeNode:
		c.collectConfig(x)
		if len(x.Decl) == 1 && len(x.Cmds) == 1 {
			// maps $site => .Site etc.
			c.decl[x.Decl[0].Ident[0]] = x.Cmds[0].String()
			c.decln[x.Decl[0].Ident[0]] = x.Cmds[0]
		}

		for i, cmd := range x.Cmds {
			keep, _ := c.applyTransformations(cmd)
			if !keep {
				x.Cmds = append(x.Cmds[:i], x.Cmds[i+1:]...)
			}
		}

	case *parse.CommandNode:
		c.wrapDot(false, x)
		if false || len(x.Args) > 1 {
			first := x.Args[0]
			var id string
			switch v := first.(type) {
			case *parse.IdentifierNode:
				id = v.Ident
			case *parse.ChainNode:
				id = v.String()
			}

			if partialRe.MatchString(id) {
				partialName := strings.Trim(x.Args[1].String(), "\"")
				if !strings.Contains(partialName, ".") {
					partialName += ".html"
				}
				partialName = "partials/" + partialName
				info := c.lookupFn(partialName)
				if info != nil {
					c.Info.Add(info.info)
				} else {
					// Delay for later
					c.identityNotFound[partialName] = true
				}
			}
		}

		c.collectInner(x)
		keep := c.collectReturnNode(x)

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
		return keep, c.err
	}

	return true, c.err
}

func (c *templateContext) applyTransformationsToNodes(nodes ...parse.Node) {
	for _, node := range nodes {
		c.applyTransformations(node)
	}
}

func (c *templateContext) updateIdentsIfNeeded(idents []string) {
	if true {
		return // TODO1 remove all this .Params stuff.
	}
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
			continue
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
