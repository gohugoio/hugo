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

package resources

import (
	"github.com/gohugoio/hugo/resources/postpub"
	"github.com/gohugoio/hugo/resources/resource"
)

type transformationKeyer interface {
	TransformationKey() string
}

// PostProcess wraps the given Resource for later processing.
func (spec *Spec) PostProcess(r resource.Resource) (postpub.PostPublishedResource, error) {
	key := r.(transformationKeyer).TransformationKey()
	spec.postProcessMu.RLock()
	result, found := spec.PostProcessResources[key]
	spec.postProcessMu.RUnlock()
	if found {
		return result, nil
	}

	spec.postProcessMu.Lock()
	defer spec.postProcessMu.Unlock()

	// Double check
	result, found = spec.PostProcessResources[key]
	if found {
		return result, nil
	}

	result = postpub.NewPostPublishResource(spec.incr.Incr(), r)
	if result == nil {
		panic("got nil result")
	}
	spec.PostProcessResources[key] = result

	return result, nil
}
