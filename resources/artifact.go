// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"
)

// Artifact represents an additional output file published as part of a
// resource transformation, e.g. a source map or a file emitted by ESBuild's
// file loader. It is exposed to templates via Resource.Data.Artifacts.
type Artifact interface {
	resource.MediaTypeProvider
	resource.ResourceLinksProvider
}

// NewArtifact creates a new Artifact with the given permalinks and media type.
func NewArtifact(permalink, relPermalink string, mediaType media.Type) Artifact {
	return &artifact{permalink: permalink, relPermalink: relPermalink, mediaType: mediaType}
}

type artifact struct {
	permalink    string
	relPermalink string
	mediaType    media.Type
}

func (a *artifact) MediaType() media.Type {
	return a.mediaType
}

func (a *artifact) Permalink() string {
	return a.permalink
}

func (a *artifact) RelPermalink() string {
	return a.relPermalink
}
