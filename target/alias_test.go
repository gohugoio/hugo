package target

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestHTMLRedirectAlias(t *testing.T) {
	var o Translator
	o = new(HTMLRedirectAlias)

	errIsNilForThisOS := runtime.GOOS != "windows"

	tests := []struct {
		value    string
		expected string
		errIsNil bool
	}{
		{"", "", false},
		{"s", filepath.FromSlash("s/index.html"), true},
		{"/", "", false},
		{"alias 1", filepath.FromSlash("alias 1/index.html"), true},
		{"alias 2/", filepath.FromSlash("alias 2/index.html"), true},
		{"alias 3.html", "alias 3.html", true},
		{"alias4.html", "alias4.html", true},
		{"/alias 5.html", "alias 5.html", true},
		{"/трям.html", "трям.html", true},
		{"../../../../tmp/passwd", "", false},
		{"/foo/../../../../tmp/passwd", filepath.FromSlash("tmp/passwd/index.html"), true},
		{"foo/../../../../tmp/passwd", "", false},
		{"C:\\Windows", filepath.FromSlash("C:\\Windows/index.html"), errIsNilForThisOS},
		{"/trailing-space /", filepath.FromSlash("trailing-space /index.html"), errIsNilForThisOS},
		{"/trailing-period./", filepath.FromSlash("trailing-period./index.html"), errIsNilForThisOS},
		{"/tab\tseparated/", filepath.FromSlash("tab\tseparated/index.html"), errIsNilForThisOS},
		{"/chrome/?p=help&ctx=keyboard#topic=3227046", filepath.FromSlash("chrome/?p=help&ctx=keyboard#topic=3227046/index.html"), errIsNilForThisOS},
		{"/LPT1/Printer/", filepath.FromSlash("LPT1/Printer/index.html"), errIsNilForThisOS},
	}

	for _, test := range tests {
		path, err := o.Translate(test.value)
		if (err == nil) != test.errIsNil {
			t.Errorf("Expected err == nil => %t, got: %t. err: %s", test.errIsNil, err == nil, err)
			continue
		}
		if err == nil && path != test.expected {
			t.Errorf("Expected: \"%s\", got: \"%s\"", test.expected, path)
		}
	}
}
