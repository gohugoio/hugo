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
	"html/template"
	"strings"
	texttemplate "text/template"
	"text/template/parse"
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
}

func (c templateContext) getIfNotVisited(name string) *parse.Tree {
	if c.visited[name] {
		return nil
	}
	c.visited[name] = true
	return c.lookupFn(name)
}

func newTemplateContext(lookupFn func(name string) *parse.Tree) *templateContext {
	return &templateContext{lookupFn: lookupFn, decl: make(map[string]string), visited: make(map[string]bool)}

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

func applyTemplateTransformersToHMLTTemplate(templ *template.Template) error {
	return applyTemplateTransformers(templ.Tree, createParseTreeLookup(templ))
}

func applyTemplateTransformersToTextTemplate(templ *texttemplate.Template) error {
	return applyTemplateTransformers(templ.Tree,
		func(nn string) *parse.Tree {
			tt := templ.Lookup(nn)
			if tt != nil {
				return tt.Tree
			}
			return nil
		})
}

func applyTemplateTransformers(templ *parse.Tree, lookupFn func(name string) *parse.Tree) error {
	if templ == nil {
		return errors.New("expected template, but none provided")
	}

	c := newTemplateContext(lookupFn)

	c.paramsKeysToLower(templ.Root)

	return nil
}

// paramsKeysToLower is made purposely non-generic to make it not so tempting
// to do more of these hard-to-maintain AST transformations.
func (c *templateContext) paramsKeysToLower(n parse.Node) {
	switch x := n.(type) {
	case *parse.ListNode:
		if x != nil {
			c.paramsKeysToLowerForNodes(x.Nodes...)
		}
	case *parse.ActionNode:
		c.paramsKeysToLowerForNodes(x.Pipe)
	case *parse.IfNode:
		c.paramsKeysToLowerForNodes(x.Pipe, x.List, x.ElseList)
	case *parse.WithNode:
		c.paramsKeysToLowerForNodes(x.Pipe, x.List, x.ElseList)
	case *parse.RangeNode:
		c.paramsKeysToLowerForNodes(x.Pipe, x.List, x.ElseList)
	case *parse.TemplateNode:
		subTempl := c.getIfNotVisited(x.Name)
		if subTempl != nil {
			c.paramsKeysToLowerForNodes(subTempl.Root)
		}
	case *parse.PipeNode:
		if len(x.Decl) == 1 && len(x.Cmds) == 1 {
			// maps $site => .Site etc.
			c.decl[x.Decl[0].Ident[0]] = x.Cmds[0].String()
		}

		for _, cmd := range x.Cmds {
			c.paramsKeysToLower(cmd)
		}

	case *parse.CommandNode:
		for _, elem := range x.Args {
			switch an := elem.(type) {
			case *parse.FieldNode:
				c.updateIdentsIfNeeded(an.Ident)
			case *parse.VariableNode:
				c.updateIdentsIfNeeded(an.Ident)
			case *parse.PipeNode:
				c.paramsKeysToLower(an)
			}

		}
	}
}

func (c *templateContext) paramsKeysToLowerForNodes(nodes ...parse.Node) {
	for _, node := range nodes {
		c.paramsKeysToLower(node)
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
