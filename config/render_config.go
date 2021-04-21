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

package config

import (
	"fmt"
)

type RenderDest string

const (
	RenderDestDisk      RenderDest = "disk"
	RenderDestComposite            = "composite"
	RenderDestMemory               = "memory"
	RenderDestHybrid               = "hybrid"
	RenderDestUnset                = ""
)

func RenderDestFrom(s string) RenderDest {
	switch s {
	case "disk":
		return RenderDestDisk
	case "composite":
		return RenderDestComposite
	case "memory":
		return RenderDestMemory
	case "hybrid":
		return RenderDestHybrid
	case "":
		return RenderDestUnset
	default:
		panic(fmt.Errorf(`renderTo must be "disk", "composite", or "memory" but "%s" is set`, s))
	}
}
