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

`allowContent`
: {{< new-in 0.162.0 />}}
: (`[]string`) A slice of [regular expressions](g) matching the [media type](g) of [content formats](g) allowed in the `content` directory. By default, Hugo denies the following content formats:

  - Emacs Org Mode (media type `text/org`): Hugo renders export blocks and `@@html:...@@` snippets verbatim, which can allow arbitrary JavaScript execution. See [details](https://orgmode.org/manual/Quoting-HTML-tags.html#Quoting-HTML-tags-1).
  - HTML (media type `text/html`): Hugo renders HTML file content verbatim, which can allow arbitrary JavaScript execution.

  See the [classification][] table for a mapping of content formats to media types.

`enableInlineShortcodes`
: (`bool`) Whether to enable [inline shortcodes][]. Default is `false`.

`exec.allow`
: (`[]string`) A slice of [regular expressions](g) matching the names of external executables that Hugo is allowed to run.

`exec.osEnv`
: (`[]string`) A slice of [regular expressions](g) matching the names of operating system environment variables that Hugo is allowed to access.

`funcs.getenv`
: (`[]string`) A slice of [regular expressions](g) matching the names of operating system environment variables that Hugo is allowed to access with the [`os.Getenv`][] function.

`http.methods`
: (`[]string`) A slice of [regular expressions](g) matching the HTTP methods that the [`resources.GetRemote`][] function is allowed to use.

`http.mediaTypes`
: (`[]string`) Applicable to the `resources.GetRemote` function, a slice of [regular expressions](g) matching the `Content-Type` in HTTP responses that Hugo trusts, bypassing file content analysis for media type detection.

`http.proxyFromEnvironment`
: {{< new-in 0.166.0 />}}
: (`bool`) Whether the `resources.GetRemote` function honors the `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` environment variables. Default is `false`. When a proxy is used, Hugo connects to the proxy rather than the destination, so the resolved address validation described under `http.urls` does not apply.

`http.urls`
: (`[]string`) A slice of [regular expressions](g) matching the URLs that the `resources.GetRemote` function is allowed to access.

  The default allowlist denies URLs with an IP address or `localhost` as the host name. In addition, with the default allowlist, Hugo validates the resolved address when connecting and rejects connections to loopback, private, link-local, and other non-public addresses, such as a host name that resolves to a cloud metadata endpoint. This validation is disabled if you override `http.urls`, as the override may intentionally allow access to hosts on your local network, such as a development server.

`node.permissions.disable`
: {{< new-in 0.161.0 />}}
: (`bool`) Whether to disable the Node.js [permission model][]. When `false`, Hugo runs Node tools with the `--permission` flag, restricting their file system and resource access to what is explicitly allowed below. Default is `false`.

`node.permissions.allowAddons`
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node tool names permitted to load native addons. For these tools, Hugo passes the `--allow-addons` flag to `node`.

`node.permissions.allowChildProcess`
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node tool names permitted to spawn child processes. For these tools, Hugo passes the `--allow-child-process` flag to `node`.

`node.permissions.allowRead`
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of file system paths that Node tools, such as those used by [`js.Babel`][], [`css.PostCSS`][], and [`css.TailwindCSS`][], are allowed to read. Hugo passes these paths to `node` with the `--allow-fs-read` flag. Paths are relative to the working directory, where `"."` means the working directory itself. Use `"*"` to allow all paths. Default is `["."]`.

  Hugo also allows reading from the `assets` directories of the project and its modules. Node tools can also read from the `node_modules` directory where they are installed.

  The Node permission model enforces these limits, so they do not apply to Node tools when `node.permissions.disable` is `true`.

  Hugo applies the same limits to files that [`js.Build`][], [`js.Batch`][], [`css.Build`][], and [`css.Sass`][] import from outside the `assets` directory, such as files in `node_modules`. These functions do not run in Node, so Hugo performs this check itself, regardless of the `node.permissions.disable` setting.

  Node follows symbolic links even when they point outside the allowed paths. When the permission model is enabled, Hugo fails the build if an allowed path contains a symbolic link whose target resolves outside the allowed paths. To permit such a link, add its target to the list. Hugo scans each allowed path once, the first time a Node tool runs. If you also grant write access to an allowed path, a Node tool can create links after the scan that Hugo does not detect.

`node.permissions.allowWorker`
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of Node tool names permitted to spawn worker threads. For these tools, Hugo passes the `--allow-worker` flag to `node`.

`node.permissions.allowWrite`
: {{< new-in 0.161.0 />}}
: (`[]string`) A slice of file system paths that Node tools are allowed to write. Hugo passes these paths to `node` with the `--allow-fs-write` flag. Paths are relative to the working directory, where `"."` means the working directory itself. Use `"*"` to allow all paths.

## Negation rules

{{< new-in 0.161.0 />}}

Any pattern in an allowlist can be negated by prefixing it with an exclamation mark (`!`) and one space to turn it into a deny rule. Deny rules take precedence over allow rules. An allowlist composed entirely of deny rules implicitly allows everything it does not deny. An empty allowlist rejects everything.

For example, to allow all URLs except those pointing to `evil.example.org`:

```toml
[security.http]
urls = ['.*', '! ^https?://evil\.example\.org']
```

Setting an allowlist to the string `none` will completely disable the associated feature.

## Environment variables

You can also override your project configuration with environment variables. For example, to block `resources.GetRemote` from accessing any URL:

```txt
export HUGO_SECURITY_HTTP_URLS=none
```

Learn more about [using environment variables][] to configure your site.

[`css.Build`]: /functions/css/build/
[`css.PostCSS`]: /functions/css/postcss/
[`css.Sass`]: /functions/css/sass/
[`css.TailwindCSS`]: /functions/css/tailwindcss/
[`js.Babel`]: /functions/js/babel/
[`js.Batch`]: /functions/js/batch/
[`js.Build`]: /functions/js/build/
[`os.Getenv`]: /functions/os/getenv/
[`resources.GetRemote`]: /functions/resources/getremote/
[classification]: /content-management/formats/#classification
[inline shortcodes]: /content-management/shortcodes/#inline
[permission model]: https://nodejs.org/api/permissions.html#permission-model
[using environment variables]: /configuration/introduction/#environment-variables
