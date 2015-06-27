package tpl

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type tstNoStringer struct {
}

type tstCompareType int

const (
	tstEq tstCompareType = iota
	tstNe
	tstGt
	tstGe
	tstLt
	tstLe
)

func tstIsEq(tp tstCompareType) bool {
	return tp == tstEq || tp == tstGe || tp == tstLe
}

func tstIsGt(tp tstCompareType) bool {
	return tp == tstGt || tp == tstGe
}

func tstIsLt(tp tstCompareType) bool {
	return tp == tstLt || tp == tstLe
}

func TestCompare(t *testing.T) {
	for _, this := range []struct {
		tstCompareType
		funcUnderTest func(a, b interface{}) bool
	}{
		{tstGt, Gt},
		{tstLt, Lt},
		{tstGe, Ge},
		{tstLe, Le},
		{tstEq, Eq},
		{tstNe, Ne},
	} {
		doTestCompare(t, this.tstCompareType, this.funcUnderTest)
	}

}

func doTestCompare(t *testing.T, tp tstCompareType, funcUnderTest func(a, b interface{}) bool) {
	for i, this := range []struct {
		left            interface{}
		right           interface{}
		expectIndicator int
	}{
		{5, 8, -1},
		{8, 5, 1},
		{5, 5, 0},
		{int(5), int64(5), 0},
		{int32(5), int(5), 0},
		{int16(4), int(5), -1},
		{uint(15), uint64(15), 0},
		{-2, 1, -1},
		{2, -5, 1},
		{0.0, 1.23, -1},
		{1.1, 1.1, 0},
		{float32(1.0), float64(1.0), 0},
		{1.23, 0.0, 1},
		{"5", "5", 0},
		{"8", "5", 1},
		{"5", "0001", 1},
		{[]int{100, 99}, []int{1, 2, 3, 4}, -1},
	} {
		result := funcUnderTest(this.left, this.right)
		success := false

		if this.expectIndicator == 0 {
			if tstIsEq(tp) {
				success = result
			} else {
				success = !result
			}
		}

		if this.expectIndicator < 0 {
			success = result && (tstIsLt(tp) || tp == tstNe)
			success = success || (!result && !tstIsLt(tp))
		}

		if this.expectIndicator > 0 {
			success = result && (tstIsGt(tp) || tp == tstNe)
			success = success || (!result && (!tstIsGt(tp) || tp != tstNe))
		}

		if !success {
			t.Errorf("[%d][%s] %v compared to %v: %t", i, path.Base(runtime.FuncForPC(reflect.ValueOf(funcUnderTest).Pointer()).Name()), this.left, this.right, result)
		}
	}
}

func TestArethmic(t *testing.T) {
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		op     rune
		expect interface{}
	}{
		{1, 2, '+', int64(3)},
		{1, 2, '-', int64(-1)},
		{2, 2, '*', int64(4)},
		{4, 2, '/', int64(2)},
		{uint8(1), uint8(3), '+', uint64(4)},
		{uint8(3), uint8(2), '-', uint64(1)},
		{uint8(2), uint8(2), '*', uint64(4)},
		{uint16(4), uint8(2), '/', uint64(2)},
		{4, 2, '¤', false},
		{4, 0, '/', false},
	} {
		// TODO(bep): Take precision into account.
		result, err := doArithmetic(this.a, this.b, this.op)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] doArethmic didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] doArethmic got %v (%T) but expected %v (%T)", i, result, result, this.expect, this.expect)
			}
		}
	}
}

