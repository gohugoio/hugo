package css

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types/css"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/gohugoio/hugo/resources/resource_transformers/babel"
	"github.com/gohugoio/hugo/resources/resource_transformers/cssjs"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/dartsass"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/gohugoio/hugo/tpl/internal/resourcehelpers"
	"github.com/spf13/cast"
)

const name = "css"

// Namespace provides template functions for the "css" namespace.
type Namespace struct {
	d                 *deps.Deps
	scssClientLibSass *scss.Client
	postcssClient     *cssjs.PostCSSClient
	tailwindcssClient *cssjs.TailwindCSSClient
	babelClient       *babel.Client

	// The Dart Client requires a os/exec process, so  only
	// create it if we really need it.
	// This is mostly to avoid creating one per site build test.
	scssClientDartSassInit sync.Once
	scssClientDartSass     *dartsass.Client
}

// Quoted returns a string that needs to be quoted in CSS.
func (ns *Namespace) Quoted(v any) css.QuotedString {
	s := cast.ToString(v)
	return css.QuotedString(s)
}

// Unquoted returns a string that does not need to be quoted in CSS.
func (ns *Namespace) Unquoted(v any) css.UnquotedString {
	s := cast.ToString(v)
	return css.UnquotedString(s)
}

// PostCSS processes the given Resource with PostCSS.
func (ns *Namespace) PostCSS(args ...any) (resource.Resource, error) {
	if len(args) > 2 {
		return nil, errors.New("must not provide more arguments than resource object and options")
	}

	r, m, err := resourcehelpers.ResolveArgs(args)
	if err != nil {
		return nil, err
	}

	return ns.postcssClient.Process(r, m)
}

// TailwindCSS processes the given Resource with tailwindcss.
func (ns *Namespace) TailwindCSS(args ...any) (resource.Resource, error) {
	if len(args) > 2 {
		return nil, errors.New("must not provide more arguments than resource object and options")
	}

	r, m, err := resourcehelpers.ResolveArgs(args)
	if err != nil {
		return nil, err
	}

	return ns.tailwindcssClient.Process(r, m)
}

// Sass processes the given Resource with SASS.
func (ns *Namespace) Sass(args ...any) (resource.Resource, error) {
	if len(args) > 2 {
		return nil, errors.New("must not provide more arguments than resource object and options")
	}

	const (
		// Transpiler implementation can be controlled from the client by
		// setting the 'transpiler' option.
		// Default is currently 'libsass', but that may change.
		transpilerDart    = "dartsass"
		transpilerLibSass = "libsass"
	)

	var (
		r          resources.ResourceTransformer
		m          map[string]any
		targetPath string
		err        error
		ok         bool
		transpiler = transpilerLibSass
	)

	r, targetPath, ok = resourcehelpers.ResolveIfFirstArgIsString(args)

	if !ok {
		r, m, err = resourcehelpers.ResolveArgs(args)
		if err != nil {
			return nil, err
		}
	}

	if m != nil {
		if t, _, found := maps.LookupEqualFold(m, "transpiler"); found {
			switch t {
			case transpilerDart, transpilerLibSass:
				transpiler = cast.ToString(t)
			default:
				return nil, fmt.Errorf("unsupported transpiler %q; valid values are %q or %q", t, transpilerLibSass, transpilerDart)
			}
		}
	}

	if transpiler == transpilerLibSass {
		var options scss.Options
		if targetPath != "" {
			options.TargetPath = paths.ToSlashTrimLeading(targetPath)
		} else if m != nil {
			options, err = scss.DecodeOptions(m)
			if err != nil {
				return nil, err
			}
		}

		return ns.scssClientLibSass.ToCSS(r, options)
	}

	if m == nil {
		m = make(map[string]any)
	}
	if targetPath != "" {
		m["targetPath"] = targetPath
	}

	client, err := ns.getscssClientDartSass()
	if err != nil {
		return nil, err
	}

	return client.ToCSS(r, m)
}

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		scssClient, err := scss.New(d.BaseFs.Assets, d.ResourceSpec)
		if err != nil {
			panic(err)
		}
		ctx := &Namespace{
			d:                 d,
			scssClientLibSass: scssClient,
			postcssClient:     cssjs.NewPostCSSClient(d.ResourceSpec),
			tailwindcssClient: cssjs.NewTailwindCSSClient(d.ResourceSpec),
			babelClient:       babel.New(d.ResourceSpec),
		}

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Sass,
			[]string{"toCSS"},
			[][2]string{},
		)

		ns.AddMethodMapping(ctx.PostCSS,
			[]string{"postCSS"},
			[][2]string{},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}

func (ns *Namespace) getscssClientDartSass() (*dartsass.Client, error) {
	var err error
	ns.scssClientDartSassInit.Do(func() {
		ns.scssClientDartSass, err = dartsass.New(ns.d.BaseFs.Assets, ns.d.ResourceSpec)
		if err != nil {
			return
		}
		ns.d.BuildClosers.Add(ns.scssClientDartSass)
	})

	return ns.scssClientDartSass, err
}
