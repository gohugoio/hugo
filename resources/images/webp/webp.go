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

// +build extended

package webp

import (
	"image"
	"io"

	"github.com/bep/gowebp/libwebp"
	"github.com/bep/gowebp/libwebp/webpoptions"
)

// Encode writes the Image m to w in Webp format with the given
// options.
func Encode(w io.Writer, m image.Image, o webpoptions.EncodingOptions) error {
	return libwebp.Encode(w, m, o)
}

// Supports returns whether webp encoding is supported in this build.
func Supports() bool {
	return true
}
