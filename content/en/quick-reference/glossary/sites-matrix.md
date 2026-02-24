---
title: sites matrix
---

A _sites matrix_ is a configuration object defined in content front matter or a file mount to precisely control which [_sites_](g) the content should be generated for. When defined in a file mount for templates, it controls which sites the template will be applied to. In Hugo multidimensional content model, the matrix defines the intersection of three dimensions: [_language_](g), [_role_](g), and [_version_](g). The configuration is structured as a map of [_glob slices_](g).

  See also [_sites complements_](g), [front matter: sites](/content-management/front-matter/#sites), [module mounts: sites](/configuration/module/#sites), and [segments: sites](/configuration/segments/#sites).
