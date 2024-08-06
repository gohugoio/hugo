// Copyright 2024 The Hugo Authors. All rights reserved.
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

package page

import "github.com/spf13/cast"

// PagesPartition is a slice of Pages objects.
type PagesPartition []Pages

func (p Pages) PartitionWith(n int) (PagesPartition, error) {
	if len(p) < 1 {
		return nil, nil
	}
	nv, err := cast.ToIntE(n)
	if err != nil {
		return nil, err
	}
	return p.partition(nv)
}

func (p Pages) partition(n int) (PagesPartition, error) {
	if n < 1 {
		return nil, nil
	}

	all := len(p)
	wholes := all / n
	remainder := all % n

	parts := wholes
	if remainder > 0 {
		parts = parts + 1
	}

	partition := make([]Pages, parts)

	for i := 0; i < wholes; i++ {
		partition[i] = make([]Page, n)
		for j := 0; j < n; j++ {
			partition[i][j] = p[i*n+j]
		}
	}

	if remainder > 0 {
		partition[wholes] = make([]Page, remainder)
		for j := 0; j < remainder; j++ {
			partition[wholes][j] = p[wholes*n+j]
		}

	}

	return partition, nil
}
