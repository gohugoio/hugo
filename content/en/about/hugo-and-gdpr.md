

---
title: Hugo and the  General Data Protection Regulation (GDPR)
linktitle: Hugo and GDPR
description: About how to configure your Hugo site to meet the new regulations.
date: 2018-05-25
layout: single
keywords: ["GDPR", "Privacy", "Data Protection"]
menu:
  docs:
    parent: "about"
    weight: 5
weight: 5
sections_weight: 5
draft: false
aliases: [/privacy/,/gdpr/]
toc: true
---


 General Data Protection Regulation ([GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation)) is a regulation in EU law on data protection and privacy for all individuals within the European Union and the European Economic Area. It became enforceable on 25 May 2018.

 **Hugo is a static site generator. By using Hugo you are already standing on very solid ground. Static HTML files on disk are much easier to reason about compared to server and database driven web sites.**

 But even static websites can integrate with external services, so from version `0.41`, Hugo provides a **Privacy Config** that covers the relevant built-in templates.

 Note that:

 * These settings have their defaults setting set to _off_, i.e. how it worked before Hugo `0.41`. You must do your own evaluation of your site and apply the appropriate settings.
 * These settings work with the [internal templates](/templates/internal/). Some theme may contain custom templates for embedding services like Google Analytics. In that case these options have no effect.
 * We will continue this work and improve this further in future Hugo versions.

## All Privacy Settings

Below are all privacy settings and their default value. These settings need to be put in your site config (e.g. `config.toml`).

 {{< code-toggle file="config">}}
[privacy]
[privacy.disqus]
disable = false
[privacy.googleAnalytics]
disable = false
respectDoNotTrack = false
anonymizeIP = false
useSessionStorage = false
[privacy.instagram]
disable = false
simple = false
[privacy.twitter]
disable = false
enableDNT = false
simple = false
[privacy.vimeo]
disable = false
simple = false
[privacy.youtube]
disable = false
privacyEnhanced = false
{{< /code-toggle >}}


## Disable All Services

An example Privacy Config that disables all the relevant services in Hugo. With this configuration, the other settings will not matter.

 {{< code-toggle file="config">}}
[privacy]
[privacy.disqus]
disable = true
[privacy.googleAnalytics]
disable = true
[privacy.instagram]
disable = true
[privacy.twitter]
disable = true
[privacy.vimeo]
disable = true
[privacy.youtube]
disable = true
{{< /code-toggle >}}

## The Privacy Settings Explained

### GoogleAnalytics

anonymizeIP
: Enabling this will make it so the users' IP addresses are anonymized within Google Analytics.

respectDoNotTrack
: Enabling this will make the GA templates respect the "Do Not Track" HTTP header.

useSessionStorage
: Enabling this will disable the use of Cookies and use Session Storage to Store the GA Client ID.

### Instagram

simple
: If simple mode is enabled, a static and no-JS version of the Instagram image card will be built. Note that this only supports image cards and the image itself will be fetched from Instagram's servers.

**Note:** If you use the _simple mode_ for Instagram and a site styled with Bootstrap 4, you may want to disable the inlines styles provided by Hugo:

 {{< code-toggle file="config">}}
[services]
[services.instagram]
disableInlineCSS = true
{{< /code-toggle >}}

### Twitter

enableDNT
: Enabling this for the twitter/tweet shortcode, the tweet and its embedded page on your site are not used for purposes that include personalized suggestions and personalized ads.

simple
: If simple mode is enabled, a static and no-JS version of a tweet will be built.


**Note:** If you use the _simple mode_ for Twitter, you may want to disable the inlines styles provided by Hugo:

 {{< code-toggle file="config">}}
[services]
[services.twitter]
disableInlineCSS = true
{{< /code-toggle >}}

### YouTube

privacyEnhanced
: When you turn on privacy-enhanced mode, YouTube won’t store information about visitors on your website unless the user plays the embedded video.

### Vimeo

simple
: If simple mode is enabled, the video thumbnail is fetched from Vimeo's servers and it is overlayed with a play button. If the user clicks to play the video, it will open in a new tab directly on Vimeo's website.