func TestMod(t *testing.T) {
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		expect interface{}
	}{
		{3, 2, int64(1)},
		{3, 1, int64(0)},
		{3, 0, false},
		{0, 3, int64(0)},
		{3.1, 2, false},
		{3, 2.1, false},
		{3.1, 2.1, false},
		{int8(3), int8(2), int64(1)},
		{int16(3), int16(2), int64(1)},
		{int32(3), int32(2), int64(1)},
		{int64(3), int64(2), int64(1)},
	} {
		result, err := Mod(this.a, this.b)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] modulo didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] modulo got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestModBool(t *testing.T) {
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		expect interface{}
	}{
		{3, 3, true},
		{3, 2, false},
		{3, 1, true},
		{3, 0, nil},
		{0, 3, true},
		{3.1, 2, nil},
		{3, 2.1, nil},
		{3.1, 2.1, nil},
		{int8(3), int8(3), true},
		{int8(3), int8(2), false},
		{int16(3), int16(3), true},
		{int16(3), int16(2), false},
		{int32(3), int32(3), true},
		{int32(3), int32(2), false},
		{int64(3), int64(3), true},
		{int64(3), int64(2), false},
	} {
		result, err := ModBool(this.a, this.b)
		if this.expect == nil {
			if err == nil {
				t.Errorf("[%d] modulo didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] modulo got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestFirst(t *testing.T) {
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"a", "b"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{100, 200}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{100}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
	} {
		results, err := First(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestLast(t *testing.T) {
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c"}, []string{"b", "c"}},
		{int32(3), []string{"a", "b"}, []string{"a", "b"}},
		{int64(2), []int{100, 200, 300}, []int{200, 300}},
		{100, []int{100, 200}, []int{100, 200}},
		{"1", []int{100, 200, 300}, []int{300}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
	} {
		results, err := Last(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestAfter(t *testing.T) {
	for i, this := range []struct {
		count    interface{}
		sequence interface{}
		expect   interface{}
	}{
		{int(2), []string{"a", "b", "c", "d"}, []string{"c", "d"}},
		{int32(3), []string{"a", "b"}, false},
		{int64(2), []int{100, 200, 300}, []int{300}},
		{100, []int{100, 200}, false},
		{"1", []int{100, 200, 300}, []int{200, 300}},
		{int64(-1), []int{100, 200, 300}, false},
		{"noint", []int{100, 200, 300}, false},
		{1, nil, false},
		{nil, []int{100}, false},
		{1, t, false},
	} {
		results, err := After(this.count, this.sequence)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] First didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
			}
		}
	}
}

func TestIn(t *testing.T) {
	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		expect bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "12", "c"}, 12, false},
		{[]int{1, 2, 4}, 2, true},
		{[]int{1, 2, 4}, 3, false},
		{[]float64{1.23, 2.45, 4.67}, 1.23, true},
		{[]float64{1.234567, 2.45, 4.67}, 1.234568, false},
		{"this substring should be found", "substring", true},
		{"this substring should not be found", "subseastring", false},
	} {
		result := In(this.v1, this.v2)

		if result != this.expect {
			t.Errorf("[%d] got %v but expected %v", i, result, this.expect)
		}
	}
}

func TestSlicestr(t *testing.T) {
	for i, this := range []struct {
		v1     interface{}
		v2     []int
		expect interface{}
	}{
		{"abc", []int{1, 2}, "b"},
		{"abc", []int{1, 3}, "bc"},
		{"abc", []int{0, 1}, "a"},
		{"abcdef", []int{}, "abcdef"},
		{"abcdef", []int{0, 6}, "abcdef"},
		{"abcdef", []int{0, 2}, "ab"},
		{"abcdef", []int{2}, "cdef"},
		{123, []int{1, 3}, "23"},
		{123, []int{1, 2, 3}, false},
		{"abcdef", []int{6}, false},
		{"abcdef", []int{4, 7}, false},
		{"abcdef", []int{-1}, false},
		{"abcdef", []int{-1, 7}, false},
		{"abcdef", []int{1, -1}, false},
		{tstNoStringer{}, []int{0, 1}, false},
	} {
		result, err := Slicestr(this.v1, this.v2...)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Slice didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}
}

func TestSubstr(t *testing.T) {
	var err error
	var n int
	for i, this := range []struct {
		v1     interface{}
		v2     interface{}
		v3     interface{}
		expect interface{}
	}{
		{"abc", 1, 2, "bc"},
		{"abc", 0, 1, "a"},
		{"abcdef", -1, 2, "ef"},
		{"abcdef", -3, 3, "bcd"},
		{"abcdef", 0, -1, "abcde"},
		{"abcdef", 2, -1, "cde"},
		{"abcdef", 4, -4, false},
		{"abcdef", 7, 1, false},
		{"abcdef", 1, 100, "bcdef"},
		{"abcdef", -100, 3, "abc"},
		{"abcdef", -3, -1, "de"},
		{"abcdef", 2, nil, "cdef"},
		{"abcdef", int8(2), nil, "cdef"},
		{"abcdef", int16(2), nil, "cdef"},
		{"abcdef", int32(2), nil, "cdef"},
		{"abcdef", int64(2), nil, "cdef"},
		{"abcdef", 2, int8(3), "cde"},
		{"abcdef", 2, int16(3), "cde"},
		{"abcdef", 2, int32(3), "cde"},
		{"abcdef", 2, int64(3), "cde"},
		{123, 1, 3, "23"},
		{1.2e3, 0, 4, "1200"},
		{tstNoStringer{}, 0, 1, false},
		{"abcdef", 2.0, nil, false},
		{"abcdef", 2.0, 2, false},
		{"abcdef", 2, 2.0, false},
	} {
		var result string
		n = i

		if this.v3 == nil {
			result, err = Substr(this.v1, this.v2)
		} else {
			result, err = Substr(this.v1, this.v2, this.v3)
		}

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Substr didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}

	n++
	_, err = Substr("abcdef")
	if err == nil {
		t.Errorf("[%d] Substr didn't return an expected error", n)
	}

	n++
	_, err = Substr("abcdef", 1, 2, 3)
	if err == nil {
		t.Errorf("[%d] Substr didn't return an expected error", n)
	}
}

func TestSplit(t *testing.T) {
	for i, this := range []struct {
		v1     interface{}
		v2     string
		expect interface{}
	}{
		{"a, b", ", ", []string{"a", "b"}},
		{"a & b & c", " & ", []string{"a", "b", "c"}},
		{"http://exmaple.com", "http://", []string{"", "exmaple.com"}},
		{123, "2", []string{"1", "3"}},
		{tstNoStringer{}, ",", false},
	} {
		result, err := Split(this.v1, this.v2)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Split didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
			}
		}
	}

}

