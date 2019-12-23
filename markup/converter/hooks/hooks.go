// Copyright 2019 The Hugo Authors. All rights reserved.
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

package hooks

import (
	"io"

	"github.com/gohugoio/hugo/identity"
)

type LinkContext interface {
	Page() interface{}
	Destination() string
	Title() string
	Text() string
	PlainText() string
}

type Render struct {
	LinkRenderer  LinkRenderer
	ImageRenderer LinkRenderer
}

func (r *Render) Eq(other interface{}) bool {
	ro, ok := other.(*Render)
	if !ok {
		return false
	}
	if r == nil || ro == nil {
		return r == nil
	}

	if r.ImageRenderer.GetIdentity() != ro.ImageRenderer.GetIdentity() {
		return false
	}

	if r.LinkRenderer.GetIdentity() != ro.LinkRenderer.GetIdentity() {
		return false
	}

	return true
}

type LinkRenderer interface {
	Render(w io.Writer, ctx LinkContext) error
	identity.Provider
}
