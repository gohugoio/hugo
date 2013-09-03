// Copyright © 2013 Steve Francia <spf@spf13.com>.
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

package template

import (
	"regexp"
	"strings"
)

var sanitizeRegexp = regexp.MustCompile("[^a-zA-Z0-9./_-]")

func Urlize(url string) string {
	return Sanitize(strings.ToLower(strings.Replace(strings.TrimSpace(url), " ", "-", -1)))
}

func Sanitize(s string) string {
	return sanitizeRegexp.ReplaceAllString(s, "")
}
