package tplimpl

import (
	"io"
	"iter"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/tpl"
	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
)

func (t *templateNamespace) readTemplateInto(templ *TemplInfo) error {
	if err := func() error {
		meta := templ.Fi.Meta()
		f, err := meta.Open()
		if err != nil {
			return err
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		templ.content = removeLeadingBOM(string(b))
		if !templ.noBaseOf {
			templ.noBaseOf = !needsBaseTemplate(templ.content)
		}
		return nil
	}(); err != nil {
		return err
	}
	return nil
}

// The tweet and twitter shortcodes were deprecated in favor of the x shortcode
// in v0.141.0. We can remove these aliases in v0.155.0 or later.
var embeddedTemplatesAliases = map[string][]string{
	"_shortcodes/twitter.html": {"_shortcodes/tweet.html"},
}

func (s *TemplateStore) parseTemplate(ti *TemplInfo, replace bool) error {
	err := s.tns.doParseTemplate(ti, replace)
	if err != nil {
		return s.addFileContext(ti, "parse of template failed", err)
	}
	return err
}

func (t *templateNamespace) doParseTemplate(ti *TemplInfo, replace bool) error {
	if !ti.noBaseOf || ti.category == CategoryBaseof {
		// Delay parsing until we have the base template.
		return nil
	}
	pi := ti.PathInfo
	name := pi.PathNoLeadingSlash()

	var (
		templ tpl.Template
		err   error
	)

	if ti.D.IsPlainText {
		prototype := t.parseText
		if !replace && prototype.Lookup(name) != nil {
			name += "-" + strconv.FormatUint(t.nameCounter.Add(1), 10)
		}
		templ, err = prototype.New(name).Parse(ti.content)
		if err != nil {
			return err
		}
	} else {
		prototype := t.parseHTML
		if !replace && prototype.Lookup(name) != nil {
			name += "-" + strconv.FormatUint(t.nameCounter.Add(1), 10)
		}
		templ, err = prototype.New(name).Parse(ti.content)
		if err != nil {
			return err
		}

		if ti.subCategory == SubCategoryEmbedded {
			// In Hugo 0.146.0 we moved the internal templates around.
			// For the "_internal/twitter_cards.html" style templates, they
			// were moved to the _partials directory.
			// But we need to make them accessible from the old path for a while.
			if pi.Type() == paths.TypePartial {
				aliasName := strings.TrimPrefix(name, "_partials/")
				aliasName = "_internal/" + aliasName
				_, err = prototype.AddParseTree(aliasName, templ.(*htmltemplate.Template).Tree)
				if err != nil {
					return err
				}
			}

			// This was also possible before Hugo 0.146.0, but this should be deprecated.
			if pi.Type() == paths.TypeShortcode {
				aliasName := strings.TrimPrefix(name, "_shortcodes/")
				aliasName = "_internal/shortcodes/" + aliasName
				_, err = prototype.AddParseTree(aliasName, templ.(*htmltemplate.Template).Tree)
				if err != nil {
					return err
				}
			}
		}

		// Issue #13599.
		if ti.category == CategoryPartial && ti.Fi != nil && ti.Fi.Meta().PathInfo.Section() == "partials" {
			aliasName := strings.TrimPrefix(name, "_")
			if _, err := prototype.AddParseTree(aliasName, templ.(*htmltemplate.Template).Tree); err != nil {
				return err
			}
		}
	}

	ti.Template = templ

	return nil
}

func (t *templateNamespace) applyBaseTemplate(overlay *TemplInfo, base keyTemplateInfo) error {
	tb := &TemplWithBaseApplied{
		Overlay: overlay,
		Base:    base.Info,
	}

	base.Info.overlays = append(base.Info.overlays, overlay)

	var templ tpl.Template
	if overlay.D.IsPlainText {
		tt := texttemplate.Must(t.parseText.Clone()).New(overlay.PathInfo.PathNoLeadingSlash())
		var err error
		tt, err = tt.Parse(base.Info.content)
		if err != nil {
			return err
		}
		tt, err = tt.Parse(overlay.content)
		if err != nil {
			return err
		}
		templ = tt
		t.baseofTextClones = append(t.baseofTextClones, tt)
	} else {
		tt := htmltemplate.Must(t.parseHTML.CloneShallow()).New(overlay.PathInfo.PathNoLeadingSlash())
		var err error
		tt, err = tt.Parse(base.Info.content)
		if err != nil {
			return err
		}
		tt, err = tt.Parse(overlay.content)
		if err != nil {
			return err
		}
		templ = tt

		t.baseofHtmlClones = append(t.baseofHtmlClones, tt)

	}

	tb.Template = &TemplInfo{
		Template: templ,
		base:     base.Info,
		PathInfo: overlay.PathInfo,
		Fi:       overlay.Fi,
		D:        overlay.D,
		noBaseOf: true,
	}

	variants := overlay.baseVariants.Get(base.Key)
	if variants == nil {
		variants = make(map[TemplateDescriptor]*TemplWithBaseApplied)
		overlay.baseVariants.Insert(base.Key, variants)
	}
	variants[base.Info.D] = tb
	return nil
}

func (t *templateNamespace) templatesIn(in tpl.Template) iter.Seq[tpl.Template] {
	return func(yield func(t tpl.Template) bool) {
		switch in := in.(type) {
		case *htmltemplate.Template:
			for t := range in.All() {
				if !yield(t) {
					return
				}
			}

		case *texttemplate.Template:
			for t := range in.All() {
				if !yield(t) {
					return
				}
			}
		}
	}
}

/*


func (t *templateHandler) applyBaseTemplate(overlay, base templateInfo) (tpl.Template, error) {
	if overlay.isText {
		var (
			templ = t.main.getPrototypeText(prototypeCloneIDBaseof).New(overlay.name)
			err   error
		)

		if !base.IsZero() {
			templ, err = templ.Parse(base.template)
			if err != nil {
				return nil, base.errWithFileContext("text: base: parse failed", err)
			}
		}

		templ, err = texttemplate.Must(templ.Clone()).Parse(overlay.template)
		if err != nil {
			return nil, overlay.errWithFileContext("text: overlay: parse failed", err)
		}

		// The extra lookup is a workaround, see
		// * https://github.com/golang/go/issues/16101
		// * https://github.com/gohugoio/hugo/issues/2549
		// templ = templ.Lookup(templ.Name())

		return templ, nil
	}

	var (
		templ = t.main.getPrototypeHTML(prototypeCloneIDBaseof).New(overlay.name)
		err   error
	)

	if !base.IsZero() {
		templ, err = templ.Parse(base.template)
		if err != nil {
			return nil, base.errWithFileContext("html: base: parse failed", err)
		}
	}

	templ, err = htmltemplate.Must(templ.Clone()).Parse(overlay.template)
	if err != nil {
		return nil, overlay.errWithFileContext("html: overlay: parse failed", err)
	}

	// The extra lookup is a workaround, see
	// * https://github.com/golang/go/issues/16101
	// * https://github.com/gohugoio/hugo/issues/2549
	templ = templ.Lookup(templ.Name())

	return templ, err
}

*/

var baseTemplateDefineRe = regexp.MustCompile(`^{{-?\s*define`)

// needsBaseTemplate returns true if the first non-comment template block is a
// define block.
func needsBaseTemplate(templ string) bool {
	idx := -1
	inComment := false
	for i := 0; i < len(templ); {
		if !inComment && strings.HasPrefix(templ[i:], "{{/*") {
			inComment = true
			i += 4
		} else if !inComment && strings.HasPrefix(templ[i:], "{{- /*") {
			inComment = true
			i += 6
		} else if inComment && strings.HasPrefix(templ[i:], "*/}}") {
			inComment = false
			i += 4
		} else if inComment && strings.HasPrefix(templ[i:], "*/ -}}") {
			inComment = false
			i += 6
		} else {
			r, size := utf8.DecodeRuneInString(templ[i:])
			if !inComment {
				if strings.HasPrefix(templ[i:], "{{") {
					idx = i
					break
				} else if !unicode.IsSpace(r) {
					break
				}
			}
			i += size
		}
	}

	if idx == -1 {
		return false
	}

	return baseTemplateDefineRe.MatchString(templ[idx:])
}

func removeLeadingBOM(s string) string {
	const bom = '\ufeff'

	for i, r := range s {
		if i == 0 && r != bom {
			return s
		}
		if i > 0 {
			return s[i:]
		}
	}

	return s
}

type templateNamespace struct {
	parseText     *texttemplate.Template
	parseHTML     *htmltemplate.Template
	prototypeText *texttemplate.Template
	prototypeHTML *htmltemplate.Template

	nameCounter atomic.Uint64

	standaloneText *texttemplate.Template

	baseofTextClones []*texttemplate.Template
	baseofHtmlClones []*htmltemplate.Template
}

func (t *templateNamespace) createPrototypesParse() error {
	if t.prototypeHTML == nil {
		panic("prototypeHTML not set")
	}
	t.parseHTML = htmltemplate.Must(t.prototypeHTML.Clone())
	t.parseText = texttemplate.Must(t.prototypeText.Clone())
	return nil
}

func (t *templateNamespace) createPrototypes(init bool) error {
	if init {
		t.prototypeHTML = htmltemplate.Must(t.parseHTML.Clone())
		t.prototypeText = texttemplate.Must(t.parseText.Clone())
	}

	return nil
}

func newTemplateNamespace(funcs map[string]any) *templateNamespace {
	return &templateNamespace{
		parseHTML:      htmltemplate.New("").Funcs(funcs),
		parseText:      texttemplate.New("").Funcs(funcs),
		standaloneText: texttemplate.New("").Funcs(funcs),
	}
}

func isText(t tpl.Template) bool {
	switch t.(type) {
	case *texttemplate.Template:
		return true
	case *htmltemplate.Template:
		return false
	default:
		panic("unknown template type")
	}
}
