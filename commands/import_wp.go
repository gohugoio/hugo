// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func init() {
	importCmd.AddCommand(importWordpressCmd)
}

var importWordpressCmd = &cobra.Command{
	Use:   "wordpress",
	Short: "hugo import from WordPress",
	Long: `hugo import from WordPress.

Import from WordPress requires list of files, e.g. ` + "`hugo import wordpress export1.xml export2.xml`.",
	RunE: importFromWordpress,
}

func importFromWordpress(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError(`Import from WordPress requires list of files, e.g. ` + "`hugo import wordpress export1.xml export2.xml`.")
	}

	// WordPress export may contain multiple files
	for _, i := range args[0:] {
		fmt.Println("\tWorking on file: ", i)
		finp, err := os.Open(i) // for read access
		if err != nil {
			fmt.Println("\tCannot open file ", i, "/", err)
		} else {
			wp2hugo(finp)
		}
		finp.Close()
	}

	return nil
}

// Global variables
type attachm_t struct { // almost always images, but can be doc, pdf, etc.
	cnt       int    // counter, starts with 0
	fname     string // file name, i.e., basename
	hugoFname string // fname prefixed with /img/
}

var attachm = make(map[string]*attachm_t) // counts number of occurences of attachments

var wpconfig = make(map[string]string) // goes to config.toml

// WordPress unfortunately inserts HTML codes into bash/C/R code.
// Have to get rid of these again.
var delimList = []struct {
	start string
	stop  string
}{
	{"[googlemaps ", "]"},
	{"[code", "[/code]"},
	{"<pre>", "</pre>"},
	{"$latex ", "$"},
}

func deHTMLify(s string) string {
	// glitch for WordPress foolishly changing "&" to "&amp;", ">" to "&gt;", etc.
	for _, v := range delimList {
		tx0, tx1, lenvstart, lenvstop := 0, 0, len(v.start), len(v.stop)
		for {
			if tx1 = strings.Index(s[tx0:], v.start); tx1 < 0 {
				break
			}
			if tx2 := strings.Index(s[tx0+tx1+lenvstart:], v.stop); tx2 > 0 {
				//fmt.Println("\t\tv =",v,", tx0 =",tx0,", tx1 =",tx1,", tx2 =",tx2,"\n\t\ts =",s[tx1:tx1+70])
				t := strings.Replace(s[tx0:tx0+tx1+tx2], "&amp;", "&", -1)
				t = strings.Replace(t, "&lt;", "<", -1)
				t = strings.Replace(t, "&gt;", ">", -1)
				t = strings.Replace(t, "&quot;", "\"", -1)
				s = s[0:tx0] + t + s[tx0+tx1+tx2:]
			} else {
				u := len(s[tx0+tx1:])
				if u > 40 {
					u = 40
				}
				// Show up to 40 chars of offending string
				fmt.Println("\tClosing tag", v.stop, " in ", s[tx0+tx1:tx0+tx1+u], " not found")
			}
			tx0 += tx1 + lenvstart + lenvstop // + len(t)
		}
	}

	return s
}

