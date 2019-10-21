* A _regular_ page is a "post" page or a "content" page.
  * A _leaf bundle_ is a regular page.
* A _list_ page can list _regular_ pages and other _list_ pages. Some
  examples are: homepage, section pages, _taxonomy term_ (`/tags/`) and
  _taxonomy_ (`/tags/foo/`) pages.
  * A _branch bundle_ is a _list_ page.

`.Site.Pages`
: Collection of **all** pages of the site: _regular_ pages,
    sections, taxonomies, etc. -- Superset of everything!

`.Site.RegularPages`
: Collection of only _regular_ pages.

The above `.Site. ..` page collections can be accessed from any scope in
the templates.

Below variables return a collection of pages only from the scope of
the current _list_ page:

`.Pages`
: Collection of _regular_ pages and _only first-level_
    section pages under the current _list_ page.

`.RegularPages`
: Collection of only _regular_ pages under the
    current _list_ page.


Note
: From the scope of _regular_ pages, `.Pages` and
    `.RegularPages` return an empty slice.
