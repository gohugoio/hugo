---
title: Node.js dependencies
description: How to manage Node dependencies in Hugo Modules.
date: 2026-03-22
categories: []
keywords: []
weight: 40
---

Hugo Modules that need Node packages (e.g. for Tailwind CSS) can declare those dependencies in a standard `package.json` at the module root. Hugo consolidates dependencies from all modules into an [npm workspace], so you only need a single `npm install` at the project level.

[npm workspace]: https://docs.npmjs.com/cli/using-npm/workspaces

## Declaring dependencies

Each Hugo Module declares its Node dependencies in a `package.json` file in its root directory, using the standard `dependencies` and `devDependencies` fields.

> [!note]
> We improved this setup greatly in Hugo [v0.159.0](https://github.com/gohugoio/hugo/releases/tag/v0.159.0), but we kept the old `package.hugo.json` in the search path. Mostly to preserve as much backward compatibility as possible, but it may also be useful in some situations to reserve a separate set of Node dependencies for Hugo.

## Consolidating with `hugo mod npm pack`

Run [`hugo mod npm pack`] to collect Node dependencies from all modules and write them to `packages/hugoautogen/package.json`. Hugo also adds a `workspaces` entry to your project's root `package.json` pointing to this auto-generated package.

The resulting project structure:

```text
project/
├── package.json                      # your project's package.json (updated with workspaces entry)
├── packages/
│   └── hugoautogen/
│       ├── package.json              # auto-generated, contains consolidated module deps
│       └── hugo_packagemeta.json     # metadata and checksums for staleness detection
└── ...
```

> [!note]
In Hugo < v0.159.0 Hugo wrote the dependencies into your project's package.json, so if you have used `hugo mod npm pack` on your project using older Hugo versions, now is the time to do a spring cleaning of your project `package.json` file: Only direct Node dependencies needs to live in this file, all incoming dependencies from imported Hugo Modules gets written to `packages/hugoautogen/package.json`.

When merging, the **topmost version, starting from the project, take precedence**. If a module declares `tailwindcss@4.1` but your project already has `tailwindcss@4.0`, the project version wins and the module dependency is excluded from the generated workspace package.

## Staleness detection

When Hugo detects that the npm dependency configuration has changed in one or more of the modules in use, you will get a warning in the console:

```text
WARN  npm dependencies are out of sync, please run "hugo mod npm pack" (you may also want to run "npm install" after that)
```

This ensures you don't forget to re-run `hugo mod npm pack` after updating module versions.

[`hugo mod npm pack`]: /commands/hugo_mod_npm_pack
