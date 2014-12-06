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
	"os"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

var cpuProfilefile string
var benchmarkTimes int

var benchmark = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark hugo by building a site a number of times",
	Long: `Hugo can build a site many times over and anlyze the
    running process creating a `,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig(cmdBenchmark)
		bench(cmd, args)
	},
}

func init() {
	benchmark.Flags().StringVar(&cpuProfilefile, "outputfile", "/tmp/hugo-cpuprofile", "path/filename for the profile file")
	benchmark.Flags().IntVarP(&benchmarkTimes, "count", "n", 13, "number of times to build the site")
}

func bench(cmd *cobra.Command, args []string) {
	f, err := os.Create(cpuProfilefile)

	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < benchmarkTimes; i++ {
		_ = buildSite()
	}
}
