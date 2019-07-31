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
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/compare"

	"github.com/gohugoio/hugo/resources/resource"
)

var (
	_ collections.Slicer   = PageGroup{}
	_ compare.ProbablyEqer = PageGroup{}
	_ compare.ProbablyEqer = PagesGroup{}
)

// PageGroup represents a group of pages, grouped by the key.
// The key is typically a year or similar.
type PageGroup struct {
	Key interface{}
	Pages
}

type mapKeyValues []reflect.Value

func (v mapKeyValues) Len() int      { return len(v) }
func (v mapKeyValues) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

type mapKeyByInt struct{ mapKeyValues }

func (s mapKeyByInt) Less(i, j int) bool { return s.mapKeyValues[i].Int() < s.mapKeyValues[j].Int() }

type mapKeyByStr struct{ mapKeyValues }

func (s mapKeyByStr) Less(i, j int) bool {
	return compare.LessStrings(s.mapKeyValues[i].String(), s.mapKeyValues[j].String())
}

func sortKeys(v []reflect.Value, order string) []reflect.Value {
	if len(v) <= 1 {
		return v
	}

	switch v[0].Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if order == "desc" {
			sort.Sort(sort.Reverse(mapKeyByInt{v}))
		} else {
			sort.Sort(mapKeyByInt{v})
		}
	case reflect.String:
		if order == "desc" {
			sort.Sort(sort.Reverse(mapKeyByStr{v}))
		} else {
			sort.Sort(mapKeyByStr{v})
		}
	}
	return v
}

// PagesGroup represents a list of page groups.
// This is what you get when doing page grouping in the templates.
type PagesGroup []PageGroup

// Reverse reverses the order of this list of page groups.
func (p PagesGroup) Reverse() PagesGroup {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}

	return p
}

var (
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	pagePtrType = reflect.TypeOf((*Page)(nil)).Elem()
	pagesType   = reflect.TypeOf(Pages{})
)

// GroupBy groups by the value in the given field or method name and with the given order.
// Valid values for order is asc, desc, rev and reverse.
func (p Pages) GroupBy(key string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	direction := "asc"

	if len(order) > 0 && (strings.ToLower(order[0]) == "desc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse") {
		direction = "desc"
	}

	var ft interface{}
	m, ok := pagePtrType.MethodByName(key)
	if ok {
		if m.Type.NumOut() == 0 || m.Type.NumOut() > 2 {
			return nil, errors.New(key + " is a Page method but you can't use it with GroupBy")
		}
		if m.Type.NumOut() == 1 && m.Type.Out(0).Implements(errorType) {
			return nil, errors.New(key + " is a Page method but you can't use it with GroupBy")
		}
		if m.Type.NumOut() == 2 && !m.Type.Out(1).Implements(errorType) {
			return nil, errors.New(key + " is a Page method but you can't use it with GroupBy")
		}
		ft = m
	} else {
		ft, ok = pagePtrType.Elem().FieldByName(key)
		if !ok {
			return nil, errors.New(key + " is neither a field nor a method of Page")
		}
	}

	var tmp reflect.Value
	switch e := ft.(type) {
	case reflect.StructField:
		tmp = reflect.MakeMap(reflect.MapOf(e.Type, pagesType))
	case reflect.Method:
		tmp = reflect.MakeMap(reflect.MapOf(e.Type.Out(0), pagesType))
	}

	for _, e := range p {
		ppv := reflect.ValueOf(e)
		var fv reflect.Value
		switch ft.(type) {
		case reflect.StructField:
			fv = ppv.Elem().FieldByName(key)
		case reflect.Method:
			fv = ppv.MethodByName(key).Call([]reflect.Value{})[0]
		}
		if !fv.IsValid() {
			continue
		}
		if !tmp.MapIndex(fv).IsValid() {
			tmp.SetMapIndex(fv, reflect.MakeSlice(pagesType, 0, 0))
		}
		tmp.SetMapIndex(fv, reflect.Append(tmp.MapIndex(fv), ppv))
	}

	sortedKeys := sortKeys(tmp.MapKeys(), direction)
	r := make([]PageGroup, len(sortedKeys))
	for i, k := range sortedKeys {
		r[i] = PageGroup{Key: k.Interface(), Pages: tmp.MapIndex(k).Interface().(Pages)}
	}

	return r, nil
}

// GroupByParam groups by the given page parameter key's value and with the given order.
// Valid values for order is asc, desc, rev and reverse.
func (p Pages) GroupByParam(key string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	direction := "asc"

	if len(order) > 0 && (strings.ToLower(order[0]) == "desc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse") {
		direction = "desc"
	}

	var tmp reflect.Value
	var keyt reflect.Type
	for _, e := range p {
		param := resource.GetParamToLower(e, key)
		if param != nil {
			if _, ok := param.([]string); !ok {
				keyt = reflect.TypeOf(param)
				tmp = reflect.MakeMap(reflect.MapOf(keyt, pagesType))
				break
			}
		}
	}
	if !tmp.IsValid() {
		return nil, errors.New("there is no such a param")
	}

	for _, e := range p {
		param := resource.GetParam(e, key)

		if param == nil || reflect.TypeOf(param) != keyt {
			continue
		}
		v := reflect.ValueOf(param)
		if !tmp.MapIndex(v).IsValid() {
			tmp.SetMapIndex(v, reflect.MakeSlice(pagesType, 0, 0))
		}
		tmp.SetMapIndex(v, reflect.Append(tmp.MapIndex(v), reflect.ValueOf(e)))
	}

	var r []PageGroup
	for _, k := range sortKeys(tmp.MapKeys(), direction) {
		r = append(r, PageGroup{Key: k.Interface(), Pages: tmp.MapIndex(k).Interface().(Pages)})
	}

	return r, nil
}

