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

package pcache

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeConverters(t *testing.T) {
	assert := require.New(t)

	s := "1/500"
	to := reflect.TypeOf(big.NewRat(1, 2))
	v, b, err := defaultTypeConverters.Convert(s, to)
	assert.NoError(err)
	assert.True(b)
	assert.Equal(big.NewRat(1, 500), v)

	// This does not exist
	to = reflect.TypeOf(t)
	v, b, err = defaultTypeConverters.Convert(s, to)
	assert.NoError(err)
	assert.False(b)
	assert.Equal(s, v)
}
