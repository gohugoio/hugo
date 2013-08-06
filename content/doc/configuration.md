---
title: "Configuring Hugo"
pubdate: "2013-07-01"
---

The directory structure and templates provide the majority of the
configuration for a site. In fact a config file isn't even needed for many websites
since the defaults used follow commonly used patterns.

Hugo expects to find the config file in the root of the source directory and
will look there first for a config.yaml file. If none is present it will
then look for a config.json file, followed by a config.toml file.

**Please note the field names must be all lowercase**

## Examples

The following is an example of a yaml config file with the default values: 

    ---
    contentdir: "content"
    layoutdir: "layouts"
    publishdir: "public"
    builddrafts: false
    indexes:
       category: "categories"
       tag: "tags"
    baseurl: "http://yoursite.com/"
    ...


The following is an example of a json config file with the default values: 

    {
        "contentdir": "content",
        "layoutdir": "layouts",
        "publishdir": "public",
        "builddrafts": false,
        "indexes": {
           "category": "categories",
           "tag": "tags"
        },
        "baseurl": "http://yoursite.com/"
    }


The following is an example of a toml config file with the default values: 

    contentdir = "content"
    layoutdir = "layouts"
    publishdir = "public"
    builddrafts = false
    baseurl = "http://yoursite.com/"
    [indexes]
       category = "categories"
       tag = "tags"


