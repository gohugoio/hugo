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

// Package commands defines and implements command-line commands and flags
// used by Hugo. Commands and flags are implemented using Cobra.

package releaser

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
)

func _TestReleaseNotesWriter(t *testing.T) {
	if os.Getenv("CI") != "" {
		// Travis has an ancient git with no --invert-grep: https://github.com/travis-ci/travis-ci/issues/6328
		t.Skip("Skip git test on CI to make Travis happy.")
	}
	c := qt.New(t)

	var b bytes.Buffer

	// TODO(bep) consider to query GitHub directly for the gitlog with author info, probably faster.
	infos, err := getGitInfosBefore("HEAD", "v0.20", "hugo", "", false)
	c.Assert(err, qt.IsNil)

	c.Assert(writeReleaseNotes("0.21", infos, infos, &b), qt.IsNil)

	fmt.Println(b.String())

}
