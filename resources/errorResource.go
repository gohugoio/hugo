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
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/images/exif"
	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ error = (*errorResource)(nil)
	// Imnage covers all current Resource implementations.
	_ images.ImageResource = (*errorResource)(nil)
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
func NewErrorResource(err resource.ResourceError) resource.Resource {
	return &errorResource{ResourceError: err}
}

type errorResource struct {
	resource.ResourceError
}

func (e *errorResource) Err() resource.ResourceError {
	return e.ResourceError
}

func (e *errorResource) ReadSeekCloser() (hugio.ReadSeekCloser, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Content() (any, error) {
	panic(e.ResourceError)
}

func (e *errorResource) ResourceType() string {
	panic(e.ResourceError)
}

func (e *errorResource) MediaType() media.Type {
	panic(e.ResourceError)
}

func (e *errorResource) Permalink() string {
	panic(e.ResourceError)
}

func (e *errorResource) RelPermalink() string {
	panic(e.ResourceError)
}

func (e *errorResource) Name() string {
	panic(e.ResourceError)
}

func (e *errorResource) Title() string {
	panic(e.ResourceError)
}

func (e *errorResource) Params() maps.Params {
	panic(e.ResourceError)
}

func (e *errorResource) Data() any {
	panic(e.ResourceError)
}

func (e *errorResource) Height() int {
	panic(e.ResourceError)
}

func (e *errorResource) Width() int {
	panic(e.ResourceError)
}

func (e *errorResource) Crop(spec string) (images.ImageResource, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Fill(spec string) (images.ImageResource, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Fit(spec string) (images.ImageResource, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Resize(spec string) (images.ImageResource, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Filter(filters ...any) (images.ImageResource, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Exif() *exif.ExifInfo {
	panic(e.ResourceError)
}

func (e *errorResource) Colors() ([]string, error) {
	panic(e.ResourceError)
}

func (e *errorResource) DecodeImage() (image.Image, error) {
	panic(e.ResourceError)
}

func (e *errorResource) Transform(...ResourceTransformation) (ResourceTransformer, error) {
	panic(e.ResourceError)
}
