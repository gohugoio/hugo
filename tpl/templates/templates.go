// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package templates provides template functions for working with templates.
package templates

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"
)

// New returns a new instance of the templates-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	ns := &Namespace{
		deps: deps,
	}

	return ns
}

// Namespace provides template functions for the "templates" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Exists returns whether the template with the given name exists.
// Note that this is the Unix-styled relative path including filename suffix,
// e.g. partials/header.html
func (ns *Namespace) Exists(name string) bool {
	return ns.deps.Tmpl().HasTemplate(name)
}

// Defer defers the execution of a template block.
func (ns *Namespace) Defer(args ...any) (bool, error) {
	// Prevent defer from being used in content adapters,
	// that just doesn't work.
	ns.deps.Site.CheckReady()

	if len(args) != 0 {
		return false, fmt.Errorf("Defer does not take any arguments")
	}
	return true, nil
}

var defferedIDCounter atomic.Uint64

type DeferOpts struct {
	// Optional cache key. If set, the deferred block will be executed
	// once per unique key.
	Key string

	// Optional data context to use when executing the deferred block.
	Data any
}

// DoDefer defers the execution of a template block.
// For internal use only.
func (ns *Namespace) DoDefer(ctx context.Context, id string, optsv any) string {
	var opts DeferOpts
	if optsv != nil {
		if err := mapstructure.WeakDecode(optsv, &opts); err != nil {
			panic(err)
		}
	}

	templateName := id
	var key string
	if opts.Key != "" {
		key = hashing.MD5FromStringHexEncoded(opts.Key)
	} else {
		key = strconv.FormatUint(defferedIDCounter.Add(1), 10)
	}

	id = fmt.Sprintf("%s_%s%s", id, key, tpl.HugoDeferredTemplateSuffix)

	_, _ = ns.deps.BuildState.DeferredExecutions.Executions.GetOrCreate(id,
		func() (*tpl.DeferredExecution, error) {
			return &tpl.DeferredExecution{
				TemplateName: templateName,
				Ctx:          ctx,
				Data:         opts.Data,
				Executed:     false,
			}, nil
		})

	return id
}
