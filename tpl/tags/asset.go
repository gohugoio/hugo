package tags

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"regexp"
)

// capture the .EXT
var extRegexp = regexp.MustCompile(`(\.[^.]+)$`)

// NewAsset initialises a new Asset using the url as the target path and b as
// it's contents.
func NewAsset(url string, b []byte) Asset {
	return Asset{
		url: url,
		raw: b,
	}
}

// Asset is used to hold assets such as CSS and JS. It provides a mechanism to
// lazily calculate sha256sum's for integrity attributes and cache-busting.
type Asset struct {
	url          string
	raw          []byte
	immutableUrl string
	sum          []byte
	base64       string
}

// Raw returns the raw contents of the asset.
func (a *Asset) Raw() []byte {
	raw := a.raw

	return raw
}

// Sum returns the sha256sum of the assets contents.
func (a *Asset) Sum() []byte {
	sum := a.sum

	if len(sum) < 1 {
		shasum := sha256.Sum256(a.Raw())
		sum = shasum[:]
		a.sum = sum
	}

	return sum
}

// Base64Sum returns the base64 encoded sha256sum of the assets contents.
func (a *Asset) Base64Sum() string {
	base64sum := a.base64

	if base64sum == "" {
		base64sum = base64.StdEncoding.EncodeToString(a.Sum())
		a.base64 = base64sum
	}

	return base64sum
}

// Url returns the assets url. If immutable is set includes a URL-safe base64
// encoded cache-buster using the sha256sum of the contents.
func (a *Asset) Url(immutable bool) string {
	immutableUrl := a.immutableUrl
	url := a.url

	if immutable && immutableUrl == "" {
		b64 := base64.RawURLEncoding.EncodeToString(a.Sum())
		immutableUrl = extRegexp.ReplaceAllString(url, fmt.Sprintf("-%s${1}", b64))
		a.immutableUrl = immutableUrl
	}

	if immutable {
		return immutableUrl
	}

	return url
}
