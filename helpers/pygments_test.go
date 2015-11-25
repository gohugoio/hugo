package helpers

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParsePygmentsArgs(t *testing.T) {
	for i, this := range []struct {
		in                 string
		pygmentsStyle      string
		pygmentsUseClasses bool
		expect1            interface{}
	}{
		{"", "foo", true, "encoding=utf8,noclasses=false,style=foo"},
		{"style=boo,noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"Style=boo, noClasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=foo"},
		{"style=boo", "foo", true, "encoding=utf8,noclasses=false,style=boo"},
		{"boo=invalid", "foo", false, false},
		{"style", "foo", false, false},
	} {
		viper.Reset()
		viper.Set("PygmentsStyle", this.pygmentsStyle)
		viper.Set("PygmentsUseClasses", this.pygmentsUseClasses)

		result1, err := parsePygmentsOpts(this.in)
		if b, ok := this.expect1.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parsePygmentArgs didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
				continue
			}
			if result1 != this.expect1 {
				t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result1, this.expect1)
			}

		}
	}
}

func TestParseDefaultPygmentsArgs(t *testing.T) {
	expect := "encoding=utf8,noclasses=false,style=foo"

	for i, this := range []struct {
		in                 string
		pygmentsStyle      interface{}
		pygmentsUseClasses interface{}
		pygmentsOptions    string
	}{
		{"", "foo", true, "style=override,noclasses=override"},
		{"", nil, nil, "style=foo,noclasses=false"},
		{"style=foo,noclasses=false", nil, nil, "style=override,noclasses=override"},
		{"style=foo,noclasses=false", "override", false, "style=override,noclasses=override"},

	} {
		viper.Reset()

		viper.Set("PygmentsOptions", this.pygmentsOptions)

		if s, ok := this.pygmentsStyle.(string); ok {
			viper.Set("PygmentsStyle", s)
		}

		if b, ok := this.pygmentsUseClasses.(bool); ok {
			viper.Set("PygmentsUseClasses", b)
		}

		result, err := parsePygmentsOpts(this.in)
		if err != nil {
			t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
			continue
		}
		if result != expect {
			t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result, expect)
		}
	}
}
