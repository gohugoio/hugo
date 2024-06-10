// Copyright 2020 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package resourcehelpers

import (
	"errors"
	"fmt"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/resources"
)

// We allow string or a map as the first argument in some cases.
func ResolveIfFirstArgIsString(args []any) (resources.ResourceTransformer, string, bool) {
	if len(args) != 2 {
		return nil, "", false
	}

	v1, ok1 := args[0].(string)
	if !ok1 {
		return nil, "", false
	}
	v2, ok2 := args[1].(resources.ResourceTransformer)

	return v2, v1, ok2
}

// This roundabout way of doing it is needed to get both pipeline behavior and options as arguments.
func ResolveArgs(args []any) (resources.ResourceTransformer, map[string]any, error) {
	if len(args) == 0 {
		return nil, nil, errors.New("no Resource provided in transformation")
	}

	if len(args) == 1 {
		r, ok := args[0].(resources.ResourceTransformer)
		if !ok {
			return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
		}
		return r, nil, nil
	}

	r, ok := args[1].(resources.ResourceTransformer)
	if !ok {
		if _, ok := args[1].(map[string]any); !ok {
			return nil, nil, fmt.Errorf("no Resource provided in transformation")
		}
		return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
	}

	m, err := maps.ToStringMapE(args[0])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid options type: %w", err)
	}

	return r, m, nil
}
