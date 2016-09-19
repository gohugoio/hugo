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
	"sort"
	"strings"

	"github.com/spf13/cast"
)

// AuthorList is a list of all authors and their metadata.
type AuthorList map[string]Author

// Author contains details about the author of a page.
type Author struct {
	ID          string
	GivenName   string            // givenName OR firstName
	FirstName   string            // alias for GivenName
	FamilyName  string            // familyName OR lastName
	LastName    string            // alias for FamilyName
	DisplayName string            // displayName
	Thumbnail   string            // thumbnail
	Image       string            // image
	ShortBio    string            // shortBio
	Bio         string            // bio
	Email       string            // email
	Social      AuthorSocial      // social
	Params      map[string]string // params
	Weight      int
	languages   map[string]Author
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

// URL is a convenience function that provides the correct canonical URL
// for a specific social network given a username. If an unsupported network
// is requested, only the username is returned
func (as AuthorSocial) URL(key string) string {
	switch key {
	case "github":
		return fmt.Sprintf("https://github.com/%s", as[key])
	case "facebook":
		return fmt.Sprintf("https://www.facebook.com/%s", as[key])
	case "twitter":
		return fmt.Sprintf("https://twitter.com/%s", as[key])
	case "googleplus":
		isNumeric := onlyNumbersRegExp.Match([]byte(as[key]))
		if isNumeric {
			return fmt.Sprintf("https://plus.google.com/%s", as[key])
		}
		return fmt.Sprintf("https://plus.google.com/+%s", as[key])
	case "pinterest":
		return fmt.Sprintf("https://www.pinterest.com/%s/", as[key])
	case "instagram":
		return fmt.Sprintf("https://www.instagram.com/%s/", as[key])
	case "youtube":
		return fmt.Sprintf("https://www.youtube.com/user/%s", as[key])
	case "linkedin":
		return fmt.Sprintf("https://www.linkedin.com/in/%s", as[key])
	default:
		return as[key]
	}
}

func mapToAuthors(m map[string]interface{}) Authors {
	authors := make(Authors, 0, len(m))
	for authorID, data := range m {
		authorMap, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		a := mapToAuthor(authorID, authorMap)
		if a.ID != "" {
			authors = append(authors, a)
		}
	}
	sort.Stable(authors)
	return authors
}

func mapToAuthor(id string, m map[string]interface{}) Author {
	if id == "" {
		return Author{}
	}

	author := Author{ID: id}
	for k, data := range m {
		switch k {
		case "givenName", "firstName":
			author.GivenName = cast.ToString(data)
			author.FirstName = author.GivenName
		case "familyName", "lastName":
			author.FamilyName = cast.ToString(data)
			author.LastName = author.FamilyName
		case "displayName":
			author.DisplayName = cast.ToString(data)
		case "thumbnail":
			author.Thumbnail = cast.ToString(data)
		case "image":
			author.Image = cast.ToString(data)
		case "shortBio":
			author.ShortBio = cast.ToString(data)
		case "bio":
			author.Bio = cast.ToString(data)
		case "email":
			author.Email = cast.ToString(data)
		case "weight":
			author.Weight = cast.ToInt(data)
		case "social":
			author.Social = normalizeSocial(cast.ToStringMapString(data))
		case "params":
			author.Params = cast.ToStringMapString(data)
		case "languages":
			if author.languages == nil {
				author.languages = make(map[string]Author)
			}
			langAuthorMap := cast.ToStringMap(data)
			for lang, m := range langAuthorMap {
				authorMap, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				author.languages[lang] = mapToAuthor(id, authorMap)
			}
		}
	}

	// set a reasonable default for DisplayName
	if author.DisplayName == "" {
		author.DisplayName = author.GivenName + " " + author.FamilyName
	}

	return author
}

// languageOverride returns an author with details overridden by the provided language
func (a Authors) languageOverride(lang string) Authors {
	langAuthors := make(Authors, 0, len(a))
	for _, author := range a {
		langAuthors = append(langAuthors, author.languageOverride(lang))
	}
	return langAuthors
}

// languageOverride returns an author with details overridden by the provided language
func (a Author) languageOverride(lang string) Author {
	if a.languages == nil {
		return a
	}

	langAuthor, ok := a.languages[lang]
	if !ok {
		return a
	}

	if langAuthor.GivenName != "" {
		a.GivenName = langAuthor.GivenName
	}
	if langAuthor.FirstName != "" {
		a.FirstName = langAuthor.FirstName
	}
	if langAuthor.FamilyName != "" {
		a.FamilyName = langAuthor.FamilyName
	}
	if langAuthor.LastName != "" {
		a.LastName = langAuthor.LastName
	}
	if langAuthor.DisplayName != "" {
		a.DisplayName = langAuthor.DisplayName
	}
	if langAuthor.Thumbnail != "" {
		a.Thumbnail = langAuthor.Thumbnail
	}
	if langAuthor.Image != "" {
		a.Image = langAuthor.Image
	}
	if langAuthor.ShortBio != "" {
		a.ShortBio = langAuthor.ShortBio
	}
	if langAuthor.Bio != "" {
		a.Bio = langAuthor.Bio
	}
	if langAuthor.Email != "" {
		a.Email = langAuthor.Email
	}

	for k, v := range langAuthor.Social {
		a.Social[k] = v
	}

	for k, v := range langAuthor.Params {
		a.Params[k] = v
	}

	return a
}

// normalizeSocial makes a naive attempt to normalize social media usernames
// and strips out extraneous characters or url info
func normalizeSocial(m map[string]string) map[string]string {
	for network, username := range m {
		if !isSupportedSocialNetwork(network) {
			continue
		}

		username = strings.TrimSpace(username)
		username = strings.TrimSuffix(username, "/")
		strs := strings.Split(username, "/")
		username = strs[len(strs)-1]
		username = strings.TrimPrefix(username, "@")
		username = strings.TrimPrefix(username, "+")
		m[network] = username
	}
	return m
}

func isSupportedSocialNetwork(network string) bool {
	switch network {
	case
		"github",
		"facebook",
		"twitter",
		"googleplus",
		"pinterest",
		"instagram",
		"youtube",
		"linkedin":
		return true
	}
	return false
}

func (a Authors) Len() int           { return len(a) }
func (a Authors) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Authors) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
