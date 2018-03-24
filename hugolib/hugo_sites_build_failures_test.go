package hugolib

import (
	"fmt"
	"testing"
)

// https://github.com/gohugoio/hugo/issues/4526
func TestSiteBuildFailureInvalidPageMetadata(t *testing.T) {
	t.Parallel()

	validContentFile := `
---
title = "This is good"
---

Some content.
`

	invalidContentFile := `
---
title = "PDF EPUB: Anne Bradstreet: Poems "The Prologue Summary And Analysis EBook Full Text  "
---

Some content.
`

	var contentFiles []string
	for i := 0; i <= 30; i++ {
		name := fmt.Sprintf("valid%d.md", i)
		contentFiles = append(contentFiles, name, validContentFile)
		if i%5 == 0 {
			name = fmt.Sprintf("invalid%d.md", i)
			contentFiles = append(contentFiles, name, invalidContentFile)
		}
	}

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithContent(contentFiles...)
	b.CreateSites().BuildFail(BuildCfg{})

}
