---
title: index
linktitle: index
description: Looks up the index(es) or key(s) of the data structure passed into it.
godocref: https://golang.org/pkg/text/template/#hdr-Functions
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: []
signature: ["index COLLECTION INDEX", "index COLLECTION KEY"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: [/functions/index/]
needsexample: true
---

From the Godocs:

> Returns the result of indexing its first argument by the following arguments. Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each indexed item must be a map, slice, or array.

In Go templates, you can't access array, slice, or map elements directly the same way you would in Go. For example, `$.Site.Data.authors[.Params.authorkey]` isn't supported syntax.

Instead, you have to use `index`, a function that handles the lookup for you.

## Example: Load Data from a Path Based on Front Matter Params

Assume you want to add a `location = ""` field to your front matter for every article written in `content/vacations/`. You want to use this field to populate information about the location at the bottom of the article in your `single.html` template. You also have a directory in `data/locations/` that looks like the following:

```
.
└── data
    └── locations
        ├── abilene.toml
        ├── chicago.toml
        ├── oslo.toml
        └── provo.toml
```

Here is an example of the data inside `data/locations/oslo.toml`:

```
website = "https://www.oslo.kommune.no"
pop_city = 658390
pop_metro = 1717900
```

The example we will use will be an article on Oslo, whose front matter should be set to exactly the same name as the corresponding file name in `data/locations/`:

```
title = "My Norwegian Vacation"
location = "oslo"
```

The content of `oslo.toml` can be accessed from your template using the following node path: `.Site.Data.locations.oslo`. However, the specific file you need is going to change according to the front matter.

This is where the `index` function is needed. `index` takes 2 parameters in this use case:

1. The node path
2. A string corresponding to the desired data; e.g.&mdash;

```
{{ index .Site.Data.locations “oslo” }}
```

The variable for `.Params.location` is a string and can therefore replace `oslo` in the example above:

```
{{ index .Site.Data.locations .Params.location }}
=> map[website:https://www.oslo.kommune.no pop_city:658390 pop_metro:1717900]
```

Now the call will return the specific file according to the location specified in the content's front matter, but you will likely want to write specific properties to the template. You can do this by continuing down the node path via dot notation (`.`):

```
{{ (index .Site.Data.locations .Params.location).pop_city }}
=> 658390
```

