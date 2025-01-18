---
title: environment
---

Typically one of `development`, `staging`, or `production`, each environment may exhibit different behavior depending on configuration and template logic. For example, in a production environment you might minify and fingerprint CSS, but that probably doesn't make sense in a development environment.

When running the built-in development server with the `hugo server` command, the environment is set to `development`. When building your site with the `hugo` command, the environment is set to `production`. To override the environment value, use the `--environment` command line flag or the `HUGO_ENVIRONMENT` environment variable.

To determine the current environment within a template, use the [`hugo.Environment`] function.

[`hugo.Environment`]: /functions/hugo/environment/

{{% include "/getting-started/glossary/_link-reference-definitions" %}}
