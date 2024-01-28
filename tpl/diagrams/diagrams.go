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

package diagrams

import (
	"html/template"
)

type SVGDiagram interface {
	// Wrapped returns the diagram as an SVG, including the <svg> container.
	Wrapped() template.HTML

	// Inner returns the inner markup of the SVG.
	// This allows for the <svg> container to be created manually.
	Inner() template.HTML

	// Width returns the width of the SVG.
	Width() int

	// Height returns the height of the SVG.
	Height() int
}
