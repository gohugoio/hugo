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
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
)

type Handler interface {
	FileConvert(*source.File, *Site) HandledResult
	PageConvert(*Page, tpl.Template) HandledResult
	Read(*source.File, *Site) HandledResult
	Extensions() []string
}

type Handle struct {
	extensions []string
}

func (h Handle) Extensions() []string {
	return h.extensions
}

type HandledResult struct {
	page *Page
	file *source.File
	err  error
}

// HandledResult is an error
func (h HandledResult) Error() string {
	if h.err != nil {
		if h.page != nil {
			return "Error:" + h.err.Error() + " for " + h.page.File.LogicalName()
		}
		if h.file != nil {
			return "Error:" + h.err.Error() + " for " + h.file.LogicalName()
		}
	}
	return h.err.Error()
}

func (h HandledResult) String() string {
	return h.Error()
}

func (h HandledResult) Page() *Page {
	return h.page
}
