---
title: "Configuring Hugo"
pubdate: "2013-07-01"
---

The directory structure and templates provide the majority of the
configuration for a site. In fact a config file isn't even needed for many websites
since the defaults used follow commonly used patterns.

The following is an example of a config file with the default values: 

    SourceDir: "content"
    LayoutDir: "layouts"
    PublishDir: "public"
    BuildDrafts: false
    Tags:
       category: "categories"
       tag: "tags"
    BaseUrl: "http://yourSite.com/"
    ...

