---
title: os.Getenv
description: Returns the value of an environment variable, or an empty string if the environment variable is not set.
categories: []
keywords: []
action:
  aliases: [getenv]
  related:
    - functions/os/FileExists
    - functions/os/ReadDir
    - functions/os/ReadFile
    - functions/os/Stat
  returnType: string
  signatures: [os.Getenv VARIABLE]
aliases: [/functions/getenv]
toc: true
---

## Security

By default, when using the `os.Getenv` function Hugo allows access to:

- The `CI` environment variable
- Any environment variable beginning with `HUGO_`

To access other environment variables, adjust your site configuration. For example, to allow access to the `HOME` and `USER` environment variables:

{{< code-toggle file=hugo >}}
[security.funcs]
getenv = ['^HUGO_', '^CI$', '^USER$', '^HOME$']
{{< /code-toggle >}}

Read more about Hugo's [security policy].

[security policy]: /about/security/#security-policy

## Examples

```go-html-template
{{ getenv "HOME" }} → /home/victor
{{ getenv "USER" }} → victor
```

You can pass values when building your site:

```sh
MY_VAR1=foo MY_VAR2=bar hugo

OR

export MY_VAR1=foo
export MY_VAR2=bar
hugo
```

And then retrieve the values within a template:

```go-html-template
{{ getenv "MY_VAR1" }} → foo
{{ getenv "MY_VAR2" }} → bar
```
