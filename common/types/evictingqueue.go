// Copyright 2017-present The Hugo Authors. All rights reserved.
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

// Package types contains types shared between packages in Hugo.
package types

import (
	"sync"
)

// LIFO 队列，溢出的元素会从顶部移除
// 没有主动删除元素的方法
// EvictingStringQueue is a queue which automatically evicts elements from the head of
// the queue when attempting to add new elements onto the queue and it is full.
// This queue orders elements LIFO (last-in-first-out). It throws away duplicates.
// Note: This queue currently does not contain any remove (poll etc.) methods.
type EvictingStringQueue struct {
	size int
	vals []string // 储存真实的数据
	set  map[string]bool // 表示是否已经存在
	mu   sync.Mutex
}

// NewEvictingStringQueue creates a new queue with the given size.
func NewEvictingStringQueue(size int) *EvictingStringQueue {
	return &EvictingStringQueue{size: size, set: make(map[string]bool)}
}

// Add adds a new string to the tail of the queue if it's not already there.
func (q *EvictingStringQueue) Add(v string) {
	q.mu.Lock()
	// 已经存在
	if q.set[v] {
		q.mu.Unlock()
		return
	}

	// 数量达到最大限制
	if len(q.set) == q.size {
		// Full
		// 移除了 0 号元素的占位符
		delete(q.set, q.vals[0])
		// :0 取空数组，1:取不包含第一个元素的其余元素
		// 移除了数组 0 号元素
		q.vals = append(q.vals[:0], q.vals[1:]...)
	}
	// 表示存在
	q.set[v] = true
	// 最新插入的值在数组最后
	// 是队列结构
	q.vals = append(q.vals, v)
	q.mu.Unlock()
}

// Contains returns whether the queue contains v.
func (q *EvictingStringQueue) Contains(v string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.set[v]
}

// Peek looks at the last element added to the queue.
func (q *EvictingStringQueue) Peek() string {
	q.mu.Lock()
	l := len(q.vals)
	// 处理边界条件
	if l == 0 {
		q.mu.Unlock()
		return ""
	}
	// 取最后一个元素
	elem := q.vals[l-1]
	q.mu.Unlock()
	return elem
}

// PeekAll looks at all the elements in the queue, with the newest first.
func (q *EvictingStringQueue) PeekAll() []string {
	q.mu.Lock()
	vals := make([]string, len(q.vals))
	copy(vals, q.vals)
	q.mu.Unlock()
	// i 从头开始循环 j 从尾循环
	// 交换 i j 元素位置
	// 数组 reverse
	// 最后插入的在最前面
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		vals[i], vals[j] = vals[j], vals[i]
	}
	return vals
}

// PeekAllSet returns PeekAll as a set.
func (q *EvictingStringQueue) PeekAllSet() map[string]bool {
	all := q.PeekAll()
	set := make(map[string]bool)
	for _, v := range all {
		set[v] = true
	}

	return set
}
