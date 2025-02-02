---
title: pipeline
---

Within a [template action](g), a pipeline is a possibly chained sequence of values, [function](g) calls, or [method](g) calls. Functions and methods in the pipeline may take multiple [arguments](g).

A pipeline may be *chained* by separating a sequence of commands with pipeline characters "|". In a chained pipeline, the result of each command is passed as the last argument to the following command. The output of the final command in the pipeline is the value of the pipeline.
