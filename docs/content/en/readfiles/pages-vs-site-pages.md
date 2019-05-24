| Variable             | Current context | Pages included                                                                                                                                                                                           |
|----------------------|-----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `.Site.Pages`        | **any** page    | ALL pages of the site: content, sections, taxonomies, etc. -- Superset of everything!                                                                                                                     |
| `.Site.RegularPages` | **any** page    | Only regular (content) pages -- Subset of `.Site.Pages`                                                                                                                                                   |
| `.Pages`             | _List_ page     | Regular pages under that _list_ page representing the homepage, section, taxonomy term (`/tags`) or taxonomy (`/tags/foo`) page -- Subset of `.Site.Pages` or `.Site.RegularPages`, depending on context. |
| `.Pages`             | _Single_ page   | empty slice                                                                                                                                                                                              |

Note
: In the **home** context (`index.html`), `.Pages` is the same as `.Site.RegularPages`.
