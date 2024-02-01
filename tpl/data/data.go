// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package data provides template functions for working with external data
// sources.
package data

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config/security"

	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/common/constants"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/deps"
)

// New returns a new instance of the data-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps:         deps,
		cacheGetCSV:  deps.ResourceSpec.FileCaches.GetCSVCache(),
		cacheGetJSON: deps.ResourceSpec.FileCaches.GetJSONCache(),
		client:       http.DefaultClient,
	}
}

// Namespace provides template functions for the "data" namespace.
type Namespace struct {
	deps *deps.Deps

	cacheGetJSON *filecache.Cache
	cacheGetCSV  *filecache.Cache

	client *http.Client
}

// GetCSV expects the separator sep and one or n-parts of a URL to a resource which
// can either be a local or a remote one.
// The data separator can be a comma, semi-colon, pipe, etc, but only one character.
// If you provide multiple parts for the URL they will be joined together to the final URL.
// GetCSV returns nil or a slice slice to use in a short code.
func (ns *Namespace) GetCSV(sep string, args ...any) (d [][]string, err error) {
	hugo.Deprecate("data.GetCSV", "use resources.Get or resources.GetRemote with transform.Unmarshal.", "v0.123.0")

	url, headers := toURLAndHeaders(args)
	cache := ns.cacheGetCSV

	unmarshal := func(b []byte) (bool, error) {
		if d, err = parseCSV(b, sep); err != nil {
			err = fmt.Errorf("failed to parse CSV file %s: %w", url, err)

			return true, err
		}

		return false, nil
	}

	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for getCSV for resource %s: %w", url, err)
	}

	// Add custom user headers.
	addUserProvidedHeaders(headers, req)
	addDefaultHeaders(req, "text/csv", "text/plain")

	err = ns.getResource(cache, unmarshal, req)
	if err != nil {
		if security.IsAccessDenied(err) {
			return nil, err
		}
		ns.deps.Log.Erroridf(constants.ErrRemoteGetCSV, "Failed to get CSV resource %q: %s", url, err)
		return nil, nil
	}

	return
}

// GetJSON expects one or n-parts of a URL in args to a resource which can either be a local or a remote one.
// If you provide multiple parts they will be joined together to the final URL.
// GetJSON returns nil or parsed JSON to use in a short code.
func (ns *Namespace) GetJSON(args ...any) (any, error) {
	hugo.Deprecate("data.GetJSON", "use resources.Get or resources.GetRemote with transform.Unmarshal.", "v0.123.0")

	var v any
	url, headers := toURLAndHeaders(args)
	cache := ns.cacheGetJSON

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for getJSON resource %s: %w", url, err)
	}

	unmarshal := func(b []byte) (bool, error) {
		err := json.Unmarshal(b, &v)
		if err != nil {
			return true, err
		}
		return false, nil
	}

	addUserProvidedHeaders(headers, req)
	addDefaultHeaders(req, "application/json")

	err = ns.getResource(cache, unmarshal, req)
	if err != nil {
		if security.IsAccessDenied(err) {
			return nil, err
		}
		ns.deps.Log.Erroridf(constants.ErrRemoteGetJSON, "Failed to get JSON resource %q: %s", url, err)
		return nil, nil
	}

	return v, nil
}

func addDefaultHeaders(req *http.Request, accepts ...string) {
	for _, accept := range accepts {
		if !hasHeaderValue(req.Header, "Accept", accept) {
			req.Header.Add("Accept", accept)
		}
	}
	if !hasHeaderKey(req.Header, "User-Agent") {
		req.Header.Add("User-Agent", "Hugo Static Site Generator")
	}
}

func addUserProvidedHeaders(headers map[string]any, req *http.Request) {
	if headers == nil {
		return
	}
	for key, val := range headers {
		vals := types.ToStringSlicePreserveString(val)
		for _, s := range vals {
			req.Header.Add(key, s)
		}
	}
}

func hasHeaderValue(m http.Header, key, value string) bool {
	var s []string
	var ok bool

	if s, ok = m[key]; !ok {
		return false
	}

	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func hasHeaderKey(m http.Header, key string) bool {
	_, ok := m[key]
	return ok
}

func toURLAndHeaders(urlParts []any) (string, map[string]any) {
	if len(urlParts) == 0 {
		return "", nil
	}

	// The last argument may be a map.
	headers, err := maps.ToStringMapE(urlParts[len(urlParts)-1])
	if err == nil {
		urlParts = urlParts[:len(urlParts)-1]
	} else {
		headers = nil
	}

	return strings.Join(cast.ToStringSlice(urlParts), ""), headers
}

// parseCSV parses bytes of CSV data into a slice slice string or an error
func parseCSV(c []byte, sep string) ([][]string, error) {
	if len(sep) != 1 {
		return nil, errors.New("Incorrect length of CSV separator: " + sep)
	}
	b := bytes.NewReader(c)
	r := csv.NewReader(b)
	rSep := []rune(sep)
	r.Comma = rSep[0]
	r.FieldsPerRecord = 0
	return r.ReadAll()
}
