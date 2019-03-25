// Copyright 2019 The Hugo Authors. All rights reserved.
// Some functions in this file (see comments) is based on the Go source code,
// copyright The Go Authors and  governed by a BSD-style license.
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

// Package codegen contains helpers for code generation.
package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Make room for insertions
const weightWidth = 1000

// NewInspector creates a new Inspector given a source root.
func NewInspector(root string) *Inspector {
	return &Inspector{ProjectRootDir: root}
}

// Inspector provides methods to help code generation. It uses a combination
// of reflection and source code AST to do the heavy lifting.
type Inspector struct {
	ProjectRootDir string

	init sync.Once

	// Determines method order. Go's reflect sorts lexicographically, so
	// we must parse the source to preserve this order.
	methodWeight map[string]map[string]int
}

// MethodsFromTypes create a method set from the include slice, excluding any
// method in exclude.
func (c *Inspector) MethodsFromTypes(include []reflect.Type, exclude []reflect.Type) Methods {
	c.parseSource()

	var methods Methods

	var excludes = make(map[string]bool)

	if len(exclude) > 0 {
		for _, m := range c.MethodsFromTypes(exclude, nil) {
			excludes[m.Name] = true
		}
	}

	// There may be overlapping interfaces in types. Do a simple check for now.
	seen := make(map[string]bool)

	nameAndPackage := func(t reflect.Type) (string, string) {
		var name, pkg string

		isPointer := t.Kind() == reflect.Ptr

		if isPointer {
			t = t.Elem()
		}

		pkgPrefix := ""
		if pkgPath := t.PkgPath(); pkgPath != "" {
			pkgPath = strings.TrimSuffix(pkgPath, "/")
			_, shortPath := path.Split(pkgPath)
			pkgPrefix = shortPath + "."
			pkg = pkgPath
		}

		name = t.Name()
		if name == "" {
			// interface{}
			name = t.String()
		}

		if isPointer {
			pkgPrefix = "*" + pkgPrefix
		}

		name = pkgPrefix + name

		return name, pkg

	}

	for _, t := range include {

		for i := 0; i < t.NumMethod(); i++ {

			m := t.Method(i)
			if excludes[m.Name] || seen[m.Name] {
				continue
			}

			seen[m.Name] = true

			if m.PkgPath != "" {
				// Not exported
				continue
			}

			numIn := m.Type.NumIn()

			ownerName, _ := nameAndPackage(t)

			method := Method{Owner: t, OwnerName: ownerName, Name: m.Name}

			for i := 0; i < numIn; i++ {
				in := m.Type.In(i)

				name, pkg := nameAndPackage(in)

				if pkg != "" {
					method.Imports = append(method.Imports, pkg)
				}

				method.In = append(method.In, name)
			}

			numOut := m.Type.NumOut()

			if numOut > 0 {
				for i := 0; i < numOut; i++ {
					out := m.Type.Out(i)
					name, pkg := nameAndPackage(out)

					if pkg != "" {
						method.Imports = append(method.Imports, pkg)
					}

					method.Out = append(method.Out, name)
				}
			}

			methods = append(methods, method)
		}

	}

	sort.SliceStable(methods, func(i, j int) bool {
		mi, mj := methods[i], methods[j]

		wi := c.methodWeight[mi.OwnerName][mi.Name]
		wj := c.methodWeight[mj.OwnerName][mj.Name]

		if wi == wj {
			return mi.Name < mj.Name
		}

		return wi < wj

	})

	return methods

}

func (c *Inspector) parseSource() {
	c.init.Do(func() {

		if !strings.Contains(c.ProjectRootDir, "hugo") {
			panic("dir must be set to the Hugo root")
		}

		c.methodWeight = make(map[string]map[string]int)
		dirExcludes := regexp.MustCompile("docs|examples")
		fileExcludes := regexp.MustCompile("autogen")
		var filenames []string

		filepath.Walk(c.ProjectRootDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if dirExcludes.MatchString(info.Name()) {
					return filepath.SkipDir
				}
			}

			if !strings.HasSuffix(path, ".go") || fileExcludes.MatchString(path) {
				return nil
			}

			filenames = append(filenames, path)

			return nil

		})

		for _, filename := range filenames {

			pkg := c.packageFromPath(filename)

			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
			if err != nil {
				panic(err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				switch t := n.(type) {
				case *ast.TypeSpec:
					if t.Name.IsExported() {
						switch it := t.Type.(type) {
						case *ast.InterfaceType:
							iface := pkg + "." + t.Name.Name
							methodNames := collectMethodsRecursive(pkg, it.Methods.List)
							weights := make(map[string]int)
							weight := weightWidth
							for _, name := range methodNames {
								weights[name] = weight
								weight += weightWidth
							}
							c.methodWeight[iface] = weights
						}
					}

				}
				return true
			})

		}

		// Complement
		for _, v1 := range c.methodWeight {
			for k2, w := range v1 {
				if v, found := c.methodWeight[k2]; found {
					for k3, v3 := range v {
						v1[k3] = (v3 / weightWidth) + w
					}
				}
			}
		}

	})
}

