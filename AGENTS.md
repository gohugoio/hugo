
* Brevity is good.
* Assume that the maintainers and readers of the code you write are Go experts:
   * Don't use comments to explain the obvious.
   * Use self-explanatory variable and function names.
   * Use short variable names when the context is clear.
* If you need to add temporary debug printing, use `hdebug.Printf`.[^1]
* Never export symbols that's not needed outside of the package.
* Avoid global state at (almost) all cost.
* This is a project with a long history; assume that a similiar problem has been solved before, look hard for helper functions before creating new ones.
* In tests, almost always write end-to-end integration tests using `hugolib.Test` or one of its siblings. Write unit tests only for isolated utilities.
* In tests, use `qt` matchers (e.g. `b.Assert(err, qt.ErrorMatches, ...)`) instead of raw `if`/`t.Fatal` checks.
* In tests, always use the latest Hugo specification, e.g. for layouts, it's `layouts/page.html` and not `layouts/_default/single.html`, `layouts/list.html` and not `layouts/_default/list.html`
* Never name tests `TestIssue1234`; always give the test function a descriptive name, e.g. `TestDisablePathToLower`, and add any issue reference as a Go doc function comment, e.g. `// See issue 1234.`.
* If you're a security researcher, read @SECURITY.md carefully.
* Brevity is good. This applies to code, comments and commit messages. Don't write a novel.
* Use `./check.sh ./somepackage/...` when iterating.
* Use `./check.sh` when you're done.


[^1]: CI build fail if you forget to remove the debug printing.
