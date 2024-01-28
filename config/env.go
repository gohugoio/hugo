// Copyright 2019 The Hugo Authors. All rights reserved.
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

package config

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/pbnjay/memory"
)

const (
	gigabyte = 1 << 30
)

// GetNumWorkerMultiplier returns the base value used to calculate the number
// of workers to use for Hugo's parallel execution.
// It returns the value in HUGO_NUMWORKERMULTIPLIER OS env variable if set to a
// positive integer, else the number of logical CPUs.
func GetNumWorkerMultiplier() int {
	if gmp := os.Getenv("HUGO_NUMWORKERMULTIPLIER"); gmp != "" {
		if p, err := strconv.Atoi(gmp); err == nil && p > 0 {
			return p
		}
	}
	return runtime.NumCPU()
}

// GetMemoryLimit returns the upper memory limit in bytes for Hugo's in-memory caches.
// Note that this does not represent "all of the memory" that Hugo will use,
// so it needs to be set to a lower number than the available system memory.
// It will read from the HUGO_MEMORYLIMIT (in Gigabytes) environment variable.
// If that is not set, it will set aside a quarter of the total system memory.
func GetMemoryLimit() uint64 {
	if mem := os.Getenv("HUGO_MEMORYLIMIT"); mem != "" {
		if v := stringToGibabyte(mem); v > 0 {
			return v
		}
	}

	// There is a FreeMemory function, but as the kernel in most situations
	// will take whatever memory that is left and use for caching etc.,
	// that value is not something that we can use.
	m := memory.TotalMemory()
	if m != 0 {
		return uint64(m / 4)
	}

	return 2 * gigabyte
}

func stringToGibabyte(f string) uint64 {
	if v, err := strconv.ParseFloat(f, 32); err == nil && v > 0 {
		return uint64(v * gigabyte)
	}
	return 0
}

// SetEnvVars sets vars on the form key=value in the oldVars slice.
func SetEnvVars(oldVars *[]string, keyValues ...string) {
	for i := 0; i < len(keyValues); i += 2 {
		setEnvVar(oldVars, keyValues[i], keyValues[i+1])
	}
}

func SplitEnvVar(v string) (string, string) {
	name, value, _ := strings.Cut(v, "=")
	return name, value
}

func setEnvVar(vars *[]string, key, value string) {
	for i := range *vars {
		if strings.HasPrefix((*vars)[i], key+"=") {
			(*vars)[i] = key + "=" + value
			return
		}
	}
	// New var.
	*vars = append(*vars, key+"="+value)
}
