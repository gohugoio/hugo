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

package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/spf13/cast"
	bp "github.com/spf13/hugo/bufferpool"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Filepath separator defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// FindAvailablePort returns an available and valid TCP port.
func FindAvailablePort() (*net.TCPAddr, error) {
	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer l.Close()
		addr := l.Addr()
		if a, ok := addr.(*net.TCPAddr); ok {
			return a, nil
		}
		return nil, fmt.Errorf("Unable to obtain a valid tcp port. %v", addr)
	}
	return nil, err
}

// InStringArray checks if a string is an element of a slice of strings
// and returns a boolean value.
func InStringArray(arr []string, el string) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

// GuessType attempts to guess the type of file from a given string.
func GuessType(in string) string {
	switch strings.ToLower(in) {
	case "md", "markdown", "mdown":
		return "markdown"
	case "asciidoc", "adoc", "ad":
		return "asciidoc"
	case "mmark":
		return "mmark"
	case "rst":
		return "rst"
	case "html", "htm":
		return "html"
	}

	return "unknown"
}

// ReaderToBytes takes an io.Reader argument, reads from it
// and returns bytes.
func ReaderToBytes(lines io.Reader) []byte {
	if lines == nil {
		return []byte{}
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	b.ReadFrom(lines)

	bc := make([]byte, b.Len(), b.Len())
	copy(bc, b.Bytes())
	return bc
}

// ReaderToString is the same as ReaderToBytes, but returns a string.
func ReaderToString(lines io.Reader) string {
	if lines == nil {
		return ""
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	b.ReadFrom(lines)
	return b.String()
}

// StringToReader does the opposite of ReaderToString.
func StringToReader(in string) io.Reader {
	return strings.NewReader(in)
}

// BytesToReader does the opposite of ReaderToBytes.
func BytesToReader(in []byte) io.Reader {
	return bytes.NewReader(in)
}

// ReaderContains reports whether subslice is within r.
func ReaderContains(r io.Reader, subslice []byte) bool {

	if r == nil || len(subslice) == 0 {
		return false
	}

	bufflen := len(subslice) * 4
	halflen := bufflen / 2
	buff := make([]byte, bufflen)
	var err error
	var n, i int

	for {
		i++
		if i == 1 {
			n, err = io.ReadAtLeast(r, buff[:halflen], halflen)
		} else {
			if i != 2 {
				// shift left to catch overlapping matches
				copy(buff[:], buff[halflen:])
			}
			n, err = io.ReadAtLeast(r, buff[halflen:], halflen)
		}

		if n > 0 && bytes.Contains(buff, subslice) {
			return true
		}

		if err != nil {
			break
		}
	}
	return false
}

// ThemeSet checks whether a theme is in use or not.
func ThemeSet() bool {
	return viper.GetString("theme") != ""
}

// DistinctErrorLogger ignores duplicate log statements.
type DistinctErrorLogger struct {
	sync.RWMutex
	m map[string]bool
}

// Printf will ERROR log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
func (l *DistinctErrorLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	if !l.m[logStatement] {
		jww.ERROR.Print(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}

// NewDistinctErrorLogger creates a new DistinctErrorLogger
func NewDistinctErrorLogger() *DistinctErrorLogger {
	return &DistinctErrorLogger{m: make(map[string]bool)}
}

// Avoid spamming the logs with errors
var deprecatedLogger = NewDistinctErrorLogger()

// Deprecated logs ERROR logs about a deprecation, but only once for a given set of arguments' values.
func Deprecated(object, item, alternative string) {
	deprecatedLogger.Printf("%s's %s is deprecated and will be removed in Hugo %s. Use %s instead.", object, item, NextHugoReleaseVersion(), alternative)
}

// SliceToLower goes through the source slice and lowers all values.
func SliceToLower(s []string) []string {
	if s == nil {
		return nil
	}

	l := make([]string, len(s))
	for i, v := range s {
		l[i] = strings.ToLower(v)
	}

	return l
}

// Md5String takes a string and returns its MD5 hash.
func Md5String(f string) string {
	h := md5.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum([]byte{}))
}

// Seq creates a sequence of integers.
// It's named and used as GNU's seq.
// Examples:
// 3 => 1, 2, 3
// 1 2 4 => 1, 3
// -3 => -1, -2, -3
// 1 4 => 1, 2, 3, 4
// 1 -2 => 1, 0, -1, -2
func Seq(args ...interface{}) ([]int, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, errors.New("Seq, invalid number of args: 'first' 'increment' (optional) 'last' (optional)")
	}

	intArgs := cast.ToIntSlice(args)

	if len(intArgs) < 1 || len(intArgs) > 3 {
		return nil, errors.New("Invalid argument(s) to Seq")
	}

	var inc = 1
	var last int
	var first = intArgs[0]

	if len(intArgs) == 1 {
		last = first
		if last == 0 {
			return []int{}, nil
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(intArgs) == 2 {
		last = intArgs[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = intArgs[1]
		last = intArgs[2]
		if inc == 0 {
			return nil, errors.New("'increment' must not be 0")
		}
		if first < last && inc < 0 {
			return nil, errors.New("'increment' must be > 0")
		}
		if first > last && inc > 0 {
			return nil, errors.New("'increment' must be < 0")
		}
	}

	// sanity check
	if last < -100000 {
		return nil, errors.New("size of result exeeds limit")
	}
	size := int(((last - first) / inc) + 1)

	// sanity check
	if size <= 0 || size > 2000 {
		return nil, errors.New("size of result exeeds limit")
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq, nil
}

// DoArithmetic performs arithmetic operations (+,-,*,/) using reflection to
// determine the type of the two terms.
func DoArithmetic(a, b interface{}, op rune) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64
	var af, bf float64
	var au, bu uint64
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
			if ai >= 0 {
				au = uint64(ai)
				ai = 0
			} else {
				bi = int64(bu) // may overflow
				bu = 0
			}
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au = av.Uint()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
			if bi >= 0 {
				bu = uint64(bi)
				bi = 0
			} else {
				ai = int64(au) // may overflow
				au = 0
			}
		case reflect.Float32, reflect.Float64:
			af = float64(au) // may overflow
			au = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.String:
		as := av.String()
		if bv.Kind() == reflect.String && op == '+' {
			bs := bv.String()
			return as + bs, nil
		}
		return nil, errors.New("Can't apply the operator to the values")
	default:
		return nil, errors.New("Can't apply the operator to the values")
	}

	switch op {
	case '+':
		if ai != 0 || bi != 0 {
			return ai + bi, nil
		} else if af != 0 || bf != 0 {
			return af + bf, nil
		} else if au != 0 || bu != 0 {
			return au + bu, nil
		} else {
			return 0, nil
		}
	case '-':
		if ai != 0 || bi != 0 {
			return ai - bi, nil
		} else if af != 0 || bf != 0 {
			return af - bf, nil
		} else if au != 0 || bu != 0 {
			return au - bu, nil
		} else {
			return 0, nil
		}
	case '*':
		if ai != 0 || bi != 0 {
			return ai * bi, nil
		} else if af != 0 || bf != 0 {
			return af * bf, nil
		} else if au != 0 || bu != 0 {
			return au * bu, nil
		} else {
			return 0, nil
		}
	case '/':
		if bi != 0 {
			return ai / bi, nil
		} else if bf != 0 {
			return af / bf, nil
		} else if bu != 0 {
			return au / bu, nil
		} else {
			return nil, errors.New("Can't divide the value by 0")
		}
	default:
		return nil, errors.New("There is no such an operation")
	}
}
