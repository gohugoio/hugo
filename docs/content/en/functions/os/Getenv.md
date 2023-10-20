---
title: os.Getenv
linkTitle: getenv
description: Returns the value of an environment variable, or an empty string if the environment variable is not set.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [getenv]
  returnType: string
  signatures: [os.Getenv VARIABLE]
relatedFunctions:
  - os.FileExists
  - os.Getenv
  - os.ReadDir
  - os.ReadFile
  - os.Stat
aliases: [/functions/getenv]
---

Examples:

```go-html-template
{{ os.Getenv "HOME" }} → /home/victor
{{ os.Getenv "USER" }} → victor
```

You can pass values when building your site:

```bash
MY_VAR1=foo MY_VAR2=bar hugo

OR

export MY_VAR1=foo
export MY_VAR2=bar
hugo
```

And then retrieve the values within a template:

```go-html-template
{{ os.Getenv "MY_VAR1" }} → foo
{{ os.Getenv "MY_VAR2" }} → bar
```

With Hugo v0.91.0 and later, you must explicitly allow access to environment variables. For details, review [Hugo's Security Policy](/about/security-model/#security-policy). By default, environment variables beginning with `HUGO_` are allowed when using the `os.Getenv` function.
