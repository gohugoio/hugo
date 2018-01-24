---
title : "Resources"
description : "Resources allows one to access content and meta-data from page, section and branch bundles, including images."
date : 2018-01-24T13:10:00-05:00
lastmod : 2018-01-24T13:45:08-05:00
categories : ["content management", "bundles"]
weight : 4003
draft : false
toc : true
linktitle : "Resources"
---


## Properties

ResourceType
: The main type of the resource. For exemple a file of MIME type `image/jpg` has for type `image`.

Name
: The filename. (relative path to the bundle) It can be overwritten with the resource's Front Matter metadata.

Title
: Same as filename. It can be overwritten with the resource's Front Matter metadata.

Permalink
: The absolute URL of the resource.

RelPermalink
: The relative URL of the resource.

## Methods / Functions
ByType
: Retrieve the page resources of the passed type.

```go
{{ .Resources.ByType "images" }}
```
Match
: Retrieve all the page resources whose Name matches the [Glob pattern](https://en.wikipedia.org/wiki/Glob_(programming)) passed as parameter. The matching is case insensitive.

```go
{{ .Resources.Match "images/*" }}
```

GetMatch
: Same as Match but will only retrieve the first matching resource.

```go
//Using Match/GetMatch to find this images/sunset.jpg ?
.Resources.Match "images/sun*" ✅ 
.Resources.Match "**/Sunset.jpg" ✅
.Resources.Match "images/*.jpg" ✅
.Resources.Match "**.jpg" ✅ 
.Resources.Match "*" 🚫
.Resources.Match "sunset.jpg" 🚫
.Resources.Match "*sunset.jpg" 🚫

```

## Metadata

Resources metadata management is possible from the Front Matter's array `resources`. Batch assign is made possible using glob pattern matching.

~~~yaml
title: "Sunflare 101"
date: 2018-01-09T13:03:53-05:00
resources:
 - src: "*/sunset.jpg"
   title: "Beautiful Sunset"
 - src: "*/sunrise.jpg"
   title: "Beautiful Sunrise"
   name: "sunrise-1"
 - src: "*/*.jpg"
   params: 
     credits: Hugo Doe
     license:CC BY
 - src: "*/tokyo-sunset.jpg"
   params: 
     credits: Jiro Ashima
~~~ 

In the exemple above, we use glob to target several files and groups of file.

- `sunset.jpg` will receive a distinct Title
- `sunrise.jpg` will receive a distinct Title and Name and won't be retrieved by `.Match "*/sunrise.jpg"` anymore but something like `.Match "sunrise-*"`.
- Every jpg in the bundle will receive the same `credits` and `license` parameter.
- `tokyo-sunset.jpg` will receive a distinct `credits`

The __order matters__, each rule overwriting the params assigned by the previous one on their shared target. The rule for `*/tokyo-sunset.jpg` has to be declared after the more generic `*/*.jpg`.

## Available metadata

name
: Will overwrite Name

{{< warning >}}
The methods Match and GetMatch use Name to match the resource. Overwrite wisely.
{{</ warning >}}

title
: Will overwrite Title

params
: An array of custom params to be retrieve much like page params <br> `{{ .Params.credits }}`



