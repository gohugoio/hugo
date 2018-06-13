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

// +build !go1.11

package tplimpl

import (
	"text/template/parse"
)

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
		for i, elem := range x.Decl {
			if len(x.Cmds) > i {
				// maps $site => .Site etc.
				c.decl[elem.Ident[0]] = x.Cmds[i].String()
			}
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
