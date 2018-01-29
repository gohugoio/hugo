---
title : "Page Resources"
description : "Page Resources are files included in a page bundle. You can use them in your template and add metadata."
date : 2018-01-24T13:10:00-05:00
lastmod : 2018-01-26T13:45:08-05:00
categories : ["content management"]
weight : 4003
draft : false
toc : true
linktitle : "Page Resources"
menu :
  docs:
    parent : "content-management"
    weight : 31
---

## Properties

ResourceType
: The main type of the resource. For example, a file of MIME type `image/jpg` has for ResourceType `image`.

Name
: The filename (relative path to the bundle). It can be overwritten with the resource's Front Matter metadata.

Title
: Same as filename. It can be overwritten with the resource's Front Matter metadata.

Permalink
: The absolute URL of the resource. This gets a value only where the _ResourceType_ is **not** `page`.

RelPermalink
: The relative URL of the resource. This gets a value only where the _ResourceType_ is **not** `page`.

## Methods
ByType
: Retrieve the page resources of the passed type.

```go
{{ .Resources.ByType "image" }}
```
Match
: Retrieve all the page resources (as a slice) whose `Name` matches the Glob pattern ([examples](https://github.com/gobwas/glob/blob/master/readme.md)) passed as parameter. The matching is case-insensitive.

```go
{{ .Resources.Match "images/*" }}
```

GetMatch
: Same as `Match` but will only retrieve the first matching resource.

### Pattern Matching
```go
// Using Match/GetMatch to find this images/sunset.jpg ?
.Resources.Match "images/sun*" ✅
.Resources.Match "**/Sunset.jpg" ✅
.Resources.Match "images/*.jpg" ✅
.Resources.Match "**.jpg" ✅
.Resources.Match "*" 🚫
.Resources.Match "sunset.jpg" 🚫
.Resources.Match "*sunset.jpg" 🚫

```

## Metadata

Page Resources metadata is managed from their page's Front Matter with an array/table parameter named `resources`. Batch assign is made possible using glob pattern matching.

### Available metadata

name
: Will overwrite Name

{{% warning %}}
The methods Match and GetMatch use Name to match the resource. Overwrite wisely.
{{%/ warning %}}

title
: Will overwrite Title

params
: An array of custom params to be retrieve much like page params
`{{ .Params.credits }}`

### Example
#### `resources` parameter in YAML
~~~yaml
title: Application
date : 2018-01-25
resources :
- src : "images/header.*"
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
~~~

#### `resources` parameter in TOML
~~~toml
title = Application
date : 2018-01-25
[[resources]]
  src = "images/header.*"
  name = "header"
[[resources]]
  src = "documents/photo_specs.pdf"
  title = "Photo Specifications"
  [resources.params]
    icon = "photo"
[[resources]]
  src = "documents/guide.pdf"
  title = "Instruction Guide"
[[resources]]
  src = "documents/checklist.pdf"
  title = "Document Checklist"
[[resources]]
  src = "documents/payment.docx"
  title = "Proof of Payment"
[[resources]]
  src = "**.pdf"
  name = "pdf-file-:counter"
  [resources.params]
    icon = "pdf"
[[resources]]
  src = "**.docx"
  [resources.params]
    icon = "word"
~~~


From the metadata example above:

- `header.jpg` will receive a new `Name` and won't be retrieved by `.Match "*/header.jpg"` anymore, but something like `.Match "header"`.
- `documents/photo_specs.pdf` will get the `photo` icon.
- `documents/checklist.pdf`, `documents/guide.pdf` and `documents/payment.docx` will receive `Title` as set by `title`.
- Every pdf in the bundle except `documents/photo_specs.pdf` will receive the `pdf` icon.
- All pdf files will get a new `Name`. The `name` parameter contains a special placeholder [`:counter`](#counter). That will cause the retrieved pdf files to have names `pdf-file-1`, `pdf-file-2`, `pdf-file-3`.
- Every docx in the bundle will receive the `word` icon.

{{% warning %}}
The __order matters__ --- Only the **first set** values of the `title`, `name` and `params`-**keys** will be used. Consecutive parameters will be set only for the ones not already set. For example, in the above example, `.Params.icon` is already first set to `"photo"` in `src = "documents/photo_specs.pdf"`. So that would not get overridden to `"pdf"` by the later set `src = "**.pdf"` rule.
{{%/ warning %}}

### `:counter` placeholder in `name` and `title` {#counter}

The `:counter` is a special placeholder recognized in `name` and `title` parameters in the `resources` Front Matter.

If the `name` value contains the `":counter"` string, a "name-counter" is initialized to 1, and if the `title` value contains the same string, a separate "title-counter" is initialized to 1 as well.

For example, if a bundle has the resources `photo_specs.pdf`, `other_specs.pdf`, `guide.pdf` and `checklist.pdf`, and the Front Matter has specified the `resources` as:

~~~toml
[[resources]]
  src = "*specs.pdf"
  title = "Specifications #:counter"
[[resources]]
  src = "**.pdf"
  name = "pdf-file-:counter"
~~~

the `Name` and `Title` will be assigned to the resource files as follows:

| Resource file     | `Name`            | `Title`               |
|-------------------|-------------------|-----------------------|
| checklist.pdf     | `"pdf-file-1.pdf` | `"checklist.pdf"`     |
| guide.pdf         | `"pdf-file-2.pdf` | `"guide.pdf"`         |
| other\_specs.pdf  | `"pdf-file-3.pdf` | `"Specifications #1"` |
| photo\_specs.pdf  | `"pdf-file-4.pdf` | `"Specifications #2"` |