// Use regexp to change various WordPress specific codes and
// map them to the equivalent Hugo codes.
// This list should be put into a configuration file.
var replaceList = []struct {
	regx    *regexp.Regexp
	replace string
}{
	// convert [code lang=bash] to ```bash
	{regexp.MustCompile("(\n{0,1})\\[code\\s*lang(uage|)=(\"|)([A-Za-z\\+]+)(\"|)(.*)\\]\\w*\n"), "\n```$4$6\n"},
	// convert [/code] to ```
	{regexp.MustCompile("\n\\[/code\\]\\s*\n"), "\n```\n"},
	// handle https://www.youtube.com/watch?v=wtqfC9v0xB0
	{regexp.MustCompile("\nhttp(.|)://www\\.youtube\\.com/watch\\?v=(\\w+)(&.+|)\n"), "\n{{< youtube $2 >}}\n"},
	// handle [youtube=http://www.youtube.com/watch?v=IA8X1cXFo9oautoplay=0&start=0&end=0]
	{regexp.MustCompile(`\[youtube=http(.|)://www\.youtube\.com/watch\?v=(\w+)(&.+|)\]`), "{{< youtube $2 >}}"},
	// handle [vimeo 199882338]
	{regexp.MustCompile(`\[vimeo (\d\d\d+)\]`), "{{< vimeo $1 >}}"},
	// handle [vimeo https://vimeo.com/167845464]
	{regexp.MustCompile(`\[vimeo http(.|)://vimeo\.com/(\d\d\d+)\]`), "{{< vimeo $2 >}}"},
	// convert <code>, which is used as <pre>, to ```
	{regexp.MustCompile("\n<code>\\s*\n"), "\n```\n"},
	// handle </code> which is used as </pre>
	{regexp.MustCompile("\n</code>\\s*\n"), "\n```\n"},
	// convert <pre> and </pre> to ```
	{regexp.MustCompile("(\n{0,1})<(/|)pre>\\s*(\n{0,1})"), "\n```\n"},
	// convert $latex ...$ to `$...$`, handle multiline matches with (?s)
	{regexp.MustCompile(`\$latex\s+(?s)(.+?)\$`), "${}$1$"},
	// convert [googlemaps ...] to <iframe ...>/<iframe>
	{regexp.MustCompile(`\[googlemaps\s+(.+)\]`), `<iframe src="$1"></iframe>`},
}

func hugofy(bodyStr string) string {
	s := deHTMLify(bodyStr)

	// replace attachments with references to /img/attachment
	for k, v := range attachm {
		if strings.Contains(s, k) {
			v.cnt += 1 // attachment is used in body
			s = strings.Replace(s, k, v.hugoFname, -1)
		}
	}

	// replace hyperlinks from baseURL to root
	// this logic is not perfect if href points to page, instead of post
	s = strings.Replace(s, "<a href=\""+wpconfig["baseURL"], "<a href=\"/post", -1)
	s = strings.Replace(s, "<a href=\""+wpconfig["baseURLalt"], "<a href=\"/post/", -1)

	// finally process all regex
	for _, r := range replaceList {
		s = r.regx.ReplaceAllString(s, r.replace)
	}

	return s
}

// Take a list of strings and return concatenated string of it,
// each entry in quotes, separated by comma.
// http://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
func commaSep(list map[int]string, prefix string) string {
	lineLen := 0
	line := make([]byte, 2048)

	if len(list) <= 0 {
		return ""
	}

	lineLen += copy(line[lineLen:], prefix)
	lineLen += copy(line[lineLen:], " = [")
	sep := ""
	for key := range list {
		lineLen += copy(line[lineLen:], sep)
		lineLen += copy(line[lineLen:], "\"")
		lineLen += copy(line[lineLen:], list[key])
		lineLen += copy(line[lineLen:], "\"")
		sep = ", "
	}
	lineLen += copy(line[lineLen:], "]\n")
	return string(line[0:lineLen])
}

