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

package js

import (
	"context"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/resources/resource"
)

// BatcherClient is used to do JS batch operations.
type BatcherClient interface {
	New(id string) (Batcher, error)
	Store() *maps.Cache[string, Batcher]
}

// BatchPackage holds a group of JavaScript resources.
type BatchPackage interface {
	Groups() map[string]resource.Resources
}

// Batcher is used to build JavaScript packages.
type Batcher interface {
	Build(context.Context) (BatchPackage, error)
	Config(ctx context.Context) OptionsSetter
	Group(ctx context.Context, id string) BatcherGroup
}

// BatcherGroup is a group of scripts and instances.
type BatcherGroup interface {
	Instance(sid, iid string) OptionsSetter
	Runner(id string) OptionsSetter
	Script(id string) OptionsSetter
}

// OptionsSetter is used to set options for a batch, script or instance.
type OptionsSetter interface {
	SetOptions(map[string]any) string
}
