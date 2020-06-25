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

package identity

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/hugofs/files"
)

// NewPathIdentity creates a new Identity with the three identifiers
// type, path and lang (optional).
func NewPathIdentity(typ, path, filename, lang string) PathIdentity {
	path = strings.ToLower(strings.TrimPrefix(filepath.ToSlash(path), "/"))
	return pathIdentity{typePath: typePath{typ: typ, path: path}, filename: filename, lang: lang}
}

type PathIdentitySet map[PathIdentity]bool

func (p PathIdentitySet) ToPathIdentities() PathIdentities {
	var ids PathIdentities
	for id, _ := range p {
		ids = append(ids, id)
	}
	return ids
}

func (p PathIdentitySet) ToIdentityProviders() []Provider {
	var ids []Provider
	for id, _ := range p {
		ids = append(ids, id)
	}
	return ids
}

type PathIdentities []PathIdentity

func (pp PathIdentities) ByType(typ string) PathIdentities {
	var res PathIdentities
	for _, p := range pp {
		if p.Type() == typ {
			res = append(res, p)
		}
	}

	return res

}

func (pp PathIdentities) Sort() PathIdentities {
	sort.Slice(pp, func(i, j int) bool {
		pi, pj := pp[i], pp[j]
		if pi.Path() != pj.Path() {
			return pi.Path() < pj.Path()
		}

		if pi.Filename() != pj.Filename() {
			return pi.Filename() < pj.Filename()
		}

		if pi.Lang() != pj.Lang() {
			return pi.Lang() < pj.Lang()
		}

		return pi.Type() < pj.Type()

	})
	return pp
}

// A PathIdentity is a common identity identified by a type and a path,
// e.g. "layouts" and "_default/single.html".
type PathIdentity interface {
	Provider
	Identity
	Type() string
	Path() string
	Filename() string
	Lang() string
}

var _ IsNotDependentProvider = (*pathIdentity)(nil)

type typePath struct {
	typ  string
	path string
}
type pathIdentity struct {
	typePath
	filename string
	lang     string
}

// GetIdentity returns itself.
func (id pathIdentity) GetIdentity() Identity {
	return id
}

func (id pathIdentity) Base() interface{} {
	return id.typePath
}

func isCrossComponent(c string) bool {
	return c == files.ComponentFolderData || c == files.ComponentFolderLayouts
}

func (id pathIdentity) IsNotDependent(other Provider) bool {
	pid, ok := other.GetIdentity().(PathIdentity)
	if !ok {
		return true
	}

	if isCrossComponent(id.Type()) {
		return id.Type() == pid.Type()
	}

	if isCrossComponent(pid.Type()) {
		// Changes in /data and /layouts currently triggers full
		// content re-renders.
		return id.Type() != files.ComponentFolderContent
	}

	return true
}

func (id typePath) Type() string {
	return id.typ
}

func (id typePath) Path() string {
	return id.path
}

func (id pathIdentity) Filename() string {
	return id.filename
}

func (id pathIdentity) Lang() string {
	return id.lang
}

// Name returns the Path.
func (id pathIdentity) Name() string {
	return id.path
}

func (id pathIdentity) Eq(other interface{}) bool {
	return id == other
}

func (id pathIdentity) ProbablyEq(other interface{}) bool {
	if id == other {
		return true
	}

	oid, ok := other.(PathIdentity)
	if !ok {
		return false
	}

	return id.Base() == oid.Base()

}
