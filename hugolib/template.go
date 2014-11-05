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
	"github.com/spf13/cast"
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

func Intersect(l1, l2 interface{}) (interface{}, error) {

	if l1 == nil || l2 == nil {
		return make([]interface{}, 0), nil
	}

	l1v := reflect.ValueOf(l1)
	l2v := reflect.ValueOf(l2)

	switch l1v.Kind() {
	case reflect.Array, reflect.Slice:
		switch l2v.Kind() {
		case reflect.Array, reflect.Slice:
			r := reflect.MakeSlice(l1v.Type(), 0, 0)
			for i := 0; i < l1v.Len(); i++ {
				l1vv := l1v.Index(i)
				for j := 0; j < l2v.Len(); j++ {
					l2vv := l2v.Index(j)
					switch l1vv.Kind() {
					case reflect.String:
						if l1vv.Type() == l2vv.Type() && l1vv.String() == l2vv.String() && !In(r, l2vv) {
							r = reflect.Append(r, l2vv)
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						switch l2vv.Kind() {
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							if l1vv.Int() == l2vv.Int() && !In(r, l2vv) {
								r = reflect.Append(r, l2vv)
							}
						}
					case reflect.Float32, reflect.Float64:
						switch l2vv.Kind() {
						case reflect.Float32, reflect.Float64:
							if l1vv.Float() == l2vv.Float() && !In(r, l2vv) {
								r = reflect.Append(r, l2vv)
							}
						}
					}
				}
			}
			return r.Interface(), nil
		default:
			return nil, errors.New("can't iterate over " + reflect.ValueOf(l2).Type().String())
		}
	default:
		return nil, errors.New("can't iterate over " + reflect.ValueOf(l1).Type().String())
	}
}

func In(l interface{}, v interface{}) bool {
	lv := reflect.ValueOf(l)
	vv := reflect.ValueOf(v)

	switch lv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < lv.Len(); i++ {
			lvv := lv.Index(i)
			switch lvv.Kind() {
			case reflect.String:
				if vv.Type() == lvv.Type() && vv.String() == lvv.String() {
					return true
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch vv.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if vv.Int() == lvv.Int() {
						return true
					}
				}
			case reflect.Float32, reflect.Float64:
				switch vv.Kind() {
				case reflect.Float32, reflect.Float64:
					if vv.Float() == lvv.Float() {
						return true
					}
				}
			}
		}
	case reflect.String:
		if vv.Type() == lv.Type() && strings.Contains(lv.String(), vv.String()) {
			return true
		}
	}
	return false
}

// First is exposed to templates, to iterate over the first N items in a
// rangeable list.
func First(limit interface{}, seq interface{}) (interface{}, error) {

	limitv, err := cast.ToIntE(limit)

	if err != nil {
		return nil, err
	}

	if limitv < 1 {
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
	if limitv > seqv.Len() {
		limitv = seqv.Len()
	}
	return seqv.Slice(0, limitv).Interface(), nil
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

	return template.HTML(helpers.Highlight(html.UnescapeString(str), lang))
}

func SafeHtml(text string) template.HTML {
	return template.HTML(text)
}

func doArithmetic(a, b interface{}, op rune) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64
	var af, bf float64
	var au, bu uint64
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
			if ai >= 0 {
				au = uint64(ai)
				ai = 0
			} else {
				bi = int64(bu) // may overflow
				bu = 0
			}
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		au = av.Uint()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
			if bi >= 0 {
				bu = uint64(bi)
				bi = 0
			} else {
				ai = int64(au) // may overflow
				au = 0
			}
		case reflect.Float32, reflect.Float64:
			af = float64(au) // may overflow
			au = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bu = bv.Uint()
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.String:
		as := av.String()
		if bv.Kind() == reflect.String && op == '+' {
			bs := bv.String()
			return as + bs, nil
		} else {
			return nil, errors.New("Can't apply the operator to the values")
		}
	default:
		return nil, errors.New("Can't apply the operator to the values")
	}

	switch op {
	case '+':
		if ai != 0 || bi != 0 {
			return ai + bi, nil
		} else if af != 0 || bf != 0 {
			return af + bf, nil
		} else if au != 0 || bu != 0 {
			return au + bu, nil
		} else {
			return 0, nil
		}
	case '-':
		if ai != 0 || bi != 0 {
			return ai - bi, nil
		} else if af != 0 || bf != 0 {
			return af - bf, nil
		} else if au != 0 || bu != 0 {
			return au - bu, nil
		} else {
			return 0, nil
		}
	case '*':
		if ai != 0 || bi != 0 {
			return ai * bi, nil
		} else if af != 0 || bf != 0 {
			return af * bf, nil
		} else if au != 0 || bu != 0 {
			return au * bu, nil
		} else {
			return 0, nil
		}
	case '/':
		if bi != 0 {
			return ai / bi, nil
		} else if bf != 0 {
			return af / bf, nil
		} else if bu != 0 {
			return au / bu, nil
		} else {
			return nil, errors.New("Can't divide the value by 0")
		}
	default:
		return nil, errors.New("There is no such an operation")
	}
}

func Mod(a, b interface{}) (int64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var ai, bi int64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bi = bv.Int()
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("The number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

func ModBool(a, b interface{}) (bool, error) {
	res, err := Mod(a, b)
	if err != nil {
		return false, err
	}
	return res == int64(0), nil
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
		"in":          In,
		"intersect":   Intersect,
		"isset":       IsSet,
		"echoParam":   ReturnWhenSet,
		"safeHtml":    SafeHtml,
		"first":       First,
		"where":       Where,
		"highlight":   Highlight,
		"add":         func(a, b interface{}) (interface{}, error) { return doArithmetic(a, b, '+') },
		"sub":         func(a, b interface{}) (interface{}, error) { return doArithmetic(a, b, '-') },
		"div":         func(a, b interface{}) (interface{}, error) { return doArithmetic(a, b, '/') },
		"mod":         Mod,
		"mul":         func(a, b interface{}) (interface{}, error) { return doArithmetic(a, b, '*') },
		"modBool":     ModBool,
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
			err := localTemplates.ExecuteTemplate(buffer, layout, context)
			if err != nil {
				jww.ERROR.Println(err, "in", layout)
			}
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
	// get the suffix and switch on that
	ext := filepath.Ext(path)
	switch ext {
	case ".amber":
		compiler := amber.New()
		// Parse the input file
		if err := compiler.ParseFile(path); err != nil {
			return nil
		}

		if _, err := compiler.CompileWithTemplate(t.New(name)); err != nil {
			return err
		}
	default:
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return t.AddTemplate(name, string(b))
	}

	return nil

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

			t.AddTemplateFile(tplName, path)

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