func TestIntersect(t *testing.T) {
	for i, this := range []struct {
		sequence1 interface{}
		sequence2 interface{}
		expect    interface{}
	}{
		{[]string{"a", "b", "c"}, []string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b"}, []string{"a", "b", "c"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"d", "e"}, []string{}},
		{[]string{}, []string{}, []string{}},
		{[]string{"a", "b"}, nil, make([]interface{}, 0)},
		{nil, []string{"a", "b"}, make([]interface{}, 0)},
		{nil, nil, make([]interface{}, 0)},
		{[]string{"1", "2"}, []int{1, 2}, []string{}},
		{[]int{1, 2}, []string{"1", "2"}, []int{}},
		{[]int{1, 2, 4}, []int{2, 4}, []int{2, 4}},
		{[]int{2, 4}, []int{1, 2, 4}, []int{2, 4}},
		{[]int{1, 2, 4}, []int{3, 6}, []int{}},
		{[]float64{2.2, 4.4}, []float64{1.1, 2.2, 4.4}, []float64{2.2, 4.4}},
	} {
		results, err := Intersect(this.sequence1, this.sequence2)
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(results, this.expect) {
			t.Errorf("[%d] got %v but expected %v", i, results, this.expect)
		}
	}

	_, err1 := Intersect("not an array or slice", []string{"a"})

	if err1 == nil {
		t.Error("Expected error for non array as first arg")
	}

	_, err2 := Intersect([]string{"a"}, "not an array or slice")

	if err2 == nil {
		t.Error("Expected error for non array as second arg")
	}
}

func TestIsSet(t *testing.T) {
	aSlice := []interface{}{1, 2, 3, 5}
	aMap := map[string]interface{}{"a": 1, "b": 2}

	assert.True(t, IsSet(aSlice, 2))
	assert.True(t, IsSet(aMap, "b"))
	assert.False(t, IsSet(aSlice, 22))
	assert.False(t, IsSet(aMap, "bc"))
}

func (x *TstX) TstRp() string {
	return "r" + x.A
}

func (x TstX) TstRv() string {
	return "r" + x.B
}

func (x TstX) unexportedMethod() string {
	return x.unexported
}

func (x TstX) MethodWithArg(s string) string {
	return s
}

func (x TstX) MethodReturnNothing() {}

func (x TstX) MethodReturnErrorOnly() error {
	return errors.New("something error occured")
}

func (x TstX) MethodReturnTwoValues() (string, string) {
	return "foo", "bar"
}

func (x TstX) MethodReturnValueWithError() (string, error) {
	return "", errors.New("something error occured")
}

func (x TstX) String() string {
	return fmt.Sprintf("A: %s, B: %s", x.A, x.B)
}

type TstX struct {
	A, B       string
	unexported string
}

func TestTimeUnix(t *testing.T) {
	var sec int64 = 1234567890
	tv := reflect.ValueOf(time.Unix(sec, 0))
	i := 1

	res := timeUnix(tv)
	if sec != res {
		t.Errorf("[%d] timeUnix got %v but expected %v", i, res, sec)
	}

	i++
	func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("[%d] timeUnix didn't return an expected error", i)
			}
		}()
		iv := reflect.ValueOf(sec)
		timeUnix(iv)
	}(t)
}

