package hugolib

import (
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

func TestFirst(t *testing.T) {
	for i, this := range []struct {
		count    int
		sequence interface{}
		expect   interface{}
	}{
		{2, []string{"a", "b", "c"}, []string{"a", "b"}},
		{3, []string{"a", "b"}, []string{"a", "b"}},
		{2, []int{100, 200, 300}, []int{100, 200}},
	} {
		results, err := First(this.count, this.sequence)
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(results, this.expect) {
			t.Errorf("[%d] First %d items, got %v but expected %v", i, this.count, results, this.expect)
		}
	}
}

func TestWhere(t *testing.T) {
	type X struct {
		A, B string
	}
	for i, this := range []struct {
		sequence interface{}
		key      interface{}
		match    interface{}
		expect   interface{}
	}{
		{[]map[int]string{{1: "a", 2: "m"}, {1: "c", 2: "d"}, {1: "e", 3: "m"}}, 2, "m", []map[int]string{{1: "a", 2: "m"}}},
		{[]map[string]int{{"a": 1, "b": 2}, {"a": 3, "b": 4}, {"a": 5, "x": 4}}, "b", 4, []map[string]int{{"a": 3, "b": 4}}},
		{[]X{{"a", "b"}, {"c", "d"}, {"e", "f"}}, "B", "f", []X{{"e", "f"}}},
		{[]*map[int]string{&map[int]string{1: "a", 2: "m"}, &map[int]string{1: "c", 2: "d"}, &map[int]string{1: "e", 3: "m"}}, 2, "m", []*map[int]string{&map[int]string{1: "a", 2: "m"}}},
		{[]*X{&X{"a", "b"}, &X{"c", "d"}, &X{"e", "f"}}, "B", "f", []*X{&X{"e", "f"}}},
	} {
		results, err := Where(this.sequence, this.key, this.match)
		if err != nil {
			t.Errorf("[%d] failed: %s", i, err)
			continue
		}
		if !reflect.DeepEqual(results, this.expect) {
			t.Errorf("[%d] Where clause matching %v with %v, got %v but expected %v", i, this.key, this.match, results, this.expect)
		}
	}
}
