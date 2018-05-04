
---
date: 2017-09-26
title: "Hugo 0.29: Template Metrics"
description: "Makes it easy to find template bottle necks."
categories: ["Releases"]
images:
- images/blog/hugo-29-poster.png
---
	
Hugo `0.29` brings Template Metrics by [@moorereason](https://github.com/moorereason). Hugo is very fast, but it is still possible to write ineffective templates. Now these should be easy to identify. Just run:

```bash
hugo --templateMetrics
```
Now, that was the tasty carrot. The real reason this release comes so fast after the last one is to change the default value for the new `noHTTPCache` flag, which gives away too much performance to make sense as a default value.

Hugo now has:

* 19817+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 454+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 180+ [themes](http://themes.gohugo.io/)

## Notes
* Make `noHTTPCache` default false [e94d4f01](https://github.com/gohugoio/hugo/commit/e94d4f0177852b357f40fb9686a0ff3667d86351) [@bep](https://github.com/bep) 

## Enhancements

### Templates
* Add simple template metrics feature [b4a14c25](https://github.com/gohugoio/hugo/commit/b4a14c25fe85c41b79497be27ead128502a4dd7b) [@moorereason](https://github.com/moorereason) 
* Set Metrics at creation time [b5e1dc58](https://github.com/gohugoio/hugo/commit/b5e1dc5892f81da798d0d4e964a1f3328532f45e) [@bep](https://github.com/bep) 
* Fix sort order [d3681f51](https://github.com/gohugoio/hugo/commit/d3681f51c08fb11e8addcf9f0b484848d20d46cc) [@bep](https://github.com/bep) 
* Add math.Ceil, Floor, and Round to method mappings [8a69d235](https://github.com/gohugoio/hugo/commit/8a69d2356703d9f2fcb75bce0ae514e70ebd8e01) [@moorereason](https://github.com/moorereason) 

### Other
* Split go build in Dockerfile [d9697e27](https://github.com/gohugoio/hugo/commit/d9697e275ecb038958b3dcea2b43e11dcba28fc9) [@tjamet](https://github.com/tjamet) 
* Update Dockerfile to benefit build cache [09d960f1](https://github.com/gohugoio/hugo/commit/09d960f17396eb7fd2c8fe6527db9503d59f0b4f) [@tjamet](https://github.com/tjamet) 
* Add git to snap package for GitInfo [a3a3f5b8](https://github.com/gohugoio/hugo/commit/a3a3f5b86114213a23337499551f000662b26022) [@ghalse](https://github.com/ghalse) 








