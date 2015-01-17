---
date: 2015-01-08T16:11:23-07:00
menu:
  main:
    parent: troubleshooting
title: Strange EOF error
weight: 5
---

## Trouble: `hugo new` aborts with cryptic EOF error

> I'm running into an issue where I cannot get archetypes working, when running `hugo new showcase/test.md`, for example, I see an `EOF` error thrown by Hugo.
>
> I have set up this test repository to show exactly what I've done, but it is essentially a vanilla installation of Hugo. https://github.com/polds/hugo-archetypes-test
>
> When in that repository, using Hugo v0.12 to run `hugo new -v showcase/test.md`, I see the following output:
>
>     INFO: 2015/01/04 Using config file: /private/tmp/test/config.toml
>     INFO: 2015/01/04 attempting to create  showcase/test.md of showcase
>     INFO: 2015/01/04 curpath: /private/tmp/test/archetypes/showcase.md
>     ERROR: 2015/01/04 EOF
>
> Is there something that I am blatantly missing?

## Solution

Thank you for reporting this issue.  The solution is to add a final newline (i.e. EOL) to the end of your default.md archetype file of your theme.  More discussions happened on the forum here:

* http://discuss.gohugo.io/t/archetypes-not-properly-working-in-0-12/544
* http://discuss.gohugo.io/t/eol-f-in-archetype-files/554

Due to popular demand, Hugo's parser has been enhanced to
accommodate archetype files without final EOL,
thanks to the great work by [@tatsushid](https://github.com/tatsushid),
in the upcoming v0.13 release,

Until then, for us running the stable v0.12 release, please remember to add the final EOL diligently.  <i class="fa fa-smile-o"></i>

## References

* https://github.com/spf13/hugo/issues/776

