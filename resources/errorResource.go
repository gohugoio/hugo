// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"image"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/resources/images/exif"

	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ error = (*errorResource)(nil)
	// Imnage covers all current Resource implementations.
	_ resource.Image = (*errorResource)(nil)
	// The list of user facing and exported interfaces in resource.go
	// Note that if we're missing some interface here, the user will still
	// get an error, but not as pretty.
	_ resource.ContentResource         = (*errorResource)(nil)
	_ resource.ReadSeekCloserResource  = (*errorResource)(nil)
	_ resource.ResourcesLanguageMerger = (*resource.Resources)(nil)
	// Make sure it also fails when passed to a pipe function.
	_ ResourceTransformer = (*errorResource)(nil)
)

// NewErrorResource wraps err in a Resource where all but the Err method will panic.
func NewErrorResource(err error) resource.Resource {
	return &errorResource{error: err}
}

type errorResource struct {
	error
}

func (e *errorResource) Err() error {
	return e.error
}

func (e *errorResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	panic(e.error)
}

func (e *errorResource) Content() (interface{}, error) {
	panic(e.error)
}

func (e *errorResource) ResourceType() string {
	panic(e.error)
}

func (e *errorResource) MediaType() media.Type {
	panic(e.error)
}

func (e *errorResource) Permalink() string {
	panic(e.error)
}

func (e *errorResource) RelPermalink() string {
	panic(e.error)
}

func (e *errorResource) Name() string {
	panic(e.error)
}

func (e *errorResource) Title() string {
	panic(e.error)
}

func (e *errorResource) Params() maps.Params {
	panic(e.error)
}

func (e *errorResource) Data() interface{} {
	panic(e.error)
}

func (e *errorResource) Height() int {
	panic(e.error)
}

func (e *errorResource) Width() int {
	panic(e.error)
}

func (e *errorResource) Crop(spec string) (resource.Image, error) {
	panic(e.error)
}

func (e *errorResource) Fill(spec string) (resource.Image, error) {
	panic(e.error)
}

func (e *errorResource) Fit(spec string) (resource.Image, error) {
	panic(e.error)
}

func (e *errorResource) Resize(spec string) (resource.Image, error) {
	panic(e.error)
}

func (e *errorResource) Filter(filters ...interface{}) (resource.Image, error) {
	panic(e.error)
}

func (e *errorResource) Exif() *exif.Exif {
	panic(e.error)
}

func (e *errorResource) DecodeImage() (image.Image, error) {
	panic(e.error)
}

func (e *errorResource) Transform(...ResourceTransformation) (ResourceTransformer, error) {
	panic(e.error)
}
