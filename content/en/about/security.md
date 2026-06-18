---
title: Security model
linkTitle: Security 
description: A summary of Hugo's security model.
categories: []
keywords: []
weight: 30
aliases: [/about/security-model/]
---

## Security Boundaries

- The templates inside `layouts` are trusted.
- The assets inside `archetypes`, `assets`, `resources`, `data`, `i18n` and `static` are trusted.
- The content and the content produced by [content adapters][] inside `content` is not trusted. The one exception here is if [inline shortcodes][] is enabled. Note that for content adapters, this is scoped to the result of the adapter.
- The development server, `hugo server`, and its livereload script is trusted and meant for _local_ development only.

## Runtime security

Hugo generates static websites, meaning the final output runs directly in the browser and interacts with any integrated APIs. However, during development and site building, the `hugo` executable itself is the runtime environment.

Securing a runtime is a complex task. Hugo addresses this through a robust sandboxing approach and a strict security policy with default protections. Key features include:

- Virtual file system: Hugo employs a virtual file system, limiting file access. Only the main project, not external components, can access files or directories outside the project root.
- Read-Only access: User-defined components have read-only access to the file system, preventing unintended modifications.
- Controlled external binaries: While Hugo utilizes external binaries for features like Asciidoctor support, these are strictly predefined with specific flags and are disabled by default. The [security policy][] details these limitations.
- No arbitrary commands: To mitigate risks, Hugo intentionally avoids implementing general functions that would allow users to execute arbitrary operating system commands.
- Pragmatic defaults: The default [security policy][] aims to balance security and usability, enabling common workflows out of the box while keeping more sensitive capabilities opt-in. These defaults may be tightened in future releases, but each project is ultimately responsible for reviewing the policy and adjusting it to match its own trust model and requirements.

This combination of sandboxing and strict defaults effectively minimizes potential security vulnerabilities during the Hugo build process.

## Dependency security

Hugo utilizes [Go modules][] to manage its dependencies, compiling as a static binary. Go modules create a `go.sum` file, a critical security feature. This file acts as a database, storing the expected cryptographic checksums of all dependencies, including those required indirectly (transitive dependencies).

[Hugo modules][], which extend the functionality of Go modules, also produce a `go.sum` file. To ensure dependency integrity, commit this `go.sum` file to your version control. If Hugo detects a checksum mismatch during the build process, it will fail, indicating a possible attempt to [tamper with your project's dependencies][].

## Web application security

Hugo's security philosophy is rooted in established security standards, primarily aligning with the threats defined by [OWASP][]. For HTML output, Hugo operates under a clear trust model. This model assumes that template and configuration authors, the developers, are trustworthy. However, the data supplied to these templates is inherently considered untrusted. This distinction is crucial for understanding how Hugo handles potential security risks.

To prevent unintended escaping of data that developers know is safe, Hugo provides  [`safe`][] functions, such as [`safe.HTML`][]. These functions allow developers to explicitly mark data as trusted, bypassing the default escaping mechanisms. This is essential for scenarios where data is generated or sourced from reliable sources. However, an exception exists: enabling [inline shortcodes][]. By activating this feature, you are implicitly trusting the logic within the shortcodes and the data contained within your content files.

It's vital to remember that Hugo is a static site generator. This architectural choice significantly reduces the attack surface by eliminating the complexities and vulnerabilities associated with dynamic user input. Unlike dynamic websites, Hugo generates static HTML files, minimizing the risk of real-time attacks. Regarding content, Hugo's default Markdown renderer is [configured to sanitize][] potentially unsafe content. This default behavior ensures that potentially malicious code or scripts are removed or escaped. However, this setting can be reconfigured if you have a high degree of confidence in the safety of your content sources.

In essence, Hugo prioritizes secure output by establishing a clear trust boundary between developers and data. By default, it errs on the side of caution, sanitizing potentially unsafe content and escaping data. Developers have the flexibility to adjust these defaults through [`safe`][] functions and [configuration settings][], but they must do so with a clear understanding of the security implications. Hugo's static site generation model further strengthens its security posture by minimizing dynamic vulnerabilities.

## Configuration

See [configure security][].

[Go modules]: https://go.dev/wiki/Modules#modules
[Hugo modules]: /hugo-modules/
[OWASP]: https://en.wikipedia.org/wiki/OWASP
[`safe.HTML`]: /functions/safe/html/
[`safe`]: /functions/safe/
[configuration settings]: /configuration/security/
[configure security]: /configuration/security/
[configured to sanitize]: /configuration/markup/#rendererunsafe
[content adapters]: /content-management/content-adapters/
[inline shortcodes]: /content-management/shortcodes/#inline
[security policy]: /configuration/security/
[tamper with your project's dependencies]: https://julienrenaux.fr/2019/12/20/github-actions-security-risk/
