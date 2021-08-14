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

package commands

import (
	// For time zone lookups on Windows without Go installed.
	// See #8892
	_ "time/tzdata"

	"github.com/spf13/cobra"
)

func init() {
	// This message to show to Windows users if Hugo is opened from explorer.exe
	cobra.MousetrapHelpText = `

  Hugo is a command-line tool for generating static website.

  You need to open cmd.exe and run Hugo from there.
  
  Visit https://gohugo.io/ for more information.`
}