func (c *Inspector) packageFromPath(p string) string {
	p = filepath.ToSlash(p)
	base := path.Base(p)
	if !strings.Contains(base, ".") {
		return base
	}
	return path.Base(strings.TrimSuffix(p, base))
}

// Method holds enough information about it to recreate it.
type Method struct {
	// The interface we extracted this method from.
	Owner reflect.Type

	// String version of the above, on the form PACKAGE.NAME, e.g.
	// page.Page
	OwnerName string

	// Method name.
	Name string

	// Imports needed to satisfy the method signature.
	Imports []string

	// Argument types, including any package prefix, e.g. string, int, interface{},
	// net.Url
	In []string

	// Return types.
	Out []string
}

// Declaration creates a method declaration (without any body) for the given receiver.
func (m Method) Declaration(receiver string) string {
	return fmt.Sprintf("func (%s %s) %s%s %s", receiverShort(receiver), receiver, m.Name, m.inStr(), m.outStr())
}

// DeclarationNamed creates a method declaration (without any body) for the given receiver
// with named return values.
func (m Method) DeclarationNamed(receiver string) string {
	return fmt.Sprintf("func (%s %s) %s%s %s", receiverShort(receiver), receiver, m.Name, m.inStr(), m.outStrNamed())
}

// Delegate creates a delegate call string.
func (m Method) Delegate(receiver, delegate string) string {
	ret := ""
	if len(m.Out) > 0 {
		ret = "return "
	}
	return fmt.Sprintf("%s%s.%s.%s%s", ret, receiverShort(receiver), delegate, m.Name, m.inOutStr())
}

func (m Method) String() string {
	return m.Name + m.inStr() + " " + m.outStr() + "\n"
}

func (m Method) inOutStr() string {
	if len(m.In) == 0 {
		return "()"
	}

	args := make([]string, len(m.In))
	for i := 0; i < len(args); i++ {
		args[i] = fmt.Sprintf("arg%d", i)
	}
	return "(" + strings.Join(args, ", ") + ")"
}

func (m Method) inStr() string {
	if len(m.In) == 0 {
		return "()"
	}

	args := make([]string, len(m.In))
	for i := 0; i < len(args); i++ {
		args[i] = fmt.Sprintf("arg%d %s", i, m.In[i])
	}
	return "(" + strings.Join(args, ", ") + ")"
}

func (m Method) outStr() string {
	if len(m.Out) == 0 {
		return ""
	}
	if len(m.Out) == 1 {
		return m.Out[0]
	}

	return "(" + strings.Join(m.Out, ", ") + ")"
}

func (m Method) outStrNamed() string {
	if len(m.Out) == 0 {
		return ""
	}

	outs := make([]string, len(m.Out))
	for i := 0; i < len(outs); i++ {
		outs[i] = fmt.Sprintf("o%d %s", i, m.Out[i])
	}

	return "(" + strings.Join(outs, ", ") + ")"
}

// Methods represents a list of methods for one or more interfaces.
// The order matches the defined order in their source file(s).
type Methods []Method

// Imports returns a sorted list of package imports needed to satisfy the
// signatures of all methods.
func (m Methods) Imports() []string {
	var pkgImports []string
	for _, method := range m {
		pkgImports = append(pkgImports, method.Imports...)
	}
	if len(pkgImports) > 0 {
		pkgImports = uniqueNonEmptyStrings(pkgImports)
		sort.Strings(pkgImports)
	}
	return pkgImports
}

