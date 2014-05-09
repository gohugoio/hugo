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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	err := http.ListenAndServe(":"+strconv.Itoa(port), http.FileServer(http.Dir(helpers.AbsPathify(viper.GetString("PublishDir")))))
	if err != nil {
		jww.ERROR.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
