package hugolib

import (
	"reflect"
	"testing"
)

var (
	defaultAuthorsMap = map[string]interface{}{
		"derek": map[string]interface{}{
			"givenName":   "Derek",
			"familyName":  "Perkins",
			"displayName": "Derek Perkins",
			"thumbnail":   "https://www.google.com/images/branding/googlelogo/googlelogo_color_272x92dp.png",
			"image":       "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
			"shortBio":    "Derek loves Hugo",
			"bio":         "Derek really loves Hugo",
			"email":       "derek@email.com",
			"weight":      2,
			"social": map[string]interface{}{
				"facebook":   "derekperkins",
				"github":     "derekperkins",
				"twitter":    "@derekperkins",
				"googleplus": "+derekperkins",
				"pinterest":  "derekperkins",
				"instagram":  "derekperkins",
				"youtube":    "derekperkins",
				"linkedin":   "derekperkins",
			},
			"params": map[string]interface{}{
				"myParamInt": 1234,
				"myParamStr": "5678",
			},
		},
		// test firstName and lastName as aliases of givenName and familyName
		// leave displayName blank to test default concatenation
		// test stripping down the built in social network supported URLs to usernames
		"tanner": map[string]interface{}{
			"firstName": "Tanner",
			"lastName":  "Linsley",
			"thumbnail": "https://www.google.com/images/branding/googlelogo/googlelogo_color_272x92dp.png",
			"image":     "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
			"shortBio":  "Tanner loves Hugo",
			"bio":       "Tanner really loves Hugo",
			"email":     "tanner@email.com",
			"weight":    1,
			"social": map[string]interface{}{
				"facebook":       "https://www.facebook.com/tannerlinsley/",
				"github":         "https://github.com/tannerlinsley/",
				"twitter":        "https://twitter.com/tannerlinsley/",
				"googleplus":     "https://plus.google.com/12345678/",
				"pinterest":      "https://www.pinterest.com/tannerlinsley/",
				"instagram":      "https://www.instagram.com/tannerlinsley/",
				"youtube":        "https://www.youtube.com/user/tannerlinsley/",
				"linkedin":       "https://www.linkedin.com/in/tannerlinsley/",
				"unknownNetwork": "https://www.unknownNetwork.com/tannerlinsley/",
			},
			"params": map[string]interface{}{
				"myParamInt": 1234,
				"myParamStr": "5678",
			},
		},
	}

	authorDerek = Author{
		ID:          "derek",
		GivenName:   "Derek",
		FirstName:   "Derek",
		FamilyName:  "Perkins",
		LastName:    "Perkins",
		DisplayName: "Derek Perkins",
		Thumbnail:   "https://www.google.com/images/branding/googlelogo/googlelogo_color_272x92dp.png",
		Image:       "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		ShortBio:    "Derek loves Hugo",
		Bio:         "Derek really loves Hugo",
		Email:       "derek@email.com",
		Weight:      2,
		Social: AuthorSocial{
			"facebook":   "derekperkins",
			"github":     "derekperkins",
			"twitter":    "derekperkins",
			"googleplus": "derekperkins",
			"pinterest":  "derekperkins",
			"instagram":  "derekperkins",
			"youtube":    "derekperkins",
			"linkedin":   "derekperkins",
		},
		Params: map[string]string{
			"myParamInt": "1234",
			"myParamStr": "5678",
		},
	}

	authorTanner = Author{
		ID:          "tanner",
		GivenName:   "Tanner",
		FirstName:   "Tanner",
		FamilyName:  "Linsley",
		LastName:    "Linsley",
		DisplayName: "Tanner Linsley",
		Thumbnail:   "https://www.google.com/images/branding/googlelogo/googlelogo_color_272x92dp.png",
		Image:       "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		ShortBio:    "Tanner loves Hugo",
		Bio:         "Tanner really loves Hugo",
		Email:       "tanner@email.com",
		Weight:      1,
		Social: AuthorSocial{
			"facebook":       "tannerlinsley",
			"github":         "tannerlinsley",
			"twitter":        "tannerlinsley",
			"googleplus":     "12345678",
			"pinterest":      "tannerlinsley",
			"instagram":      "tannerlinsley",
			"youtube":        "tannerlinsley",
			"linkedin":       "tannerlinsley",
			"unknownNetwork": "https://www.unknownNetwork.com/tannerlinsley/",
		},
		Params: map[string]string{
			"myParamInt": "1234",
			"myParamStr": "5678",
		},
	}

	// ordered by Weight
	defaultAuthors = Authors{authorTanner, authorDerek}

	defaultAuthorSiteInfo = &SiteInfo{
		Authors: defaultAuthors,
	}
)

