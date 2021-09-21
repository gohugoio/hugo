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
	"strings"
)

// GetNumWorkerMultiplier returns the base value used to calculate the number
// of workers to use for Hugo's parallel execution.
// It returns the value in HUGO_NUMWORKERMULTIPLIER OS env variable if set to a
// positive integer, else the number of logical CPUs.
func GetNumWorkerMultiplier() int {
	return 1
}

// SetEnvVars sets vars on the form key=value in the oldVars slice.
func SetEnvVars(oldVars *[]string, keyValues ...string) {
	for i := 0; i < len(keyValues); i += 2 {
		setEnvVar(oldVars, keyValues[i], keyValues[i+1])
	}
}

func SplitEnvVar(v string) (string, string) {
	parts := strings.SplitN(v, "=", 2)
	return parts[0], parts[1]
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
