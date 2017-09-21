---
title: Accented Characters in URLs
linktitle: Accented Characters in URLs
description: If you're having trouble with special characters in your taxonomies or titles adding odd characters to your URLs.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [urls,multilingual,special characters]
categories: [troubleshooting]
menu:
  docs:
    parent: "troubleshooting"
weight:
draft: false
slug:
aliases: [/troubleshooting/categories-with-accented-characters/]
toc: true
---

## Trouble: Categories with accented characters

> One of my categories is named "Le-carré," but the link ends up being generated like this:
>
> ```
> categories/le-carr%C3%A9
> ```
>
> And not working. Is there an easy fix for this that I'm overlooking?

## Solution

Are you a macOS user? If so, you are likely a victim of HFS Plus file system's insistence to store the "é" (U+00E9) character in Normal Form Decomposed (NFD) mode, i.e. as "e" + "  ́" (U+0065 U+0301).

`le-carr%C3%A9` is actually correct, `%C3%A9` being the UTF-8 version of U+00E9 as expected by the web server. The problem is that OS X turns [U+00E9] into [U+0065 U+0301], and thus `le-carr%C3%A9` no longer works.  Instead, only `le-carre%CC%81` ending with `e%CC%81` would match that [U+0065 U+0301] at the end.

This is unique to OS X. The rest of the world does not do this, and most certainly not your web server which is most likely running Linux. This is not a Hugo-specific problem either. Other people have been bitten by this when they have accented characters in their HTML files.

Note that this problem is not specific to Latin scripts. Japanese Mac users often run into the same issue; e.g., with `だ` decomposing into `た` and `&#x3099;`. (Read the [Japanese Perl users article][]).

Rsync 3.x to the rescue! From [an answer posted on Server Fault][]:

> You can use rsync's `--iconv` option to convert between UTF-8 NFC & NFD, at least if you're on a Mac. There is a special `utf-8-mac` character set that stands for UTF-8 NFD. So to copy files from your Mac to your web server, you'd need to run something like:
>
> `rsync -a --iconv=utf-8-mac,utf-8 localdir/ mywebserver:remotedir/`
>
> This will convert all the local filenames from UTF-8 NFD to UTF-8 NFC on the remote server. The files' contents won't be affected. - [Server Fault][]

Please make sure you have the latest version of rsync 3.x installed. The rsync that ships with OS X is outdated. Even the version that comes packaged with 10.10 (Yosemite) is version 2.6.9 protocol version 29. The `--iconv` flag is new in rsync 3.x.

### Discussion Forum References

* http://discourse.gohugo.io/t/categories-with-accented-characters/505
* http://wiki.apache.org/subversion/NonNormalizingUnicodeCompositionAwareness
* https://en.wikipedia.org/wiki/Unicode_equivalence#Example
* http://zaiste.net/2012/07/brand_new_rsync_for_osx/
* https://gogo244.wordpress.com/2014/09/17/drived-me-crazy-convert-utf-8-mac-to-utf-8/

[an Answer posted on Server Fault]: http://serverfault.com/questions/397420/converting-utf-8-nfd-filenames-to-utf-8-nfc-in-either-rsync-or-afpd "Converting UTF-8 NFD filenames to UTF-8 NFC in either rsync or afpd, Server Fault Discussion"
[Japanese Perl users article]: http://perl-users.jp/articles/advent-calendar/2010/english/24 "Encode::UTF8Mac makes you happy while handling file names on MacOSX"
[Server Fault]: http://serverfault.com/questions/397420/converting-utf-8-nfd-filenames-to-utf-8-nfc-in-either-rsync-or-afpd "Converting UTF-8 NFD filenames to UTF-8 NFC in either rsync or afpd, Server Fault Discussion"
