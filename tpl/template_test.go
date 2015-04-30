package tpl

import (
	"errors"
	"io/ioutil"
	"testing"
)

// Test for bugs discovered by https://github.com/dvyukov/go-fuzz
func TestTplGoFuzzReports(t *testing.T) {
	for i, this := range []struct {
		data      string
		expectErr int
	}{
		// Issue #1089
		{"{{apply .C \"first\" }}", 2},
		// Issue #1090
		{"{{ slicestr \"000000\" 10}}", 2}} {
		templ := New()

		d := &Data{
			A: 42,
			B: "foo",
			C: []int{1, 2, 3},
			D: map[int]string{1: "foo", 2: "bar"},
			E: Data1{42, "foo"},
		}

		err := templ.AddTemplate("fuzz", this.data)

		if err != nil && this.expectErr == 0 {
			t.Fatalf("Test %d errored: %s", i, err)
		} else if err == nil && this.expectErr == 1 {
			t.Fatalf("#1 Test %d should have errored", i)
		}

		err = templ.ExecuteTemplate(ioutil.Discard, "fuzz", d)

		if err != nil && this.expectErr == 0 {
			t.Fatalf("Test %d errored: %s", i, err)
		} else if err == nil && this.expectErr == 2 {
			t.Fatalf("#2 Test %d should have errored", i)
		}
	}
}

type Data struct {
	A int
	B string
	C []int
	D map[int]string
	E Data1
}

type Data1 struct {
	A int
	B string
}

func (Data1) Q() string {
	return "foo"
}

func (Data1) W() (string, error) {
	return "foo", nil
}

func (Data1) E() (string, error) {
	return "foo", errors.New("Data.E error")
}

func (Data1) R(v int) (string, error) {
	return "foo", nil
}

func (Data1) T(s string) (string, error) {
	return s, nil
}
