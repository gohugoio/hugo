
---
date: 2021-08-03
title: "Localized Time and Dates and Numbers"
description: "Hugo 0.87.0 brings time zone support, localized time and dates and numbers backed by CLDR, and more."
categories: ["Releases"]
---

Hugo `0.87` brings two long sought-after features: Default time zone support (per language if needed) for dates without zone offset or location info, and localized time and dates and numbers (backed by [CLDR](https://en.wikipedia.org/wiki/Common_Locale_Data_Repository)).

For more information, see:

* The [time zone config](https://gohugo.io/getting-started/configuration/#timezone) documentation.
* The [time.Format](https://gohugo.io/functions/dateformat/) documentation. This function will now give you localized dates (with weekdays and months in the current language). It supports all of Go's layout syntax, but also some predefined constants, e.g. `{{ .Date | time.Format ":date_long" }}`.
* A set of new [localized number formatting ](https://gohugo.io/functions/lang/) functions.

Also in this release, we have switched to using [go-toml](https://github.com/pelletier/go-toml) for all things TOML in Hugo. A big thanks to [@pelletier](https://github.com/pelletier) for his work on the `v2` version. It's both faster than what we had and [TOML v1.0.0](https://toml.io/en/v1.0.0) compliant.

This release represents **40 contributions by 4 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@dependabot[bot]](https://github.com/apps/dependabot), [@digitalcraftsman](https://github.com/digitalcraftsman), and [@jmooring](https://github.com/jmooring) for their ongoing contributions.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **1 contributions by 1 contributors**.

Hugo now has:

* 53261+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 430+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 395+ [themes](http://themes.gohugo.io/)


## Notes

* Deprecate Blackfriday and fix a potential deadlock in config [c7252224](https://github.com/gohugoio/hugo/commit/c7252224c4fecfe8321f31b901e2510d98b657c4) [@bep](https://github.com/bep) [#8792](https://github.com/gohugoio/hugo/issues/8792)[#8791](https://github.com/gohugoio/hugo/issues/8791)

## Enhancements

### Templates

* Adjust tests to handle matching local time zones [9ff17c33](https://github.com/gohugoio/hugo/commit/9ff17c332405da5830cef9b3711706b1fc9a7444) [@bep](https://github.com/bep) [#8843](https://github.com/gohugoio/hugo/issues/8843)
* Add new localized versions of lang.FormatNumber etc. [7907d24b](https://github.com/gohugoio/hugo/commit/7907d24ba16fc5a80930c1aabf5144e684ff7f29) [@bep](https://github.com/bep) [#8820](https://github.com/gohugoio/hugo/issues/8820)

### Other

* Make sure module config loading errors have file positioning info [d70c4857](https://github.com/gohugoio/hugo/commit/d70c485707edfd445bcfc0e84181bc15eb146e76) [@bep](https://github.com/bep) [#8845](https://github.com/gohugoio/hugo/issues/8845)
* Remove superflous replace statement [7aaaf7e3](https://github.com/gohugoio/hugo/commit/7aaaf7e33afd05d2c74d74fbbfbd34d55e8129eb) [@bep](https://github.com/bep) 
* Reduce binary size vs locale, update to CLDR v36.1 [3a966555](https://github.com/gohugoio/hugo/commit/3a96655592d0b0db4126f20ca717d553dda9c4ed) [@bep](https://github.com/bep) [#8839](https://github.com/gohugoio/hugo/issues/8839)[#8841](https://github.com/gohugoio/hugo/issues/8841)
* Update github.com/tdewolff/minify/v2 v2.9.20 => v2.9.21 [9a7383ca](https://github.com/gohugoio/hugo/commit/9a7383caf3945b9b11db2b108003f87e2e8b6a3a) [@jmooring](https://github.com/jmooring) [#8831](https://github.com/gohugoio/hugo/issues/8831)
* Fail on invalid time zone [4d221ce4](https://github.com/gohugoio/hugo/commit/4d221ce468a1209ee9dd6cbece9d1273dad6a29b) [@bep](https://github.com/bep) [#8832](https://github.com/gohugoio/hugo/issues/8832)
* Improve handling of <nil> Params [e3dc5240](https://github.com/gohugoio/hugo/commit/e3dc5240f01fd5ec67643e40f27c026d707da110) [@bep](https://github.com/bep) [#8825](https://github.com/gohugoio/hugo/issues/8825)
* Merge branch 'release-0.86.1' [268065cb](https://github.com/gohugoio/hugo/commit/268065cb2d8339392766a23703beaf7cc49d6b5c) [@bep](https://github.com/bep) 
* bump github.com/evanw/esbuild from 0.12.16 to 0.12.17 [e90b3591](https://github.com/gohugoio/hugo/commit/e90b3591a155d1266a86c9490886720740b9d62e) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/getkin/kin-openapi from 0.67.0 to 0.68.0 [4b7da6a9](https://github.com/gohugoio/hugo/commit/4b7da6a9d720ed5ab4b45d6aa3b0b7af4683d02f) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Go back to WARNING for Page deprecations [726fe9c3](https://github.com/gohugoio/hugo/commit/726fe9c3c97a9c979dc7862e7f226fc5ec1341de) [@bep](https://github.com/bep) 
* Handle toml.LocalDate and toml.LocalDateTime in front matter [b5de37ee](https://github.com/gohugoio/hugo/commit/b5de37ee793c01f2acccdea7119be05c4182723f) [@bep](https://github.com/bep) [#8801](https://github.com/gohugoio/hugo/issues/8801)
* Upgrade github.com/pelletier/go-toml/v2 v2.0.0-beta.3 => v2.0.0-beta.3.0.20210727221244-fa0796069526 [bf301daf](https://github.com/gohugoio/hugo/commit/bf301daf158e5e9673ad5f457ea3a264315942b5) [@bep](https://github.com/bep) 
* Switch to go-toml v2 [a3701e09](https://github.com/gohugoio/hugo/commit/a3701e09313695d4a0f6fb0eb7844c1a4befc07a) [@bep](https://github.com/bep) [#8801](https://github.com/gohugoio/hugo/issues/8801)
* bump github.com/tdewolff/minify/v2 from 2.9.19 to 2.9.20 [40b6016c](https://github.com/gohugoio/hugo/commit/40b6016cf3f7aac541b042d32e3a162411fd9cd0) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Add a TOML front matter benchmark [7e130571](https://github.com/gohugoio/hugo/commit/7e1305710f08d26d9214abb5410ccd675e59a6e9) [@bep](https://github.com/bep) 
* Add timezone support for front matter dates without one [efa5760d](https://github.com/gohugoio/hugo/commit/efa5760db5ef39ae084bfccb5b8f756c7b117a2a) [@bep](https://github.com/bep) [#8810](https://github.com/gohugoio/hugo/issues/8810)
* Localize time.Format [a57dda85](https://github.com/gohugoio/hugo/commit/a57dda854b5efd3429af5f0b1564fc9d9d5439b9) [@bep](https://github.com/bep) [#8797](https://github.com/gohugoio/hugo/issues/8797)
* bump github.com/getkin/kin-openapi from 0.61.0 to 0.67.0 [f9afba93](https://github.com/gohugoio/hugo/commit/f9afba933579de07d2d2e36a457895ec5f1b7f01) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/spf13/cast from 1.3.1 to 1.4.0 [a5d2ba42](https://github.com/gohugoio/hugo/commit/a5d2ba429d34004efd3c6b82c1bcb130c85aca9c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump google.golang.org/api from 0.45.0 to 0.51.0 [31972f36](https://github.com/gohugoio/hugo/commit/31972f3647b284eea1a66a2e27ed42d04a391a7a) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/sanity-io/litter from 1.5.0 to 1.5.1 [2e58782f](https://github.com/gohugoio/hugo/commit/2e58782f96972487dc5e5ba91d0256ec6e86dad7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/mattn/go-isatty from 0.0.12 to 0.0.13 [7b68f652](https://github.com/gohugoio/hugo/commit/7b68f6524d24d450330cbe4a2380301e66abee4a) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/spf13/cobra from 1.1.3 to 1.2.1 [81265af2](https://github.com/gohugoio/hugo/commit/81265af2cccd3247df87f05eebf8907a14e978a4) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/mitchellh/mapstructure from 1.3.3 to 1.4.1 [c102c971](https://github.com/gohugoio/hugo/commit/c102c9719b3a29406ef59dc18eca6bd280e4dc43) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/google/go-cmp from 0.5.5 to 0.5.6 [7c0f904f](https://github.com/gohugoio/hugo/commit/7c0f904f29c41e8782b44a37fd4e98e441cd2b2c) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/mitchellh/hashstructure from 1.0.0 to 1.1.0 [b2fbd4d1](https://github.com/gohugoio/hugo/commit/b2fbd4d13a47ce3f6a56f08d0bda77e16793de72) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/gobuffalo/flect from 0.2.2 to 0.2.3 [90041d1b](https://github.com/gohugoio/hugo/commit/90041d1b6d4eeb91ea085f5a97b02887159a655b) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/pelletier/go-toml from 1.9.0 to 1.9.3 [05047096](https://github.com/gohugoio/hugo/commit/05047096f52e43ff09acbc50616441bb42a1c6f7) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/aws/aws-sdk-go from 1.38.23 to 1.40.8 [a469156e](https://github.com/gohugoio/hugo/commit/a469156ea4ad023aa4fda0d3fb657ce003412abb) [@dependabot[bot]](https://github.com/apps/dependabot) 
* bump github.com/tdewolff/minify/v2 from 2.9.18 to 2.9.19 [18fdd85b](https://github.com/gohugoio/hugo/commit/18fdd85bcc4ac2d9a33546dca8a0a24f63987361) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Update github.com/evanw/esbuild v0.11.16 => v0.12.16 [aeb1935d](https://github.com/gohugoio/hugo/commit/aeb1935d44eb258a794c8f055eedaf3a7655a3ad) [@bep](https://github.com/bep) 
* Update github.com/yuin/goldmark v1.3.9 => v1.4.0 [e09d7882](https://github.com/gohugoio/hugo/commit/e09d7882c88812bedb2c2e66b68c7eed21213dbc) [@bep](https://github.com/bep) 
* bump github.com/frankban/quicktest from 1.12.0 to 1.13.0 [15c0eed0](https://github.com/gohugoio/hugo/commit/15c0eed0487598ac1e5a6fff167b19031b6595bc) [@dependabot[bot]](https://github.com/apps/dependabot) 
* Bump all long-living deprecations to ERRORs [91cbb963](https://github.com/gohugoio/hugo/commit/91cbb963020ac2aead68ff2bbd7e9077d5558abd) [@bep](https://github.com/bep) 

## Fixes

### Other

* Fix error handling for the time func alias [6c70e1f2](https://github.com/gohugoio/hugo/commit/6c70e1f22f365322d5f754302e110c9ed716b215) [@bep](https://github.com/bep) [#8835](https://github.com/gohugoio/hugo/issues/8835)
* Fix a potential deadlock in config reading [94b616bd](https://github.com/gohugoio/hugo/commit/94b616bdfad177daa99f5e87535943f509198f6f) [@bep](https://github.com/bep) [#8791](https://github.com/gohugoio/hugo/issues/8791)
* Fix theme count in release notes [a352d19d](https://github.com/gohugoio/hugo/commit/a352d19d881474f53d01791be4febd305453a9d6) [@digitalcraftsman](https://github.com/digitalcraftsman) 





