// Copyright 2018 The Hugo Authors. All rights reserved.
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

package urlreplacers

import "github.com/gohugoio/hugo/transform"

var ar = newAbsURLReplacer()

// NewAbsURLTransformer replaces relative URLs with absolute ones
// in HTML files, using the baseURL setting.
func NewAbsURLTransformer(path string) transform.Transformer {
	return func(ft transform.FromTo) error {
		ar.replaceInHTML(path, ft)
		return nil
	}
}

// NewAbsURLInXMLTransformer replaces relative URLs with absolute ones
// in XML files, using the baseURL setting.
func NewAbsURLInXMLTransformer(path string) transform.Transformer {
	return func(ft transform.FromTo) error {
		ar.replaceInXML(path, ft)
		return nil
	}
}