func TestMapToAuthors(t *testing.T) {
	tests := []struct {
		desc      string
		authorMap map[string]interface{}
		expected  Authors
	}{
		{
			"valid authors",
			defaultAuthorsMap,
			defaultAuthors,
		},
		{
			"no author info",
			map[string]interface{}{},
			Authors{},
		},
		{
			"invalid author",
			map[string]interface{}{
				"derek": "",
			},
			Authors{},
		},
		{
			"blank author id",
			map[string]interface{}{
				"": map[string]interface{}{},
			},
			Authors{},
		},
	}

	for _, test := range tests {
		authors := mapToAuthors(test.authorMap)
		if !reflect.DeepEqual(authors, test.expected) {
			t.Errorf("authors: expected:\n%#v\ngot:\n%#v", test.expected, authors)
		}
	}
}

func TestAuthorsGet(t *testing.T) {
	tests := []struct {
		desc     string
		authors  Authors
		id       string
		expected Author
	}{
		{
			"valid ID",
			defaultAuthors,
			"derek",
			authorDerek,
		},
		{
			"valid ID",
			defaultAuthors,
			"tanner",
			authorTanner,
		},
		{
			"invalid ID",
			defaultAuthors,
			"abc",
			Author{},
		},
		{
			"blank ID",
			defaultAuthors,
			"",
			Author{},
		},
	}

	for _, test := range tests {
		author := test.authors.Get(test.id)
		if !reflect.DeepEqual(author, test.expected) {
			t.Errorf("author: expected:\n%#v\ngot:\n%#v", test.expected, author)
		}
	}
}

func TestPageAuthor(t *testing.T) {
	tests := []struct {
		desc     string
		page     *Page
		expected Author
	}{
		{
			"valid ID",
			&Page{
				Params: map[string]interface{}{
					"author": "derek",
				},
				Site: defaultAuthorSiteInfo,
			},
			authorDerek,
		},
		{
			"valid ID",
			&Page{
				Params: map[string]interface{}{
					"author": "tanner",
				},
				Site: defaultAuthorSiteInfo,
			},
			authorTanner,
		},
		{
			"invalid ID",
			&Page{
				Params: map[string]interface{}{
					"author": "abc",
				},
				Site: defaultAuthorSiteInfo,
			},
			Author{},
		},
		{
			"no ID",
			&Page{
				Params: map[string]interface{}{},
				Site: defaultAuthorSiteInfo,
			},
			Author{},
		},
		{
			"blank ID",
			&Page{
				Params: map[string]interface{}{
					"author": "",
				},
				Site: defaultAuthorSiteInfo,
			},
			Author{},
		},
		{
			"valid authors instead of author",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"derek"},
				},
				Site: defaultAuthorSiteInfo,
			},
			authorDerek,
		},
		{
			"invalid authors instead of author",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"abc"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Author{},
		},
		{
			"blank authors instead of author",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{""},
				},
				Site: defaultAuthorSiteInfo,
			},
			Author{},
		},
		{
			"no site authors",
			&Page{
				Params: map[string]interface{}{
					"author": "derek",
				},
				Site: &SiteInfo{Authors: Authors{}},
			},
			Author{},
		},
	}

	for _, test := range tests {
		author := test.page.Author()
		if !reflect.DeepEqual(author, test.expected) {
			t.Errorf("author: expected:\n%#v\ngot:\n%#v", test.expected, author)
		}
	}
}

