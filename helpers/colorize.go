// Copyright © 2014 Thibault Normand <me@zenithar.org>.
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

package helpers

// Extracted from strings.go

// primeRK is the prime base used in Rabin-Karp algorithm.
const primeRK = 16777619
func hashstr(sep string) (uint32, uint32) {
	hash := uint32(0)
    for i := 0; i < len(sep); i++ {
    	hash = hash*primeRK + uint32(sep[i])
    }
    
    var pow, sq uint32 = 1, primeRK
    for i := len(sep); i > 0; i >>= 1 {
    	if i&1 != 0 {
    		pow *= sq
    	}
    	sq *= sq
    }
    return hash, pow
}
// /-- Extract

func flatcolors() [16]string {
	return [...]string{
		"#1abc9c", "#16a085", "#2ecc71", "#27ae60",
		"#3498db", "#2980b9", "#9b59b6", "#8e44ad",
		"#f1c40f", "#f39c12", "#e67e22", "#d35400",
		"#e74c3c", "#c0392b", "#95a5a6", "#7f8c8d",
	}
}

func Colorize16(a string) string {
	hash, _ := hashstr(a)
	return flatcolors()[hash % 16]
}