func TestEvaluateSubElem(t *testing.T) {
	tstx := TstX{A: "foo", B: "bar"}
	var inner struct {
		S fmt.Stringer
	}
	inner.S = tstx
	interfaceValue := reflect.ValueOf(&inner).Elem().Field(0)

	for i, this := range []struct {
		value  reflect.Value
		key    string
		expect interface{}
	}{
		{reflect.ValueOf(tstx), "A", "foo"},
		{reflect.ValueOf(&tstx), "TstRp", "rfoo"},
		{reflect.ValueOf(tstx), "TstRv", "rbar"},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1, "foo"},
		{reflect.ValueOf(map[string]string{"key1": "foo", "key2": "bar"}), "key1", "foo"},
		{interfaceValue, "String", "A: foo, B: bar"},
		{reflect.Value{}, "foo", false},
		//{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), 1.2, false},
		{reflect.ValueOf(tstx), "unexported", false},
		{reflect.ValueOf(tstx), "unexportedMethod", false},
		{reflect.ValueOf(tstx), "MethodWithArg", false},
		{reflect.ValueOf(tstx), "MethodReturnNothing", false},
		{reflect.ValueOf(tstx), "MethodReturnErrorOnly", false},
		{reflect.ValueOf(tstx), "MethodReturnTwoValues", false},
		{reflect.ValueOf(tstx), "MethodReturnValueWithError", false},
		{reflect.ValueOf((*TstX)(nil)), "A", false},
		{reflect.ValueOf(tstx), "C", false},
		{reflect.ValueOf(map[int]string{1: "foo", 2: "bar"}), "1", false},
		{reflect.ValueOf([]string{"foo", "bar"}), "1", false},
	} {
		result, err := evaluateSubElem(this.value, this.key)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] evaluateSubElem didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result.Kind() != reflect.String || result.String() != this.expect {
				t.Errorf("[%d] evaluateSubElem with %v got %v but expected %v", i, this.key, result, this.expect)
			}
		}
	}
}

func TestCheckCondition(t *testing.T) {
	type expect struct {
		result  bool
		isError bool
	}

	for i, this := range []struct {
		value reflect.Value
		match reflect.Value
		op    string
		expect
	}{
		{reflect.ValueOf(123), reflect.ValueOf(123), "", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("foo"), "", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"",
			expect{true, false},
		},
		{reflect.ValueOf(nil), reflect.ValueOf(nil), "", expect{true, false}},
		{reflect.ValueOf(123), reflect.ValueOf(456), "!=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), "!=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			"!=",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(nil), "!=", expect{true, false}},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">=", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">=",
			expect{true, false},
		},
		{reflect.ValueOf(456), reflect.ValueOf(123), ">", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar"), ">", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			">",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<=", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<=", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<=",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf(456), "<", expect{true, false}},
		{reflect.ValueOf("bar"), reflect.ValueOf("foo"), "<", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			"<",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{123, 45, 678}), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"foo", "bar", "baz"}), "in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.June, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"in",
			expect{true, false},
		},
		{reflect.ValueOf(123), reflect.ValueOf([]int{45, 678}), "not in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]string{"bar", "baz"}), "not in", expect{true, false}},
		{
			reflect.ValueOf(time.Date(2015, time.May, 26, 19, 18, 56, 12345, time.UTC)),
			reflect.ValueOf([]time.Time{
				time.Date(2015, time.February, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.March, 26, 19, 18, 56, 12345, time.UTC),
				time.Date(2015, time.April, 26, 19, 18, 56, 12345, time.UTC),
			}),
			"not in",
			expect{true, false},
		},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar-foo-baz"), "in", expect{true, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf("bar--baz"), "not in", expect{true, false}},
		{reflect.Value{}, reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.Value{}, "", expect{false, false}},
		{reflect.ValueOf((*TstX)(nil)), reflect.ValueOf("foo"), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf((*TstX)(nil)), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf(map[int]string{}), "", expect{false, false}},
		{reflect.ValueOf("foo"), reflect.ValueOf([]int{1, 2}), "", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf([]int{}), "in", expect{false, false}},
		{reflect.ValueOf(123), reflect.ValueOf(123), "op", expect{false, true}},
	} {
		result, err := checkCondition(this.value, this.match, this.op)
		if this.expect.isError {
			if err == nil {
				t.Errorf("[%d] checkCondition didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if result != this.expect.result {
				t.Errorf("[%d] check condition %v %s %v, got %v but expected %v", i, this.value, this.op, this.match, result, this.expect.result)
			}
		}
	}
}