func (p Pages) groupByDateField(sorter func(p Pages) Pages, formatter func(p Page) string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	sp := sorter(p)

	if !(len(order) > 0 && (strings.ToLower(order[0]) == "asc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse")) {
		sp = sp.Reverse()
	}

	date := formatter(sp[0].(Page))
	var r []PageGroup
	r = append(r, PageGroup{Key: date, Pages: make(Pages, 0)})
	r[0].Pages = append(r[0].Pages, sp[0])

	i := 0
	for _, e := range sp[1:] {
		date = formatter(e.(Page))
		if r[i].Key.(string) != date {
			r = append(r, PageGroup{Key: date})
			i++
		}
		r[i].Pages = append(r[i].Pages, e)
	}
	return r, nil
}

// GroupByDate groups by the given page's Date value in
// the given format and with the given order.
// Valid values for order is asc, desc, rev and reverse.
// For valid format strings, see https://golang.org/pkg/time/#Time.Format
func (p Pages) GroupByDate(format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		return p.ByDate()
	}
	formatter := func(p Page) string {
		return p.Date().Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

// GroupByPublishDate groups by the given page's PublishDate value in
// the given format and with the given order.
// Valid values for order is asc, desc, rev and reverse.
// For valid format strings, see https://golang.org/pkg/time/#Time.Format
func (p Pages) GroupByPublishDate(format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		return p.ByPublishDate()
	}
	formatter := func(p Page) string {
		return p.PublishDate().Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

// GroupByExpiryDate groups by the given page's ExpireDate value in
// the given format and with the given order.
// Valid values for order is asc, desc, rev and reverse.
// For valid format strings, see https://golang.org/pkg/time/#Time.Format
func (p Pages) GroupByExpiryDate(format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		return p.ByExpiryDate()
	}
	formatter := func(p Page) string {
		return p.ExpiryDate().Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

// GroupByParamDate groups by a date set as a param on the page in
// the given format and with the given order.
// Valid values for order is asc, desc, rev and reverse.
// For valid format strings, see https://golang.org/pkg/time/#Time.Format
func (p Pages) GroupByParamDate(key string, format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		var r Pages
		for _, e := range p {
			param := resource.GetParamToLower(e, key)
			if _, ok := param.(time.Time); ok {
				r = append(r, e)
			}
		}
		pdate := func(p1, p2 Page) bool {
			p1p, p2p := p1.(Page), p2.(Page)
			return resource.GetParamToLower(p1p, key).(time.Time).Unix() < resource.GetParamToLower(p2p, key).(time.Time).Unix()
		}
		pageBy(pdate).Sort(r)
		return r
	}
	formatter := func(p Page) string {
		return resource.GetParamToLower(p, key).(time.Time).Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

// ProbablyEq wraps comare.ProbablyEqer
func (p PageGroup) ProbablyEq(other interface{}) bool {
	otherP, ok := other.(PageGroup)
	if !ok {
		return false
	}

	if p.Key != otherP.Key {
		return false
	}

	return p.Pages.ProbablyEq(otherP.Pages)

}

// Slice is not meant to be used externally. It's a bridge function
// for the template functions. See collections.Slice.
func (p PageGroup) Slice(in interface{}) (interface{}, error) {
	switch items := in.(type) {
	case PageGroup:
		return items, nil
	case []interface{}:
		groups := make(PagesGroup, len(items))
		for i, v := range items {
			g, ok := v.(PageGroup)
			if !ok {
				return nil, fmt.Errorf("type %T is not a PageGroup", v)
			}
			groups[i] = g
		}
		return groups, nil
	default:
		return nil, fmt.Errorf("invalid slice type %T", items)
	}
}

// Len returns the number of pages in the page group.
func (psg PagesGroup) Len() int {
	l := 0
	for _, pg := range psg {
		l += len(pg.Pages)
	}
	return l
}

// ProbablyEq wraps comare.ProbablyEqer
func (psg PagesGroup) ProbablyEq(other interface{}) bool {
	otherPsg, ok := other.(PagesGroup)
	if !ok {
		return false
	}

	if len(psg) != len(otherPsg) {
		return false
	}

	for i := range psg {
		if !psg[i].ProbablyEq(otherPsg[i]) {
			return false
		}
	}

	return true

}

// ToPagesGroup tries to convert seq into a PagesGroup.
func ToPagesGroup(seq interface{}) (PagesGroup, error) {
	switch v := seq.(type) {
	case nil:
		return nil, nil
	case PagesGroup:
		return v, nil
	case []PageGroup:
		return PagesGroup(v), nil
	case []interface{}:
		l := len(v)
		if l == 0 {
			break
		}
		switch v[0].(type) {
		case PageGroup:
			pagesGroup := make(PagesGroup, l)
			for i, ipg := range v {
				if pg, ok := ipg.(PageGroup); ok {
					pagesGroup[i] = pg
				} else {
					return nil, fmt.Errorf("unsupported type in paginate from slice, got %T instead of PageGroup", ipg)
				}
			}
			return pagesGroup, nil
		}
	}

	return nil, nil
}
