// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpl

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var remoteURLLock = &remoteLock{m: make(map[string]*sync.Mutex)}

type remoteLock struct {
	sync.RWMutex
	m map[string]*sync.Mutex
}

// URLLock locks an URL during download
func (l *remoteLock) URLLock(url string) {
	l.Lock()
	if _, ok := l.m[url]; !ok {
		l.m[url] = &sync.Mutex{}
	}
	l.Unlock() // call this Unlock before the next lock will be called. NFI why but defer doesn't work.
	l.m[url].Lock()
}

// URLUnlock unlocks an URL when the download has been finished. Use only in defer calls.
func (l *remoteLock) URLUnlock(url string) {
	l.RLock()
	defer l.RUnlock()
	if um, ok := l.m[url]; ok {
		um.Unlock()
	}
}

// getCacheFileID returns the cache ID for a string
func getCacheFileID(id string) string {
	return viper.GetString("CacheDir") + url.QueryEscape(id)
}

// resGetCache returns the content for an ID from the file cache or an error
// if the file is not found returns nil,nil
func resGetCache(id string, fs afero.Fs, ignoreCache bool) ([]byte, error) {
	if ignoreCache {
		return nil, nil
	}
	fID := getCacheFileID(id)
	isExists, err := helpers.Exists(fID, fs)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, nil
	}

	f, err := fs.Open(fID)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// resWriteCache writes bytes to an ID into the file cache
func resWriteCache(id string, c []byte, fs afero.Fs) error {
	fID := getCacheFileID(id)
	f, err := fs.Create(fID)
	if err != nil {
		return err
	}
	n, err := f.Write(c)
	if n == 0 {
		return errors.New("No bytes written to file: " + fID)
	}
	return err
}

// resGetRemote loads the content of a remote file. This method is thread safe.
func resGetRemote(url string, fs afero.Fs, hc *http.Client) ([]byte, error) {

	c, err := resGetCache(url, fs, viper.GetBool("IgnoreCache"))
	if c != nil && err == nil {
		return c, nil
	}
	if err != nil {
		return nil, err
	}

	// avoid race condition with locks, block other goroutines if the current url is processing
	remoteURLLock.URLLock(url)
	defer func() { remoteURLLock.URLUnlock(url) }()

	// avoid multiple locks due to calling resGetCache twice
	c, err = resGetCache(url, fs, viper.GetBool("IgnoreCache"))
	if c != nil && err == nil {
		return c, nil
	}
	if err != nil {
		return nil, err
	}

	jww.INFO.Printf("Downloading: %s ...", url)
	res, err := hc.Get(url)
	if err != nil {
		return nil, err
	}
	c, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = resWriteCache(url, c, fs)
	if err != nil {
		return nil, err
	}
	jww.INFO.Printf("... and cached to: %s", getCacheFileID(url))
	return c, nil
}

// resGetLocal loads the content of a local file
func resGetLocal(url string, fs afero.Fs) ([]byte, error) {
	p := ""
	if viper.GetString("WorkingDir") != "" {
		p = viper.GetString("WorkingDir")
		if helpers.FilePathSeparator != p[len(p)-1:] {
			p = p + helpers.FilePathSeparator
		}
	}
	jFile := p + url
	if e, err := helpers.Exists(jFile, fs); !e {
		return nil, err
	}

	f, err := fs.Open(jFile)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// resGetResource loads the content of a local or remote file
func resGetResource(url string) ([]byte, error) {
	if url == "" {
		return nil, nil
	}
	if strings.Contains(url, "://") {
		return resGetRemote(url, hugofs.SourceFs, http.DefaultClient)
	}
	return resGetLocal(url, hugofs.SourceFs)
}

// GetJSON expects one or n-parts of a URL to a resource which can either be a local or a remote one.
// If you provide multiple parts they will be joined together to the final URL.
// GetJSON returns nil or parsed JSON to use in a short code.
func GetJSON(urlParts ...string) interface{} {
	url := strings.Join(urlParts, "")
	c, err := resGetResource(url)
	if err != nil {
		jww.ERROR.Printf("Failed to get json resource %s with error message %s", url, err)
		return nil
	}

	var v interface{}
	err = json.Unmarshal(c, &v)
	if err != nil {
		jww.ERROR.Printf("Cannot read json from resource %s with error message %s", url, err)
		return nil
	}
	return v
}

// parseCSV parses bytes of CSV data into a slice slice string or an error
func parseCSV(c []byte, sep string) ([][]string, error) {
	if len(sep) != 1 {
		return nil, errors.New("Incorrect length of csv separator: " + sep)
	}
	b := bytes.NewReader(c)
	r := csv.NewReader(b)
	rSep := []rune(sep)
	r.Comma = rSep[0]
	r.FieldsPerRecord = 0
	return r.ReadAll()
}

// GetCSV expects a data separator and one or n-parts of a URL to a resource which
// can either be a local or a remote one.
// The data separator can be a comma, semi-colon, pipe, etc, but only one character.
// If you provide multiple parts for the URL they will be joined together to the final URL.
// GetCSV returns nil or a slice slice to use in a short code.
func GetCSV(sep string, urlParts ...string) [][]string {
	url := strings.Join(urlParts, "")
	c, err := resGetResource(url)
	if err != nil {
		jww.ERROR.Printf("Failed to get csv resource %s with error message %s", url, err)
		return nil
	}
	d, err := parseCSV(c, sep)
	if err != nil {
		jww.ERROR.Printf("Failed to read csv resource %s with error message %s", url, err)
		return nil
	}
	return d
}