func TestPageAuthors(t *testing.T) {
	tests := []struct {
		desc     string
		page     *Page
		expected Authors
	}{
		{
			"single author in array",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"derek"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{authorDerek},
		},
		{
			"verify author ordering",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"tanner", "derek"},
				},
				Site: defaultAuthorSiteInfo,
			},
			defaultAuthors,
		},
		{
			"verify author ordering",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"derek", "tanner"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{authorDerek, authorTanner},
		},
		{
			"invalid ID authors param",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"abc"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"blank string ID authors param",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{""},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"empty authors param",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"no authors param",
			&Page{
				Params: map[string]interface{}{},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"author instead of authors param with valid ID",
			&Page{
				Params: map[string]interface{}{
					"author": "derek",
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{authorDerek},
		},
		{
			"author instead of authors param with invalid ID",
			&Page{
				Params: map[string]interface{}{
					"author": "abc",
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"author instead of authors param with blank ID",
			&Page{
				Params: map[string]interface{}{
					"author": "",
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"author AND authors param with valid IDs",
			&Page{
				Params: map[string]interface{}{
					"author":  "tanner",
					"authors": []string{"derek"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{authorTanner},
		},
		{
			"author AND authors param with invalid author ID",
			&Page{
				Params: map[string]interface{}{
					"author":  "abc",
					"authors": []string{"derek"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{authorDerek},
		},
		{
			"author AND authors param with invalid author IDs",
			&Page{
				Params: map[string]interface{}{
					"author":  "abc",
					"authors": []string{"abc"},
				},
				Site: defaultAuthorSiteInfo,
			},
			Authors{},
		},
		{
			"no site authors",
			&Page{
				Params: map[string]interface{}{
					"authors": []string{"derek"},
				},
				Site: &SiteInfo{Authors: Authors{}},
			},
			Authors{},
		},
	}

	for _, test := range tests {
		authors := test.page.Authors()
		if !reflect.DeepEqual(authors, test.expected) {
			t.Errorf("author: expected:\n%#v\ngot:\n%#v", test.expected, authors)
		}
	}
}

func TestAuthorSocialURL(t *testing.T) {
	tests := []struct {
		authorSocial AuthorSocial
		network      string
		expected     string
	}{
		{
			authorTanner.Social,
			"facebook",
			"https://www.facebook.com/tannerlinsley",
		},
		{
			authorTanner.Social,
			"github",
			"https://github.com/tannerlinsley",
		},
		{
			authorTanner.Social,
			"twitter",
			"https://twitter.com/tannerlinsley",
		},
		{
			authorTanner.Social,
			"googleplus",
			"https://plus.google.com/12345678",
		},
		{
			authorDerek.Social,
			"googleplus",
			"https://plus.google.com/+derekperkins",
		},
		{
			authorTanner.Social,
			"pinterest",
			"https://www.pinterest.com/tannerlinsley/",
		},
		{
			authorTanner.Social,
			"instagram",
			"https://www.instagram.com/tannerlinsley/",
		},
		{
			authorTanner.Social,
			"youtube",
			"https://www.youtube.com/user/tannerlinsley",
		},
		{
			authorTanner.Social,
			"linkedin",
			"https://www.linkedin.com/in/tannerlinsley",
		},
		{
			authorTanner.Social,
			"unknownNetwork",
			"https://www.unknownNetwork.com/tannerlinsley/",
		},
	}

	for _, test := range tests {
		u := test.authorSocial.URL(test.network)
		if test.expected != u {
			t.Errorf("url: expected:\n%#v\ngot:\n%#v", test.expected, u)
		}
	}
}
