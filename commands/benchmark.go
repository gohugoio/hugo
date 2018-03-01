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

type benchmarkCmd struct {
	benchmarkTimes int
	cpuProfileFile string
	memProfileFile string

	*baseBuilderCmd
}

func (b *commandsBuilder) newBenchmarkCmd() *benchmarkCmd {
	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "Benchmark Hugo by building a site a number of times.",
		Long: `Hugo can build a site many times over and analyze the running process
creating a benchmark.`,
	}

	c := &benchmarkCmd{baseBuilderCmd: b.newBuilderCmd(cmd)}

	cmd.Flags().StringVar(&c.cpuProfileFile, "cpuprofile", "", "path/filename for the CPU profile file")
	cmd.Flags().StringVar(&c.memProfileFile, "memprofile", "", "path/filename for the memory profile file")
	cmd.Flags().IntVarP(&c.benchmarkTimes, "count", "n", 13, "number of times to build the site")
	cmd.Flags().Bool("renderToMemory", false, "render to memory (only useful for benchmark testing)")

	cmd.RunE = c.benchmark

	return c
}

func (c *benchmarkCmd) benchmark(cmd *cobra.Command, args []string) error {
	cfgInit := func(c *commandeer) error {
		return nil
	}

	comm, err := initializeConfig(true, false, &c.hugoBuilderCommon, c, cfgInit)
	if err != nil {
		return err
	}

	var memProf *os.File
	if c.memProfileFile != "" {
		memProf, err = os.Create(c.memProfileFile)
		if err != nil {
			return err
		}
	}

	var cpuProf *os.File
	if c.cpuProfileFile != "" {
		cpuProf, err = os.Create(c.cpuProfileFile)
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
	for i := 0; i < c.benchmarkTimes; i++ {
		if err = comm.resetAndBuildSites(); err != nil {
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
	jww.FEEDBACK.Printf("Average time per operation: %vms\n", int(1000*totalTime.Seconds()/float64(c.benchmarkTimes)))
	jww.FEEDBACK.Printf("Average memory allocated per operation: %vkB\n", totalMemAllocated/uint64(c.benchmarkTimes)/1024)
	jww.FEEDBACK.Printf("Average allocations per operation: %v\n", totalMallocs/uint64(c.benchmarkTimes))

	return nil
}
