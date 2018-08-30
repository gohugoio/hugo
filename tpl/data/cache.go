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

package data

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
)

var cacheMu sync.RWMutex

// getCacheFileID returns the cache ID for a string.
func getCacheFileID(cfg config.Provider, id string) string {
	hash := md5.Sum([]byte(id))
	return cfg.GetString("cacheDir") + hex.EncodeToString(hash[:])
}

// getCache returns the content for an ID from the file cache or an error.
// If the ID is not found, return nil,nil.
func getCache(id string, fs afero.Fs, cfg config.Provider, ignoreCache bool) ([]byte, error) {
	if ignoreCache {
		return nil, nil
	}

	cacheMu.RLock()
	defer cacheMu.RUnlock()

	fID := getCacheFileID(cfg, id)
	isExists, err := helpers.Exists(fID, fs)
	if err != nil {
		return nil, err
	}
	if !isExists {
		return nil, nil
	}

	return afero.ReadFile(fs, fID)
}

// writeCache writes bytes associated with an ID into the file cache.
func writeCache(id string, c []byte, fs afero.Fs, cfg config.Provider, ignoreCache bool) error {
	if ignoreCache {
		return nil
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	fID := getCacheFileID(cfg, id)
	f, err := fs.Create(fID)
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to create file: " + fID)
	}
	defer f.Close()

	n, err := f.Write(c)
	if err != nil {
		return errors.New("Error: " + err.Error() + ". Failed to write to file: " + fID)
	}
	if n == 0 {
		return errors.New("No bytes written to file: " + fID)
	}
	return nil
}

func deleteCache(id string, fs afero.Fs, cfg config.Provider) error {
	return fs.Remove(getCacheFileID(cfg, id))
}
