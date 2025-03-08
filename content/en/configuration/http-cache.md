---
title: Configure the HTTP cache
linkTitle: HTTP cache
description: Configure the HTTP cache.
categories: []
keywords: []
---

> [!note]
> This configuration is only relevant when using the [`resources.GetRemote`] function.

## Layered caching

Hugo employs a layered caching system.

```goat {.w-40}
 .-----------.
|  dynacache  |
 '-----+-----'
       |
       v
 .----------.
| HTTP cache |
 '-----+----'
       |
       v
 .----------.
| file cache |
 '-----+----'
```

Dynacache
: An in-memory cache employing a Least Recently Used (LRU) eviction policy. Entries are removed from the cache when changes occur, when they match [cache-busting] patterns, or under low-memory conditions.

HTTP Cache
: An HTTP cache for remote resources as specified in [RFC 9111]. Optimal performance is achieved when resources include appropriate HTTP cache headers. The HTTP cache utilizes the file cache for storage and retrieval of cached resources.

File cache
: See [configure file caches].

The HTTP cache involves two key aspects: determining which content to cache (the caching process itself) and defining the frequency with which to check for updates (the polling strategy).

## HTTP caching

The HTTP cache behavior is defined for a configured set of resources. Stale resources will be refreshed from the file cache, even if their configured Time-To-Live (TTL) has not expired. If HTTP caching is disabled for a resource, Hugo will bypass the cache and access the file directly.

The default configuration disables everything:

{{< code-toggle file=hugo >}}
[HTTPCache.cache.for]
excludes = ['**']
includes = []
{{< /code-toggle >}}

cache.for.excludes
: (`string`) A list of [glob](g) patterns to exclude from caching.

cache.for.includes
: (`string`) A list of [glob](g) patterns to cache.

## HTTP polling

Polling is used in watch mode (e.g., `hugo server`) to detect changes in remote resources. Polling can be enabled even if HTTP caching is disabled. Detected changes trigger a rebuild of pages using the affected resource. Polling can be disabled for specific resources, typically those known to be static.

The default configuration disables everything:

{{< code-toggle file=hugo >}}
[[HTTPCache.polls]]
disable = true
high = '0s'
low = '0s'
[HTTPCache.polls.for]
includes = ['**']
excludes = []
{{< /code-toggle >}}

polls
: A slice of polling configurations.

polls.disable
: (`bool`) Whether to disable polling for this configuration.

polls.low
: (`string`) The minimum polling interval expressed as a [duration](g). This is used after a recent change and gradually increases towards `polls.high`.

polls.high
: (`string`) The maximum polling interval expressed as a [duration](g). This is used when the resource is considered stable.

polls.for.excludes
: (`string`) A list of [glob](g) patterns to exclude from polling for this configuration.

polls.for.includes
: (`string`) A list of [glob](g) patterns to include in polling for this configuration.

## Behavior

Polling and HTTP caching interact as follows:

- With polling enabled, rebuilds are triggered only by actual changes, detected via `eTag` changes (Hugo generates an MD5 hash if the server doesn't provide one).
- If polling is enabled but HTTP caching is disabled, the remote is checked for changes only after the file cache's TTL expires (e.g., a `maxAge` of `10h` with a `1s` polling interval is inefficient).
- If both polling and HTTP caching are enabled, changes are checked for even before the file cache's TTL expires. Cached `eTag` and `last-modified` values are sent in `if-none-match` and `if-modified-since` headers, respectively, and a cached response is returned on HTTP [304].

[`resources.GetRemote`]: /functions/resources/getremote/
[304]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304
[cache-busting]: /configuration/build/#cache-busters
[configure file caches]: /configuration/caches/
[RFC 9111]: https://datatracker.ietf.org/doc/html/rfc9111
