---
title: Configure file caches
linkTitle: Caches
description: Configure file caches.
categories: []
keywords: []
---

This is the default configuration:

{{< code-toggle config=caches />}}

## Keys

dir
: (`string`) The absolute file system path where the cached files will be stored. You can begin the path with the `:cacheDir` or `:resourceDir` token. These tokens will be replaced with the actual configured cache directory and resource directory paths, respectively.

maxAge
: (`string`) The [duration](g) a cached entry remains valid before being evicted. A value of `0` disables the cache. A value of `-1` means the cache entry never expires (the default).

## Tokens

`:cacheDir`
: (`string`) The designated cache directory. See&nbsp;[details](/configuration/all/#cachedir).

`:project`
: (`string`) The base directory name of the current Hugo project. By default, this ensures each project has isolated file caches, so running `hugo --gc` will only affect the current project's cache and not those of other Hugo projects on the same machine.

`:resourceDir`
: (`string`) The designated directory for caching output from [asset pipelines](g). See&nbsp;[details](/configuration/all/#resourcedir).
