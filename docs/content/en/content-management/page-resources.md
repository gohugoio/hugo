---
title : "Page Resources"
description : "Page Resources -- images, other pages, documents etc. -- have page-relative URLs and their own metadata."
date: 2018-01-24
categories: ["content management"]
keywords: [bundle,content,resources]
weight: 4003
draft: false
toc: true
linktitle: "Page Resources"
menu:
  docs:
    parent: "content-management"
    weight: 31
---

## Properties

ResourceType
: The main type of the resource. For example, a file of MIME type `image/jpeg` has the ResourceType `image`.

Name
: Default value is the filename (relative to the owning page). Can be set in front matter.

Title
: Default value is the same as `.Name`. Can be set in front matter.

Permalink
: The absolute URL to the resource. Resources of type `page` will have no value.

RelPermalink
: The relative URL to the resource. Resources of type `page` will have no value.

Content
: The content of the resource itself. For most resources, this returns a string with the contents of the file. This can be used to inline some resources, such as `<script>{{ (.Resources.GetMatch "myscript.js").Content | safeJS }}</script>` or `<img src="{{ (.Resources.GetMatch "mylogo.png").Content | base64Encode }}">`.

MediaType
: The MIME type of the resource, such as `image/jpeg`.

MediaType.MainType
: The main type of the resource's MIME type. For example, a file of MIME type `application/pdf` has for MainType `application`.

MediaType.SubType
: The subtype of the resource's MIME type. For example, a file of MIME type `application/pdf` has for SubType `pdf`. Note that this is not the same as the file extension - PowerPoint files have a subtype of `vnd.mspowerpoint`.

MediaType.Suffixes
: A slice of possible suffixes for the resource's MIME type.

## Methods
ByType
: Returns the page resources of the given type.

```go
{{ .Resources.ByType "image" }}
```
Match
: Returns all the page resources (as a slice) whose `Name` matches the given Glob pattern ([examples](https://github.com/gobwas/glob/blob/master/readme.md)). The matching is case-insensitive.

```go
{{ .Resources.Match "images/*" }}
```

GetMatch
: Same as `Match` but will return the first match.

### Pattern Matching
```go
// Using Match/GetMatch to find this images/sunset.jpg ?
.Resources.Match "images/sun*" ✅
.Resources.Match "**/sunset.jpg" ✅
.Resources.Match "images/*.jpg" ✅
.Resources.Match "**.jpg" ✅
.Resources.Match "*" 🚫
.Resources.Match "sunset.jpg" 🚫
.Resources.Match "*sunset.jpg" 🚫

```

## Page Resources Metadata

The page resources' metadata is managed from the corresponding page's front matter with an array/table parameter named `resources`. You can batch assign values using [wildcards](http://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm).

{{% note %}}
Resources of type `page` get `Title` etc. from their own front matter.
{{% /note %}}

name
: Sets the value returned in `Name`.

{{% warning %}}
The methods `Match` and `GetMatch` use `Name` to match the resources.
{{%/ warning %}}

title
: Sets the value returned in `Title`

params
: A map of custom key/values.


###  Resources metadata example

{{< code-toggle copy="false">}}
title: Application
date : 2018-01-25
resources :
- src : "images/sunset.jpg"
  name : "header"
- src : "documents/photo_specs.pdf"
  title : "Photo Specifications"
  params:
    icon : "photo"
- src : "documents/guide.pdf"
  title : "Instruction Guide"
- src : "documents/checklist.pdf"
  title : "Document Checklist"
- src : "documents/payment.docx"
  title : "Proof of Payment"
- src : "**.pdf"
  name : "pdf-file-:counter"
  params :
    icon : "pdf"
- src : "**.docx"
  params :
    icon : "word"
{{</ code-toggle >}}

From the example above:

- `sunset.jpg` will receive a new `Name` and can now be found with `.GetMatch "header"`.
- `documents/photo_specs.pdf` will get the `photo` icon.
- `documents/checklist.pdf`, `documents/guide.pdf` and `documents/payment.docx` will get `Title` as set by `title`.
- Every `PDF` in the bundle except `documents/photo_specs.pdf` will get the `pdf` icon.
- All `PDF` files will get a new `Name`. The `name` parameter contains a special placeholder [`:counter`](#the-counter-placeholder-in-name-and-title), so the `Name` will be `pdf-file-1`, `pdf-file-2`, `pdf-file-3`.
- Every docx in the bundle will receive the `word` icon.

{{% warning %}}
The __order matters__ --- Only the **first set** values of the `title`, `name` and `params`-**keys** will be used. Consecutive parameters will be set only for the ones not already set. In the above example, `.Params.icon` is first set to `"photo"` in `src = "documents/photo_specs.pdf"`. So that would not get overridden to `"pdf"` by the later set `src = "**.pdf"` rule.
{{%/ warning %}}

### The `:counter` placeholder in `name` and `title`

The `:counter` is a special placeholder recognized in `name` and `title` parameters `resources`.

The counter starts at 1 the first time they are used in either `name` or `title`.

For example, if a bundle has the resources `photo_specs.pdf`, `other_specs.pdf`, `guide.pdf` and `checklist.pdf`, and the front matter has specified the `resources` as:

{{< code-toggle copy="false">}}
[[resources]]
  src = "*specs.pdf"
  title = "Specification #:counter"
[[resources]]
  src = "**.pdf"
  name = "pdf-file-:counter"
{{</ code-toggle >}}

the `Name` and `Title` will be assigned to the resource files as follows:

| Resource file     | `Name`            | `Title`               |
|-------------------|-------------------|-----------------------|
| checklist.pdf     | `"pdf-file-1.pdf` | `"checklist.pdf"`     |
| guide.pdf         | `"pdf-file-2.pdf` | `"guide.pdf"`         |
| other\_specs.pdf  | `"pdf-file-3.pdf` | `"Specification #1"` |
| photo\_specs.pdf  | `"pdf-file-4.pdf` | `"Specification #2"` |
