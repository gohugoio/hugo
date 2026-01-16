package tplimpl

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"

	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
)

type templateTransformContext struct {
	visited          map[string]bool
	templateNotFound map[string]bool
	deferNodes       map[string]*parse.ListNode
	lookupFn         func(name string, in *TemplInfo) *TemplInfo
	store            *TemplateStore

	// The last error encountered.
	err error

	// Set when we're done checking for config header.
	configChecked bool

	t *TemplInfo

	// Store away the return node in partials.
	returnNode *parse.CommandNode
}

func (c templateTransformContext) getIfNotVisited(name string) *TemplInfo {
	if c.visited[name] {
		return nil
	}
	c.visited[name] = true
	templ := c.lookupFn(name, c.t)
	if templ == nil {
		// This may be a inline template defined outside of this file
		// and not yet parsed. Unusual, but it happens.
		// Store the name to try again later.
		c.templateNotFound[name] = true
	}

	return templ
}

func newTemplateTransformContext(
	t *TemplInfo,
	store *TemplateStore,
	lookupFn func(name string, in *TemplInfo) *TemplInfo,
) *templateTransformContext {
	return &templateTransformContext{
		t:                t,
		lookupFn:         lookupFn,
		store:            store,
		visited:          make(map[string]bool),
		templateNotFound: make(map[string]bool),
		deferNodes:       make(map[string]*parse.ListNode),
	}
}

func applyTemplateTransformers(
	t *TemplInfo,
	store *TemplateStore,
	lookupFn func(name string, in *TemplInfo) *TemplInfo,
) (*templateTransformContext, error) {
	if t == nil {
		return nil, errors.New("expected template, but none provided")
	}

	c := newTemplateTransformContext(t, store, lookupFn)
	c.t.ParseInfo = defaultParseInfo
	tree := getParseTree(t.Template)
	if tree == nil {
		panic(fmt.Errorf("template %s not parsed", t))
	}

	if err := c.applyTransformationsAndSetReturnWrapper(tree); err != nil {
		return c, fmt.Errorf("failed to transform template %q: %w", t.Name(), err)
	}

	return c, c.err
}

func getParseTree(templ tpl.Template) *parse.Tree {
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

	// _pushPartialDecorator is always falsy.
	pushPartialDecoratorTempl = `{{ if or (_pushPartialDecorator ("PLACEHOLDER")) }}{{ end }}`
	popPartialDecoratorTempl  = `{{ if (_popPartialDecorator ("PLACEHOLDER1")) }}{{ . }}{{ else }}("PLACEHOLDER2"){{ end }}`
)

var (
	partialReturnWrapper *parse.ListNode
	doDefer              *parse.ListNode
	popPartialDecorator  *parse.ListNode
	pushPartialDecorator *parse.ListNode
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

	templ, err = texttemplate.New("").Funcs(texttemplate.FuncMap{"_popPartialDecorator": func(string) string { return "" }}).Parse(popPartialDecoratorTempl)
	if err != nil {
		panic(err)
	}
	popPartialDecorator = templ.Tree.Root

	templ, err = texttemplate.New("").Funcs(texttemplate.FuncMap{"_pushPartialDecorator": func(string) string { return "" }}).Parse(pushPartialDecoratorTempl)
	if err != nil {
		panic(err)
	}
	pushPartialDecorator = templ.Tree.Root
}

// wrapInPartialReturnWrapper copies and modifies the parsed nodes of a
// predefined partial return wrapper to insert those of a user-defined partial.
func (c *templateTransformContext) wrapInPartialReturnWrapper(n *parse.ListNode) *parse.ListNode {
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

func (c *templateTransformContext) applyTransformationsAndSetReturnWrapper(tree *parse.Tree) error {
	_, err := c.applyTransformations(tree.Root)
	if err != nil {
		return err
	}
	if c.returnNode != nil {
		// This is a partial with a return statement.
		c.t.ParseInfo.HasReturn = true
		tree.Root = c.wrapInPartialReturnWrapper(tree.Root)
	}
	return nil
}

// applyTransformations does 2 things:
// 1) Parses partial return statement.
// 2) Tracks template (partial) dependencies and some other info.
func (c *templateTransformContext) applyTransformations(n parse.Node) (bool, error) {
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
		c.handleWith(x)
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
				x.Cmds = slices.Delete(x.Cmds, i, i+1)
			}
		}

	case *parse.CommandNode:
		if x == nil {
			return true, nil
		}
		c.collectInnerInShortcode(x)
		c.collectInnerInPartial(x)
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

