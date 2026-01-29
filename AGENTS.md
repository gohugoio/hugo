* Brevity is good.
* Assume that the maintainers and readers of the code you write are Go experts:
   * Don't use comments to explain the obvious.
   * Use self-explanatory variable and function names.
   * Use short variable names when the context is clear.
* If you need to add temporary debug printing, use `hdebug.Printf`.[^1]
* Never export symbols that's not needed outside of the package.
* Avoid global state at (almost) all cost.
* This is a project with a long history; assume that a similiar problem has been solved before, look hard for helper functions before creating new ones.
* Use `./check.sh ./somepackage/...` when iterating.
* Use `./check.sh` when you're done.


[^1]: CI build fail if you forget to remove the debug printing.
