// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"errors"

	"fmt"

	"github.com/spf13/hugo/source"
	"regexp"
)

var handlers []Handler

type MetaHandler interface {
	// Read the Files in and register
	Read(*source.File, *Site, HandleResults)

	// Generic Convert Function with coordination
	Convert(interface{}, *Site, HandleResults)

	Handle() Handler
}

type HandleResults chan<- HandledResult

func NewMetaHandler(ext string) *MetaHandle {
	return NewFilenameMetaHandler(ext, "*")
}

func NewFilenameMetaHandler(ext string, filename string) *MetaHandle {
	x := &MetaHandle{ext: ext, filename: filename}
	x.Handler()
	return x
}

type MetaHandle struct {
	handler  Handler
	ext      string
	filename string
}

func (mh *MetaHandle) Read(f *source.File, s *Site, results HandleResults) {
	if h := mh.Handler(); h != nil {
		results <- h.Read(f, s)
		return
	}

	results <- HandledResult{err: errors.New("No handler found"), file: f}
}

func (mh *MetaHandle) Convert(i interface{}, s *Site, results HandleResults) {
	h := mh.Handler()

	if f, ok := i.(*source.File); ok {
		results <- h.FileConvert(f, s)
		return
	}

	if p, ok := i.(*Page); ok {
		if p == nil {
			results <- HandledResult{err: errors.New("file resulted in a nil page")}
			return
		}

		if h == nil {
			results <- HandledResult{err: fmt.Errorf("No handler found for page '%s'. Verify the markup is supported by Hugo.", p.FullFilePath())}
			return
		}

		results <- h.PageConvert(p)
	}
}

func (mh *MetaHandle) Handler() Handler {
	if mh.handler == nil {
		mh.handler = FindHandler(mh.ext, mh.filename)

		// if no handler found, use default handler
		if mh.handler == nil {
			mh.handler = FindHandler("*", "*")
		}
	}
	return mh.handler
}

func FindHandler(ext string, filename string) Handler {
	for _, h := range Handlers() {
		if HandlerMatch(h, ext, filename) {
			return h
		}
	}
	return nil
}

func HandlerMatch(h Handler, ext string, filename string) bool {
	match := true
	if len(h.FilenamePatterns()) > 0 && filename != "*" {
		match = checkFileNameMatch(filename, h.FilenamePatterns())
	}
	if len(h.IgnoreFilenamePatterns()) > 0 && filename != "*" {
		match = match && !checkFileNameMatch(filename, h.IgnoreFilenamePatterns())
	}
	return match && checkFileExtension(h, ext)
}

func checkFileNameMatch(filename string, filenamePatterns []string) bool {
	for _, filenamePattern := range filenamePatterns {
		if match, err := regexp.MatchString(filenamePattern, filename); err == nil && match {
			return true
		}
	}
	return false
}

func checkFileExtension(h Handler, ext string) bool {
	for _, x := range h.Extensions() {
		if ext == x {
			return true
		}
	}
	return false
}

func RegisterHandler(h Handler) {
	handlers = append(handlers, h)
}

func Handlers() []Handler {
	return handlers
}
