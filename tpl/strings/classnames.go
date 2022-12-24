package strings

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ClassNames the method ported npm package. see https://npm.im/classnames
func (ns *Namespace) ClassNames(inputs ...any) string {
	var classes []string
	for _, input := range inputs {
		v := reflect.ValueOf(input)
		if isZeroValue(v) {
			continue
		}
		switch v.Kind() {
		case reflect.String:
			classes = append(classes, v.String())
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				if inner := ns.ClassNames(v.Index(i).Interface()); inner != "" {
					classes = append(classes, inner)
				}
			}
		case reflect.Map:
			var keys []string
			iter := v.MapRange()
			for iter.Next() {
				key := reflect.ValueOf(iter.Key().Interface())
				value := reflect.ValueOf(iter.Value().Interface())
				if key.Kind() != reflect.String || key.IsZero() {
					continue
				}
				if isZeroValue(value) {
					continue
				}
				keys = append(keys, key.String())
			}
			sort.Strings(keys)
			classes = append(classes, keys...)
		default:
		}
		switch {
		case v.CanInt():
			classes = append(classes, strconv.FormatInt(v.Int(), 10))
		case v.CanUint():
			classes = append(classes, strconv.FormatUint(v.Uint(), 10))
		case v.CanFloat():
			classes = append(classes, strconv.FormatFloat(v.Float(), 'f', 32, 64))
		}
	}
	return strings.Join(classes, " ")
}

func isZeroValue(value reflect.Value) bool {
	if !value.IsValid() || value.IsZero() {
		return true
	}
	switch value.Kind() {
	case reflect.Slice:
		return value.Len() == 0
	}
	return false
}
