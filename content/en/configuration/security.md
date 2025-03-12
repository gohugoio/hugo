---
title: Configure security
linkTitle: Security
description: Configure security.
categories: []
keywords: []
---

Hugo's built-in security policy, which restricts access to `os/exec`, remote communication, and similar operations, is configured via allow lists. By default, access is restricted. If a build attempts to use a feature not included in the allow list, it will fail, providing a detailed message.

This is the default security configuration:

{{< code-toggle config=security />}}

enableInlineShortcodes
: (`bool`) Whether to enable [inline shortcodes]. Default is `false`.

exec.allow
: (`[]string`) A slice of [regular expressions](g) matching the names of external executables that Hugo is allowed to run.

exec.osEnv
: (`[]string`) A slice of [regular expressions](g) matching the names of operating system environment variables that Hugo is allowed to access.

funcs.getenv
: (`[]string`) A slice of [regular expressions](g) matching the names of operating system environment variables that Hugo is allowed to access with the [`os.Getenv`] function.

http.methods
: (`[]string`) A slice of [regular expressions](g) matching the HTTP methods that the [`resources.GetRemote`] function is allowed to use.

http.mediaTypes
: (`[]string`) Applicable to the `resources.GetRemote` function, a slice of [regular expressions](g) matching the `Content-Type` in HTTP responses that Hugo trusts, bypassing file content analysis for media type detection.

http.urls
: (`[]string`) A slice of [regular expressions](g) matching the URLs that the `resources.GetRemote` function is allowed to access.

> [!note]
> Setting an allow list to the string `none` will completely disable the associated feature.

You can also override the site configuration with environment variables. For example, to block `resources.GetRemote` from accessing any URL:

```txt
export HUGO_SECURITY_HTTP_URLS=none
```

Learn more about [using environment variables] to configure your site.

[`os.Getenv`]: /functions/os/getenv
[`resources.GetRemote`]: /functions/resources/getremote
[inline shortcodes]: /content-management/shortcodes/#inline
[using environment variables]: /configuration/introduction/#environment-variables
