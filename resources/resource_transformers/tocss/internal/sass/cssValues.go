// Copyright 2022 The Hugo Authors. All rights reserved.
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

package sass

type cssValue struct {
	prefix []string
	sufix  []string
}

var (
	cssValues = cssValue{
		prefix: []string{
			"#",
			"rgb(",
			"hsl(",
			"hwb(",
			"lch(",
			"lab(",
			"calc(",
			"min(",
			"max(",
			"minmax(",
			"clamp(",
			"attr(",
		},
		sufix: []string{
			"em",
			"ex",
			"cap",
			"ch",
			"ic",
			"rem",
			"lh",
			"rlh",
			"vw",
			"vh",
			"vi",
			"vb",
			"vmin",
			"vmax",
			"cqw",
			"cqh",
			"cqi",
			"cqb",
			"cqmin",
			"cqmax",
			"cm",
			"mm",
			"Q",
			"in",
			"pc",
			"pt",
			"px",
			"deg",
			"grad",
			"rad",
			"turn",
			"s",
			"ms",
			"fr",
			"dpi",
			"dpcm",
			"dppx",
			"x",
			"%",
		},
	}
)
