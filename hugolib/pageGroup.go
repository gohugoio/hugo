// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"time"
)

type PageGroup struct {
	Key   interface{}
	Pages Pages
}

type mapKeyValues []reflect.Value

func (v mapKeyValues) Len() int      { return len(v) }
func (v mapKeyValues) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

type mapKeyByInt struct{ mapKeyValues }

func (s mapKeyByInt) Less(i, j int) bool { return s.mapKeyValues[i].Int() < s.mapKeyValues[j].Int() }

type mapKeyByStr struct{ mapKeyValues }

func (s mapKeyByStr) Less(i, j int) bool {
	return s.mapKeyValues[i].String() < s.mapKeyValues[j].String()
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

type PagesGroup []PageGroup

func (p PagesGroup) Reverse() PagesGroup {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}

	return p
}

var (
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
	pagePtrType = reflect.TypeOf((*Page)(nil))
)

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
		if m.Type.NumIn() != 1 || m.Type.NumOut() == 0 || m.Type.NumOut() > 2 {
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
		tmp = reflect.MakeMap(reflect.MapOf(e.Type, reflect.SliceOf(pagePtrType)))
	case reflect.Method:
		tmp = reflect.MakeMap(reflect.MapOf(e.Type.Out(0), reflect.SliceOf(pagePtrType)))
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
			tmp.SetMapIndex(fv, reflect.MakeSlice(reflect.SliceOf(pagePtrType), 0, 0))
		}
		tmp.SetMapIndex(fv, reflect.Append(tmp.MapIndex(fv), ppv))
	}

	var r []PageGroup
	for _, k := range sortKeys(tmp.MapKeys(), direction) {
		r = append(r, PageGroup{Key: k.Interface(), Pages: tmp.MapIndex(k).Interface().([]*Page)})
	}

	return r, nil
}

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
		param := e.GetParam(key)
		if param != nil {
			if _, ok := param.([]string); !ok {
				keyt = reflect.TypeOf(param)
				tmp = reflect.MakeMap(reflect.MapOf(keyt, reflect.SliceOf(pagePtrType)))
				break
			}
		}
	}
	if !tmp.IsValid() {
		return nil, errors.New("There is no such a param")
	}

	for _, e := range p {
		param := e.GetParam(key)
		if param == nil || reflect.TypeOf(param) != keyt {
			continue
		}
		v := reflect.ValueOf(param)
		if !tmp.MapIndex(v).IsValid() {
			tmp.SetMapIndex(v, reflect.MakeSlice(reflect.SliceOf(pagePtrType), 0, 0))
		}
		tmp.SetMapIndex(v, reflect.Append(tmp.MapIndex(v), reflect.ValueOf(e)))
	}

	var r []PageGroup
	for _, k := range sortKeys(tmp.MapKeys(), direction) {
		r = append(r, PageGroup{Key: k.Interface(), Pages: tmp.MapIndex(k).Interface().([]*Page)})
	}

	return r, nil
}

func (p Pages) groupByDateField(sorter func(p Pages) Pages, formatter func(p *Page) string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	sp := sorter(p)

	if !(len(order) > 0 && (strings.ToLower(order[0]) == "asc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse")) {
		sp = sp.Reverse()
	}

	date := formatter(sp[0])
	var r []PageGroup
	r = append(r, PageGroup{Key: date, Pages: make(Pages, 0)})
	r[0].Pages = append(r[0].Pages, sp[0])

	i := 0
	for _, e := range sp[1:] {
		date = formatter(e)
		if r[i].Key.(string) != date {
			r = append(r, PageGroup{Key: date})
			i++
		}
		r[i].Pages = append(r[i].Pages, e)
	}
	return r, nil
}

func (p Pages) GroupByDate(format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		return p.ByDate()
	}
	formatter := func(p *Page) string {
		return p.Date.Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

func (p Pages) GroupByPublishDate(format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		return p.ByPublishDate()
	}
	formatter := func(p *Page) string {
		return p.PublishDate.Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}

func (p Pages) GroupByParamDate(key string, format string, order ...string) (PagesGroup, error) {
	sorter := func(p Pages) Pages {
		var r Pages
		for _, e := range p {
			param := e.GetParam(key)
			if param != nil {
				if _, ok := param.(time.Time); ok {
					r = append(r, e)
				}
			}
		}
		pdate := func(p1, p2 *Page) bool {
			return p1.GetParam(key).(time.Time).Unix() < p2.GetParam(key).(time.Time).Unix()
		}
		PageBy(pdate).Sort(r)
		return r
	}
	formatter := func(p *Page) string {
		return p.GetParam(key).(time.Time).Format(format)
	}
	return p.groupByDateField(sorter, formatter, order...)
}
