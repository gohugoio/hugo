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

func (p Pages) GroupBy(key string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	direction := "asc"

	if len(order) > 0 && (strings.ToLower(order[0]) == "desc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse") {
		direction = "desc"
	}

	ppt := reflect.TypeOf(&Page{}) // *hugolib.Page

	ft, ok := ppt.Elem().FieldByName(key)

	if !ok {
		return nil, errors.New("No such field in Page struct")
	}

	tmp := reflect.MakeMap(reflect.MapOf(ft.Type, reflect.SliceOf(ppt)))

	for _, e := range p {
		ppv := reflect.ValueOf(e)
		fv := ppv.Elem().FieldByName(key)
		if !fv.IsNil() {
			if !tmp.MapIndex(fv).IsValid() {
				tmp.SetMapIndex(fv, reflect.MakeSlice(reflect.SliceOf(ppt), 0, 0))
			}
			tmp.SetMapIndex(fv, reflect.Append(tmp.MapIndex(fv), ppv))
		}
	}

	var r []PageGroup
	for _, k := range sortKeys(tmp.MapKeys(), direction) {
		r = append(r, PageGroup{Key: k.Interface(), Pages: tmp.MapIndex(k).Interface().([]*Page)})
	}

	return r, nil
}

func (p Pages) GroupByDate(format string, order ...string) (PagesGroup, error) {
	if len(p) < 1 {
		return nil, nil
	}

	sp := p.ByDate()

	if !(len(order) > 0 && (strings.ToLower(order[0]) == "asc" || strings.ToLower(order[0]) == "rev" || strings.ToLower(order[0]) == "reverse")) {
		sp = sp.Reverse()
	}

	date := sp[0].Date.Format(format)
	var r []PageGroup
	r = append(r, PageGroup{Key: date, Pages: make(Pages, 0)})
	r[0].Pages = append(r[0].Pages, sp[0])

	i := 0
	for _, e := range sp[1:] {
		date = e.Date.Format(format)
		if r[i].Key.(string) != date {
			r = append(r, PageGroup{Key: date})
			i++
		}
		r[i].Pages = append(r[i].Pages, e)
	}
	return r, nil
}
