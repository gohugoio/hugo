.CurrentSection
: The page's current section. The value can be the page itself if it is a section or the homepage.

.InSection $anotherPage
: Whether the given page is in the current section. Note that this will always return false for pages that are not either regular, home or section pages.

.IsAncestor $anotherPage
: Whether the current page is an ancestor of the given page. Note that this method is not relevant for taxonomy lists and taxonomy terms pages.

.IsDescendant $anotherPage
: Whether the current page is a descendant of the given page. Note that this method is not relevant for taxonomy lists and taxonomy terms pages.

.Parent
: A section's parent section or a page's section.

.Section
: The [section](/content-management/sections/) this content belongs to. **Note:** For nested sections, this is the first path element in the directory, for example, `/blog/funny/mypost/ => blog`.

.Sections
: The [sections](/content-management/sections/) below this content.
