---
title: crypto.Hash
description: Hashes the given input with the given algorithm and returns its checksum encoded to a hexadecimal string.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: ['crypto.Hash [ALGORITHM] INPUT']
---

The `ALGORITHM` is one of `md5`, `sha1`, `sha256` (the default), `sha384`, or `sha512`:

```go-html-template
{{ crypto.Hash "sha256" "Hello world" }} → 64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c
{{ "Hello world" | crypto.Hash "sha512" }} → b7f783baed8297f0db917462184ff4f08e69c2d5e5f79a942600f9725f58ce1f29c18139bf80b06c0fff2bdd34738452ecf40c488c22a7e3d80cdf6f9c1c0d47
```

If you omit the algorithm, it defaults to `sha256`:

```go-html-template
{{ "Hello world" | crypto.Hash }} → 64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c
```

The supported algorithms match those used for the [Subresource Integrity] hash in [`.Data.Integrity`] on a fingerprinted resource. Combine `crypto.Hash` with [`encoding.HexDecode`] and [`encoding.Base64Encode`] to construct an SRI hash from a string:

```go-html-template
{{ $algo := "sha256" }}
{{ $integrity := printf "%s-%s" $algo ("Hello world" | crypto.Hash $algo | encoding.HexDecode | encoding.Base64Encode) }}
```

[Subresource Integrity]: https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
[`.Data.Integrity`]: /methods/resource/data/
[`encoding.HexDecode`]: /functions/encoding/hexdecode/
[`encoding.Base64Encode`]: /functions/encoding/base64encode/
