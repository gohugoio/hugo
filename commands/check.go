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
    "github.com/spf13/cobra"
    "github.com/spf13/hugo/hugolib"
    "syscall"
)

func init() {
    check.AddCommand(limit)
}

var check = &cobra.Command{
    Use:   "check",
    Short: "Check content in the source directory",
    Long: `Hugo will perform some basic analysis on the
    content provided and will give feedback.`,
    Run: func(cmd *cobra.Command, args []string) {
        InitializeConfig()
        site := hugolib.Site{Config: *Config}
        site.Analyze()
    },
}

var limit = &cobra.Command{
    Use:   "ulimit",
    Short: "Check system ulimit settings",
    Long: `Hugo will inspect the current ulimit settings on the system.
    This is primarily to ensure that Hugo can watch enough files on some OSs`,
    Run: func(cmd *cobra.Command, args []string) {
        var rLimit syscall.Rlimit
        err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
        if err != nil {
            fmt.Println("Error Getting Rlimit ", err)
        }
        fmt.Println("Current rLimit:", rLimit)

        fmt.Println("Attempting to increase limit")
        rLimit.Max = 999999
        rLimit.Cur = 999999
        err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
        if err != nil {
            fmt.Println("Error Setting rLimit ", err)
        }
        err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
        if err != nil {
            fmt.Println("Error Getting rLimit ", err)
        }
        fmt.Println("rLimit after change:", rLimit)
    },
}
