
---
date: 2021-07-05
title: "Poll based watching"
description: "Hugo 0.85.0: Polled based alternative when watching for changes and some other nice improvements."
categories: ["Releases"]
---

Hugo `0.85.0` is on the smaller side of releases, but the main new thing it brings should be important to those who need it: Poll based watching the filesystem for changes.

Hugo uses [Fsnotify](https://github.com/fsnotify/fsnotify) to provide native file system notifications. This is still the default, but there may situations where this isn't working. The file may not support it (e.g. NFS), or you get the "too many open files" error and cannot or do not want to increase the `ulimit`. Enable polling by passing the `--poll` flag with an interval:

```bash
hugo server --poll 700ms
```

You can even do "long polling" by passing a long interval:

```bash
hugo server --poll 24h
```

This release represents **23 contributions by 6 contributors** to the main Hugo code base.[@bep](https://github.com/bep) leads the Hugo development with a significant amount of contributions, but also a big shoutout to [@raoulb](https://github.com/raoulb), [@jmooring](https://github.com/jmooring), and [@digitalcraftsman](https://github.com/digitalcraftsman) for their ongoing contributions.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **1 contributions by 1 contributors**.

Hugo now has:

* 52755+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 431+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 370+ [themes](http://themes.gohugo.io/)

## Enhancements

### Other

* Move time notification to after any build errors [04dc469f](https://github.com/gohugoio/hugo/commit/04dc469fbd78d9fe784829f2cba61c8cce982bdb) [@jhollowe](https://github.com/jhollowe) [#8403](https://github.com/gohugoio/hugo/issues/8403)
* Log warning for metadata decode error [07919d1c](https://github.com/gohugoio/hugo/commit/07919d1ccb01733f4c6c5952e59228cecc9b26c8) [@IveGotNorto](https://github.com/IveGotNorto) [#8519](https://github.com/gohugoio/hugo/issues/8519)
* Make the --poll flag a duration [e31b1d19](https://github.com/gohugoio/hugo/commit/e31b1d194655ac3a38fe903ff3995806b129b88a) [@bep](https://github.com/bep) [#8720](https://github.com/gohugoio/hugo/issues/8720)
* Regen CLI docs [43a23239](https://github.com/gohugoio/hugo/commit/43a23239b2e3ad602c06d9af0b648e0304fc8744) [@bep](https://github.com/bep) 
* Add polling as a fallback to native filesystem events in server watch [24ce98b6](https://github.com/gohugoio/hugo/commit/24ce98b6d10b2088af61c15112f5c5ed915a0c35) [@bep](https://github.com/bep) [#8720](https://github.com/gohugoio/hugo/issues/8720)[#6849](https://github.com/gohugoio/hugo/issues/6849)[#7930](https://github.com/gohugoio/hugo/issues/7930)
* Bump github.com/yuin/goldmark v1.3.9 [0019d60f](https://github.com/gohugoio/hugo/commit/0019d60f67b6c4dde085753641a917fcd0aa4c76) [@bep](https://github.com/bep) [#8727](https://github.com/gohugoio/hugo/issues/8727)
* Add module.import.noMounts config [40dfdd09](https://github.com/gohugoio/hugo/commit/40dfdd09521bcb8f56150e6791d60445198f27ab) [@bep](https://github.com/bep) [#8708](https://github.com/gohugoio/hugo/issues/8708)
* Use value type for module.Time [3a6dc6d3](https://github.com/gohugoio/hugo/commit/3a6dc6d3f423c4acb79ef21b5a76e616fa2c9477) [@bep](https://github.com/bep) 
* Add version time to "hugo config mounts" [6cd2110a](https://github.com/gohugoio/hugo/commit/6cd2110ab295f598907a18da91e34d31407c1d9d) [@bep](https://github.com/bep) 
* Add some more info to "hugo config mounts" [6a365c27](https://github.com/gohugoio/hugo/commit/6a365c2712c7607e067e192d213b266f0c88d0f3) [@bep](https://github.com/bep) 
* Update to Minify v2.9.18 [d9bdd37d](https://github.com/gohugoio/hugo/commit/d9bdd37d35ccd436b4dd470ef99efa372a6a086b) [@bep](https://github.com/bep) [#8693](https://github.com/gohugoio/hugo/issues/8693)
* Remove credit from release notes [b2eaf4c8](https://github.com/gohugoio/hugo/commit/b2eaf4c8c2e31aa1c1bc4a2c0061f661e01d2de1) [@digitalcraftsman](https://github.com/digitalcraftsman) 
* Rename Header(s) to Heading(s) in ToC struct [a7e3da24](https://github.com/gohugoio/hugo/commit/a7e3da242f98d4799dad013d7ba2f285717640d6) [@bep](https://github.com/bep) 

## Fixes

### Other

* Fix tab selection of disabled items in internal pagination template [f75f9007](https://github.com/gohugoio/hugo/commit/f75f90079a6f2a239c8186faba5db5dbe6e36cb6) [@raoulb](https://github.com/raoulb) 
* Fix panic when theme has permalinks config [e451b984](https://github.com/gohugoio/hugo/commit/e451b984cfb45b54a3972cefa59a02d50b0b0fd2) [@bep](https://github.com/bep) [#8724](https://github.com/gohugoio/hugo/issues/8724)
* Fix Cloudflare vs Netlify cache dir issue [4c8552b1](https://github.com/gohugoio/hugo/commit/4c8552b11477141777101e0e0609dd1f32d191e9) [@bep](https://github.com/bep) [#8714](https://github.com/gohugoio/hugo/issues/8714)
* Fix date format in schema and opengraph templates [34e4742f](https://github.com/gohugoio/hugo/commit/34e4742f0caab0d3eb9efd00fce4157d112617b5) [@jmooring](https://github.com/jmooring) [#8671](https://github.com/gohugoio/hugo/issues/8671)
* Fix Netlify default cache dir logic [6c8c0c8b](https://github.com/gohugoio/hugo/commit/6c8c0c8b6a0b39b91de44d72a7bd1cd49534a0f1) [@bep](https://github.com/bep) [#8710](https://github.com/gohugoio/hugo/issues/8710)
* Fix handling of invalid OS env config overrides [49fedbc5](https://github.com/gohugoio/hugo/commit/49fedbc51cafa64e4eb0eae9fb79ccbe2d4c6774) [@bep](https://github.com/bep) [#8709](https://github.com/gohugoio/hugo/issues/8709)
* Fix config handling with empty config entries after merge [19aa95fc](https://github.com/gohugoio/hugo/commit/19aa95fc7f4cd58dcc8a8ff075762cfc86d41dc3) [@bep](https://github.com/bep) [#8701](https://github.com/gohugoio/hugo/issues/8701)
* Fix config loading for "hugo mod init" [923dd9d1](https://github.com/gohugoio/hugo/commit/923dd9d1c1f649142f3f377109318b07e0f44d5d) [@bep](https://github.com/bep) [#8697](https://github.com/gohugoio/hugo/issues/8697)
* Fix language menu config regression [093dacab](https://github.com/gohugoio/hugo/commit/093dacab29a3c6fc363408453d0bc3b1fc159ad5) [@bep](https://github.com/bep) [#8672](https://github.com/gohugoio/hugo/issues/8672)
* Fix merge of config with map[string]string values. [4a9d408f](https://github.com/gohugoio/hugo/commit/4a9d408fe0bbf4c563546e35d2be7ade4e920c4c) [@bep](https://github.com/bep) [#8679](https://github.com/gohugoio/hugo/issues/8679)





