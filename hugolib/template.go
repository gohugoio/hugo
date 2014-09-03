package hugolib

import (
	"bytes"
	"errors"
	"html"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/eknkc/amber"
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
)

var localTemplates *template.Template

func Eq(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

func Ne(x, y interface{}) bool {
	return !Eq(x, y)
}

func Ge(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left >= right
}

func Gt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left > right
}

func Le(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left <= right
}

func Lt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left < right
}

func compareGetFloat(a interface{}, b interface{}) (float64, float64) {
	var left, right float64
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = float64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = float64(av.Int())
	case reflect.Float32, reflect.Float64:
		left = av.Float()
	case reflect.String:
		left, _ = strconv.ParseFloat(av.String(), 64)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = float64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = float64(bv.Int())
	case reflect.Float32, reflect.Float64:
		right = bv.Float()
	case reflect.String:
		right, _ = strconv.ParseFloat(bv.String(), 64)
	}

	return left, right
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

func Where(seq, key, match interface{}) (interface{}, error) {
	seqv := reflect.ValueOf(seq)
	kv := reflect.ValueOf(key)
	mv := reflect.ValueOf(match)

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
	case reflect.Array, reflect.Slice:
		r := reflect.MakeSlice(seqv.Type(), 0, 0)
		for i := 0; i < seqv.Len(); i++ {
			var vvv reflect.Value
			vv := seqv.Index(i)
			switch vv.Kind() {
			case reflect.Map:
				if kv.Type() == vv.Type().Key() && vv.MapIndex(kv).IsValid() {
					vvv = vv.MapIndex(kv)
				}
			case reflect.Struct:
				if kv.Kind() == reflect.String && vv.FieldByName(kv.String()).IsValid() {
					vvv = vv.FieldByName(kv.String())
				}
			case reflect.Ptr:
				if !vv.IsNil() {
					ev := vv.Elem()
					switch ev.Kind() {
					case reflect.Map:
						if kv.Type() == ev.Type().Key() && ev.MapIndex(kv).IsValid() {
							vvv = ev.MapIndex(kv)
						}
					case reflect.Struct:
						if kv.Kind() == reflect.String && ev.FieldByName(kv.String()).IsValid() {
							vvv = ev.FieldByName(kv.String())
						}
					}
				}
			}

			if vvv.IsValid() && mv.Type() == vvv.Type() {
				switch mv.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if mv.Int() == vvv.Int() {
						r = reflect.Append(r, vv)
					}
				case reflect.String:
					if mv.String() == vvv.String() {
						r = reflect.Append(r, vv)
					}
				}
			}
		}
		return r.Interface(), nil
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(seq).Type().String())
	}
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
	LoadTemplatesWithPrefix(absPath, prefix string)
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

	localTemplates = &templates.Template

	funcMap := template.FuncMap{
		"urlize":      helpers.Urlize,
		"sanitizeurl": helpers.SanitizeUrl,
		"eq":          Eq,
		"ne":          Ne,
		"gt":          Gt,
		"ge":          Ge,
		"lt":          Lt,
		"le":          Le,
		"isset":       IsSet,
		"echoParam":   ReturnWhenSet,
		"safeHtml":    SafeHtml,
		"first":       First,
		"where":       Where,
		"highlight":   Highlight,
		"add":         func(a, b int) int { return a + b },
		"sub":         func(a, b int) int { return a - b },
		"div":         func(a, b int) int { return a / b },
		"mod":         func(a, b int) int { return a % b },
		"mul":         func(a, b int) int { return a * b },
		"modBool":     func(a, b int) bool { return a%b == 0 },
		"lower":       func(a string) string { return strings.ToLower(a) },
		"upper":       func(a string) string { return strings.ToUpper(a) },
		"title":       func(a string) string { return strings.Title(a) },
		"partial":     Partial,
	}

	templates.Funcs(funcMap)

	templates.LoadEmbedded()
	return templates
}

func Partial(name string, context_list ...interface{}) template.HTML {
	if strings.HasPrefix("partials/", name) {
		name = name[8:]
	}
	var context interface{}

	if len(context_list) == 0 {
		context = nil
	} else {
		context = context_list[0]
	}
	return ExecuteTemplateToHTML(context, "partials/"+name, "theme/partials/"+name)
}

func ExecuteTemplate(context interface{}, layouts ...string) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	worked := false
	for _, layout := range layouts {
		if localTemplates.Lookup(layout) != nil {
			localTemplates.ExecuteTemplate(buffer, layout, context)
			worked = true
			break
		}
	}
	if !worked {
		jww.ERROR.Println("Unable to render", layouts)
		jww.ERROR.Println("Expecting to find a template in either the theme/layouts or /layouts in one of the following relative locations", layouts)
	}

	return buffer
}

func ExecuteTemplateToHTML(context interface{}, layouts ...string) template.HTML {
	b := ExecuteTemplate(context, layouts...)
	return template.HTML(string(b.Bytes()))
}

func (t *GoHtmlTemplate) LoadEmbedded() {
	t.EmbedShortcodes()
	t.EmbedTemplates()
}

func (t *GoHtmlTemplate) AddInternalTemplate(prefix, name, tpl string) error {
	if prefix != "" {
		return t.AddTemplate("_internal/"+prefix+"/"+name, tpl)
	} else {
		return t.AddTemplate("_internal/"+name, tpl)
	}
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
	return t.AddTemplate(name, string(b))
}

func (t *GoHtmlTemplate) generateTemplateNameFrom(base, path string) string {
	return filepath.ToSlash(path[len(base)+1:])
}

func ignoreDotFile(path string) bool {
	return filepath.Base(path)[0] == '.'
}

func (t *GoHtmlTemplate) loadTemplates(absPath string, prefix string) {
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !fi.IsDir() {
			if ignoreDotFile(path) {
				return nil
			}

			tplName := t.generateTemplateNameFrom(absPath, path)

			if prefix != "" {
				tplName = strings.Trim(prefix, "/") + "/" + tplName
			}

			// TODO move this into the AddTemplateFile function
			if strings.HasSuffix(path, ".amber") {
				compiler := amber.New()
				// Parse the input file
				if err := compiler.ParseFile(path); err != nil {
					return nil
				}

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

func (t *GoHtmlTemplate) LoadTemplatesWithPrefix(absPath string, prefix string) {
	t.loadTemplates(absPath, prefix)
}

func (t *GoHtmlTemplate) LoadTemplates(absPath string) {
	t.loadTemplates(absPath, "")
}
