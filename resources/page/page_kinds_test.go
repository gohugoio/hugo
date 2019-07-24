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

package page

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKind(t *testing.T) {
	t.Parallel()
	// Add tests for these constants to make sure they don't change
	require.Equal(t, "page", KindPage)
	require.Equal(t, "home", KindHome)
	require.Equal(t, "section", KindSection)
	require.Equal(t, "taxonomy", KindTaxonomy)
	require.Equal(t, "taxonomyTerm", KindTaxonomyTerm)
}
