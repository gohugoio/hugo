package tags_test

import (
	"io"
	"strings"
)

func equalTo(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}

func str2reader(a []string) []io.Reader {
	var arr []io.Reader
	for _, s := range a {
		arr = append(arr, strings.NewReader(s))
	}
	return arr
}
