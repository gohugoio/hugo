// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"html/template"

	"github.com/spf13/hugo/helpers"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

var hugoInfo *HugoInfo

// HugoInfo contains information about the current Hugo environment
type HugoInfo struct {
	Version    string
	Generator  template.HTML
	CommitHash string
	BuildDate  string
}

func init() {
	hugoInfo = &HugoInfo{
		Version:    helpers.HugoVersion(),
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		Generator:  template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, helpers.HugoVersion())),
	}
}
