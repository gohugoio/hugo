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

package hugolib

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestNewFrontmatterConfig(t *testing.T) {
	t.Parallel()

	v := viper.New()

	v.Set("frontmatter", map[string]interface{}{
		"defaultDate": []string{"filename"},
	})

	assert := require.New(t)

	fc, err := newFrontmatterConfig(newWarningLogger(), v)

	assert.NoError(err)
	assert.Equal(2, len(fc.dateHandlers))

}