func (c *templateTransformContext) isWithPartial(args []parse.Node) bool {
	if len(args) == 0 {
		return false
	}

	first := args[0]

	if pn, ok := first.(*parse.PipeNode); ok {
		if len(pn.Cmds) == 0 || pn.Cmds[0] == nil {
			return false
		}
		return c.isWithPartial(pn.Cmds[0].Args)
	}

	if id1, ok := first.(*parse.IdentifierNode); ok && (id1.Ident == "partial" || id1.Ident == "partialCached") {
		return true
	}

	if chain, ok := first.(*parse.ChainNode); ok {
		if id2, ok := chain.Node.(*parse.IdentifierNode); !ok || (id2.Ident != "partials") {
			return false
		}
		if len(chain.Field) != 1 {
			return false
		}
		if chain.Field[0] != "Include" && chain.Field[0] != "IncludeCached" {
			return false
		}
		return true
	}
	return false
}

func (c *templateTransformContext) isWithDefer(idArg parse.Node) bool {
	id, ok := idArg.(*parse.ChainNode)
	if !ok || len(id.Field) != 1 || id.Field[0] != "Defer" {
		return false
	}
	if id2, ok := id.Node.(*parse.IdentifierNode); !ok || id2.Ident != "templates" {
		return false
	}
	return true
}

// PartialDecoratorPrefix is the prefix used for internal partial decorator templates.
const PartialDecoratorPrefix = "_internal/decorator_"

var templatesInnerRe = regexp.MustCompile(`{{\s*(templates\.Inner\b|inner\b)`)

// hasBreakOrContinueOutsideRange returns true if the given list node contains a break or continue statement without being nested in a range.
func (c *templateTransformContext) hasBreakOrContinueOutsideRange(n *parse.ListNode) bool {
	if n == nil {
		return false
	}
	for _, node := range n.Nodes {
		switch x := node.(type) {
		case *parse.ListNode:
			if c.hasBreakOrContinueOutsideRange(x) {
				return true
			}
		case *parse.RangeNode:
			// skip
		case *parse.IfNode:
			if c.hasBreakOrContinueOutsideRange(x.List) {
				return true
			}
			if c.hasBreakOrContinueOutsideRange(x.ElseList) {
				return true
			}
		case *parse.WithNode:
			if c.hasBreakOrContinueOutsideRange(x.List) {
				return true
			}
			if c.hasBreakOrContinueOutsideRange(x.ElseList) {
				return true
			}
		case *parse.BreakNode, *parse.ContinueNode:
			return true

		}
	}
	return false
}

