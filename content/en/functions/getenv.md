---
title: getenv
description: Returns the value of an environment variable, or an empty string if the environment variable is not set.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
signature: ["os.Getenv VARIABLE", "getenv VARIABLE"]
relatedfuncs: []
---
Examples:

```go-html-template
{{ os.Getenv "HOME" }} --> /home/victor
{{ os.Getenv "USER" }} --> victor
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
{{ os.Getenv "MY_VAR1" }} --> foo
{{ os.Getenv "MY_VAR2" }} --> bar
```

With Hugo v0.91.0 and later, you must explicitly allow access to environment variables. For details, review [Hugo's Security Policy](/about/security-model/#security-policy). By default, environment variables beginning with `HUGO_` are allowed when using the `os.Getenv` function.
