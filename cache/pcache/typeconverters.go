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
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var defaultTypeConverters typeConverters = make(typeConverters)

func init() {
	defaultTypeConverters[typeConverterKey{reflect.TypeOf(""), reflect.TypeOf(big.NewRat(1, 2))}] = new(stringToRatConverter)
	defaultTypeConverters[typeConverterKey{reflect.TypeOf(""), reflect.TypeOf(time.Now())}] = new(stringToTimeConverter)

}

type typeConverter interface {
	Convert(v interface{}) (interface{}, error)
}

type typeConverters map[typeConverterKey]typeConverter

func (t typeConverters) ConvertTypes(v interface{}, from, to reflect.Type) (interface{}, bool, error) {
	converter, found := t[typeConverterKey{from, to}]
	if !found {
		return v, false, nil
	}

	c, err := converter.Convert(v)

	return c, true, err
}

func (t typeConverters) Convert(v interface{}, to reflect.Type) (interface{}, bool, error) {
	return t.ConvertTypes(v, reflect.TypeOf(v), to)
}

type typeConverterKey struct {
	From reflect.Type
	To   reflect.Type
}

type stringToRatConverter int

func (s stringToRatConverter) Convert(v interface{}) (interface{}, error) {
	vs := v.(string)
	parts := strings.Split(vs, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("does not look like a Rat: %v", v)
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return big.NewRat(int64(a), int64(b)), nil

}

type stringToTimeConverter int

func (s stringToTimeConverter) Convert(v interface{}) (interface{}, error) {
	vs := v.(string)

	t, err := time.Parse(time.RFC3339, vs)

	return t, err
}
