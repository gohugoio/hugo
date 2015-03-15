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

package hugolib

// AuthorList is a list of all authors and their metadata.
type AuthorList map[string]Author

// Author contains details about the author of a page.
type Author struct {
	GivenName   string
	FamilyName  string
	DisplayName string
	Thumbnail   string
	Image       string
	ShortBio    string
	LongBio     string
	Email       string
	Social      AuthorSocial
}

// AuthorSocial is a place to put social details per author. These are the
// standard keys that themes will expect to have available, but can be
// expanded to any others on a per site basis
// - website
// - github
// - facebook
// - twitter
// - googleplus
// - pinterest
// - instagram
// - youtube
// - linkedin
// - skype
type AuthorSocial map[string]string
