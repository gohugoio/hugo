// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package create contains functions for to create Resource objects. This will
// typically non-files.
package create

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/pkg/errors"
)

// Client contains methods to create Resource objects.
// tasks to Resource objects.
type Client struct {
	rs         *resources.Spec
	httpClient *http.Client
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{
		rs: rs,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get creates a new Resource by opening the given filename in the assets filesystem.
func (c *Client) Get(filename string) (resource.Resource, error) {
	filename = filepath.Clean(filename)
	return c.rs.ResourceCache.GetOrCreate(resources.ResourceCacheKey(filename), func() (resource.Resource, error) {
		return c.rs.New(resources.ResourceSourceDescriptor{
			Fs:             c.rs.BaseFs.Assets.Fs,
			LazyPublish:    true,
			SourceFilename: filename,
		})
	})
}

// Match gets the resources matching the given pattern from the assets filesystem.
func (c *Client) Match(pattern string) (resource.Resources, error) {
	return c.match(pattern, false)
}

// GetMatch gets first resource matching the given pattern from the assets filesystem.
func (c *Client) GetMatch(pattern string) (resource.Resource, error) {
	res, err := c.match(pattern, true)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res[0], err
}

func (c *Client) match(pattern string, firstOnly bool) (resource.Resources, error) {
	var name string
	if firstOnly {
		name = "__get-match"
	} else {
		name = "__match"
	}

	pattern = glob.NormalizePath(pattern)
	partitions := glob.FilterGlobParts(strings.Split(pattern, "/"))
	if len(partitions) == 0 {
		partitions = []string{resources.CACHE_OTHER}
	}
	key := path.Join(name, path.Join(partitions...))
	key = path.Join(key, pattern)

	return c.rs.ResourceCache.GetOrCreateResources(key, func() (resource.Resources, error) {
		var res resource.Resources

		handle := func(info hugofs.FileMetaInfo) (bool, error) {
			meta := info.Meta()
			r, err := c.rs.New(resources.ResourceSourceDescriptor{
				LazyPublish: true,
				FileInfo:    info,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return meta.Open()
				},
				RelTargetFilename: meta.Path,
			})
			if err != nil {
				return true, err
			}

			res = append(res, r)

			return firstOnly, nil
		}

		if err := hugofs.Glob(c.rs.BaseFs.Assets.Fs, pattern, handle); err != nil {
			return nil, err
		}

		return res, nil
	})
}

// FromString creates a new Resource from a string with the given relative target path.
func (c *Client) FromString(targetPath, content string) (resource.Resource, error) {
	return c.rs.ResourceCache.GetOrCreate(path.Join(resources.CACHE_OTHER, targetPath), func() (resource.Resource, error) {
		return c.rs.New(
			resources.ResourceSourceDescriptor{
				Fs:          c.rs.FileCaches.AssetsCache().Fs,
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloserFromString(content), nil
				},
				RelTargetFilename: filepath.Clean(targetPath),
			})
	})
}

// FromRemote expects one or n-parts of a URL to a resource
// If you provide multiple parts they will be joined together to the final URL.
func (c *Client) FromRemote(uri string, options map[string]interface{}) (resource.Resource, error) {
	rURL, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse URL for resource %s", uri)
	}

	resourceID := helpers.HashString(uri, options)

	return c.rs.ResourceCache.GetOrCreate(resources.ResourceCacheKey(resourceID), func() (resource.Resource, error) {
		method, reqBody, err := getMethodAndBody(options)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get method or body for resource %s", uri)
		}

		req, err := http.NewRequest(method, uri, reqBody)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create request for resource %s", uri)
		}
		addDefaultHeaders(req)

		if _, ok := options["headers"]; ok {
			headers, err := maps.ToStringMapE(options["headers"])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse request headers for resource %s", uri)
			}
			addUserProvidedHeaders(headers, req)
		}
		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return nil, errors.Errorf("failed to retrieve remote resource: %s", http.StatusText(res.StatusCode))
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read remote resource %s", uri)
		}

		filename := path.Base(rURL.Path)
		if _, params, _ := mime.ParseMediaType(res.Header.Get("Content-Disposition")); params != nil {
			if _, ok := params["filename"]; ok {
				filename = params["filename"]
			}
		}

		var contentType string
		if arr, _ := mime.ExtensionsByType(res.Header.Get("Content-Type")); len(arr) == 1 {
			contentType = arr[0]
		}

		// If content type was not determined by header, look for a file extention
		if contentType == "" {
			if ext := path.Ext(filename); ext != "" {
				contentType = ext
			}
		}

		// If content type was not determined by header or file extention, try using content itself
		if contentType == "" {
			if ct := http.DetectContentType(body); ct != "application/octet-stream" {
				if arr, _ := mime.ExtensionsByType(ct); arr != nil {
					contentType = arr[0]
				}
			}
		}

		resourceID = filename[:len(filename)-len(path.Ext(filename))] + "_" + resourceID + contentType

		return c.rs.New(
			resources.ResourceSourceDescriptor{
				Fs:          c.rs.FileCaches.AssetsCache().Fs,
				LazyPublish: true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloser(bytes.NewReader(body)), nil
				},
				RelTargetFilename: filepath.Clean(resourceID),
			})
	})
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

func addUserProvidedHeaders(headers map[string]interface{}, req *http.Request) {
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

func getMethodAndBody(options map[string]interface{}) (string, io.Reader, error) {
	if options == nil {
		return "GET", nil, nil
	}

	if method, ok := options["method"].(string); ok {
		method = strings.ToUpper(method)
		switch method {
		case "GET", "DELETE", "HEAD", "OPTIONS":
			return method, nil, nil
		case "POST", "PUT", "PATCH":
			var body []byte
			if _, ok := options["body"]; ok {
				switch b := options["body"].(type) {
				case string:
					body = []byte(b)
				case []byte:
					body = b
				}
			}
			return method, bytes.NewBuffer(body), nil
		}

		return "", nil, fmt.Errorf("invalid HTTP method %q", method)
	}

	return "GET", nil, nil
}
