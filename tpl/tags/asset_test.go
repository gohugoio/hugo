package tags_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/tpl/tags"
)

const helloWorld = "hello world"

const helloWorldSum = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
const helloWorldBase64 = "uU0nuZNNPgilLlLX2n2r+sSE7+N6U4DukIj3rOLvzek="

func TestAssetRaw(t *testing.T) {
	t.Parallel()
	// explicitly not using const here
	asset := tags.NewAsset("/js/test/js", []byte("hello world"))

	actual := asset.Raw()

	if string(actual) != helloWorld {
		t.Errorf("want Raw() = <%s>, got <%s>", helloWorld, actual)
	}
}

func TestAssetSum(t *testing.T) {
	t.Parallel()
	asset := tags.NewAsset("/js/test.js", []byte("hello world"))

	actual := asset.Sum()

	if !strings.EqualFold(helloWorldSum, hex.EncodeToString(actual)) {
		t.Errorf("want Sum() = <%s>, got <%s>", helloWorldSum, hex.EncodeToString(asset.Sum()))
	}
}

func TestAssetBase64Sum(t *testing.T) {
	t.Parallel()
	asset := tags.NewAsset("/js/test.js", []byte("hello world"))

	actual := asset.Base64Sum()

	if helloWorldBase64 != actual {
		t.Errorf("want Base64Sum() = <%v>, got <%v>", helloWorldBase64, actual)
	}
}

func TestAssetUrl(t *testing.T) {
	t.Parallel()

	td := []struct {
		msg       string
		immutable bool
		filename  string
		contents  string
		expected  string
	}{
		{"non-immutable URL",
			false, "/js/test.js", "hello world", "/js/test.js"},
		{"immutable URL",
			true, "/js/test.js", "hello world", "/js/test-uU0nuZNNPgilLlLX2n2r-sSE7-N6U4DukIj3rOLvzek.js"},
	}

	for _, v := range td {
		asset := tags.NewAsset(v.filename, []byte(v.contents))

		actual := asset.Url(v.immutable)

		if v.expected != actual {
			t.Errorf("%v | want Url(%v) = <%v>, got <%v>", v.msg, v.immutable, v.expected, actual)
		}
	}
}
