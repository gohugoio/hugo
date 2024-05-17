// Copyright 2021 The Hugo Authors. All rights reserved.
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

package create

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	gmaps "maps"

	"github.com/gohugoio/httpcache"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/tasks"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/mitchellh/mapstructure"
)

type HTTPError struct {
	error
	Data map[string]any

	StatusCode int
	Body       string
}

func responseToData(res *http.Response, readBody bool) map[string]any {
	var body []byte
	if readBody {
		body, _ = io.ReadAll(res.Body)
	}

	m := map[string]any{
		"StatusCode":       res.StatusCode,
		"Status":           res.Status,
		"TransferEncoding": res.TransferEncoding,
		"ContentLength":    res.ContentLength,
		"ContentType":      res.Header.Get("Content-Type"),
	}

	if readBody {
		m["Body"] = string(body)
	}

	return m
}

func toHTTPError(err error, res *http.Response, readBody bool) *HTTPError {
	if err == nil {
		panic("err is nil")
	}
	if res == nil {
		return &HTTPError{
			error: err,
			Data:  map[string]any{},
		}
	}

	return &HTTPError{
		error: err,
		Data:  responseToData(res, readBody),
	}
}

var temporaryHTTPStatusCodes = map[int]bool{
	408: true,
	429: true,
	500: true,
	502: true,
	503: true,
	504: true,
}

func (c *Client) configurePollingIfEnabled(uri, optionsKey string, getRes func() (*http.Response, error)) {
	if c.remoteResourceChecker == nil {
		return
	}

	// Set up polling for changes to this resource.
	pollingConfig := c.httpCacheConfig.PollConfigFor(uri)
	if pollingConfig.IsZero() || pollingConfig.Config.Disable {
		return
	}

	if c.remoteResourceChecker.Has(optionsKey) {
		return
	}

	var lastChange time.Time
	c.remoteResourceChecker.Add(optionsKey,
		tasks.Func{
			IntervalLow:  pollingConfig.Config.Low,
			IntervalHigh: pollingConfig.Config.High,
			F: func(interval time.Duration) (time.Duration, error) {
				start := time.Now()
				defer func() {
					duration := time.Since(start)
					c.rs.Logger.Debugf("Polled remote resource for changes in %13s. Interval: %4s (low: %4s high: %4s) resource: %q ", duration, interval, pollingConfig.Config.Low, pollingConfig.Config.High, uri)
				}()
				// TODO(bep) figure out a ways to remove unused tasks.
				res, err := getRes()
				if err != nil {
					return pollingConfig.Config.High, err
				}
				// The caching is delayed until the body is read.
				io.Copy(io.Discard, res.Body)
				res.Body.Close()
				x1, x2 := res.Header.Get(httpcache.XETag1), res.Header.Get(httpcache.XETag2)
				if x1 != x2 {
					lastChange = time.Now()
					c.remoteResourceLogger.Logf("detected change in remote resource %q", uri)
					c.rs.Rebuilder.SignalRebuild(identity.StringIdentity(optionsKey))
				}

				if time.Since(lastChange) < 10*time.Second {
					// The user is typing, check more often.
					return 0, nil
				}

				// Increase the interval to avoid hammering the server.
				interval += 1 * time.Second

				return interval, nil
			},
		})
}

// FromRemote expects one or n-parts of a URL to a resource
// If you provide multiple parts they will be joined together to the final URL.
func (c *Client) FromRemote(uri string, optionsm map[string]any) (resource.Resource, error) {
	rURL, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL for resource %s: %w", uri, err)
	}

	method := "GET"
	if s, _, ok := maps.LookupEqualFold(optionsm, "method"); ok {
		method = strings.ToUpper(s.(string))
	}
	isHeadMethod := method == "HEAD"

	optionsm = gmaps.Clone(optionsm)
	userKey, optionsKey := remoteResourceKeys(uri, optionsm)

	// A common pattern is to use the key in the options map as
	// a way to control cache eviction,
	// so make sure we use any user provided kehy as the file cache key,
	// but the auto generated and more stable key for everything else.
	filecacheKey := userKey

	return c.rs.ResourceCache.CacheResourceRemote.GetOrCreate(optionsKey, func(key string) (resource.Resource, error) {
		options, err := decodeRemoteOptions(optionsm)
		if err != nil {
			return nil, fmt.Errorf("failed to decode options for resource %s: %w", uri, err)
		}

		if err := c.validateFromRemoteArgs(uri, options); err != nil {
			return nil, err
		}

		getRes := func() (*http.Response, error) {
			ctx := context.Background()
			ctx = c.resourceIDDispatcher.Set(ctx, filecacheKey)

			req, err := options.NewRequest(uri)
			if err != nil {
				return nil, fmt.Errorf("failed to create request for resource %s: %w", uri, err)
			}

			req = req.WithContext(ctx)

			return c.httpClient.Do(req)
		}

		res, err := getRes()
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		c.configurePollingIfEnabled(uri, optionsKey, getRes)

		if res.StatusCode == http.StatusNotFound {
			// Not found. This matches how lookups for local resources work.
			return nil, nil
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return nil, toHTTPError(fmt.Errorf("failed to fetch remote resource: %s", http.StatusText(res.StatusCode)), res, !isHeadMethod)
		}

		var (
			body      []byte
			mediaType media.Type
		)
		// A response to a HEAD method should not have a body. If it has one anyway, that body must be ignored.
		// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/HEAD
		if !isHeadMethod && res.Body != nil {
			body, err = io.ReadAll(res.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read remote resource %q: %w", uri, err)
			}
		}

		filename := path.Base(rURL.Path)
		if _, params, _ := mime.ParseMediaType(res.Header.Get("Content-Disposition")); params != nil {
			if _, ok := params["filename"]; ok {
				filename = params["filename"]
			}
		}

		contentType := res.Header.Get("Content-Type")

		// For HEAD requests we have no body to work with, so we need to use the Content-Type header.
		if isHeadMethod || c.rs.ExecHelper.Sec().HTTP.MediaTypes.Accept(contentType) {
			var found bool
			mediaType, found = c.rs.MediaTypes().GetByType(contentType)
			if !found {
				// A media type not configured in Hugo, just create one from the content type string.
				mediaType, _ = media.FromString(contentType)
			}
		}

		if mediaType.IsZero() {

			var extensionHints []string

			// mime.ExtensionsByType gives a long list of extensions for text/plain,
			// just use ".txt".
			if strings.HasPrefix(contentType, "text/plain") {
				extensionHints = []string{".txt"}
			} else {
				exts, _ := mime.ExtensionsByType(contentType)
				if exts != nil {
					extensionHints = exts
				}
			}

			// Look for a file extension. If it's .txt, look for a more specific.
			if extensionHints == nil || extensionHints[0] == ".txt" {
				if ext := path.Ext(filename); ext != "" {
					extensionHints = []string{ext}
				}
			}

			// Now resolve the media type primarily using the content.
			mediaType = media.FromContent(c.rs.MediaTypes(), extensionHints, body)

		}

		if mediaType.IsZero() {
			return nil, fmt.Errorf("failed to resolve media type for remote resource %q", uri)
		}

		userKey = filename[:len(filename)-len(path.Ext(filename))] + "_" + userKey + mediaType.FirstSuffix.FullSuffix
		data := responseToData(res, false)

		return c.rs.NewResource(
			resources.ResourceSourceDescriptor{
				MediaType:     mediaType,
				Data:          data,
				GroupIdentity: identity.StringIdentity(optionsKey),
				LazyPublish:   true,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloser(bytes.NewReader(body)), nil
				},
				TargetPath: userKey,
			})
	})
}

