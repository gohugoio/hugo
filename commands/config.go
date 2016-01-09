package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
)

var translationsConfigs []*viper.Viper

func init() {
	translationsConfigs = make([]*viper.Viper, 0)
}

// readInMultilingualConfig reads the configuration, and loads any translations
// sub-documents in the main config file.
func readInMultilingualConfig(cfgFile, source string) error {
	viperSetConfigFile(cfgFile, source)

	// TODO: ignore error in ReadInConfig, simply use that for
	// discovery of the file.  Then split the config file in chunks,
	// identify the chunks and feed to ReadConfig().  This could be
	// replaced if `viper` exposed the File without reading it.
	_ = viper.ReadInConfig()
	filename := viper.ConfigFileUsed()
	cnt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Initial support for YAML only:
	for i, configDocument := range strings.Split(string(cnt), "\n---") {
		if len(strings.TrimSpace(string(configDocument))) == 0 {
			continue
		}
		viper.SetDefaultConfig(viper.New())

		viperSetConfigFile(cfgFile, source)

		err := viper.ReadConfig(bytes.NewBuffer([]byte(configDocument)))
		if err != nil {
			return fmt.Errorf("in config document %d: %s", i+1, err)
		}

		viper.RegisterAlias("indexes", "taxonomies")

		LoadDefaultSettings()

		translationsConfigs = append(translationsConfigs, viper.DefaultConfig())
	}

	if len(translationsConfigs) == 0 {
		return fmt.Errorf("found 0 config blocks in %q", filename)
	}

	viper.SetDefaultConfig(translationsConfigs[0])

	return nil
}

func viperSetConfigFile(cfgFile, source string) {
	viper.SetConfigFile(CfgFile)
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if Source == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(Source)
	}
}

// viperSetAll calls `viper.Set` on the all translations' configs.
func viperSetAll(key string, value interface{}) {
	viper.Set(key, value)
	for _, conf := range translationsConfigs {
		conf.Set(key, value)
	}
}

// viperSetDefaultAll calls `viper.SetDefault` on all translations'
// configs.
func viperSetDefaultAll(key string, value interface{}) {
	viper.SetDefault(key, value)
	for _, conf := range translationsConfigs {
		conf.SetDefault(key, value)
	}
}
