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

package livereload

import (
	"bytes"
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// There is a potential data race, especially visible with large files.
	// This is protected by synchronization of the send channel's close.
	closer sync.Once
}

func (c *connection) close() {
	c.closer.Do(func() {
		close(c.send)
	})
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		if bytes.Contains(message, []byte(`"command":"hello"`)) {
			c.send <- []byte(`{
				"command": "hello",
				"protocols": [ "http://livereload.com/protocols/official-7" ],
				"serverName": "Hugo"
			}`)
		}
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}
