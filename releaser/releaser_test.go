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
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/stretchr/testify/require"
)

func TestCalculateVersions(t *testing.T) {
	startVersion := helpers.HugoVersion{Number: 0.20, Suffix: "-DEV"}

	tests := []struct {
		handler *ReleaseHandler
		version helpers.HugoVersion
		v1      string
		v2      string
	}{
		{
			New(0, 0, true),
			startVersion,
			"0.20",
			"0.21-DEV",
		},
		{
			New(2, 0, true),
			startVersion,
			"0.20.2",
			"0.20-DEV",
		},
		{
			New(0, 1, true),
			startVersion,
			"0.20",
			"0.21-DEV",
		},
		{
			New(0, 2, true),
			startVersion,
			"0.20",
			"0.21-DEV",
		},
		{
			New(3, 1, true),
			startVersion,
			"0.20.3",
			"0.20-DEV",
		},
		{
			New(3, 2, true),
			startVersion.Next(),
			"0.21",
			"0.21-DEV",
		},
	}

	for _, test := range tests {
		v1, v2 := test.handler.calculateVersions(test.version)
		require.Equal(t, test.v1, v1.String(), "Release version")
		require.Equal(t, test.v2, v2.String(), "Final version")
	}
}
