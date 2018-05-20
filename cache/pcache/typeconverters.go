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
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	stringToRatConverter       = new(stringToRatConverterTp)
	stringToTimeConverter      = new(stringToTimeConverterTp)
	jsonNumberToFloatConverter = new(jsonNumberToFloatConverterTp)

	strType = reflect.TypeOf("")
)

// TODO(bep) this may be general useful. But let us keep this inside this package for now.
var defaultTypeConverters typeConverters = make(typeConverters)

func init() {

	defaultTypeConverters[typeConverterKey{strType, reflect.TypeOf(big.NewRat(1, 2))}] = stringToRatConverter
	defaultTypeConverters[typeConverterKey{strType, reflect.TypeOf(time.Now())}] = stringToTimeConverter
	defaultTypeConverters[typeConverterKey{strType, reflect.TypeOf(float64(0))}] = jsonNumberToFloatConverter

	// Add the named variants.
	for _, v := range defaultTypeConverters {
		defaultTypeConverters[v.Name()] = v
	}

}

type typeConverter interface {
	Convert(v interface{}) (interface{}, error)
	Name() string
}

// The key is either a string (named converter) or a typeConverterKey.
type typeConverters map[interface{}]typeConverter

// GetByTypes returns nil if no converter could be found for the given from/to.
func (t typeConverters) GetByTypes(from, to reflect.Type) typeConverter {
	converter, _ := t[typeConverterKey{from, to}]
	return converter
}

// GetByTypes returns nil if no converter could be found for the given from/to.
func (t typeConverters) GetByName(name string) typeConverter {
	converter, _ := t[name]
	return converter
}

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

type stringToRatConverterTp int

func (s stringToRatConverterTp) Convert(v interface{}) (interface{}, error) {
	vs, ok := v.(string)
	if !ok {
		return v, nil
	}

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

func (s stringToRatConverterTp) Name() string {
	return "stringToRatConverter"
}

type stringToTimeConverterTp int

func (s stringToTimeConverterTp) Convert(v interface{}) (interface{}, error) {
	vs, ok := v.(string)
	if !ok {
		// TODO(bep)
		return v, nil
	}

	t, err := time.Parse(time.RFC3339, vs)

	return t, err
}

func (s stringToTimeConverterTp) Name() string {
	return "stringToTimeConverter"
}

type jsonNumberToFloatConverterTp int

func (s jsonNumberToFloatConverterTp) Convert(v interface{}) (interface{}, error) {
	vs, ok := v.(json.Number)
	if !ok {
		return v, nil
	}

	f, err := strconv.ParseFloat(string(vs), 64)
	if err != nil {
		return v, err
	}

	return f, nil
}

func (s jsonNumberToFloatConverterTp) Name() string {
	return "jsonNumberToFloatConverter"
}

type typeMapper struct {
	converters typeConverters
	from       reflect.Type

	// The result.
	namedConverters map[string]string
}

func newDefaultTypeMapper() *typeMapper {
	// This is the current JSON use case. We get strings and have some special
	// converters to time.Time etc.
	return &typeMapper{from: strType, converters: defaultTypeConverters, namedConverters: make(map[string]string)}
}

func (t *typeMapper) mapTypes(in interface{}) {
	t.mapRecursive("", reflect.ValueOf(in))
}

func (t *typeMapper) mapRecursive(name string, v reflect.Value) {
	switch v.Kind() {

	case reflect.Ptr:
		vv := v.Elem()
		if !vv.IsValid() {
			// Nil
			return
		}
		// A Ptr type can have a mapping (*big.Rat)
		t.checkMapping(name, v)
		t.mapRecursive(name, vv)

	case reflect.Interface:
		vv := v.Elem()
		t.mapRecursive(name, vv)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i += 1 {
			f := v.Field(i)
			if f.CanSet() {
				// Exported field.
				fieldName := v.Type().Field(i).Name
				if name != "" {
					fieldName = name + "." + fieldName
				}
				t.mapRecursive(fieldName, f)
			}
		}
		// A struct can have a type mapping (time.Time).
		t.checkMapping(name, v)

	case reflect.Slice:
		for i := 0; i < v.Len(); i += 1 {
			t.mapRecursive(name, v.Index(i))
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			fieldName := key.Interface().(string)
			if name != "" {
				// This is the format mapstructure delivers
				fieldName = fmt.Sprintf("%s[%s]", name, fieldName)
			}

			vv := v.MapIndex(key)

			t.mapRecursive(fieldName, vv)
		}

	default:
		// This should now be a non-container type with a potential type
		// mapper.
		t.checkMapping(name, v)

	}

}

func (t *typeMapper) checkMapping(name string, v reflect.Value) {
	// TODO(bep) t.from vs json.Number
	converter := t.converters.GetByTypes(t.from, v.Type())
	if converter != nil {
		t.namedConverters[name] = converter.Name()
	}
}