func TestWhere(t *testing.T) {
	// TODO(spf): Put these page tests back in
	//page1 := &Page{contentType: "v", Source: Source{File: *source.NewFile("/x/y/z/source.md")}}
	//page2 := &Page{contentType: "w", Source: Source{File: *source.NewFile("/y/z/a/source.md")}}

	type Mid struct {
		Tst TstX
	}

	for i, this := range []struct {
		sequence interface{}
		key      interface{}
		op       string
		match    interface{}
		expect   interface{}
	}{
		{
			sequence: []map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []map[int]string{
				{1: "a", 2: "m"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4},
			},
			key: "b", match: 4,
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []TstX{
				{A: "e", B: "f"},
			},
		},
		{
			sequence: []*map[int]string{
				{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"},
			},
			key: 2, match: "m",
			expect: []*map[int]string{
				{1: "a", 2: "m"},
			},
		},
		{
			sequence: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", match: "f",
			expect: []*TstX{
				{A: "e", B: "f"},
			},
		},
		{
			sequence: []*TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRp", match: "rc",
			expect: []*TstX{
				{A: "c", B: "d"},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "c"},
			},
			key: "TstRv", match: "rc",
			expect: []TstX{
				{A: "e", B: "c"},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: ".foo.B", match: "d",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]TstX{
				{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRv", match: "rd",
			expect: []map[string]TstX{
				{"foo": TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]*TstX{
				{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}},
			},
			key: "foo.TstRp", match: "rc",
			expect: []map[string]*TstX{
				{"foo": &TstX{A: "c", B: "d"}},
			},
		},
		{
			sequence: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.B", match: "d",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRv", match: "rd",
			expect: []map[string]Mid{
				{"foo": Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": &Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": &Mid{Tst: TstX{A: "e", B: "f"}}},
			},
			key: "foo.Tst.TstRp", match: "rc",
			expect: []map[string]*Mid{
				{"foo": &Mid{Tst: TstX{A: "c", B: "d"}}},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: 3,
			expect: []map[string]int{
				{"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "!=", match: "f",
			expect: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "b": 6},
			},
			key: "b", op: "in", match: []int{3, 4, 5},
			expect: []map[string]int{
				{"a": 3, "b": 4},
			},
		},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "not in", match: []string{"c", "d", "e"},
			expect: []TstX{
				{A: "a", B: "b"}, {A: "e", B: "f"},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "", match: nil,
			expect: []map[string]int{
				{"a": 3},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: "!=", match: nil,
			expect: []map[string]int{
				{"a": 1, "b": 2}, {"a": 5, "b": 6},
			},
		},
		{
			sequence: []map[string]int{
				{"a": 1, "b": 2}, {"a": 3}, {"a": 5, "b": 6},
			},
			key: "b", op: ">", match: nil,
			expect: []map[string]int{},
		},
		{sequence: (*[]TstX)(nil), key: "A", match: "a", expect: false},
		{sequence: TstX{A: "a", B: "b"}, key: "A", match: "a", expect: false},
		{sequence: []map[string]*TstX{{"foo": nil}}, key: "foo.B", match: "d", expect: false},
		{
			sequence: []TstX{
				{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"},
			},
			key: "B", op: "op", match: "f",
			expect: false,
		},
		//{[]*Page{page1, page2}, "Type", "v", []*Page{page1}},
		//{[]*Page{page1, page2}, "Section", "y", []*Page{page2}},
	} {
		var results interface{}
		var err error
		if len(this.op) > 0 {
			results, err = Where(this.sequence, this.key, this.op, this.match)
		} else {
			results, err = Where(this.sequence, this.key, this.match)
		}
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Where didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(results, this.expect) {
				t.Errorf("[%d] Where clause matching %v with %v, got %v but expected %v", i, this.key, this.match, results, this.expect)
			}
		}
	}

	var err error
	_, err = Where(map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1)
	if err == nil {
		t.Errorf("Where called with none string op value didn't return an expected error")
	}

	_, err = Where(map[string]int{"a": 1, "b": 2}, "a", []byte("="), 1, 2)
	if err == nil {
		t.Errorf("Where called with more than two variable arguments didn't return an expected error")
	}

	_, err = Where(map[string]int{"a": 1, "b": 2}, "a")
	if err == nil {
		t.Errorf("Where called with no variable arguments didn't return an expected error")
	}
}

func TestDelimit(t *testing.T) {
	for i, this := range []struct {
		sequence  interface{}
		delimiter interface{}
		last      interface{}
		expect    template.HTML
	}{
		{[]string{"class1", "class2", "class3"}, " ", nil, "class1 class2 class3"},
		{[]int{1, 2, 3, 4, 5}, ",", nil, "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", nil, "1, 2, 3, 4, 5"},
		{[]string{"class1", "class2", "class3"}, " ", " and ", "class1 class2 and class3"},
		{[]int{1, 2, 3, 4, 5}, ",", ",", "1,2,3,4,5"},
		{[]int{1, 2, 3, 4, 5}, ", ", ", and ", "1, 2, 3, 4, and 5"},
		// test maps with and without sorting required
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", nil, "10--20--30--40--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", nil, "10--20--30--40--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", nil, "30--20--10--40--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", nil, "50--40--10--30--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", nil, "10--20--30--40--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", nil, "30--20--10--40--50"},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, "--", nil, "30--20--10--40--50"},
		// test maps with a last delimiter
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, "--", "--and--", "50--40--10--30--and--20"},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, "--", "--and--", "10--20--30--40--and--50"},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
		{map[float64]string{3.5: "10", 2.5: "20", 1.5: "30", 4.5: "40", 5.5: "50"}, "--", "--and--", "30--20--10--40--and--50"},
	} {
		var result template.HTML
		var err error
		if this.last == nil {
			result, err = Delimit(this.sequence, this.delimiter)
		} else {
			result, err = Delimit(this.sequence, this.delimiter, this.last)
		}
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] Delimit called on sequence: %v | delimiter: `%v` | last: `%v`, got %v but expected %v", i, this.sequence, this.delimiter, this.last, result, this.expect)
		}
	}
}

