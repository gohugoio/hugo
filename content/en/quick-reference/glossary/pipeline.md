---
title: pipeline
---

Within a [_template action_](g), a _pipeline_ is a possibly chained sequence of values, [_function_](g) calls, or [_method_](g) calls. Functions and methods in the pipeline may take multiple [_arguments_](g).

A pipeline may be *chained* by separating a sequence of commands with pipeline characters (`|`). In a chained pipeline, the result of each command is passed as the last argument to the following command. The output of the final command in the pipeline is the value of the pipeline.
