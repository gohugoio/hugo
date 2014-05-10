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
