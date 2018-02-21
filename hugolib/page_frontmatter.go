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

package hugolib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
)

// TODO(bep) should probably make the date handling chain complete to give people the flexibility they want.

type frontmatterConfig struct {
	// Ordered chain.
	dateHandlers []frontmatterFieldHandler

	logger *jww.Notepad
}

func (f frontmatterConfig) handleField(handlers []frontmatterFieldHandler, frontmatter map[string]interface{}, p *Page) {
	for _, h := range handlers {
		handled, err := h(frontmatter, p)
		if err != nil {
			f.logger.ERROR.Println(err)
		}
		if handled {
			break
		}
	}
}

func (f frontmatterConfig) handleDate(frontmatter map[string]interface{}, p *Page) {
	f.handleField(f.dateHandlers, frontmatter, p)
}

type frontmatterFieldHandler func(frontmatter map[string]interface{}, p *Page) (bool, error)

func newFrontmatterConfig(logger *jww.Notepad, cfg config.Provider) (frontmatterConfig, error) {

	if logger == nil {
		logger = jww.NewNotepad(jww.LevelWarn, jww.LevelWarn, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	f := frontmatterConfig{logger: logger}

	handlers := &frontmatterFieldHandlers{logger: logger}

	f.dateHandlers = []frontmatterFieldHandler{handlers.defaultDateHandler}

	if cfg.IsSet("frontmatter") {
		fm := cfg.GetStringMap("frontmatter")
		if fm != nil {
			dateFallbacks, found := fm["defaultdate"]
			if found {
				slice, err := cast.ToStringSliceE(dateFallbacks)
				if err != nil {
					return f, fmt.Errorf("invalid value for dataCallbacks, expeced a string slice, got %T", dateFallbacks)
				}

				for _, v := range slice {
					if strings.EqualFold(v, "filename") {
						f.dateHandlers = append(f.dateHandlers, handlers.fileanameFallbackDateHandler)
						// No more for now.
						break
					}
				}
			}
		}
	}

	return f, nil
}

type frontmatterFieldHandlers struct {
	logger *jww.Notepad
}

// TODO(bep) modtime

func (f *frontmatterFieldHandlers) defaultDateHandler(frontmatter map[string]interface{}, p *Page) (bool, error) {
	loki := "date"
	v, found := frontmatter[loki]
	if !found {
		return false, nil
	}

	var err error
	p.Date, err = cast.ToTimeE(v)
	if err != nil {
		return false, fmt.Errorf("Failed to parse date %q in page %s", v, p.File.Path())
	}

	p.params[loki] = p.Date

	return true, nil
}

func (f *frontmatterFieldHandlers) fileanameFallbackDateHandler(frontmatter map[string]interface{}, p *Page) (bool, error) {
	return true, nil
}
