---
date: 2016-02-06
linktitle: Analytics
menu:
  main:
    parent: extras
next: /extras/builders
prev: /extras/aliases
title: Analytics in Hugo
---

Hugo ships with prebuilt internal templates for Google Analytics tracking, including both synchronous and asynchronous tracking codes.

## Configuring Google Analytics

Provide your tracking id in your configuration file, e.g. config.yaml.

    googleAnalytics = "UA-123-45"

## Example

Include the internal template in your templates like so:

    {{ template "_internal/google_analytics.html" . }}

For async include the async template:

    {{ template "_internal/google_analytics_async.html" . }}
