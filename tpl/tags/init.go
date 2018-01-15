package tags

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const Name = "tags"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		ns := &internal.TemplateFuncsNamespace{
			Name:    Name,
			Context: func(args ...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Script,
			[]string{"script"},
			[][2]string{
				{
					`{{ script "/js/bootstrap.js" nil "bootstrap.js" }}`,
					`<script src="/js/bootstrap-pa0caWEhSlXeU8HObmDSe2t2H1SFH6ZeMwZkYN-moNs.js" integrity="sha256-pa0caWEhSlXeU8HObmDSe2t2H1SFH6ZeMwZkYN+moNs="></script>`,
				},
			},
		)

		ns.AddMethodMapping(ctx.Style,
			[]string{"style"},
			[][2]string{
				{
					`{{ style "/css/bootstrap.css" nil "bootstrap.css" }}`,
					`<link rel="stylesheet" href="/css/bootstrap-pa0caWEhSlXeU8HObmDSe2t2H1SFH6ZeMwZkYN-moNs.css" integrity="sha256-pa0caWEhSlXeU8HObmDSe2t2H1SFH6ZeMwZkYN+moNs="/>`},
			},
		)

		return ns
	}

	internal.AddTemplateFuncsNamespace(f)
}