func (c *Client) validateFromRemoteArgs(uri string, options fromRemoteOptions) error {
	if err := c.rs.ExecHelper.Sec().CheckAllowedHTTPURL(uri); err != nil {
		return err
	}

	if err := c.rs.ExecHelper.Sec().CheckAllowedHTTPMethod(options.Method); err != nil {
		return err
	}

	return nil
}

func remoteResourceKeys(uri string, optionsm map[string]any) (string, string) {
	var userKey string
	if key, k, found := maps.LookupEqualFold(optionsm, "key"); found {
		userKey = identity.HashString(key)
		delete(optionsm, k)
	}
	optionsKey := identity.HashString(uri, optionsm)
	if userKey == "" {
		userKey = optionsKey
	}
	return userKey, optionsKey
}

func addDefaultHeaders(req *http.Request) {
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

func hasHeaderKey(m http.Header, key string) bool {
	_, ok := m[key]
	return ok
}

type fromRemoteOptions struct {
	Method  string
	Headers map[string]any
	Body    []byte
}

func (o fromRemoteOptions) BodyReader() io.Reader {
	if o.Body == nil {
		return nil
	}
	return bytes.NewBuffer(o.Body)
}

func (o fromRemoteOptions) NewRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(o.Method, url, o.BodyReader())
	if err != nil {
		return nil, err
	}

	// First add any user provided headers.
	if o.Headers != nil {
		addUserProvidedHeaders(o.Headers, req)
	}

	// Then add default headers not provided by the user.
	addDefaultHeaders(req)

	return req, nil
}

func decodeRemoteOptions(optionsm map[string]any) (fromRemoteOptions, error) {
	options := fromRemoteOptions{
		Method: "GET",
	}

	err := mapstructure.WeakDecode(optionsm, &options)
	if err != nil {
		return options, err
	}
	options.Method = strings.ToUpper(options.Method)

	return options, nil
}

var _ http.RoundTripper = (*transport)(nil)

type transport struct {
	Cfg    config.AllProvider
	Logger loggers.Logger
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	defer func() {
		if resp != nil && resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusNotModified {
			t.Logger.Debugf("Fetched remote resource: %s", req.URL.String())
		}
	}()

	var (
		start          time.Time
		nextSleep      = time.Duration((rand.Intn(1000) + 100)) * time.Millisecond
		nextSleepLimit = time.Duration(5) * time.Second
		retry          bool
	)

	for {
		resp, retry, err = func() (*http.Response, bool, error) {
			resp2, err := http.DefaultTransport.RoundTrip(req)
			if err != nil {
				return resp2, false, err
			}

			if resp2.StatusCode != http.StatusNotFound && resp2.StatusCode != http.StatusNotModified {
				if resp2.StatusCode < 200 || resp2.StatusCode > 299 {
					return resp2, temporaryHTTPStatusCodes[resp2.StatusCode], nil
				}
			}
			return resp2, false, nil
		}()

		if retry {
			if start.IsZero() {
				start = time.Now()
			} else if d := time.Since(start) + nextSleep; d >= t.Cfg.Timeout() {
				msg := "<nil>"
				if resp != nil {
					msg = resp.Status
				}
				err := toHTTPError(fmt.Errorf("retry timeout (configured to %s) fetching remote resource: %s", t.Cfg.Timeout(), msg), resp, req.Method != "HEAD")
				return resp, err
			}
			time.Sleep(nextSleep)
			if nextSleep < nextSleepLimit {
				nextSleep *= 2
			}
			continue
		}

		return
	}
}

// We need to send the redirect responses back to the HTTP client from RoundTrip,
// but we don't want to cache them.
func shouldCache(statusCode int) bool {
	switch statusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect, http.StatusPermanentRedirect:
		return false
	}
	return true
}
