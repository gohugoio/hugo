package hugolib

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func interfaceToStringToDate(i interface{}) time.Time {
	s := interfaceToString(i)

	if d, e := parseDateWith(s, []string{
		time.RFC3339,
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
	}); e == nil {
		return d
	}

	return time.Unix(0, 0)
}

// TODO remove this and return a proper error.
func errorf(str string, a ...interface{}) {
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
	default:
		errorf("Only Boolean values are supported for this YAML key")
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

func interfaceToString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	default:
		errorf("Only Strings are supported for this YAML key")
	}

	return ""
}
