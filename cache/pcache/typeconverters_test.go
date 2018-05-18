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

	assert.Equal(stringToRatConverter, defaultTypeConverters.GetByName(stringToRatConverter.Name()))
	assert.Equal(stringToTimeConverter, defaultTypeConverters.GetByName(stringToTimeConverter.Name()))

}

type nestedObject struct {
	TopRat      *big.Rat
	TestObjectV testObject
	TestObjectP *testObject
	MyMap       map[string]interface{}
}

func TestTypeMapper(t *testing.T) {
	assert := require.New(t)

	top := testObject{
		MyString:  "hi",
		MyRat:     big.NewRat(1, 100),
		MyInt64:   int64(64),
		MyFloat64: float64(3.14159264),
	}

	no := &nestedObject{
		TopRat:      big.NewRat(1, 100),
		TestObjectV: top,
		TestObjectP: &top,
		MyMap: map[string]interface{}{
			"MyString": "A string",
			"MyRat":    big.NewRat(1, 100),
		},
	}

	mapper := newDefaultTypeMapper()

	mapper.mapTypes(no)

	assert.Equal("stringToTimeConverter", mapper.namedConverters["TestObjectV.MyDate"])
	assert.Equal("stringToTimeConverter", mapper.namedConverters["TestObjectP.MyDate"])
	assert.Equal("stringToRatConverter", mapper.namedConverters["TopRat"])
	assert.Equal("stringToRatConverter", mapper.namedConverters["TestObjectV.MyRat"])

	assert.Equal("stringToRatConverter", mapper.namedConverters["MyMap[MyRat]"])

}
