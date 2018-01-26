---
title : "Page Resources"
description : "Page Resources are files included in a page bundle. You can use them in your template and add metadata"
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
: The main type of the resource. For exemple a file of MIME type `image/jpg` has for ResourceType `image`.

Name
: The filename (relative path to the bundle). It can be overwritten with the resource's Front Matter metadata.

Title
: Same as filename. It can be overwritten with the resource's Front Matter metadata.

Permalink
: The absolute URL of the resource.

RelPermalink
: The relative URL of the resource.

## Methods
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

### Pattern Matching
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

Page Resources metadata is managed from their page's Front Matter with an array named `resources`. Batch assign is made possible using glob pattern matching.

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
~~~yaml
title: Application
date : 2018-01-25
resources :
- src : "images/header.*"
  name : "header"
- src : "**.pdf"
  title = "PDF file #:counter"
  params :
    icon : "pdf"
- src : "**.docx"
  title : "Word file #:counter"
  params :
    icon : "word"
- src : "documents/photo_specs.pdf"
  title : "Photo Specifications"
  params:
    icon : "image"
- src : "documents/guide.pdf"
  title : "Instruction Guide"
- src : "documents/checklist.pdf"
  title : "Document Checklist"
- src : "documents/payment.docx"
  title : "Proof of Payment"
~~~ 

~~~toml
title = Application
date : 2018-01-25
[[resources]]
  src = "images/header.*"
  name = "header"
[[resources]]
  src = "**.pdf"
  title = "PDF file #:counter"
  [resources.params]
    icon = "pdf"
[[resources]]
  src = "**.docx"
  title = "Word file #:counter"
  [resources.params]
    icon = "word"
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
 ~~~


From the metadata example above:

- `header.jpg` will receive a new `Name` and won't be retrieved by `.Match "*/header.jpg"` anymore but something like `.Match "header"`.
- `documents/photo_specs.pdf` will get the `image` icon
- `documents/checklist.pdf`, `documents/guide.pdf` and `documents/payment.docx` will receive a unique Title
- Every pdf in the bundle exepct documents/photo_specs.pdf` will receive the `pdf` icon along with a Title using the keyword `:counter`
- Every docx in the bundle will receive the `word` icon along with a Title using the keyword `:counter`

{{% warning %}}
The __order matters__, every metadata key/value pair assigned overwrites any previous ones assigned to the same `src` target. As in the example above broad targets rules will usually be defined before the narrower ones.
{{%/ warning %}}


