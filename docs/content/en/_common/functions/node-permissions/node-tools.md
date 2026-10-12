---
_comment: Do not remove front matter.
---

Hugo runs the Node.js tool used by this function with a permission model that limits the files the tool can read to the paths in [`security.node.permissions.allowRead`][]. By default, this is the project directory. If the tool needs to read files from other locations, such as configuration files, plugins, or imported files in the `node_modules` directory of a parent directory in a monorepo, add those paths to the list.

[`security.node.permissions.allowRead`]: /configuration/security/#nodepermissionsallowread
