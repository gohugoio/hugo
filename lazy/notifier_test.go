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

package lazy

import (
	"testing"
	"time"
)

func TestNotifier(t *testing.T) {
	type foo struct {
		value   int
		created *Notifier
	}

	f := foo{
		created: NewNotifier(),
		value:   3,
	}
	f.value = 3
	go func() {
		time.Sleep(time.Duration(100) * time.Millisecond)
		f.value = 5
		f.created.Close()
	}()
	f.created.Wait()
	if f.value != 5 {
		t.Errorf("expecting a value of 5 but got %v", f.value)

	}
	f.created.Reset()
	go func() {
		time.Sleep(time.Duration(100) * time.Millisecond)
		f.value = 6
		f.created.Close()
	}()
	f.created.Wait()
	if f.value != 6 {
		t.Errorf("expecting a value of 6 but got %v", f.value)
	}
}
