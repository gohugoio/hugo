// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNodeSimpleMethods(t *testing.T) {
	for i, this := range []struct {
		assertFunc func(n *Node) bool
	}{
		{func(n *Node) bool { return n.IsNode() }},
		{func(n *Node) bool { return !n.IsPage() }},
		{func(n *Node) bool { return n.RSSlink() == "rssLink" }},
		{func(n *Node) bool { return n.Scratch() != nil }},
		{func(n *Node) bool { return n.Hugo() != nil }},
		{func(n *Node) bool { return n.Now().Unix() == time.Now().Unix() }},
	} {

		n := &Node{}
		n.RSSLink = "rssLink"

		if !this.assertFunc(n) {
			t.Errorf("[%d] Node method error", i)
		}
	}
}

func TestNodeID(t *testing.T) {
	t.Parallel()

	n1 := &Node{}
	n2 := &Node{}

	assert.True(t, n1.ID() > 0)
	assert.Equal(t, n1.ID(), n1.ID())
	assert.True(t, n2.ID() > n1.ID())

	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(j int) {
			for k := 0; k < 10; k++ {
				n := &Node{}
				assert.True(t, n.ID() > 0)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
