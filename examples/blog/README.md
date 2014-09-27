Hugo Example Blog
=================

This repository provides a fully-working example of a [Hugo](https://github.com/spf13/hugo)-powered blog. Many
Hugo-specific features are used as a way to see them in action, and hopefully ease the learning curve for creating your
very own site with Hugo.

Features
--------

- Recent Posts at main index
- Indexes for `tags` and `categories`
- Post information block, with links for all `tags` and `categories` post belongs to
- [Bootstrap 3](http://getbootstrap.com/) ready
  - Currently using the [Yeti](http://bootswatch.com/yeti/) theme from http://bootswatch.com/

Common things that should be added in the near future *(pull requests are welcome!)*:

- Disqus integration
- More content types to demonstrate different layout methods
  - About Me
  - Contact

Getting Started
---------------

To get started, you should simply fork or clone this repository! That's definitely an important first step.

[Install Hugo](http://gohugo.io/overview/installing) in a way that best suits your environment and comfort level.

Edit `config.toml` and change the default properties to suit your own information. This is not required to run the
example, but this is the global configuration file and you're going to need to use it eventually. Start here!

In a command prompt or terminal, navigate to the path that contains your `config.toml` file and run `hugo`. That's it!
You should now have a `public` directory with a complete blog! Open `public/index.html` in your browser and bask.

If that wasn't amazing enough, from the same terminal, run `hugo server -w`. This will watch your directories for changes
and rebuild the site immediately, *and* it will make these changes available at http://localhost:1313/ so you can view
your finished site in your browser. Go on, try it. This is one of the best ways to preview your site while working on it.

To further learn Hugo and learn more, read through the Hugo [documentation](http://gohugo.io/overview/introduction)
or browse around the files in this repository. Have fun!
