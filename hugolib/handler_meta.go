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

	"github.com/gohugoio/hugo/source"
)

var handlers []Handler

// MetaHandler abstracts reading and converting functionality of a Handler.
type MetaHandler interface {
	// Read the Files in and register
	Read(*source.File, *Site, HandleResults)

	// Generic Convert Function with coordination
	Convert(interface{}, *Site, HandleResults)

	Handle() Handler
}

// HandleResults is a channel for HandledResult.
type HandleResults chan<- HandledResult

// NewMetaHandler creates a MetaHandle for a given extensions.
func NewMetaHandler(in string) *MetaHandle {
	x := &MetaHandle{ext: in}
	x.Handler()
	return x
}

// MetaHandle is a generic MetaHandler that internally uses
// the globally registered handlers for handling specific file types.
type MetaHandle struct {
	handler Handler
	ext     string
}

func (mh *MetaHandle) Read(f *source.File, s *Site, results HandleResults) {
	if h := mh.Handler(); h != nil {
		results <- h.Read(f, s)
		return
	}

	results <- HandledResult{err: errors.New("No handler found"), file: f}
}

// Convert handles the conversion of files and pages.
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

// Handler finds the registered handler for the used extensions.
func (mh *MetaHandle) Handler() Handler {
	if mh.handler == nil {
		mh.handler = FindHandler(mh.ext)

		// if no handler found, use default handler
		if mh.handler == nil {
			mh.handler = FindHandler("*")
		}
	}
	return mh.handler
}

// FindHandler finds a Handler in the globally registered handlers.
func FindHandler(ext string) Handler {
	for _, h := range Handlers() {
		if HandlerMatch(h, ext) {
			return h
		}
	}
	return nil
}

// HandlerMatch checks if the given extensions matches.
func HandlerMatch(h Handler, ext string) bool {
	for _, x := range h.Extensions() {
		if ext == x {
			return true
		}
	}
	return false
}

// RegisterHandler adds a handler to the globally registered ones.
func RegisterHandler(h Handler) {
	handlers = append(handlers, h)
}

// Handlers returns the globally registered handlers.
func Handlers() []Handler {
	return handlers
}
