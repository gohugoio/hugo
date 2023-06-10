package external

import (
	"context"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl/internal"
)

// TODO: Pull functions location from config.toml, with this as default
const defaultFunctionsPath = "functions/"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx, err := LoadFunctionFiles(defaultFunctionsPath)
		if err != nil {
			helpers.DistinctErrorLog.Warnf("Unable to load any function files: %v\n", err)
			return nil
		}

		ns := &internal.TemplateFuncsNamespace{
			Name:    "external",
			Context: func(cctx context.Context, args ...any) (any, error) { return ctx, nil },
		}

		ns.AddMethodMapping(ctx.Function, []string{"fn"}, nil)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
