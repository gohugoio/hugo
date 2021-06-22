
---
date: 2021-06-18
title: "Config Revamp"
description: "Hugo 0.84.0: Deep merge of theme configuration, config dir support now also in themes/modules, HTTP header support in getJSON, and more."
categories: ["Releases"]
---

**This release brings several configuration fixes and improvements that will be especially useful for themes.**

## Deep merge of theme Params

One of the most common complaints from Hugo theme owners/users has been about the configuration handling. Hugo has up until now only performed a shallow merge of theme `params` into the configuration.

With that, given this example from a theme configuration:

```toml
[params]
[params.colours]
blue="#337DFF"
green="#68FF33"
red="#FF3358"
```

If you would like to use the above theme, but want a different shade of red, you earlier had to copy the entire block, even the colours you're totally happy with. This was painful with even the simplest setup.

Now you can just override the `params` keys you want to change, e.g.:

```toml
[params]
[params.colours]
red="#fc0f03"
```

For more information, and especially about the way you can opt out of the above behaviour, see [Merge Configuration from Themes](https://gohugo.io/getting-started/configuration/#merge-configuration-from-themes).

## Themes now support the config directory

Now both the project and themes/modules can store its configuration in both the top level config file (e.g. `config.toml`) or in the `config` directory. See [Configuration Directory](https://gohugo.io/getting-started/configuration/#configuration-directory).

## HTTP headers in getJSON/getCSV

`getJSON` now supports custom HTTP headers. This has been a big limitation in Hugo, especially considering the [Authorization](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization) header.

We have updated the internal Instagram shortcode to pass the access token in a header:

```
{{ $hideCaption := cond (eq (.Get 1) "hidecaption") "1" "0" }}
{{ $headers := dict "Authorization" (printf "Bearer %s" $accessToken) }}
{{ with getJSON "https://graph.facebook.com/v8.0/instagram_oembed/?url=https://instagram.com/p/" $id "/&hidecaption=" $hideCaption $headers }}
    {{ .html | safeHTML }}
{{ end }}
```

 Also see the discussion [this issue](https://github.com/gohugoio/hugo/issues/7879) about the access token above.

## New erroridf template func

Sometimes, especially when creating themes, it is useful to be able to let the user decide if an error situation is critical enough to fail the build. The new `erroridf` produces `ERROR` log statements that can be toggled off:

```html
{{ erroridf "some-custom-id" "Some error message." }}
```

Will log:

```
ERROR: Some error message.
If you feel that this should not be logged as an ERROR, you can ignore it by adding this to your site config:
ignoreErrors = ["some-custom-id"]
```
## Stats

This release represents **46 contributions by 11 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@jmooring](https://github.com/jmooring), [@satotake](https://github.com/satotake), and [@Seirdy](https://github.com/Seirdy) for their ongoing contributions.
And a big thanks to [@digitalcraftsman](https://github.com/digitalcraftsman) for his relentless work on keeping the themes site in pristine condition.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **20 contributions by 10 contributors**. A special thanks to [@salim-b](https://github.com/salim-b), [@bep](https://github.com/bep), [@thomasjsn](https://github.com/thomasjsn), and [@lucasew](https://github.com/lucasew) for their work on the documentation site.


Hugo now has:

* 52487+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 432+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)


## Notes

* We now do deep merging of `params` from theme config(s). That is you most likely what you want, but [Merge Configuration from Themes](https://gohugo.io/getting-started/configuration/#merge-configuration-from-themes) describes how you can get the old behaviour back.

## Enhancements

### Templates

* Rename err-missing-instagram-accesstoken => error-missing-instagram-accesstoken [9096842b](https://github.com/gohugoio/hugo/commit/9096842b0494166e401cc08a70b93ae2ee19a198) [@bep](https://github.com/bep) 
* Add a terse pagination template variant to improve performance [73483d0f](https://github.com/gohugoio/hugo/commit/73483d0f9eb46838d41640f88cc05c1d16811dc5) [@jmooring](https://github.com/jmooring) [#8599](https://github.com/gohugoio/hugo/issues/8599)
* Add erroridf template func [f55d2f43](https://github.com/gohugoio/hugo/commit/f55d2f43769053b80b419a690554e747dc5dcede) [@bep](https://github.com/bep) [#8613](https://github.com/gohugoio/hugo/issues/8613)
* Print response body on HTTP errors [282f1aa3](https://github.com/gohugoio/hugo/commit/282f1aa3db9f6420fdd360e46db1ffadd5b083a1) [@bep](https://github.com/bep) 
* Misc header improvements, tests, allow multiple headers of same key [fcd63de3](https://github.com/gohugoio/hugo/commit/fcd63de3a54fadcd30972654d8eb86dc4d889784) [@bep](https://github.com/bep) [#5617](https://github.com/gohugoio/hugo/issues/5617)
* Allows user-defined HTTP headers with getJSON and getCSV [150d7573](https://github.com/gohugoio/hugo/commit/150d75738b54acddc485d363436757189144da6a) [@chamberlainpj](https://github.com/chamberlainpj) [#5617](https://github.com/gohugoio/hugo/issues/5617)
* Allow 'Querify' to take lone slice/interface argument [c46fc838](https://github.com/gohugoio/hugo/commit/c46fc838a9320adfc6532b1b543e903c48b3b4cb) [@importhuman](https://github.com/importhuman) [#6735](https://github.com/gohugoio/hugo/issues/6735)

### Output

* Make WebAppManifestFormat NotAlternative=true [643b6719](https://github.com/gohugoio/hugo/commit/643b671931ed5530855e7d4819896790bf3f6c28) [@bep](https://github.com/bep) [#8624](https://github.com/gohugoio/hugo/issues/8624)
* Adjust  test assertion [ab4e1dfa](https://github.com/gohugoio/hugo/commit/ab4e1dfa4eebe0ac18f1d1f60a9647cbb7b41d7f) [@bep](https://github.com/bep) [#8625](https://github.com/gohugoio/hugo/issues/8625)
* support application/manifest+json [02f31897](https://github.com/gohugoio/hugo/commit/02f31897b4f7252154850a65c900e88e0b237fa3) [@Seirdy](https://github.com/Seirdy) [#8624](https://github.com/gohugoio/hugo/issues/8624)

### Other

* Regenerate docs helper [be6b901c](https://github.com/gohugoio/hugo/commit/be6b901cf7d07238337334e6b6d886a7b039f5e6) [@bep](https://github.com/bep) 
* Regenerate docshelper [402da3f8](https://github.com/gohugoio/hugo/commit/402da3f8f327f97302c4b5d69cd4832a94bd189b) [@bep](https://github.com/bep) 
* Implement configuration in a directory for modules [bb2aa087](https://github.com/gohugoio/hugo/commit/bb2aa08709c812a5be29922a1a7f4d814e200cab) [@bep](https://github.com/bep) [#8654](https://github.com/gohugoio/hugo/issues/8654)
* Update github.com/alecthomas/chroma v0.9.1 => v0.9.2 [3aa7f0b2](https://github.com/gohugoio/hugo/commit/3aa7f0b27fc736b4c32adbb1fc1fc7fbefd6efd9) [@bep](https://github.com/bep) [#8658](https://github.com/gohugoio/hugo/issues/8658)
* Run go mod tidy [9b870aa7](https://github.com/gohugoio/hugo/commit/9b870aa788ab1b5159bc836fbac6e60a29bee329) [@bep](https://github.com/bep) 
* Split out the puthe path/filepath functions into common/paths [93aad3c5](https://github.com/gohugoio/hugo/commit/93aad3c543828efca2adeb7f96cf50ae29878593) [@bep](https://github.com/bep) [#8654](https://github.com/gohugoio/hugo/issues/8654)
* Update to Goldmark v1.3.8 [8eafe084](https://github.com/gohugoio/hugo/commit/8eafe0845d66efd3cf442a8ed89a6da5c1d3117b) [@jmooring](https://github.com/jmooring) [#8648](https://github.com/gohugoio/hugo/issues/8648)
* Do not read config from os.Environ when running tests [31fb29fb](https://github.com/gohugoio/hugo/commit/31fb29fb3f306678f3697e05bbccefb2078d7f78) [@bep](https://github.com/bep) [#8655](https://github.com/gohugoio/hugo/issues/8655)
* Set a dummy Instagram token [a886dd53](https://github.com/gohugoio/hugo/commit/a886dd53b80322e1edf924f2ede4d4ea037c5baf) [@bep](https://github.com/bep) 
* Regenerate docs helper [a91cd765](https://github.com/gohugoio/hugo/commit/a91cd7652f7559492b070dbe02fe558348f3d0b6) [@bep](https://github.com/bep) 
* Update to Go 1.16.5, Goreleaser 0.169.0 [552cef5c](https://github.com/gohugoio/hugo/commit/552cef5c576ae4dbf4626f77f3c8b15b42a9e7f3) [@bep](https://github.com/bep) [#8619](https://github.com/gohugoio/hugo/issues/8619)[#8263](https://github.com/gohugoio/hugo/issues/8263)
* Upgrade Instagram shortcode [9b5debe4](https://github.com/gohugoio/hugo/commit/9b5debe4b820132759cfdf7bff7fe9c1ad0a6bb1) [@bep](https://github.com/bep) [#7879](https://github.com/gohugoio/hugo/issues/7879)
* Set modTime at creation time [06d29542](https://github.com/gohugoio/hugo/commit/06d295427f798da85de469924fd10f58c0de9a58) [@bep](https://github.com/bep) [#6161](https://github.com/gohugoio/hugo/issues/6161)
* Add math.Max and math.Min [01758f99](https://github.com/gohugoio/hugo/commit/01758f99b915f34fe7ca4621e4d1ee09efe385b1) [@jmooring](https://github.com/jmooring) [#8583](https://github.com/gohugoio/hugo/issues/8583)
* Catch incomplete shortcode error [845a7ba4](https://github.com/gohugoio/hugo/commit/845a7ba4fc30c61842148d67d31d0fa3db8f40b9) [@satotake](https://github.com/satotake) [#6866](https://github.com/gohugoio/hugo/issues/6866)
* Use SPDX license identifier [10f60de8](https://github.com/gohugoio/hugo/commit/10f60de89a5a53528f1e3a47a77224e5c7915e4e) [@jmooring](https://github.com/jmooring) [#8555](https://github.com/gohugoio/hugo/issues/8555)
* Cache and copy Menu for sorting [785a31b5](https://github.com/gohugoio/hugo/commit/785a31b5b84643f4769f9bd363599cbcce86f098) [@satotake](https://github.com/satotake) [#7594](https://github.com/gohugoio/hugo/issues/7594)
* Update to LibSASS 3.6.5 [bc1e0528](https://github.com/gohugoio/hugo/commit/bc1e05286a96d08ad02ad200d6a4076bb01c486e) [@bep](https://github.com/bep) 
* Make the HTML element collector more robust [f518b4f7](https://github.com/gohugoio/hugo/commit/f518b4f71e1a61b09d660b5c284121ebf3b3b86b) [@bep](https://github.com/bep) [#8530](https://github.com/gohugoio/hugo/issues/8530)
* Make the HTML element collector more robust" [dc6b7a75](https://github.com/gohugoio/hugo/commit/dc6b7a75ff5b7fcb8a0b0e3f7ed406422d847624) [@bep](https://github.com/bep) 
* Get the collector in line with the io.Writer interface" [3f515f0e](https://github.com/gohugoio/hugo/commit/3f515f0e3395b24776ae24045b846ff2b33b8906) [@bep](https://github.com/bep) 
* Get the collector in line with the io.Writer interface [a9bcd381](https://github.com/gohugoio/hugo/commit/a9bcd38181ceb79afba82adcd4de1aebf571e74c) [@bep](https://github.com/bep) 
* Make the HTML element collector more robust [ef0f1a72](https://github.com/gohugoio/hugo/commit/ef0f1a726901d6c614040cfc2d7e8f9a2ca97816) [@bep](https://github.com/bep) [#8530](https://github.com/gohugoio/hugo/issues/8530)
* Add Scratch.DeleteInMap [abbc99d4](https://github.com/gohugoio/hugo/commit/abbc99d4c60b102e2779e4362ceb433095719384) [@meehawk](https://github.com/meehawk) [#8504](https://github.com/gohugoio/hugo/issues/8504)
* Display version when building site (#8533) [76c95f55](https://github.com/gohugoio/hugo/commit/76c95f55a5d18290baa7f23667161d4af9fb9b53) [@jmooring](https://github.com/jmooring) [#8531](https://github.com/gohugoio/hugo/issues/8531)
* Update querify function description and examples [2c7f5b62](https://github.com/gohugoio/hugo/commit/2c7f5b62f6c1fa1c7b3cf2c1f3a1663b18e75004) [@jmooring](https://github.com/jmooring) 
* Change SetEscapeHTML to false [504c78da](https://github.com/gohugoio/hugo/commit/504c78da4b5020e1fd13a1195ad38a9e85f8289a) [@peaceiris](https://github.com/peaceiris) [#8512](https://github.com/gohugoio/hugo/issues/8512)
* Add a benchmark [b660ea8d](https://github.com/gohugoio/hugo/commit/b660ea8d545d6ba5479dd28a670044d57e5d196f) [@bep](https://github.com/bep) 
* Update dependency list [64f88f30](https://github.com/gohugoio/hugo/commit/64f88f3011de5a510d8e6d6bad8ac4a091b11c0c) [@bep](https://github.com/bep) 

## Fixes

### Templates

* Fix countwords to handle special chars [7a2c10ae](https://github.com/gohugoio/hugo/commit/7a2c10ae60f096dacee4b44e0c8ae0a1b66ae033) [@ResamVi](https://github.com/ResamVi) [#8479](https://github.com/gohugoio/hugo/issues/8479)

### Other

* Fix fill with smartcrop sometimes returning 0 bytes images [5af045eb](https://github.com/gohugoio/hugo/commit/5af045ebab109d3e5501b8b6d9fd448840c96c9a) [@bep](https://github.com/bep) [#7955](https://github.com/gohugoio/hugo/issues/7955)
* Misc config loading fixes [d392893c](https://github.com/gohugoio/hugo/commit/d392893cd73dc00c927f342778f6dca9628d328e) [@bep](https://github.com/bep) [#8633](https://github.com/gohugoio/hugo/issues/8633)[#8618](https://github.com/gohugoio/hugo/issues/8618)[#8630](https://github.com/gohugoio/hugo/issues/8630)[#8591](https://github.com/gohugoio/hugo/issues/8591)[#6680](https://github.com/gohugoio/hugo/issues/6680)[#5192](https://github.com/gohugoio/hugo/issues/5192)
* Fix nested OS env config override when parent does not exist [12530519](https://github.com/gohugoio/hugo/commit/12530519d8fb4513c9c18a6494099b7dff8e4fd4) [@bep](https://github.com/bep) [#8618](https://github.com/gohugoio/hugo/issues/8618)
* Fix invalid timestamp of the "public" folder [26ae12c0](https://github.com/gohugoio/hugo/commit/26ae12c0c64b847d24bde60d7d710ea2efcb40d4) [@anthonyfok](https://github.com/anthonyfok) [#6161](https://github.com/gohugoio/hugo/issues/6161)
* Fix env split to allow = character in  values [ee733085](https://github.com/gohugoio/hugo/commit/ee733085b7f5d3f2aef1667901ab6ecb8041d699) [@xqbumu](https://github.com/xqbumu) [#8589](https://github.com/gohugoio/hugo/issues/8589)
* Fix warning regression in i18n [ececd1b1](https://github.com/gohugoio/hugo/commit/ececd1b122c741567a80acd8d60ccd6356fa5323) [@bep](https://github.com/bep) [#8492](https://github.com/gohugoio/hugo/issues/8492)





