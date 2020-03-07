// Copyright 2016 The Hugo Authors. All rights reserved.
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

package data

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/cache/filecache"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
)

var (
	resSleep   = time.Second * 2 // if JSON decoding failed sleep for n seconds before retrying
	resRetries = 1               // number of retries to load the JSON from URL
)

// getRemote loads the content of a remote file. This method is thread safe.
func (ns *Namespace) getRemote(cache *filecache.Cache, unmarshal func([]byte) (bool, error), req *http.Request) error {
	url := req.URL.String()
	id := helpers.MD5String(url)
	var handled bool
	var retry bool

	_, b, err := cache.GetOrCreateBytes(id, func() ([]byte, error) {
		var err error
		handled = true
		for i := 0; i <= resRetries; i++ {
			ns.deps.Log.INFO.Printf("Downloading: %s ...", url)
			var res *http.Response
			res, err = ns.client.Do(req)
			if err != nil {
				return nil, err
			}

			if isHTTPError(res) {
				return nil, errors.Errorf("Failed to retrieve remote file: %s", http.StatusText(res.StatusCode))
			}

			var b []byte
			b, err = ioutil.ReadAll(res.Body)

			if err != nil {
				return nil, err
			}
			res.Body.Close()

			retry, err = unmarshal(b)

			if err == nil {
				// Return it so it can be cached.
				return b, nil
			}

			if !retry {
				return nil, err
			}

			ns.deps.Log.INFO.Printf("Cannot read remote resource %s: %s", url, err)
			ns.deps.Log.INFO.Printf("Retry #%d for %s and sleeping for %s", i+1, url, resSleep)
			time.Sleep(resSleep)
		}

		return nil, err

	})

	if !handled {
		// This is cached content and should be correct.
		_, err = unmarshal(b)
	}

	return err
}

// getLocal loads the content of a local file
func getLocal(url string, fs afero.Fs, cfg config.Provider) ([]byte, error) {
	filename := filepath.Join(cfg.GetString("workingDir"), url)
	if e, err := helpers.Exists(filename, fs); !e {
		return nil, err
	}

	return afero.ReadFile(fs, filename)

}

// getResource loads the content of a local or remote file and returns its content and the
// cache ID used, if relevant.
func (ns *Namespace) getResource(cache *filecache.Cache, unmarshal func(b []byte) (bool, error), req *http.Request) error {
	switch req.URL.Scheme {
	case "":
		url, err := url.QueryUnescape(req.URL.String())
		if err != nil {
			return err
		}
		b, err := getLocal(url, ns.deps.Fs.Source, ns.deps.Cfg)
		if err != nil {
			return err
		}
		_, err = unmarshal(b)
		return err
	default:
		return ns.getRemote(cache, unmarshal, req)
	}
}

func isHTTPError(res *http.Response) bool {
	return res.StatusCode < 200 || res.StatusCode > 299
}
