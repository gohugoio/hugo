// Copyright 2024 The Hugo Authors. All rights reserved.
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
//
// Contains an embedded version of livereload.js
//
// Copyright (c) 2010-2015 Andrey Tarantsov
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package livereload

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	_ "embed"

	"github.com/gohugoio/hugo/media"
	"github.com/gorilla/websocket"
)

// Prefix to signal to LiveReload that we need to navigate to another path.
// Do not change this.
const hugoNavigatePrefix = "__hugo_navigate"

var upgrader = &websocket.Upgrader{
	// Hugo may potentially spin up multiple HTTP servers, so we need to exclude the
	// port when checking the origin.
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}

		rHost := r.Host
		// For Github codespace in browser #9936
		if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
			rHost = forwardedHost
		}

		if u.Host == rHost {
			return true
		}

		h1, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			return false
		}
		h2, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			return false
		}

		return h1 == h2
	},
	ReadBufferSize: 1024, WriteBufferSize: 1024,
}

// Handler is a HandlerFunc handling the livereload
// Websocket interaction.
func Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	wsHub.register <- c
	defer func() { wsHub.unregister <- c }()
	go c.writer()
	c.reader()
}

// Initialize starts the Websocket Hub handling live reloads.
func Initialize() {
	go wsHub.run()
}

// ForceRefresh tells livereload to force a hard refresh.
func ForceRefresh() {
	RefreshPath("/x.js")
}

// NavigateToPathForPort is similar to NavigateToPath but will also
// set window.location.port to the given port value.
func NavigateToPathForPort(path string, port int) {
	refreshPathForPort(hugoNavigatePrefix+path, port)
}

// RefreshPath tells livereload to refresh only the given path.
// If that path points to a CSS stylesheet or an image, only the changes
// will be updated in the browser, not the entire page.
func RefreshPath(s string) {
	refreshPathForPort(s, -1)
}

func refreshPathForPort(s string, port int) {
	// Tell livereload a file has changed - will force a hard refresh if not CSS or an image
	urlPath := filepath.ToSlash(s)
	portStr := ""
	if port > 0 {
		portStr = fmt.Sprintf(`, "overrideURL": %d`, port)
	}
	msg := fmt.Sprintf(`{"command":"reload","path":%q,"originalPath":"","liveCSS":true,"liveImg":true%s}`, urlPath, portStr)
	wsHub.broadcast <- []byte(msg)
}

// ServeJS serves the livereload.js who's reference is injected into the page.
func ServeJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", media.Builtin.JavascriptType.Type)
	w.Write(liveReloadJS())
}

func liveReloadJS() []byte {
	return []byte(livereloadJS)
}

//go:embed livereload.min.js
var livereloadJS string
