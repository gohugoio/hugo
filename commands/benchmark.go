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
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	benchmarkTimes int
	cpuProfileFile string
	memProfileFile string
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark Hugo by building a site a number of times.",
	Long: `Hugo can build a site many times over and analyze the running process
creating a benchmark.`,
}

func init() {
	initHugoBuilderFlags(benchmarkCmd)
	initBenchmarkBuildingFlags(benchmarkCmd)

	benchmarkCmd.Flags().StringVar(&cpuProfileFile, "cpuprofile", "", "path/filename for the CPU profile file")
	benchmarkCmd.Flags().StringVar(&memProfileFile, "memprofile", "", "path/filename for the memory profile file")
	benchmarkCmd.Flags().IntVarP(&benchmarkTimes, "count", "n", 13, "number of times to build the site")

	benchmarkCmd.RunE = benchmark
}

func benchmark(cmd *cobra.Command, args []string) error {
	cfg, err := InitializeConfig(benchmarkCmd)
	if err != nil {
		return err
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	var memProf *os.File
	if memProfileFile != "" {
		memProf, err = os.Create(memProfileFile)
		if err != nil {
			return err
		}
	}

	var cpuProf *os.File
	if cpuProfileFile != "" {
		cpuProf, err = os.Create(cpuProfileFile)
		if err != nil {
			return err
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memAllocated := memStats.TotalAlloc
	mallocs := memStats.Mallocs
	if cpuProf != nil {
		pprof.StartCPUProfile(cpuProf)
	}

	t := time.Now()
	for i := 0; i < benchmarkTimes; i++ {
		if err = c.resetAndBuildSites(false); err != nil {
			return err
		}
	}
	totalTime := time.Since(t)

	if memProf != nil {
		pprof.WriteHeapProfile(memProf)
		memProf.Close()
	}
	if cpuProf != nil {
		pprof.StopCPUProfile()
		cpuProf.Close()
	}

	runtime.ReadMemStats(&memStats)
	totalMemAllocated := memStats.TotalAlloc - memAllocated
	totalMallocs := memStats.Mallocs - mallocs

	jww.FEEDBACK.Println()
	jww.FEEDBACK.Printf("Average time per operation: %vms\n", int(1000*totalTime.Seconds()/float64(benchmarkTimes)))
	jww.FEEDBACK.Printf("Average memory allocated per operation: %vkB\n", totalMemAllocated/uint64(benchmarkTimes)/1024)
	jww.FEEDBACK.Printf("Average allocations per operation: %v\n", totalMallocs/uint64(benchmarkTimes))

	return nil
}
