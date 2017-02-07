---
title: Configuration
linktitle:
description: Scaffolding new projects, configuration, and source organization.
date: 2017-01-02
publishdate: 2017-01-02
lastmod: 2017-01-02
weight: 10
draft: false
slug:
aliases:
notes:
---

Here is `some code` in the body copy.

{{% input "my-first-template.html" %}}
```go
<ul>
  {{range $.Site.Pages}}
    <li>{{.Title}}</li>
  {{end}}
</ul>
```
{{% /input %}}

{{% input "first-script.js" %}}
```js
function myFunction(arg){
  console.log(arg);
}
```
{{% /input %}}

{{% input "first-data.json" %}}
```json
function myFunction(arg){
  console.log(arg);
}
```
{{% /input %}}

{{% output "layouts/_default/list.html" %}}
```go
<ul>
  {{range $.Site.Pages}}
    <li>{{.Title}}</li>
  {{end}}
</ul>
```
{{% /output %}}

```javascript
function myFunction(val){
  console.log(val);
}
```
