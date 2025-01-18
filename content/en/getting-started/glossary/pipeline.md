---
title: pipeline
---

Within a [template action], a pipeline is a possibly chained sequence of values, [function] calls, or [method] calls. Functions and methods in the pipeline may take multiple [arguments][argument].

A pipeline may be *chained* by separating a sequence of commands with pipeline characters "|". In a chained pipeline, the result of each command is passed as the last argument to the following command. The output of the final command in the pipeline is the value of the pipeline. See the [Go&nbsp;documentation](https://pkg.go.dev/text/template#hdr-Pipelines) for details.

{{% include "/getting-started/glossary/_link-reference-definitions" %}}
