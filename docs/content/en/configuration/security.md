---
title: Configure security
linkTitle: Security
description: Configure security.
categories: []
keywords: []
---

Hugo's built-in security policy, which restricts access to `os/exec`, remote communication, and similar operations, is configured via allowlists. By default, access is restricted. If a build attempts to use a feature not included in the allowlist, it will fail, providing a detailed message.

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

node.permissions.disable
: {{< new-in 0.161.0 />}}
: (`bool`) Whether to disable the Node.js [permission model]. When `false`, Hugo runs Node.js tools with the `--permission` flag, restricting their file system and resource access to what is explicitly allowed below. Default is `false`.

node.permissions.allowAddons
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node.js tool names permitted to load native addons (`--allow-addons`).

node.permissions.allowChildProcess
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node.js tool names permitted to spawn child processes (`--allow-child-process`).

node.permissions.allowRead
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of file system paths that Node.js tools are allowed to read (`--allow-fs-read`). Paths are relative to the working directory; `"."` means the working directory itself. Use `"*"` to allow all paths.

node.permissions.allowWorker
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node.js tool names permitted to spawn worker threads (`--allow-worker`).

node.permissions.allowWrite
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of file system paths that Node.js tools are allowed to write (`--allow-fs-write`). Paths are relative to the working directory; `"."` means the working directory itself. Use `"*"` to allow all paths.

## Negation rules

{{< new-in 0.161.0 />}}

Any pattern in an allowlist can be negated by prefixing it with an exclamation mark (`!`) and one space to turn it into a deny rule. Deny rules take precedence over allow rules. An allowlist composed entirely of deny rules implicitly allows everything it does not deny. An empty allowlist rejects everything.

For example, to allow all URLs except those pointing to `evil.example.com`:

```toml
[security.http]
urls = ['.*', '! ^https?://evil\.example\.com']
```

Setting an allowlist to the string `none` will completely disable the associated feature.

## Environment variables

You can also override your project configuration with environment variables. For example, to block `resources.GetRemote` from accessing any URL:

```txt
export HUGO_SECURITY_HTTP_URLS=none
```

Learn more about [using environment variables] to configure your site.

[`os.Getenv`]: /functions/os/getenv
[`resources.GetRemote`]: /functions/resources/getremote
[inline shortcodes]: /content-management/shortcodes/#inline
[permission model]: https://nodejs.org/api/permissions.html#permission-model
[using environment variables]: /configuration/introduction/#environment-variables
