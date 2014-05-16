// Copyright © 2013 Steve Francia <spf@spf13.com>.
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

package commands

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	//"code.google.com/p/go.net/websocket"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var serverPort int
var serverWatch bool
var serverAppend bool

func init() {
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port to run the server on")
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	serverCmd.Flags().BoolVarP(&serverAppend, "append-port", "", true, "append port to baseurl")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Hugo runs it's own a webserver to render the files",
	Long: `Hugo is able to run it's own high performance web server.
Hugo will render all the files defined in the source directory and
Serve them up.`,
	Run: server,
}

func server(cmd *cobra.Command, args []string) {
	InitializeConfig()

	if BaseUrl == "" {
		BaseUrl = "http://localhost"
	}

	if !strings.HasPrefix(BaseUrl, "http://") {
		BaseUrl = "http://" + BaseUrl
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(serverPort))
	if err == nil {
		l.Close()
	} else {
		jww.ERROR.Println("port", serverPort, "already in use, attempting to use an available port")
		sp, err := helpers.FindAvailablePort()
		if err != nil {
			jww.ERROR.Println("Unable to find alternative port to use")
			jww.ERROR.Fatalln(err)
		}
		serverPort = sp.Port
	}

	if serverAppend {
		viper.Set("BaseUrl", strings.TrimSuffix(BaseUrl, "/")+":"+strconv.Itoa(serverPort))
	} else {
		viper.Set("BaseUrl", strings.TrimSuffix(BaseUrl, "/"))
	}

	build(serverWatch)

	// Watch runs its own server as part of the routine
	if serverWatch {
		jww.FEEDBACK.Println("Watching for changes in", helpers.AbsPathify(viper.GetString("ContentDir")))
		err := NewWatcher(serverPort)
		if err != nil {
			fmt.Println(err)
		}
	}

	serve(serverPort)
}

func serve(port int) {
	jww.FEEDBACK.Println("Serving pages from " + helpers.AbsPathify(viper.GetString("PublishDir")))

	jww.FEEDBACK.Printf("Web Server is available at %s\n", viper.GetString("BaseUrl"))

	fmt.Println("Press ctrl+c to stop")

	http.Handle("/", http.FileServer(http.Dir(helpers.AbsPathify(viper.GetString("PublishDir")))))
	go wsHub.run()
	http.HandleFunc("/livereload", wsHandler)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		jww.ERROR.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var wsHub = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println(string(message))
		switch true {
		case bytes.Contains(message, []byte(`"command":"hello"`)):
			wsHub.broadcast <- []byte(`{
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

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
