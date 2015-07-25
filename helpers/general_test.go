package helpers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestGuessType(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect string
	}{
		{"md", "markdown"},
		{"markdown", "markdown"},
		{"mdown", "markdown"},
		{"asciidoc", "asciidoc"},
		{"adoc", "asciidoc"},
		{"ad", "asciidoc"},
		{"rst", "rst"},
		{"mmark", "mmark"},
		{"html", "html"},
		{"htm", "html"},
		{"excel", "unknown"},
	} {
		result := GuessType(this.in)
		if result != this.expect {
			t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
		}
	}
}

func TestFirstUpper(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect string
	}{
		{"foo", "Foo"},
		{"foo bar", "Foo bar"},
		{"Foo Bar", "Foo Bar"},
		{"", ""},
		{"å", "Å"},
	} {
		result := FirstUpper(this.in)
		if result != this.expect {
			t.Errorf("[%d] got %s but expected %s", i, result, this.expect)
		}
	}
}

func TestBytesToReader(t *testing.T) {
	asBytes := ReaderToBytes(strings.NewReader("Hello World!"))
	asReader := BytesToReader(asBytes)
	assert.Equal(t, []byte("Hello World!"), asBytes)
	assert.Equal(t, asBytes, ReaderToBytes(asReader))
}

func TestStringToReader(t *testing.T) {
	asString := ReaderToString(strings.NewReader("Hello World!"))
	assert.Equal(t, "Hello World!", asString)
	asReader := StringToReader(asString)
	assert.Equal(t, asString, ReaderToString(asReader))
}

var containsTestText = (`На берегу пустынных волн
Стоял он, дум великих полн,
И вдаль глядел. Пред ним широко
Река неслася; бедный чёлн
По ней стремился одиноко.
По мшистым, топким берегам
Чернели избы здесь и там,
Приют убогого чухонца;
И лес, неведомый лучам
В тумане спрятанного солнца,
Кругом шумел.

Τη γλώσσα μου έδωσαν ελληνική
το σπίτι φτωχικό στις αμμουδιές του Ομήρου.
Μονάχη έγνοια η γλώσσα μου στις αμμουδιές του Ομήρου.

από το Άξιον Εστί
του Οδυσσέα Ελύτη

Sîne klâwen durh die wolken sint geslagen,
er stîget ûf mit grôzer kraft,
ich sih in grâwen tägelîch als er wil tagen,
den tac, der im geselleschaft
erwenden wil, dem werden man,
den ich mit sorgen în verliez.
ich bringe in hinnen, ob ich kan.
sîn vil manegiu tugent michz leisten hiez.
`)

var containsBenchTestData = []struct {
	v1     string
	v2     []byte
	expect bool
}{
	{"abc", []byte("a"), true},
	{"abc", []byte("b"), true},
	{"abcdefg", []byte("efg"), true},
	{"abc", []byte("d"), false},
	{containsTestText, []byte("стремился"), true},
	{containsTestText, []byte(containsTestText[10:80]), true},
	{containsTestText, []byte(containsTestText[100:111]), true},
	{containsTestText, []byte(containsTestText[len(containsTestText)-100 : len(containsTestText)-10]), true},
	{containsTestText, []byte(containsTestText[len(containsTestText)-20:]), true},
	{containsTestText, []byte("notfound"), false},
}

// some corner cases
var containsAdditionalTestData = []struct {
	v1     string
	v2     []byte
	expect bool
}{
	{"", nil, false},
	{"", []byte("a"), false},
	{"a", []byte(""), false},
	{"", []byte(""), false},
}

func TestReaderContains(t *testing.T) {
	for i, this := range append(containsBenchTestData, containsAdditionalTestData...) {
		result := ReaderContains(StringToReader(this.v1), this.v2)
		if result != this.expect {
			t.Errorf("[%d] got %t but expected %t", i, result, this.expect)
		}
	}

	assert.False(t, ReaderContains(nil, []byte("a")))
	assert.False(t, ReaderContains(nil, nil))
}

func BenchmarkReaderContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, this := range containsBenchTestData {
			result := ReaderContains(StringToReader(this.v1), this.v2)
			if result != this.expect {
				b.Errorf("[%d] got %t but expected %t", i, result, this.expect)
			}
		}
	}
}

// kept to compare the above
func _BenchmarkReaderContains(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, this := range containsBenchTestData {
			bs, err := ioutil.ReadAll(StringToReader(this.v1))
			if err != nil {
				b.Fatalf("Failed %s", err)
			}
			result := bytes.Contains(bs, this.v2)
			if result != this.expect {
				b.Errorf("[%d] got %t but expected %t", i, result, this.expect)
			}
		}
	}
}

func TestFindAvailablePort(t *testing.T) {
	addr, err := FindAvailablePort()
	assert.Nil(t, err)
	assert.NotNil(t, addr)
	assert.True(t, addr.Port > 0)
}

func TestInStringArrayCaseSensitive(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}
	array := []string{
		"Albert",
		"Joe",
		"Francis",
	}
	data := []test{
		{"Albert", true},
		{"ALBERT", false},
	}
	for i, in := range data {
		output := InStringArray(array, in.input)
		if output != in.expected {
			t.Errorf("[%d] got %t but expected %t", i, output, in.expected)
		}
	}
}

func TestSliceToLowerStable(t *testing.T) {
	input := []string{
		"New York",
		"BARCELONA",
		"COffEE",
		"FLOWer",
		"CanDY",
	}

	output := SliceToLower(input)

	for i, e := range output {
		if e != strings.ToLower(input[i]) {
			t.Errorf("Expected %s. Found %s.", strings.ToLower(input[i]), e)
		}
	}
}

func TestSliceToLowerNil(t *testing.T) {
	var input []string

	output := SliceToLower(input)

	if output != nil {
		t.Errorf("Expected nil to be returned. Had %s.", output)
	}
}

func TestSliceToLowerNonDestructive(t *testing.T) {
	input := []string{
		"New York",
		"BARCELONA",
		"COffEE",
		"FLOWer",
		"CanDY",
	}

	// This assignment actually copies the content
	// of input into a new object.
	// Otherwise, the test would not make sense...
	inputCopy := input

	SliceToLower(input)

	for i, e := range input {
		if e != inputCopy[i] {
			t.Errorf("TestSliceToLowerNonDestructive failed. Expected element #%d of input slice to be %s. Found %s.", i, inputCopy[i], input[i])
		}
	}
}

// Just make sure there is no error for empty-like strings
func TestMd5StringEmpty(t *testing.T) {
	inputs := []string{"", " ", "   "}

	for _, in := range inputs {
		Md5String(in)
	}
}

func TestSeq(t *testing.T) {
	for i, this := range []struct {
		in     []interface{}
		expect interface{}
	}{
		{[]interface{}{-2, 5}, []int{-2, -1, 0, 1, 2, 3, 4, 5}},
		{[]interface{}{1, 2, 4}, []int{1, 3}},
		{[]interface{}{1}, []int{1}},
		{[]interface{}{3}, []int{1, 2, 3}},
		{[]interface{}{3.2}, []int{1, 2, 3}},
		{[]interface{}{0}, []int{}},
		{[]interface{}{-1}, []int{-1}},
		{[]interface{}{-3}, []int{-1, -2, -3}},
		{[]interface{}{3, -2}, []int{3, 2, 1, 0, -1, -2}},
		{[]interface{}{6, -2, 2}, []int{6, 4, 2}},
		{[]interface{}{1, 0, 2}, false},
		{[]interface{}{1, -1, 2}, false},
		{[]interface{}{2, 1, 1}, false},
		{[]interface{}{2, 1, 1, 1}, false},
		{[]interface{}{2001}, false},
		{[]interface{}{}, false},
		// TODO(bep) {[]interface{}{t}, false},
		{nil, false},
	} {

		result, err := Seq(this.in...)

		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] TestSeq didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, this.expect) {
				t.Errorf("[%d] TestSeq got %v but expected %v", i, result, this.expect)
			}
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
		result, err := DoArithmetic(this.a, this.b, this.op)
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
