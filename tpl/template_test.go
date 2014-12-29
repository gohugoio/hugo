package tpl

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"testing"
)

func TestGt(t *testing.T) {
	for i, this := range []struct {
		left          interface{}
		right         interface{}
		leftShouldWin bool
	}{
		{5, 8, false},
		{8, 5, true},
		{5, 5, false},
		{-2, 1, false},
		{2, -5, true},
		{0.0, 1.23, false},
		{1.23, 0.0, true},
		{"8", "5", true},
		{"5", "0001", true},
		{[]int{100, 99}, []int{1, 2, 3, 4}, false},
	} {
		leftIsBigger := Gt(this.left, this.right)
		if leftIsBigger != this.leftShouldWin {
			var which string
			if leftIsBigger {
				which = "expected right to be bigger, but left was"
			} else {
				which = "expected left to be bigger, but right was"
			}
			t.Errorf("[%d] %v compared to %v: %s", i, this.left, this.right, which)
		}
	}
}

func TestDoArithmetic(t *testing.T) {
	for i, this := range []struct {
		a      interface{}
		b      interface{}
		op     rune
		expect interface{}
	}{
		{3, 2, '+', int64(5)},
		{3, 2, '-', int64(1)},
		{3, 2, '*', int64(6)},
		{3, 2, '/', int64(1)},
		{3.0, 2, '+', float64(5)},
		{3.0, 2, '-', float64(1)},
		{3.0, 2, '*', float64(6)},
		{3.0, 2, '/', float64(1.5)},
		{3, 2.0, '+', float64(5)},
		{3, 2.0, '-', float64(1)},
		{3, 2.0, '*', float64(6)},
		{3, 2.0, '/', float64(1.5)},
		{3.0, 2.0, '+', float64(5)},
		{3.0, 2.0, '-', float64(1)},
		{3.0, 2.0, '*', float64(6)},
		{3.0, 2.0, '/', float64(1.5)},
		{uint(3), uint(2), '+', uint64(5)},
		{uint(3), uint(2), '-', uint64(1)},
		{uint(3), uint(2), '*', uint64(6)},
		{uint(3), uint(2), '/', uint64(1)},
		{uint(3), 2, '+', uint64(5)},
		{uint(3), 2, '-', uint64(1)},
		{uint(3), 2, '*', uint64(6)},
		{uint(3), 2, '/', uint64(1)},
		{3, uint(2), '+', uint64(5)},
		{3, uint(2), '-', uint64(1)},
		{3, uint(2), '*', uint64(6)},
		{3, uint(2), '/', uint64(1)},
		{uint(3), -2, '+', int64(1)},
		{uint(3), -2, '-', int64(5)},
		{uint(3), -2, '*', int64(-6)},
		{uint(3), -2, '/', int64(-1)},
		{-3, uint(2), '+', int64(-1)},
		{-3, uint(2), '-', int64(-5)},
		{-3, uint(2), '*', int64(-6)},
		{-3, uint(2), '/', int64(-1)},
		{uint(3), 2.0, '+', float64(5)},
		{uint(3), 2.0, '-', float64(1)},
		{uint(3), 2.0, '*', float64(6)},
		{uint(3), 2.0, '/', float64(1.5)},
		{3.0, uint(2), '+', float64(5)},
		{3.0, uint(2), '-', float64(1)},
		{3.0, uint(2), '*', float64(6)},
		{3.0, uint(2), '/', float64(1.5)},
		{0, 0, '+', 0},
		{0, 0, '-', 0},
		{0, 0, '*', 0},
		{"foo", "bar", '+', "foobar"},
		{3, 0, '/', false},
		{3.0, 0, '/', false},
		{3, 0.0, '/', false},
		{uint(3), uint(0), '/', false},
		{3, uint(0), '/', false},
		{-3, uint(0), '/', false},
		{uint(3), 0, '/', false},
		{3.0, uint(0), '/', false},
		{uint(3), 0.0, '/', false},
		{3, "foo", '+', false},
		{3.0, "foo", '+', false},
		{uint(3), "foo", '+', false},
		{"foo", 3, '+', false},
		{"foo", "bar", '-', false},
		{3, 2, '%', false},
	} {
		result, err := doArithmetic(this.a, this.b, this.op)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] doArithmetic didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] doArithmetic got %v but expected %v", i, result, this.expect)
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
			t.Errorf("[%d] Got %v but expected %v", i, result, this.expect)
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
			t.Errorf("[%d] Got %v but expected %v", i, results, this.expect)
		}
	}

	_, err1 := Intersect("not an array or slice", []string{"a"})

	if err1 == nil {
		t.Error("Excpected error for non array as first arg")
	}

	_, err2 := Intersect([]string{"a"}, "not an array or slice")

	if err2 == nil {
		t.Error("Excpected error for non array as second arg")
	}
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
	A, B string
	unexported string
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
		match    interface{}
		expect   interface{}
	}{
		{[]map[int]string{{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"}}, 2, "m", []map[int]string{{1: "a", 2: "m"}}},
		{[]map[string]int{{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4}}, "b", 4, []map[string]int{{"a": 3, "b": 4}}},
		{[]TstX{{A: "a", B: "b"}, {A: "c", B: "d"}, {A: "e", B: "f"}}, "B", "f", []TstX{{A: "e", B: "f"}}},
		{[]*map[int]string{&map[int]string{1: "a", 2: "m"}, &map[int]string{1: "c", 2: "d"}, &map[int]string{1: "e", 3: "m"}}, 2, "m", []*map[int]string{&map[int]string{1: "a", 2: "m"}}},
		{[]*TstX{&TstX{A: "a", B: "b"}, &TstX{A: "c", B: "d"}, &TstX{A: "e", B: "f"}}, "B", "f", []*TstX{&TstX{A: "e", B: "f"}}},
		{[]*TstX{&TstX{A: "a", B: "b"}, &TstX{A: "c", B: "d"}, &TstX{A: "e", B: "c"}}, "TstRp", "rc", []*TstX{&TstX{A: "c", B: "d"}}},
		{[]TstX{TstX{A: "a", B: "b"}, TstX{A: "c", B: "d"}, TstX{A: "e", B: "c"}}, "TstRv", "rc", []TstX{TstX{A: "e", B: "c"}}},
		{[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}}, "foo.B", "d", []map[string]TstX{{"foo": TstX{A: "c", B: "d"}}}},
		{[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}}, ".foo.B", "d", []map[string]TstX{{"foo": TstX{A: "c", B: "d"}}}},
		{[]map[string]TstX{{"foo": TstX{A: "a", B: "b"}}, {"foo": TstX{A: "c", B: "d"}}, {"foo": TstX{A: "e", B: "f"}}}, "foo.TstRv", "rd", []map[string]TstX{{"foo": TstX{A: "c", B: "d"}}}},
		{[]map[string]*TstX{{"foo": &TstX{A: "a", B: "b"}}, {"foo": &TstX{A: "c", B: "d"}}, {"foo": &TstX{A: "e", B: "f"}}}, "foo.TstRp", "rc", []map[string]*TstX{{"foo": &TstX{A: "c", B: "d"}}}},
		{[]map[string]Mid{{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}}}, "foo.Tst.B", "d", []map[string]Mid{{"foo": Mid{Tst: TstX{A: "c", B: "d"}}}}},
		{[]map[string]Mid{{"foo": Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": Mid{Tst: TstX{A: "e", B: "f"}}}}, "foo.Tst.TstRv", "rd", []map[string]Mid{{"foo": Mid{Tst: TstX{A: "c", B: "d"}}}}},
		{[]map[string]*Mid{{"foo": &Mid{Tst: TstX{A: "a", B: "b"}}}, {"foo": &Mid{Tst: TstX{A: "c", B: "d"}}}, {"foo": &Mid{Tst: TstX{A: "e", B: "f"}}}}, "foo.Tst.TstRp", "rc", []map[string]*Mid{{"foo": &Mid{Tst: TstX{A: "c", B: "d"}}}}},
		{(*[]TstX)(nil), "A", "a", false},
		{TstX{A: "a", B: "b"}, "A", "a", false},
		{[]map[string]*TstX{{"foo": nil}}, "foo.B", "d", false},
		//{[]*Page{page1, page2}, "Type", "v", []*Page{page1}},
		//{[]*Page{page1, page2}, "Section", "y", []*Page{page2}},
	} {
		results, err := Where(this.sequence, this.key, this.match)
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
			map[string]ts{"1": ts{10, 10.5, "ten"}, "2": ts{20, 20.5, "twenty"}, "3": ts{30, 30.5, "thirty"}, "4": ts{40, 40.5, "forty"}, "5": ts{50, 50.5, "fifty"}},
			"MyInt",
			"asc",
			[]interface{}{ts{10, 10.5, "ten"}, ts{20, 20.5, "twenty"}, ts{30, 30.5, "thirty"}, ts{40, 40.5, "forty"}, ts{50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": ts{10, 10.5, "ten"}, "2": ts{20, 20.5, "twenty"}, "3": ts{30, 30.5, "thirty"}, "4": ts{40, 40.5, "forty"}, "5": ts{50, 50.5, "fifty"}},
			"MyFloat",
			"asc",
			[]interface{}{ts{10, 10.5, "ten"}, ts{20, 20.5, "twenty"}, ts{30, 30.5, "thirty"}, ts{40, 40.5, "forty"}, ts{50, 50.5, "fifty"}},
		},
		{
			map[string]ts{"1": ts{10, 10.5, "ten"}, "2": ts{20, 20.5, "twenty"}, "3": ts{30, 30.5, "thirty"}, "4": ts{40, 40.5, "forty"}, "5": ts{50, 50.5, "fifty"}},
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

func TestMarkdownify(t *testing.T) {

	result := Markdownify("Hello **World!**")

	expect := template.HTML("<p>Hello <strong>World!</strong></p>\n")

	if result != expect {
		t.Errorf("Markdownify: got '%s', expected '%s'", result, expect)
	}
}
