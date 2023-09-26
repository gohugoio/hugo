---
title: singularize
description: Converts a word according to a set of common English singularization rules.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: inflect
relatedFuncs:
  - inflect.Humanize
  - inflect.Pluralize
  - inflect.Singularize
signature:
  - inflect.Singularize INPUT
  - singularize INPUT
---

`{{ "cats" | singularize }}` → "cat"

See also the `.Data.Singular` [taxonomy variable](/variables/taxonomy/) for singularizing taxonomy names.