func (c *templateTransformContext) handleWithPartial(withNode *parse.WithNode) {
	withNodeInnerString := withNode.List.String()
	if templatesInnerRe.MatchString(withNodeInnerString) {
		c.err = fmt.Errorf("inner cannot be used inside a with block that wraps a partial decorator")
		return
	}

	// See #14333. That is a very odd construct, but we need to guard against it.
	if c.hasBreakOrContinueOutsideRange(withNode.List) {
		return
	}
	innerHash := hashing.XxHashFromStringHexEncoded(c.t.Name() + withNodeInnerString)
	internalPartialName := fmt.Sprintf("_partials/%s%s", PartialDecoratorPrefix, innerHash)

	if c.lookupFn(internalPartialName, c.t) == nil {
		innerCopy := withNode.List.CopyList()
		ti, err := c.store.addTransformedTemplateInsert(internalPartialName, SubCategoryInline)
		if err != nil {
			c.err = fmt.Errorf("failed to create internal partial decorator template %q: %w", internalPartialName, err)
			return
		}
		if ti == nil {
			c.err = fmt.Errorf("failed to find internal partial decorator template %q after insertion", internalPartialName)
			return
		}

		cc := newTemplateTransformContext(ti, c.store, c.lookupFn)

		tree, err := c.store.addTransformedTemplateSetTree(ti, innerCopy)
		if err != nil {
			c.err = fmt.Errorf("failed to add internal partial decorator template %q: %w", internalPartialName, err)
			return
		}

		if err := cc.applyTransformationsAndSetReturnWrapper(tree); err != nil {
			c.err = fmt.Errorf("failed to transform internal partial decorator template %q: %w", internalPartialName, err)
			return
		}

	}

	newInner := popPartialDecorator.CopyList()
	ifNode := newInner.Nodes[0].(*parse.IfNode)

	placeholderPipe := ifNode.Pipe.Cmds[0].Args[0].(*parse.PipeNode)

	// Set PLACEHOLDER1 to the unique ID for this partial decorator.
	sn1 := placeholderPipe.Cmds[0].Args[1].(*parse.PipeNode).Cmds[0].Args[0].(*parse.StringNode)
	sn1.Text = innerHash
	sn1.Quoted = fmt.Sprintf("%q", sn1.Text)

	ifNode.ElseList = withNode.List.CopyList()

	newPipe := pushPartialDecorator.CopyList()
	orNode := newPipe.Nodes[0].(*parse.IfNode)
	setContext := orNode.Pipe.Cmds[0].Args[1]
	// Replace PLACEHOLDER with the unique ID for this partial decorator.
	sn2 := setContext.(*parse.PipeNode).Cmds[0].Args[1].(*parse.PipeNode).Cmds[0].Args[0].(*parse.StringNode)
	sn2.Text = innerHash
	sn2.Quoted = fmt.Sprintf("%q", sn2.Text)
	if pn, ok := withNode.Pipe.Cmds[0].Args[0].(*parse.PipeNode); ok {
		withNode.Pipe.Cmds[0].Args = pn.Cmds[0].Args
	}
	withNode.Pipe.Cmds = append(orNode.Pipe.Cmds, withNode.Pipe.Cmds...)

	withNode.List = newInner
}

func (c *templateTransformContext) handleWith(withNode *parse.WithNode) {
	if len(withNode.Pipe.Cmds) != 1 {
		return
	}

	if c.isWithPartial(withNode.Pipe.Cmds[0].Args) {
		c.handleWithPartial(withNode)
		return
	}

	cmd := withNode.Pipe.Cmds[0]

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

	if !c.isWithDefer(idArg) {
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

func (c *templateTransformContext) applyTransformationsToNodes(nodes ...parse.Node) {
	for _, node := range nodes {
		c.applyTransformations(node)
	}
}

func (c *templateTransformContext) hasIdent(idents []string, ident string) bool {
	return slices.Contains(idents, ident)
}

// collectConfig collects and parses any leading template config variable declaration.
// This will be the first PipeNode in the template, and will be a variable declaration
// on the form:
//
//	{{ $_hugo_config:= `{ "version": 1 }` }}
func (c *templateTransformContext) collectConfig(n *parse.PipeNode) {
	if c.t.category != CategoryShortcode {
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
		m, err := hmaps.ToStringMapE(s.Text)
		if err != nil {
			c.err = fmt.Errorf(errMsg, err)
			return
		}
		if err := mapstructure.WeakDecode(m, &c.t.ParseInfo.Config); err != nil {
			c.err = fmt.Errorf(errMsg, err)
		}
	}
}

// collectInnerInShortcode determines if the given CommandNode represents a
// shortcode call to its .Inner.
func (c *templateTransformContext) collectInnerInShortcode(n *parse.CommandNode) {
	if c.t.category != CategoryShortcode {
		return
	}
	if c.t.ParseInfo.IsInner || len(n.Args) == 0 {
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
			c.t.ParseInfo.IsInner = true
			break
		}
	}
}

func (c *templateTransformContext) collectInnerInPartial(n *parse.CommandNode) {
	if c.t.category != CategoryPartial {
		return
	}

	if c.t.ParseInfo.HasPartialInner || len(n.Args) == 0 {
		return
	}

	switch v := n.Args[0].(type) {
	case *parse.IdentifierNode:
		if v.Ident == "inner" {
			c.t.ParseInfo.HasPartialInner = true
		}
	case *parse.ChainNode:
		if v.Field[0] == "Inner" {
			if id, ok := v.Node.(*parse.IdentifierNode); ok && id.Ident == "templates" {
				c.t.ParseInfo.HasPartialInner = true
			}
		}
	}
}

func (c *templateTransformContext) collectReturnNode(n *parse.CommandNode) bool {
	if c.t.category != CategoryPartial || c.returnNode != nil {
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
