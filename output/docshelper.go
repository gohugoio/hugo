package output

import (

	//	"fmt"

	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() docshelper.DocProvider {
		return docshelper.DocProvider{
			"output": map[string]any{
				// TODO(bep), maybe revisit this later, but I hope this isn't needed.
				// "layouts": createLayoutExamples(),
				"layouts": map[string]any{},
			},
		}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}
