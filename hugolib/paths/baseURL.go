// Copyright 2018 The Hugo Authors. All rights reserved.
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

package paths

import (
	"fmt"
	"net/url"
	"strings"
)

// A BaseURL in Hugo is normally on the form scheme://path, but the
// form scheme: is also valid (mailto:hugo@rules.com).
type BaseURL struct {
	url    *url.URL
	urlStr string
}

func (b BaseURL) String() string {
	if b.urlStr != "" {
		return b.urlStr
	}
	return b.url.String()
}

func (b BaseURL) Path() string {
	return b.url.Path
}

// HostURL returns the URL to the host root without any path elements.
func (b BaseURL) HostURL() string {
	return strings.TrimSuffix(b.String(), b.Path())
}

// WithProtocol returns the BaseURL prefixed with the given protocol.
// The Protocol is normally of the form "scheme://", i.e. "webcal://".
func (b BaseURL) WithProtocol(protocol string) (string, error) {
	u := b.URL()

	scheme := protocol
	isFullProtocol := strings.HasSuffix(scheme, "://")
	isOpaqueProtocol := strings.HasSuffix(scheme, ":")

	if isFullProtocol {
		scheme = strings.TrimSuffix(scheme, "://")
	} else if isOpaqueProtocol {
		scheme = strings.TrimSuffix(scheme, ":")
	}

	u.Scheme = scheme

	if isFullProtocol && u.Opaque != "" {
		u.Opaque = "//" + u.Opaque
	} else if isOpaqueProtocol && u.Opaque == "" {
		return "", fmt.Errorf("Cannot determine BaseURL for protocol %q", protocol)
	}

	return u.String(), nil
}

// URL returns a copy of the internal URL.
// The copy can be safely used and modified.
func (b BaseURL) URL() *url.URL {
	c := *b.url
	return &c
}

func newBaseURLFromString(b string) (BaseURL, error) {
	var result BaseURL

	base, err := url.Parse(b)
	if err != nil {
		return result, err
	}

	return BaseURL{url: base, urlStr: base.String()}, nil
}
