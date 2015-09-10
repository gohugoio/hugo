package commands

// TODO Support Mac Encoding (\r)

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/spf13/hugo/parser"
)

var (
	jsonFM      = "{\n \"date\": \"12-04-06\",\n \"title\": \"test json\"\n}"
	jsonDraftFM = "{\n \"draft\": true,\n \"date\": \"12-04-06\",\n \"title\":\"test json\"\n}"
	tomlFM      = "+++\n date= \"12-04-06\"\n title= \"test toml\"\n+++"
	tomlDraftFM = "+++\n draft= true\n date= \"12-04-06\"\n title=\"test toml\"\n+++"
	yamlFM      = "---\n date: \"12-04-06\"\n title: \"test yaml\"\n---"
	yamlDraftFM = "---\n draft: true\n date: \"12-04-06\"\n title: \"test yaml\"\n---"
)

func TestUndraftContent(t *testing.T) {
	tests := []struct {
		fm          string
		expectedErr string
	}{
		{jsonFM, "not a Draft: nothing was done"},
		{jsonDraftFM, ""},
		{tomlFM, "not a Draft: nothing was done"},
		{tomlDraftFM, ""},
		{yamlFM, "not a Draft: nothing was done"},
		{yamlDraftFM, ""},
	}

	for _, test := range tests {
		r := bytes.NewReader([]byte(test.fm))
		p, _ := parser.ReadFrom(r)
		res, err := undraftContent(p)
		if test.expectedErr != "" {
			if err == nil {
				t.Error("Expected error, got none")
				continue
			}
			if err.Error() != test.expectedErr {
				t.Errorf("Expected %q, got %q", test.expectedErr, err)
				continue
			}
		} else {
			r = bytes.NewReader(res.Bytes())
			p, _ = parser.ReadFrom(r)
			meta, err := p.Metadata()
			if err != nil {
				t.Errorf("unexpected error %q", err)
				continue
			}
			for k, v := range meta.(map[string]interface{}) {
				if k == "draft" {
					if v.(bool) {
						t.Errorf("Expected %q to be \"false\", got \"true\"", k)
						continue
					}
				}
				if k == "date" {
					if !strings.HasPrefix(v.(string), time.Now().Format("2006-01-02")) {
						t.Errorf("Expected %v to start with %v", v.(string), time.Now().Format("2006-01-02"))
					}
				}
			}
		}
	}
}
