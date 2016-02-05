---
date: 2015-09-21T16:14:03+02:00
menu:
  main:
    parent: "extras"
title: Dynamic Stylesheets
weight: 120
---

Dynamic stylesheets allow you to change properties of CSS elements directly
in the configs. Every variable that is defined under `params` will be available
within a stylesheet. Let's look at an example.

## Defining properties

Perhaps you want to change the global font family or main colors of your theme easily.
First of all, we need to define these properties in the config file.

#### config.toml

```toml
[params]
    fontFamily = "Arial, Verdana, sans-serif"
    textColor  = "#6BC825"
```

***

## Make the stylesheets dynamic

After defining the CSS properties we need to access them in the stylesheet(s). Like
in other content templates we can access variables with the `.`. Remember that `params`
is the scope that `.` can access.

#### style.css

```css
body {
    font-family: {{ .fontFamily }};
    color: {{ .textColor }};
}
```

***

## The rendering

Now we need to run the command `hugo` or `hugo server` as usual and the output should
look like this if everything gone right:

```css
body {
    font-family: Arial, Verdana, sans-serif;
    color: #6BC825;
}
```

How does it work? After all files are copied into the `public` directory, Hugo searches for stylesheets and treats them as a template that will be rendered. The original stylesheets remain unchanged since we only work with copies of them. It also doesn't matter if the stylesheets are predefined by the theme creator or if you create a custom CSS file under `static` in the root of your Hugo site.
