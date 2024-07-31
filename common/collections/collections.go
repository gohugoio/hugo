// Copyright 2018 The Hugo Authors. All rights reserved.
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

// Package collections contains common Hugo functionality related to collection
// handling.
package collections

// Grouper defines a very generic way to group items by a given key.
type Grouper interface {
	Group(key any, items any) (any, error)
}

// Partitioner defines a very generic way to partition sets of items into
// chunks of given maximal length.
type Partitioner interface {
	Partition(n any, items any) (any, error)
}