func TestSort(t *testing.T) {
	type ts struct {
		MyInt    int
		MyFloat  float64
		MyString string
	}
	for i, this := range []struct {
		sequence    interface{}
		sortByField interface{}
		sortAsc     string
		expect      []interface{}
	}{
		{[]string{"class1", "class2", "class3"}, nil, "asc", []interface{}{"class1", "class2", "class3"}},
		{[]string{"class3", "class1", "class2"}, nil, "asc", []interface{}{"class1", "class2", "class3"}},
		{[]int{1, 2, 3, 4, 5}, nil, "asc", []interface{}{1, 2, 3, 4, 5}},
		{[]int{5, 4, 3, 1, 2}, nil, "asc", []interface{}{1, 2, 3, 4, 5}},
		// test map sorting by keys
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, nil, "asc", []interface{}{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, nil, "asc", []interface{}{30, 20, 10, 40, 50}},
		{map[string]string{"1": "10", "2": "20", "3": "30", "4": "40", "5": "50"}, nil, "asc", []interface{}{"10", "20", "30", "40", "50"}},
		{map[string]string{"3": "10", "2": "20", "1": "30", "4": "40", "5": "50"}, nil, "asc", []interface{}{"30", "20", "10", "40", "50"}},
		{map[string]string{"one": "10", "two": "20", "three": "30", "four": "40", "five": "50"}, nil, "asc", []interface{}{"50", "40", "10", "30", "20"}},
		{map[int]string{1: "10", 2: "20", 3: "30", 4: "40", 5: "50"}, nil, "asc", []interface{}{"10", "20", "30", "40", "50"}},
		{map[int]string{3: "10", 2: "20", 1: "30", 4: "40", 5: "50"}, nil, "asc", []interface{}{"30", "20", "10", "40", "50"}},
		{map[float64]string{3.3: "10", 2.3: "20", 1.3: "30", 4.3: "40", 5.3: "50"}, nil, "asc", []interface{}{"30", "20", "10", "40", "50"}},
		// test map sorting by value
		{map[string]int{"1": 10, "2": 20, "3": 30, "4": 40, "5": 50}, "value", "asc", []interface{}{10, 20, 30, 40, 50}},
		{map[string]int{"3": 10, "2": 20, "1": 30, "4": 40, "5": 50}, "value", "asc", []interface{}{10, 20, 30, 40, 50}},
		// test map sorting by field value
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyInt",
			"asc",
			[]interface{}{ts{10, 10.5, "ten"}, ts{20, 20.5, "twenty"}, ts{30, 30.5, "thirty"}, ts{40, 40.5, "forty"}, ts{50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyFloat",
			"asc",
			[]interface{}{ts{10, 10.5, "ten"}, ts{20, 20.5, "twenty"}, ts{30, 30.5, "thirty"}, ts{40, 40.5, "forty"}, ts{50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": {10, 10.5, "ten"}, "2": {20, 20.5, "twenty"}, "3": {30, 30.5, "thirty"}, "4": {40, 40.5, "forty"}, "5": {50, 50.5, "fifty"}},
			"MyString",
			"asc",
			[]interface{}{ts{50, 50.5, "fifty"}, ts{40, 40.5, "forty"}, ts{10, 10.5, "ten"}, ts{30, 30.5, "thirty"}, ts{20, 20.5, "twenty"}},
		},
		// Test sort desc
		{[]string{"class1", "class2", "class3"}, "value", "desc", []interface{}{"class3", "class2", "class1"}},
		{[]string{"class3", "class1", "class2"}, "value", "desc", []interface{}{"class3", "class2", "class1"}},
	} {
		var result []interface{}
		var err error
		if this.sortByField == nil {
			result, err = Sort(this.sequence)
		} else {
			result, err = Sort(this.sequence, this.sortByField, this.sortAsc)
		}
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] Sort called on sequence: %v | sortByField: `%v` | got %v but expected %v", i, this.sequence, this.sortByField, result, this.expect)
		}
	}
}

func TestReturnWhenSet(t *testing.T) {
	for i, this := range []struct {
		data   interface{}
		key    interface{}
		expect interface{}
	}{
		{[]int{1, 2, 3}, 1, int64(2)},
		{[]uint{1, 2, 3}, 1, uint64(2)},
		{[]float64{1.1, 2.2, 3.3}, 1, float64(2.2)},
		{[]string{"foo", "bar", "baz"}, 1, "bar"},
		{[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}}, 1, ""},
		{map[string]int{"foo": 1, "bar": 2, "baz": 3}, "bar", int64(2)},
		{map[string]uint{"foo": 1, "bar": 2, "baz": 3}, "bar", uint64(2)},
		{map[string]float64{"foo": 1.1, "bar": 2.2, "baz": 3.3}, "bar", float64(2.2)},
		{map[string]string{"foo": "FOO", "bar": "BAR", "baz": "BAZ"}, "bar", "BAR"},
		{map[string]TstX{"foo": {A: "a", B: "b"}, "bar": {A: "c", B: "d"}, "baz": {A: "e", B: "f"}}, "bar", ""},
		{(*[]string)(nil), "bar", ""},
	} {
		result := ReturnWhenSet(this.data, this.key)
		if !reflect.DeepEqual(result, this.expect) {
			t.Errorf("[%d] ReturnWhenSet got %v (type %v) but expected %v (type %v)", i, result, reflect.TypeOf(result), this.expect, reflect.TypeOf(this.expect))
		}
	}
}

