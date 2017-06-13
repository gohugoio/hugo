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

package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printHugoVersion()
		return nil
	},
}

func printHugoVersion() {
	if hugolib.BuildDate == "" {
		setBuildDate() // set the build date from executable's mdate
	} else {
		formatBuildDate() // format the compile time
	}
	if hugolib.CommitHash == "" {
		jww.FEEDBACK.Printf("Hugo Static Site Generator v%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, runtime.GOOS, runtime.GOARCH, hugolib.BuildDate)
	} else {
		jww.FEEDBACK.Printf("Hugo Static Site Generator v%s-%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, strings.ToUpper(hugolib.CommitHash), runtime.GOOS, runtime.GOARCH, hugolib.BuildDate)
	}
}

// setBuildDate checks the ModTime of the Hugo executable and returns it as a
// formatted string.  This assumes that the executable name is Hugo, if it does
// not exist, an empty string will be returned.  This is only called if the
// hugolib.BuildDate wasn't set during compile time.
//
// osext is used for cross-platform.
func setBuildDate() {
	fname, _ := osext.Executable()
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	t := fi.ModTime()
	hugolib.BuildDate = t.Format(time.RFC3339)
}

// formatBuildDate formats the hugolib.BuildDate according to the value in
// .Params.DateFormat, if it's set.
func formatBuildDate() {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", hugolib.BuildDate)
	hugolib.BuildDate = t.Format(time.RFC3339)
}
