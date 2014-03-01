package hugolib

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

func interfaceToTime(i interface{}) time.Time {
	switch s := i.(type) {
	case time.Time:
		return s
	case string:
		d, e := stringToDate(s)
		if e == nil {
			return d
		}
		errorln("Could not parse Date/Time format:", e)
	default:
		errorln("Only Time is supported for this key")
	}

	return *new(time.Time)
}

func interfaceToStringToDate(i interface{}) time.Time {
	s := interfaceToString(i)

	if d, e := stringToDate(s); e == nil {
		return d
	}

	return time.Unix(0, 0)
}

func stringToDate(s string) (time.Time, error) {
	return parseDateWith(s, []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05Z07:00",
		"02 Jan 06 15:04 MST",
		"2006-01-02",
		"02 Jan 2006",
	})
}

// TODO remove this and return a proper error.
func errorln(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}

func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, errors.New(fmt.Sprintf("Unable to parse date: %s", s))
}

func interfaceToBool(i interface{}) bool {
	switch b := i.(type) {
	case bool:
		return b
	case int:
		if i.(int) > 0 {
			return true
		}
		return false
	default:
		errorln("Only Boolean values are supported for this YAML key")
	}

	return false

}

func interfaceArrayToStringArray(i interface{}) []string {
	var a []string

	switch vv := i.(type) {
	case []interface{}:
		for _, u := range vv {
			a = append(a, interfaceToString(u))
		}
	}

	return a
}

func interfaceToFloat64(i interface{}) float64 {
	switch s := i.(type) {
	case float64:
		return s
	case float32:
		return float64(s)

	case string:
		v, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return float64(v)
		} else {
			errorln("Only Floats are supported for this key\nErr:", err)
		}

	default:
		errorln("Only Floats are supported for this key")
	}

	return 0.0
}

func interfaceToInt(i interface{}) int {
	switch s := i.(type) {
	case int:
		return s
	case int64:
		return int(s)
	case int32:
		return int(s)
	case int16:
		return int(s)
	case int8:
		return int(s)
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err == nil {
			return int(v)
		} else {
			errorln("Only Ints are supported for this key\nErr:", err)
		}
	default:
		errorln("Only Ints are supported for this key")
	}

	return 0
}

func interfaceToString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(i.(int)), 10)
	default:
		errorln(fmt.Sprintf("Only Strings are supported for this key (got type '%T'): %s", s, s))
	}

	return ""
}