func TestMarkdownify(t *testing.T) {

	result := Markdownify("Hello **World!**")

	expect := template.HTML("Hello <strong>World!</strong>")

	if result != expect {
		t.Errorf("Markdownify: got '%s', expected '%s'", result, expect)
	}
}

func TestApply(t *testing.T) {
	strings := []interface{}{"a\n", "b\n"}
	noStringers := []interface{}{tstNoStringer{}, tstNoStringer{}}

	var nilErr *error = nil

	chomped, _ := Apply(strings, "chomp", ".")
	assert.Equal(t, []interface{}{"a", "b"}, chomped)

	chomped, _ = Apply(strings, "chomp", "c\n")
	assert.Equal(t, []interface{}{"c", "c"}, chomped)

	chomped, _ = Apply(nil, "chomp", ".")
	assert.Equal(t, []interface{}{}, chomped)

	_, err := Apply(strings, "apply", ".")
	if err == nil {
		t.Errorf("apply with apply should fail")
	}

	_, err = Apply(nilErr, "chomp", ".")
	if err == nil {
		t.Errorf("apply with nil in seq should fail")
	}

	_, err = Apply(strings, "dobedobedo", ".")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

	_, err = Apply(noStringers, "chomp", ".")
	if err == nil {
		t.Errorf("apply when func fails should fail")
	}

	_, err = Apply(tstNoStringer{}, "chomp", ".")
	if err == nil {
		t.Errorf("apply with non-sequence should fail")
	}

}

