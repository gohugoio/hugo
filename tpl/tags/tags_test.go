package tags_test

import (
	"html/template"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/tags"

	"testing"
)

const bootstrapJs = "./bootstrap.js"
const jqueryJs = "./jquery-1.11.3.min.js"

func TestAsStrings(t *testing.T) {
	// TODO (nfisher - 2018-01-23): test that works with Windows paths.
	t.Parallel()

	td := []struct {
		msg         string
		wd          string
		srcs        []interface{}
		expected    []string
		expectedErr error
	}{
		{"invalid type",
			"/home/gohugo", []interface{}{true},
			nil, tags.ErrNotString},
		{"single src",
			"/home/gohugo", []interface{}{"js/jquery.js"},
			[]string{"js/jquery.js"}, nil},
		{"multiple srcs",
			"/home/gohugo", []interface{}{"js/react.js", "js/redux.js"},
			[]string{"js/react.js", "js/redux.js"}, nil},
	}

	for _, v := range td {
		actual, err := tags.AsStrings(v.srcs)

		if err != v.expectedErr {
			t.Errorf("%v | err = %#v, want %#v", v.msg, err, v.expectedErr)
		}

		if !equalTo(actual, v.expected) {
			t.Errorf("%v | AsStrings = %v, want %v", v.msg, actual, v.expected)
		}
	}
}

var scriptTable = []struct {
	msg       string
	dest      string
	incSri    bool
	immutable bool
	srcs      []interface{}
	expected  string
}{
	/*
		{
			"single src w/out sri",
			"/js/postman.js", false, false, []interface{}{bootstrapJs},
			`<script src="/js/postman.js"></script>`,
		},

		{
			"single src w/ sri",
			"/js/postman.js", true, false, []interface{}{bootstrapJs},
			`<script src="/js/postman.js" integrity="sha256-3vw5dArBhZ2OJ4XtRzIIQJYn6Hrd1fePLeqsuToS1R0="></script>`,
		},

		{
			"single src is immutable",
			"/js/postman.js", false, true, []interface{}{bootstrapJs},
			`<script src="/js/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.js"></script>`,
		},
	*/

	{
		"single src w/ sri, and is immutable",
		"/js/postman.js", true, true, []interface{}{bootstrapJs},
		`<script src="/js/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.js" integrity="sha256-3vw5dArBhZ2OJ4XtRzIIQJYn6Hrd1fePLeqsuToS1R0="></script>`,
	},

	/*
		{
			"multiple src w/ sri",
			"/js/postman.js", true, false, []interface{}{bootstrapJs, jqueryJs},
			`<script src="/js/postman.js" integrity="sha256-969whtXXli9J7gPWBE6D2wK9oFbNULHIPg4e53HSBKQ="></script>`,
		},
	*/

	{
		"multiple src w/ sri, and is immutable",
		"/js/postman.js", true, true, []interface{}{bootstrapJs, jqueryJs},
		`<script src="/js/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.js" integrity="sha256-969whtXXli9J7gPWBE6D2wK9oFbNULHIPg4e53HSBKQ="></script>`,
	},
}

func TestTagsScript(t *testing.T) {
	t.Parallel()
	t.Skip()

	ns := tags.New(&deps.Deps{})

	for _, v := range scriptTable {
		tag, err := ns.Script(v.dest, nil, v.srcs...)
		if err != nil {
			t.Errorf("%s\n | err = %v, want nil", v.msg, err)
		}

		if tag != template.HTML(v.expected) {
			t.Errorf("%s\n | tag = %v, want %v", v.msg, tag, v.expected)
		}
	}
}

// reusing JS files as I'm only interested in the tag output (for now).
// TODO (nfisher 2018-01-25): Replace JS file refs w/ CSS files.
var styleTable = []struct {
	msg       string
	dest      string
	incSri    bool
	immutable bool
	srcs      []interface{}
	expected  string
}{
	{
		"single src w/out sri",
		"/css/postman.css", false, false, []interface{}{bootstrapJs},
		`<link rel="stylesheet" href="/css/postman.css"/>`,
	},

	{
		"single src w/ sri",
		"/css/postman.css", true, false, []interface{}{bootstrapJs},
		`<link rel="stylesheet" href="/css/postman.css" integrity="sha256-3vw5dArBhZ2OJ4XtRzIIQJYn6Hrd1fePLeqsuToS1R0="/>`,
	},

	{
		"single src is immutable",
		"/css/postman.css", false, true, []interface{}{bootstrapJs},
		`<link rel="stylesheet" href="/css/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.css"/>`,
	},

	{
		"single src w/ sri, and is immutable",
		"/css/postman.css", true, true, []interface{}{bootstrapJs},
		`<link rel="stylesheet" href="/css/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.css" integrity="sha256-3vw5dArBhZ2OJ4XtRzIIQJYn6Hrd1fePLeqsuToS1R0="/>`,
	},

	{
		"multiple src w/ sri",
		"/css/postman.css", true, false, []interface{}{bootstrapJs, jqueryJs},
		`<link rel="stylesheet" href="/css/postman.css" integrity="sha256-969whtXXli9J7gPWBE6D2wK9oFbNULHIPg4e53HSBKQ="/>`,
	},

	{
		"multiple src w/ sri, and is immutable",
		"/js/postman.js", true, true, []interface{}{bootstrapJs, jqueryJs},
		`<link rel="stylesheet" src="/css/postman-3PVPh0f01T2hWni0XcU5x0jU8Kqdgb64v1avoKGC4G2O.css" integrity="sha256-969whtXXli9J7gPWBE6D2wK9oFbNULHIPg4e53HSBKQ="/>`,
	},
}

func TestTagsStyle(t *testing.T) {
	t.Parallel()
	t.Skip()

	ns := tags.New(&deps.Deps{})

	for _, v := range styleTable {
		tag, err := ns.Style(v.dest, nil, v.srcs...)
		if err != nil {
			t.Errorf("%s\n | err = %v, want nil", v.msg, err)
		}

		if tag != template.HTML(v.expected) {
			t.Errorf("%s\n | tag = %v, want %v", v.msg, tag, v.expected)
		}
	}
}
