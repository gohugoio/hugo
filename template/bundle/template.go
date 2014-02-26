package bundle

import (
	"errors"
	"github.com/eknkc/amber"
	"github.com/spf13/hugo/helpers"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

func Gt(a interface{}, b interface{}) bool {
	var left, right int64
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, 64)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, 64)
	}

	return left > right
}

// First is exposed to templates, to iterate over the first N items in a
// rangeable list.
func First(limit int, seq interface{}) (interface{}, error) {
	if limit < 1 {
		return nil, errors.New("can't return negative/empty count of items from sequence")
	}

	seqv := reflect.ValueOf(seq)
	// this is better than my first pass; ripped from text/template/exec.go indirect():
	for ; seqv.Kind() == reflect.Ptr || seqv.Kind() == reflect.Interface; seqv = seqv.Elem() {
		if seqv.IsNil() {
			return nil, errors.New("can't iterate over a nil value")
		}
		if seqv.Kind() == reflect.Interface && seqv.NumMethod() > 0 {
			break
		}
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		// okay
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}
	if limit > seqv.Len() {
		limit = seqv.Len()
	}
	return seqv.Slice(0, limit).Interface(), nil
}

func IsSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

func ReturnWhenSet(a interface{}, index int) interface{} {
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Slice:
		if av.Len() > index {

			avv := av.Index(index)
			switch avv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return avv.Int()
			case reflect.String:
				return avv.String()
			}
		}
	}

	return ""
}

func Highlight(in interface{}, lang string) template.HTML {
	var str string
	av := reflect.ValueOf(in)
	switch av.Kind() {
	case reflect.String:
		str = av.String()
	}

	if strings.HasPrefix(strings.TrimSpace(str), "<pre><code>") {
		str = str[strings.Index(str, "<pre><code>")+11:]
	}
	if strings.HasSuffix(strings.TrimSpace(str), "</code></pre>") {
		str = str[:strings.LastIndex(str, "</code></pre>")]
	}
	return template.HTML(helpers.Highlight(html.UnescapeString(str), lang))
}

func SafeHtml(text string) template.HTML {
	return template.HTML(text)
}

type Template interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	Lookup(name string) *template.Template
	Templates() []*template.Template
	New(name string) *template.Template
	LoadTemplates(absPath string)
	AddTemplate(name, tpl string) error
	AddInternalTemplate(prefix, name, tpl string) error
	AddInternalShortcode(name, tpl string) error
}

type templateErr struct {
	name string
	err  error
}

type GoHtmlTemplate struct {
	template.Template
	errors []*templateErr
}

func NewTemplate() Template {
	var templates = &GoHtmlTemplate{
		Template: *template.New(""),
		errors:   make([]*templateErr, 0),
	}

	funcMap := template.FuncMap{
		"urlize":    helpers.Urlize,
		"gt":        Gt,
		"isset":     IsSet,
		"echoParam": ReturnWhenSet,
		"safeHtml":  SafeHtml,
		"first":     First,
		"highlight": Highlight,
		"add":       func(a, b int) int { return a + b },
		"sub":       func(a, b int) int { return a - b },
		"div":       func(a, b int) int { return a / b },
		"mod":       func(a, b int) int { return a % b },
		"mul":       func(a, b int) int { return a * b },
		"modBool":   func(a, b int) bool { return a%b == 0 },
		"lower":     func(a string) string { return strings.ToLower(a) },
		"upper":     func(a string) string { return strings.ToUpper(a) },
		"title":     func(a string) string { return strings.Title(a) },
	}

	templates.Funcs(funcMap)

	templates.LoadEmbedded()
	return templates
}

func (t *GoHtmlTemplate) LoadEmbedded() {
	t.EmbedShortcodes()
}

func (t *GoHtmlTemplate) AddInternalTemplate(prefix, name, tpl string) error {
	return t.AddTemplate("_internal/"+prefix+"/"+name, tpl)
}

func (t *GoHtmlTemplate) AddInternalShortcode(name, content string) error {
	return t.AddInternalTemplate("shortcodes", name, content)
}

func (t *GoHtmlTemplate) AddTemplate(name, tpl string) error {
	_, err := t.New(name).Parse(tpl)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	}
	return err
}

func (t *GoHtmlTemplate) AddTemplateFile(name, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	_, err = t.New(name).Parse(s)
	if err != nil {
		t.errors = append(t.errors, &templateErr{name: name, err: err})
	}
	return err
}

func (t *GoHtmlTemplate) generateTemplateNameFrom(base, path string) string {
	return filepath.ToSlash(path[len(base)+1:])
}

func ignoreDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}

func (t *GoHtmlTemplate) LoadTemplates(absPath string) {
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !fi.IsDir() {
			if ignoreDotFile(path) {
				return nil
			}

			tplName := t.generateTemplateNameFrom(absPath, path)

			if strings.HasSuffix(path, ".amber") {
				compiler := amber.New()
				// Parse the input file
				if err := compiler.ParseFile(path); err != nil {
					return nil
				}

				// note t.New(tplName)
				if _, err := compiler.CompileWithTemplate(t.New(tplName)); err != nil {
					return err
				}

			} else {
				t.AddTemplateFile(tplName, path)
			}
		}
		return nil
	}

	filepath.Walk(absPath, walker)
}
