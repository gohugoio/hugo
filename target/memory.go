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

package target

import (
	"bytes"
	"io"
)

type InMemoryTarget struct {
	Files map[string][]byte
}

func (t *InMemoryTarget) Publish(label string, reader io.Reader) (err error) {
	if t.Files == nil {
		t.Files = make(map[string][]byte)
	}
	bytes := new(bytes.Buffer)
	bytes.ReadFrom(reader)
	t.Files[label] = bytes.Bytes()
	return
}

func (t *InMemoryTarget) Translate(label string) (dest string, err error) {
	return label, nil
}
