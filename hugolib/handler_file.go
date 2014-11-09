// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"github.com/dchest/cssmin"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
)

func init() {
	RegisterHandler(css)
}

var css = Handle{
	extensions: []string{"css"},
	read: func(f *source.File, s *Site, results HandleResults) {
		results <- HandledResult{file: f}
	},
	fileConvert: func(f *source.File, s *Site, results HandleResults) {
		x := cssmin.Minify(f.Bytes())
		s.WriteDestFile(f.Path(), helpers.BytesToReader(x))
	},
}
