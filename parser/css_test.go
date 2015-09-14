package parser

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

const (
	SAMPLE_CSS = `
* {
  margin: 0;
  padding: 0;
}

body {
    font-family: {{ .mainFont }};
    color: {{ .mainTextColor }};
}

a {
   font-size: {{ .fontSizeLinks }};
}
`

	EXPECTED_SAMPLE_CSS = `
* {
  margin: 0;
  padding: 0;
}

body {
    font-family: Arial, Verdana, sans-serif;
    color: #adcdef;
}

a {
   font-size: 16px;
}
`
)

func TestParse(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("params.css", map[string]interface{}{
		"mainFont":      "Arial, Verdana, sans-serif",
		"mainTextColor": "#adcdef",
		"fontSizeLinks": "16px",
	})

	testFile := "./Testparse.css"

	if err := ioutil.WriteFile(testFile, []byte(SAMPLE_CSS), 0644); err != nil {
		t.Errorf("Unable to create [%s] as temporary file.", testFile)
	}

	parse(testFile)

	b, err := ioutil.ReadFile(testFile)

	if err != nil {
		t.Errorf("Unable to read from [%s].", testFile)
	}

	if string(b) != EXPECTED_SAMPLE_CSS {
		t.Errorf("[%s] wasn't parsed correctly. Expected:%s\n====\nGot:\n====\n%s",
			testFile,
			EXPECTED_SAMPLE_CSS,
			string(b),
		)
	}

	if err := os.Remove(testFile); err != nil {
		t.Errorf("Unable to remove temporary file [%s].", testFile)
	}
}
