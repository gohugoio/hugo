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

package hugolib

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPageBundler(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	var (
		cfg, fs = newTestCfg()
		//th      = testHelper{cfg, fs, t}
	)

	writeSource(t, fs, filepath.Join("content", "aaa.md"), `---
title: Home Page
---`)

	writeSource(t, fs, filepath.Join("content", "a", "page.md"), `---
title: A Page
---`)
	writeSource(t, fs, filepath.Join("content", "b", "page.md"), `---
title: A Page
---`)

	writeSource(t, fs, filepath.Join("content", "bundle", "page.md"), `---
title: A Page
---`)

	writeSource(t, fs, filepath.Join("content", "bundle", "index.md"), `---
title: A Page
---`)

	writeSource(t, fs, filepath.Join("content", "bundle", "images", "cartoon.jpg"), `THIS IS A JPG`)

	b := newBundler(newTestPathSpec(fs, cfg))

	assert.Nil(b.captureFiles())

	fmt.Println(">>>", b.bundles)

	assert.Nil(b.processFiles())

}