func TestChomp(t *testing.T) {
	base := "\n This is\na story "
	for i, item := range []string{
		"\n", "\n\n",
		"\r", "\r\r",
		"\r\n", "\r\n\r\n",
	} {
		chomped, _ := Chomp(base + item)

		if chomped != base {
			t.Errorf("[%d] Chomp failed, got '%v'", i, chomped)
		}

		_, err := Chomp(tstNoStringer{})

		if err == nil {
			t.Errorf("Chomp should fail")
		}
	}
}

func TestReplace(t *testing.T) {
	v, _ := Replace("aab", "a", "b")
	assert.Equal(t, "bbb", v)
	v, _ = Replace("11a11", 1, 2)
	assert.Equal(t, "22a22", v)
	v, _ = Replace(12345, 1, 2)
	assert.Equal(t, "22345", v)
	_, e := Replace(tstNoStringer{}, "a", "b")
	assert.NotNil(t, e, "tstNoStringer isn't trimmable")
	_, e = Replace("a", tstNoStringer{}, "b")
	assert.NotNil(t, e, "tstNoStringer cannot be converted to string")
	_, e = Replace("a", "b", tstNoStringer{})
	assert.NotNil(t, e, "tstNoStringer cannot be converted to string")
}

func TestTrim(t *testing.T) {
	v, _ := Trim("1234 my way 13", "123")
	assert.Equal(t, "4 my way ", v)
	v, _ = Trim("   my way    ", " ")
	assert.Equal(t, "my way", v)
	v, _ = Trim(1234, "14")
	assert.Equal(t, "23", v)
	_, e := Trim(tstNoStringer{}, " ")
	assert.NotNil(t, e, "tstNoStringer isn't trimmable")
}

func TestDateFormat(t *testing.T) {
	for i, this := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
		{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
		{"Monday, Jan 2, 2006", time.Date(2015, time.January, 21, 0, 0, 0, 0, time.UTC), "Wednesday, Jan 21, 2015"},
		{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
		{"Monday, Jan 2, 2006", 1421733600, false},
		{"Monday, Jan 2, 2006", 1421733600.123, false},
	} {
		result, err := DateFormat(this.layout, this.value)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] DateFormat didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] DateFormat failed: %s", i, err)
				continue
			}
			if result != this.expect {
				t.Errorf("[%d] DateFormat got %v but expected %v", i, result, this.expect)
			}
		}
	}
}

func TestSafeHTML(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`<div></div>`, `{{ . }}`, `&lt;div&gt;&lt;/div&gt;`, `<div></div>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		err = tmpl.Execute(buf, SafeHTML(this.str))
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by SafeHTML returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by SafeHTML, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestSafeHTMLAttr(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`href="irc://irc.freenode.net/#golang"`, `<a {{ . }}>irc</a>`, `<a ZgotmplZ>irc</a>`, `<a href="irc://irc.freenode.net/#golang">irc</a>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		err = tmpl.Execute(buf, SafeHTMLAttr(this.str))
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by SafeHTMLAttr returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by SafeHTMLAttr, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestSafeCSS(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`width: 60px;`, `<div style="{{ . }}"></div>`, `<div style="ZgotmplZ"></div>`, `<div style="width: 60px;"></div>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		err = tmpl.Execute(buf, SafeCSS(this.str))
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by SafeCSS returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by SafeCSS, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}

func TestSafeURL(t *testing.T) {
	for i, this := range []struct {
		str                 string
		tmplStr             string
		expectWithoutEscape string
		expectWithEscape    string
	}{
		{`irc://irc.freenode.net/#golang`, `<a href="{{ . }}">IRC</a>`, `<a href="#ZgotmplZ">IRC</a>`, `<a href="irc://irc.freenode.net/#golang">IRC</a>`},
	} {
		tmpl, err := template.New("test").Parse(this.tmplStr)
		if err != nil {
			t.Errorf("[%d] unable to create new html template %q: %s", i, this.tmplStr, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, this.str)
		if err != nil {
			t.Errorf("[%d] execute template with a raw string value returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithoutEscape {
			t.Errorf("[%d] execute template with a raw string value, got %v but expected %v", i, buf.String(), this.expectWithoutEscape)
		}

		buf.Reset()
		err = tmpl.Execute(buf, SafeURL(this.str))
		if err != nil {
			t.Errorf("[%d] execute template with an escaped string value by SafeURL returns unexpected error: %s", i, err)
		}
		if buf.String() != this.expectWithEscape {
			t.Errorf("[%d] execute template with an escaped string value by SafeURL, got %v but expected %v", i, buf.String(), this.expectWithEscape)
		}
	}
}