// Write post/page to appropriate directory
func wrtPostFile(frontmatter map[string]string, cats, tags map[int]string, body []byte, bodyLen int) {
	var dirname string
	if frontmatter["wp:status"] == "draft" {
		dirname = filepath.Join("content", "draft")
	} else if frontmatter["wp:status"] == "private" {
		dirname = filepath.Join("content", "private")
	} else if frontmatter["wp:post_type"] == "post" {
		dirname = filepath.Join("content", "post")
	} else if frontmatter["wp:post_type"] == "page" {
		dirname = filepath.Join("content", "page")
	} else {
		return // do not write anything if attachment, or nav_menu_item
	}
	// create directory content/post, content/page, etc.
	if os.MkdirAll(dirname, 0775) != nil {
		fmt.Println("Cannot create directory", dirname)
		os.Exit(11)
	}

	link := frontmatter["link"]
	fname := frontmatter["wp:post_name"]
	// Create directory, e.g., content/post/2015/05/14/
	if tx2 := strings.Index(link, fname); tx2 > 0 {
		if tx1 := strings.Index(link, wpconfig["baseURL"]); tx1 >= 0 {
			//fmt.Println("\t\ttx1 =",tx1,", tx2 =",tx2,", link =",link,", fname =",fname,", baseURL =",wpconfig["baseURL"])
			dirname = filepath.Join(dirname, link[tx1+len(wpconfig["baseURL"])+1:tx2])
			//fmt.Println("\t\tdirname =",dirname)
			if len(dirname) > 0 && os.MkdirAll(dirname, 0775) != nil {
				fmt.Println("Cannot create directory", dirname)
				os.Exit(12)
			}
		}
	}

	// Create file in above directory, e.g., content/post/2015/05/14/my-post.md
	fname = filepath.Join(dirname, fname) + ".md"
	fout, err := os.Create(fname)
	if err != nil {
		fmt.Println("Cannot open ", fname, "for writing")
		os.Exit(13)
	}

	bodyStr := string(body[0:bodyLen])

	catLine := commaSep(cats, "categories")
	tagLine := commaSep(tags, "tags")
	authorLine := ""
	if len(frontmatter["dc:creator"]) > 0 {
		authorLine = "author = \"" + frontmatter["dc:creator"] + "\"\n"
	}
	mathLine := ""
	if strings.Contains(bodyStr, "$latex ") {
		mathLine = "math = true\n"
	}

	// Write frontmatter + body
	w := bufio.NewWriter(fout)
	fmt.Fprintf(w, "+++\n"+
		"date = \"%s\"\n"+
		"title = \"%s\"\n"+
		"draft = \"%t\"\n"+
		"%s"+
		"%s"+
		"%s"+
		"%s"+
		"+++\n\n"+
		"%s\n",
		frontmatter["wp:post_date"],
		frontmatter["title"],
		frontmatter["wp:status"] == "draft",
		catLine,
		tagLine,
		authorLine,
		mathLine,
		hugofy(bodyStr))

	w.Flush()
	fout.Close()
}

func wrtConfigToml(configCats map[int]string) {
	fout, err := os.Create("config.toml")
	if err != nil {
		fmt.Println("Cannot open config.toml for writing")
		os.Exit(14)
	}

	//catLine := commaSep(configCats,"categories");

	w := bufio.NewWriter(fout)
	fmt.Fprintf(w, "\n"+
		"title = \"%s\"\n"+
		"languageCode = \"%s\"\n"+
		"baseURL = \"%s\"\n"+
		"paginate = 20\n"+
		"\n"+
		"[taxonomies]\n"+
		"   tag = \"tags\"\n"+
		"   category = \"categories\"\n"+
		"   archive = \"archives\"\n"+
		"\n"+
		"[params]\n"+
		"   description = \"%s\"\n"+
		"\n\n",
		wpconfig["title"],
		wpconfig["language"],
		wpconfig["baseURL"],
		wpconfig["description"])

	w.Flush()
	fout.Close()
}

// Write file attachm.txt with all used attachments.
// This can be used to wget/curl these attachments from the WordPress server.
// For example:
//     perl -ane '`curl $F[0] -o $F[1]\n`' attachm.txt
func wrtAttachm() {
	fout, err := os.Create("attachm.txt")
	if err != nil {
		fmt.Println("Cannot open attachm.txt for writing")
		os.Exit(15)
	}

	w := bufio.NewWriter(fout)
	for k, v := range attachm {
		if v.cnt > 0 && strings.Contains(k, "https://") {
			fmt.Fprintf(w, "%s\t%s\n", k, v.fname)
		}
	}

	w.Flush()
	fout.Close()
}

