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
	"net/http"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/cache/filecache"
	"github.com/gohugoio/hugo/deps"
	_errors "github.com/pkg/errors"
)

// New returns a new instance of the data-namespaced template functions.
func New(deps *deps.Deps) *Namespace {

	return &Namespace{
		deps:         deps,
		cacheGetCSV:  deps.FileCaches.GetCSVCache(),
		cacheGetJSON: deps.FileCaches.GetJSONCache(),
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

// GetCSV expects a data separator and one or n-parts of a URL to a resource which
// can either be a local or a remote one.
// The data separator can be a comma, semi-colon, pipe, etc, but only one character.
// If you provide multiple parts for the URL they will be joined together to the final URL.
// GetCSV returns nil or a slice slice to use in a short code.
func (ns *Namespace) GetCSV(sep string, urlParts ...interface{}) (d [][]string, err error) {
	url := joinURL(urlParts)
	cache := ns.cacheGetCSV

	unmarshal := func(b []byte) (bool, error) {
		if !bytes.Contains(b, []byte(sep)) {
			return false, _errors.Errorf("cannot find separator %s in CSV for %s", sep, url)
		}

		if d, err = parseCSV(b, sep); err != nil {
			err = _errors.Wrapf(err, "failed to parse CSV file %s", url)

			return true, err
		}

		return false, nil
	}

	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, _errors.Wrapf(err, "failed to create request for getCSV for resource %s", url)
	}

	req.Header.Add("Accept", "text/csv")
	req.Header.Add("Accept", "text/plain")

	err = ns.getResource(cache, unmarshal, req)
	if err != nil {
		ns.deps.Log.ERROR.Printf("Failed to get CSV resource %q: %s", url, err)
		return nil, nil
	}

	return
}

// GetJSON expects one or n-parts of a URL to a resource which can either be a local or a remote one.
// If you provide multiple parts they will be joined together to the final URL.
// GetJSON returns nil or parsed JSON to use in a short code.
func (ns *Namespace) GetJSON(urlParts ...interface{}) (interface{}, error) {
	var v interface{}
	url := joinURL(urlParts)
	cache := ns.cacheGetJSON

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, _errors.Wrapf(err, "Failed to create request for getJSON resource %s", url)
	}

	unmarshal := func(b []byte) (bool, error) {
		err := json.Unmarshal(b, &v)
		if err != nil {
			return true, err
		}
		return false, nil
	}

	req.Header.Add("Accept", "application/json")

	err = ns.getResource(cache, unmarshal, req)
	if err != nil {
		ns.deps.Log.ERROR.Printf("Failed to get JSON resource %q: %s", url, err)
		return nil, nil
	}

	return v, nil
}

func joinURL(urlParts []interface{}) string {
	return strings.Join(cast.ToStringSlice(urlParts), "")
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
