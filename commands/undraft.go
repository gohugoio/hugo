// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package commands

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/parser"
	jww "github.com/spf13/jwalterweatherman"
)

var undraftCmd = &cobra.Command{
	Use:   "undraft path/to/content",
	Short: "Undraft changes the content's draft status from 'True' to 'False'",
	Long:  `Undraft changes the content's draft status from 'True' to 'False' and updates the date to the current date and time. If the content's draft status is 'False', nothing is done`,
	Run:   Undraft,
}

// Publish publishes the specified content by setting its draft status
// to false and setting its publish date to now. If the specified content is
// not a draft, it will log an error.
func Undraft(cmd *cobra.Command, args []string) {
	InitializeConfig()

	if len(args) < 1 {
		cmd.Usage()
		jww.FATAL.Fatalln("a piece of content needs to be specified")
	}

	location := args[0]
	// open the file
	f, err := os.Open(location)
	if err != nil {
		jww.ERROR.Print(err)
		return
	}

	// get the page from file
	p, err := parser.ReadFrom(f)
	f.Close()
	if err != nil {
		jww.ERROR.Print(err)
		return
	}

	w, err := undraftContent(p)
	if err != nil {
		jww.ERROR.Printf("an error occurred while undrafting %q: %s", location, err)
		return
	}

	f, err = os.OpenFile(location, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		jww.ERROR.Printf("%q not be undrafted due to error opening file to save changes: %q\n", location, err)
		return
	}
	defer f.Close()
	_, err = w.WriteTo(f)
	if err != nil {
		jww.ERROR.Printf("%q not be undrafted due to save error: %q\n", location, err)
	}
	return
}

// undraftContent: if the content is a draft, change it's draft status to
// 'false' and set the date to time.Now(). If the draft status is already
// 'false', don't do anything.
func undraftContent(p parser.Page) (bytes.Buffer, error) {
	var buff bytes.Buffer
	// get the metadata; easiest way to see if it's a draft
	meta, err := p.Metadata()
	if err != nil {
		return buff, err
	}
	// since the metadata was obtainable, we can also get the key/value separator for
	// Front Matter
	fm := p.FrontMatter()
	if fm == nil {
		err := fmt.Errorf("Front Matter was found, nothing was finalized")
		return buff, err
	}

	var isDraft, gotDate bool
	var date string
L:
	for k, v := range meta.(map[string]interface{}) {
		switch k {
		case "draft":
			if !v.(bool) {
				return buff, fmt.Errorf("not a Draft: nothing was done")
			}
			isDraft = true
			if gotDate {
				break L
			}
		case "date":
			date = v.(string) // capture the value to make replacement easier
			gotDate = true
			if isDraft {
				break L
			}
		}
	}

	// if draft wasn't found in FrontMatter, it isn't a draft.
	if !isDraft {
		return buff, fmt.Errorf("not a Draft: nothing was done")
	}

	// get the front matter as bytes and split it into lines
	var lineEnding []byte
	fmLines := bytes.Split(fm, parser.UnixEnding)
	if len(fmLines) == 1 { // if the result is only 1 element, try to split on dos line endings
		fmLines = bytes.Split(fm, parser.DosEnding)
		if len(fmLines) == 1 {
			return buff, fmt.Errorf("unable to split FrontMatter into lines")
		}
		lineEnding = append(lineEnding, parser.DosEnding...)
	} else {
		lineEnding = append(lineEnding, parser.UnixEnding...)
	}

	// Write the front matter lines to the buffer, replacing as necessary
	for _, v := range fmLines {
		pos := bytes.Index(v, []byte("draft"))
		if pos != -1 {
			v = bytes.Replace(v, []byte("true"), []byte("false"), 1)
			goto write
		}
		pos = bytes.Index(v, []byte("date"))
		if pos != -1 { // if date field wasn't found, add it
			v = bytes.Replace(v, []byte(date), []byte(time.Now().Format(time.RFC3339)), 1)
		}
	write:
		buff.Write(v)
		buff.Write(lineEnding)
	}

	// append the actual content
	buff.Write([]byte(p.Content()))

	return buff, nil
}
