.CurrentSection
: The page's current section. The value can be the page itself if it is a section or the homepage.

.FirstSection
: The page's first section below root, e.g. `/docs`, `/blog` etc.

.InSection $anotherPage
: Whether the given page is in the current section.

.IsAncestor $anotherPage
: Whether the current page is an ancestor of the given page.

.IsDescendant $anotherPage
: Whether the current page is a descendant of the given page.

.Parent
: A section's parent section or a page's section.

.Section
: The [section](/content-management/sections/) this content belongs to. **Note:** For nested sections, this is the first path element in the directory, for example, `/blog/funny/mypost/ => blog`.

.Sections
: The [sections](/content-management/sections/) below this content.
