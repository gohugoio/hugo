// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package provides ways to identify values in Hugo. Used for dependency tracking etc.
package identity

import (
	"fmt"
	"sync/atomic"

	hglob "github.com/gohugoio/hugo/hugofs/glob"
)

// NewGlobIdentity creates a new Identity that
// is probably dependent on any other Identity
// that matches the given pattern.
func NewGlobIdentity(pattern string) Identity {
	glob, err := hglob.GetGlob(pattern)
	if err != nil {
		panic(err)
	}

	predicate := func(other Identity) bool {
		return glob.Match(other.IdentifierBase())
	}

	return NewPredicateIdentity(predicate, nil)
}

var predicateIdentityCounter = &atomic.Uint32{}

type predicateIdentity struct {
	id                 string
	probablyDependent  func(Identity) bool
	probablyDependency func(Identity) bool
}

var (
	_ IsProbablyDependencyProvider = &predicateIdentity{}
	_ IsProbablyDependentProvider  = &predicateIdentity{}
)

// NewPredicateIdentity creates a new Identity that implements both IsProbablyDependencyProvider and IsProbablyDependentProvider
// using the provided functions, both of which are optional.
func NewPredicateIdentity(
	probablyDependent func(Identity) bool,
	probablyDependency func(Identity) bool,
) *predicateIdentity {
	if probablyDependent == nil {
		probablyDependent = func(Identity) bool { return false }
	}
	if probablyDependency == nil {
		probablyDependency = func(Identity) bool { return false }
	}
	return &predicateIdentity{probablyDependent: probablyDependent, probablyDependency: probablyDependency, id: fmt.Sprintf("predicate%d", predicateIdentityCounter.Add(1))}
}

func (id *predicateIdentity) IdentifierBase() string {
	return id.id
}

func (id *predicateIdentity) IsProbablyDependent(other Identity) bool {
	return id.probablyDependent(other)
}

func (id *predicateIdentity) IsProbablyDependency(other Identity) bool {
	return id.probablyDependency(other)
}