// ToMarshalJSON creates a MarshalJSON method for these methods. Any method name
// matchin any of the regexps in excludes will be ignored.
func (m Methods) ToMarshalJSON(receiver, pkgPath string, excludes ...string) (string, []string) {
	var sb strings.Builder

	r := receiverShort(receiver)
	what := firstToUpper(trimAsterisk(receiver))
	pgkName := path.Base(pkgPath)

	fmt.Fprintf(&sb, "func Marshal%sToJSON(%s %s) ([]byte, error) {\n", what, r, receiver)

	var methods Methods
	var excludeRes = make([]*regexp.Regexp, len(excludes))

	for i, exclude := range excludes {
		excludeRes[i] = regexp.MustCompile(exclude)
	}

	for _, method := range m {
		// Exclude methods with arguments and incompatible return values
		if len(method.In) > 0 || len(method.Out) == 0 || len(method.Out) > 2 {
			continue
		}

		if len(method.Out) == 2 {
			if method.Out[1] != "error" {
				continue
			}
		}

		for _, re := range excludeRes {
			if re.MatchString(method.Name) {
				continue
			}
		}

		methods = append(methods, method)
	}

	for _, method := range methods {
		varn := varName(method.Name)
		if len(method.Out) == 1 {
			fmt.Fprintf(&sb, "\t%s := %s.%s()\n", varn, r, method.Name)
		} else {
			fmt.Fprintf(&sb, "\t%s, err := %s.%s()\n", varn, r, method.Name)
			fmt.Fprint(&sb, "\tif err != nil {\n\t\treturn nil, err\n\t}\n")
		}
	}

	fmt.Fprint(&sb, "\n\ts := struct {\n")

	for _, method := range methods {
		fmt.Fprintf(&sb, "\t\t%s %s\n", method.Name, typeName(method.Out[0], pgkName))
	}

	fmt.Fprint(&sb, "\n\t}{\n")

	for _, method := range methods {
		varn := varName(method.Name)
		fmt.Fprintf(&sb, "\t\t%s: %s,\n", method.Name, varn)
	}

	fmt.Fprint(&sb, "\n\t}\n\n")
	fmt.Fprint(&sb, "\treturn json.Marshal(&s)\n}")

	pkgImports := append(methods.Imports(), "encoding/json")

	if pkgPath != "" {
		// Exclude self
		for i, pkgImp := range pkgImports {
			if pkgImp == pkgPath {
				pkgImports = append(pkgImports[:i], pkgImports[i+1:]...)
			}
		}
	}

	return sb.String(), pkgImports

}

func collectMethodsRecursive(pkg string, f []*ast.Field) []string {
	var methodNames []string
	for _, m := range f {
		if m.Names != nil {
			methodNames = append(methodNames, m.Names[0].Name)
			continue
		}

		if ident, ok := m.Type.(*ast.Ident); ok && ident.Obj != nil {
			// Embedded interface
			methodNames = append(
				methodNames,
				collectMethodsRecursive(
					pkg,
					ident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List)...)
		} else {
			// Embedded, but in a different file/package. Return the
			// package.Name and deal with that later.
			name := packageName(m.Type)
			if !strings.Contains(name, ".") {
				// Assume current package
				name = pkg + "." + name
			}
			methodNames = append(methodNames, name)
		}
	}

	return methodNames

}

func firstToLower(name string) string {
	return strings.ToLower(name[:1]) + name[1:]
}

func firstToUpper(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

func packageName(e ast.Expr) string {
	switch tp := e.(type) {
	case *ast.Ident:
		return tp.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", packageName(tp.X), packageName(tp.Sel))
	}
	return ""
}

func receiverShort(receiver string) string {
	return strings.ToLower(trimAsterisk(receiver))[:1]
}

func trimAsterisk(name string) string {
	return strings.TrimPrefix(name, "*")
}

func typeName(name, pkg string) string {
	return strings.TrimPrefix(name, pkg+".")
}

func uniqueNonEmptyStrings(s []string) []string {
	var unique []string
	set := map[string]interface{}{}
	for _, val := range s {
		if val == "" {
			continue
		}
		if _, ok := set[val]; !ok {
			unique = append(unique, val)
			set[val] = val
		}
	}
	return unique
}

func varName(name string) string {
	name = firstToLower(name)

	// Adjust some reserved keywords, see https://golang.org/ref/spec#Keywords
	switch name {
	case "type":
		name = "typ"
	case "package":
		name = "pkg"
		// Not reserved, but syntax highlighters has it as a keyword.
	case "len":
		name = "length"
	}

	return name

}
