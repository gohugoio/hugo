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

package transform

var ar = newAbsURLReplacer()

// AbsURL replaces relative URLs with absolute ones
// in HTML files, using the baseURL setting.
var AbsURL = func(ct contentTransformer) {
	ar.replaceInHTML(ct)
}

// AbsURLInXML replaces relative URLs with absolute ones
// in XML files, using the baseURL setting.
var AbsURLInXML = func(ct contentTransformer) {
	ar.replaceInXML(ct)
}
