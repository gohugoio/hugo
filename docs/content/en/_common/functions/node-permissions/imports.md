---
_comment: Do not remove front matter.
---

Hugo limits the files that this function can import from outside the `assets` directory to the paths in [`security.node.permissions.allowRead`][]. By default, this is the project directory, including its `node_modules` directory. To import files from other locations, such as the `node_modules` directory of a parent directory in a monorepo, add those paths to the list.

[`security.node.permissions.allowRead`]: /configuration/security/#nodepermissionsallowread