func wp2hugo(finp *os.File) {
	scanner := bufio.NewScanner(finp)
	inItem, inImage, empty, inBody := false, false, false, false
	bodyLen, tx1, tx2 := 0, 0, 0
	body := make([]byte, 524288) // 2^19 as a somewhat arbitrary limit on body length
	frontmatter := make(map[string]string)
	configCatCnt := 0
	configCats := make(map[int]string) // later: categories/taxonomies in config.toml
	catCnt := 0
	cats := make(map[int]string) // goes to categories = ["a","b",...]
	tagCnt := 0
	tags := make(map[int]string) // goes to tags = ["u","v",...]

	for scanner.Scan() { // read each line
		//fmt.Println(scanner.Text())
		s := scanner.Text()
		if !inItem && !inImage {
			if strings.Contains(s, "<image>") {
				inImage = true
			} else if tx1 = strings.Index(s, "<title>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</title>"); tx2 >= tx1 {
					wpconfig["title"] = strings.TrimSpace(s[tx1+7 : tx2])
					//fmt.Println("\t\tinItem =",inItem,", config-title =",wpconfig["title"])
				}
			} else if tx1 = strings.Index(s, "<link>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</link>"); tx2 >= tx1 {
					t := s[tx1+6 : tx2]
					wpconfig["baseURL"] = t // usually this will contain https://
					wpconfig["baseURLalt"] = strings.Replace(t, "http://", "https://", -1)
					wpconfig["baseURLalt"] = strings.Replace(t, "https://", "http://", -1)
				}
			} else if tx1 = strings.Index(s, "<description>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</description>"); tx2 >= tx1 {
					wpconfig["description"] = s[tx1+13 : tx2]
				}
			} else if tx1 = strings.Index(s, "<language>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</language>"); tx2 >= tx1 {
					wpconfig["language"] = s[tx1+10 : tx2]
				}
			} else if tx1 = strings.Index(s, "<wp:cat_name><![CDATA["); tx1 >= 0 {
				if tx2 = strings.Index(s, "]]></wp:cat_name>"); tx2 >= tx1 {
					configCats[configCatCnt] = s[tx1+22 : tx2]
					configCatCnt++
				}
			} else if strings.Contains(s, "<item>") {
				// For each new post frontmatter, categories, etc. are cleared
				inItem, empty, inBody, bodyLen = true, false, false, 0
				for key := range frontmatter {
					delete(frontmatter, key)
				}
				for key := range cats {
					delete(cats, key)
				}
				for key := range tags {
					delete(tags, key)
				}
				//fmt.Println("\t\tOpening <item> found: inItem =",inItem)
				continue
			}
		} else if strings.Contains(s, "</image>") {
			inImage = false
		} else if strings.Contains(s, "</item>") {
			//fmt.Println("tile=",frontmatter["title"],"link=",frontmatter["link"],"bodyLen=",bodyLen,"wp:post_name=",frontmatter["wp:post_name"])
			inItem = false
			if !empty && bodyLen > 0 && len(frontmatter["wp:post_name"]) > 0 {
				wrtPostFile(frontmatter, cats, tags, body, bodyLen)
			}
			//fmt.Println("\t\t</item> found: inItem =",inItem)
			continue
		}
		if inItem && !empty {
			if tx1 = strings.Index(s, "<title>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</title>"); tx2 >= tx1 {
					t := strings.TrimSpace(s[tx1+7 : tx2])
					t = strings.Replace(t, "\\", "\\\\", -1) // replace backslash with double backslash
					t = strings.Replace(t, "\"", "\\\"", -1) // replace quote with backslash quote
					frontmatter["title"] = t
				}
			} else if tx1 = strings.Index(s, "<link>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</link>"); tx2 >= tx1 {
					frontmatter["link"] = s[tx1+6 : tx2]
				}
			} else if tx1 = strings.Index(s, "<pubDate>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</pubDate>"); tx2 >= tx1 {
					frontmatter["pubDate"] = s[tx1+9 : tx2]
				}
			} else if tx1 = strings.Index(s, "<dc:creator>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</dc:creator>"); tx2 >= tx1 {
					frontmatter["dc:creator"] = s[tx1+12 : tx2]
				}
			} else if strings.Contains(s, "<wp:post_name/>") && len(frontmatter["title"]) > 0 {
				// convert file name with spaces and special chars to something without that
				t := strings.ToLower(frontmatter["title"])
				nameLen := 0
				name := make([]byte, 256)
				flip := false
				for _, elem := range t {
					if unicode.IsLetter(elem) {
						nameLen += copy(name[nameLen:], string(elem))
						flip = true
					} else if flip {
						nameLen += copy(name[nameLen:], "-")
						flip = false
					}
				}
				if name[nameLen-1] == '-' { // file name ending in '-' looks ugly
					nameLen--
				}
				frontmatter["wp:post_name"] = string(name[0:nameLen])
			} else if tx1 = strings.Index(s, "<wp:post_name>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:post_name>"); tx2 >= tx1 {
					frontmatter["wp:post_name"] = s[tx1+14 : tx2]
				}
			} else if tx1 = strings.Index(s, "<wp:post_date>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:post_date>"); tx2 >= tx1 {
					frontmatter["wp:post_date"] = s[tx1+14 : tx2]
				}
			} else if tx1 = strings.Index(s, "<wp:status>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:status>"); tx2 >= tx1 {
					frontmatter["wp:status"] = s[tx1+11 : tx2] // either: draft, future, inherit, private, publish
				}
			} else if tx1 = strings.Index(s, "<wp:post_id>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:post_id>"); tx2 >= tx1 {
					frontmatter["wp:post_id"] = s[tx1+12 : tx2]
				}
			} else if tx1 = strings.Index(s, "<wp:post_type>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:post_type>"); tx2 >= tx1 {
					frontmatter["wp:post_type"] = s[tx1+14 : tx2] // either: attachment, nav_menu_item, page, post
				}
			} else if tx1 = strings.Index(s, "<wp:attachment_url>"); tx1 >= 0 {
				if tx2 = strings.Index(s, "</wp:attachment_url>"); tx2 >= tx1 {
					url := s[tx1+19 : tx2] // e.g., https://eklausmeier.files.wordpress.com/2013/06/cisco1.png
					//fmt.Println("\t\t",url)
					a := new(attachm_t)
					a.cnt = 0                                     // count number of occurences of this attachment
					a.fname = url[strings.LastIndex(url, "/")+1:] // e.g., cisco1.png
					a.hugoFname = "/img/" + a.fname               // e.g., /img/cisco1.png
					attachm[url] = a
					urlAlt := strings.Replace(url, "https://", "http://", -1)
					attachm[urlAlt] = a
				}
			} else if strings.Contains(s, "<category domain=\"category\"") {
				if tx1 = strings.Index(s, "<![CDATA["); tx1 >= 0 {
					if tx2 = strings.Index(s, "]]></category>"); tx2 >= tx1 {
						cats[catCnt] = s[tx1+9 : tx2]
						catCnt++
					}
				}
			} else if strings.Contains(s, "<category domain=\"post_tag\"") {
				if tx1 = strings.Index(s, "<![CDATA["); tx1 >= 0 {
					if tx2 = strings.Index(s, "]]></category>"); tx2 >= tx1 {
						tags[tagCnt] = s[tx1+9 : tx2]
						tagCnt++
					}
				}
				//} else if strings.Contains(s,"<content:encoded><![CDATA[]]></content:encoded>") {
				//	empty = true
			} else if tx2 = strings.Index(s, "]]></content:encoded>"); tx2 >= 0 {
				if tx1 = strings.Index(s, "<content:encoded><![CDATA["); tx1 < 0 {
					tx1 = 0 // here we accumulate body text from previous lines
				} else {
					tx1 += 26 // content is on single line
				}
				if inBody || !strings.Contains(s, "jpg]]></content:encoded>") {
					bodyLen += copy(body[bodyLen:], s[tx1:tx2])
					bodyLen += copy(body[bodyLen:], "\n")
					inBody = false
				}
			} else if tx1 = strings.Index(s, "<content:encoded><![CDATA["); tx1 >= 0 {
				bodyLen += copy(body[bodyLen:], s[tx1+26:])
				bodyLen += copy(body[bodyLen:], "\n")
				inBody = true
			} else if inBody {
				bodyLen += copy(body[bodyLen:], s)
				bodyLen += copy(body[bodyLen:], "\n")
			}
		}
	}

	wrtConfigToml(configCats) // write config.toml
	wrtAttachm()

	os.MkdirAll("layouts", 0775)
	os.MkdirAll("archetypes", 0775)
	os.MkdirAll("static", 0775)
	os.MkdirAll("data", 0775)
	os.MkdirAll("themes", 0775)
}
