package parser

import (
	"os"
	p "path"
	"path/filepath"
	"text/template"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// parse parses the file of the given path as templates,
// injects variables of the config safes the modified file.
func parse(path string) {
	t, err := template.ParseFiles(path)

	if err != nil {
		jww.WARN.Println(err)
		return
	}

	// Overwrite the original stylesheet
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		jww.WARN.Println("Can't overwrite the CSS file [%s]: ", path, err)
		return
	}

	params := viper.GetStringMap("params.css")

	if len(params) == 0 {
		jww.ERROR.Println("Can't find variables to insert them into", path)
	}

	err = t.Execute(f, viper.GetStringMap("params.css"))

	if err != nil {
		jww.WARN.Println("Can't interpret the CSS file [%s]:", path, err)
		return
	}
}

// filerFiles checks if the paths found by ParseCSSFiles pointing
// to a .css file. It returns always nil. Otherwise would this stop
// the file search in ParseCSSFiles.
func filterFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		jww.WARN.Println(err)
	}

	// If the path points to a .css file
	if !info.IsDir() && p.Ext(path) == ".css" {
		parse(path)
	}

	return nil
}

// ParseCSSFiles searches inside the folder defined by publishDir and it's
// subfodlers and uses filerFiles to find actual .css files.
func ParseCSSFiles() {
	err := filepath.Walk(viper.GetString("PublishDir"), filterFiles)

	if err != nil {
		jww.ERROR.Println("Can't parse CSS files:", err)
	}
}
