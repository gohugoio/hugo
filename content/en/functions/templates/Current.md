---
title: templates.Current
description: Returns information about the currently executing template.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: tpl.CurrentTemplateInfo
    signatures: [templates.Current]
---

> [!note]
> This function is experimental and subject to change.

{{< new-in 0.146.0 />}}

The `templates.Current` function provides introspection capabilities, allowing you to access details about the currently executing templates. This is useful for debugging complex template hierarchies and understanding the flow of execution during rendering.

## Methods

Ancestors
: (`tpl.CurrentTemplateInfos`) Returns a slice containing information about each template in the current execution chain, starting from the parent of the current template and going up towards the initial template called. It excludes any base template applied via `define` and `block`. You can chain the `Reverse` method to this result to get the slice in chronological execution order.

Base
: (`tplimpl.TemplInfo`) Returns an object representing the base template that was applied to the current template, if any. This may be `nil`.

Filename
: (`string`) Returns the absolute path of the current template. This will be empty for embedded templates.

Name
: (`string`) Returns the path of the current template, relative to the `layouts` directory.

Parent
: (`tpl.CurrentTemplateInfo`) Returns an object representing the parent of the current template, if any. This may be `nil`.

## Examples

The examples below help visualize template execution and require a `debug` parameter set to `true` in your site configuration:

{{< code-toggle file=hugo >}}
[params]
debug = true
{{< /code-toggle >}}

### Template boundaries

To visually mark where a template begins and ends execution:

```go-html-template {file="layouts/_default/single.html"}
{{ define "main" }}
  {{ if site.Params.debug }}
    <div class="debug">[entering {{ templates.Current.Filename }}]</div>
  {{ end }}

  <h1>{{ .Title }}</h1>
  {{ .Content }}

  {{ if site.Params.debug }}
    <div class="debug">[leaving {{ templates.Current.Filename }}]</div>
  {{ end }}
{{ end }}
```

### Template call stack

To show the chain of templates that led to the current one, create a partial template that iterates through its ancestors:

```go-html-template {file="layouts/partials/render-template-stack.html" copy=true}
<div class="debug">
  {{ with templates.Current }}
    {{ range .Ancestors }}
      {{ .Filename }}<br>
      {{ with .Base }}
        {{ .Filename }}<br>
      {{ end }}
    {{ end }}
  {{ end }}
</div>
```

Then call the partial from any template:

```go-html-template {file="layouts/partials/footer/copyright.html" copy=true}
{{ if site.Params.debug }}
  {{ partial "render-template-stack.html" . }}
{{ end }}
```

The rendered template stack would look something like this:

```text
/home/user/project/layouts/partials/footer/copyright.html
/home/user/project/themes/foo/layouts/partials/footer.html
/home/user/project/themes/foo/layouts/_default/single.html
/home/user/project/themes/foo/layouts/_default/baseof.html
```

To reverse the order of the entries chain the `Reverse` method to the `Ancestors` method:

```go-html-template {file="layouts/partials/render-template-stack.html" copy=true}
<div class="debug">
  {{ with templates.Current }}
    {{ range .Ancestors.Reverse }}
      {{ with .Base }}
        {{ .Filename }}<br>
      {{ end }}
      {{ .Filename }}<br>
    {{ end }}
  {{ end }}
</div>
```

To render each entry as an anchor element that, when clicked, will open the template in VS Code:

```go-html-template {file="layouts/partials/render-template-stack.html" copy=true}
<div class="debug">
  {{ with templates.Current }}
    {{ range .Ancestors }}
      <a href="vscode://file/{{ .Filename }}">{{ .Filename }}</a><br>
      {{ with .Base }}
        <a href="vscode://file/{{ .Filename }}">{{ .Filename }}</a><br>
      {{ end }}
    {{ end }}
  {{ end }}
</div>
```
