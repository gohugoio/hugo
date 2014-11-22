package commands

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// config json
var JSONConfig = []byte(`{
	"params": {
		"DateFormat": "Jan 2 2006"
	}		
}`)

// config toml
var TOMLConfig = []byte(`
[params]
DateFormat =  "Jan 2 2006"
`)

// config yaml
var YAMLConfig = []byte(`
params:
  DateFormat: "Jan 2 2006"
`)

var config map[string]interface{} = make(map[string]interface{})

func TestGetDateFormatJSON(t *testing.T) {
	jsonFile, _ := ioutil.TempFile("", "config.json")
	fname := jsonFile.Name()
	jsonFile.Write(JSONConfig)
	jsonFile.Close()
	viper.SetConfigFile(fname)
	viper.SetConfigType("json")
	viper.ReadInConfig()

	dateFmt := getDateFormat()
	assert.Equal(t, "Jan 2 2006", dateFmt)
}

func TestGetDateFormatTOML(t *testing.T) {
	viper.Reset()
	tomlFile, _ := ioutil.TempFile("", "config.toml")
	fname := tomlFile.Name()
	tomlFile.Write(TOMLConfig)
	tomlFile.Close()
	viper.SetConfigFile(fname)
	viper.SetConfigType("toml")
	viper.ReadInConfig()

	dateFmt := getDateFormat()
	assert.Equal(t, "Jan 2 2006", dateFmt)
}

func TestGetDateFormatYAML(t *testing.T) {
	viper.Reset()
	yamlFile, _ := ioutil.TempFile("", "config.yaml")
	fname := yamlFile.Name()
	yamlFile.Write(YAMLConfig)
	yamlFile.Close()
	viper.SetConfigFile(fname)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	dateFmt := getDateFormat()
	assert.Equal(t, "Jan 2 2006", dateFmt)
}
