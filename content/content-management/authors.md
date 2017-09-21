---
title: Authors
linktitle: Authors
description:
date: 2016-08-22
publishdate: 2017-03-12
lastmod: 2017-03-12
keywords: [authors]
categories: ["content management"]
menu:
  docs:
    parent: "content-management"
    weight: 55
weight: 55	#rem
draft: true
aliases: [/content/archetypes/]
toc: true
comments: Before this page is published, need to also update both site- and page-level variables documentation.
---



Larger sites often have multiple content authors. Hugo provides standardized author profiles to organize relationships between content and content creators for sites operating under a distributed authorship model.

## Author Profiles

You can create a profile containing metadata for each author on your website. These profiles have to be saved under `data/_authors/`. The filename of the profile will later be used as an identifier. This way Hugo can associate content with one or multiple authors. An author's profile can be defined in the JSON, YAML, or TOML format.

### Example: Author Profile

Let's suppose Alice Allison is a blogger. A simple unique identifier would be `alice`. Now, we have to create a file called `alice.toml` in the `data/_authors/` directory. The following example is the standardized template written in TOML:

{{< code  file="data/_authors/alice.toml" >}}
givenName      = "Alice"   # or firstName as alias
familyName     = "Allison" # or lastName as alias
displayName    = "Alice Allison"
thumbnail      = "static/authors/alice-thumb.jpg"
image          = "static/authors/alice-full.jpg"
shortBio       = "My name is Alice and I'm a blogger."
bio            = "My name is Alice and I'm a blogger... some other stuff"
email          = "alice.allison@email.com"
weight         = 10

[social]
    facebook   = "alice.allison"
    twitter    = "alice"
    googleplus = "aliceallison1"
    website    = "www.example.com"

[params]
    random     = "whatever you want"
{{< /code >}}

All variables are optional but it's advised to fill all important ones (e.g. names and biography) because themes can vary in their usage.

You can store files for the `thumbnail` and `image` attributes in the `static` folder. Then add the path to the photos relative to `static`; e.g., `/static/path/to/thumbnail.jpg`.

`weight` allows you to define the order of an author in an `.Authors` list and can be accessed on list or via the `.Site.Authors` variable.

The `social` section contains all the links to the social network accounts of an author. Hugo is able to generate the account links for the most popular social networks automatically. This way, you only have to enter your username. You can find a list of all supported social networks [here](#linking-social-network-accounts-automatically). All other variables, like `website` in the example above remain untouched.

The `params` section can contain arbitrary data much like the same-named section in the config file. What it contains is up to you.

## Associate Content Through Identifiers

Earlier it was mentioned that content can be associated with an author through their corresponding identifier. In our case, blogger Alice has the identifier `alice`. In the front matter of a content file, you can create a list of identifiers and assign it to the `authors` variable. Here are examples for `alice` using YAML and TOML, respectively.

```
---
title: Why Hugo is so Awesome
date: 2016-08-22T14:27:502:00
authors: ["alice"]
---

Nothing to read here. Move along...
```

```
+++
title = Why Hugo is so Awesome
date = "2016-08-22T14:27:502:00"
authors: ["alice"]
+++

Nothing to read here. Move along...
```

Future authors who might work on this blog post can append their identifiers to the `authors` array in the front matter as well.

## Work with Templates

After a successful setup it's time to give some credit to the authors by showing them on the website. Within the templates Hugo provides a list of the author's profiles if they are listed in the `authors` variable within the front matter.

The list is accessible via the `.Authors` template variable. Printing all authors of a the blog post is straight forward:

```
{{ range .Authors }}
    {{ .DisplayName }}
{{ end }}
=> Alice Allison
```

Even if there are co-authors you may only want to show the main author. For this case you can use the `.Author` template variable **(note the singular form)**. The template variable contains the profile of the author that is first listed with his identifier in the front matter.

{{% note %}}
You can find a list of all template variables to access the profile information in [Author Variables](/variables/authors/).
{{% /note %}}

### Link Social Network Accounts

As aforementioned, Hugo is able to generate links to profiles of the most popular social networks. The following social networks with their corrersponding identifiers are supported:  `github`, `facebook`, `twitter`, `googleplus`, `pinterest`, `instagram`, `youtube` and `linkedin`.

This is can be done with the `.Social.URL` function. Its only parameter is the name of the social network as they are defined in the profile (e.g. `facebook`, `googleplus`). Custom variables like `website` remain as they are.

Most articles feature a small section with information about the author at the end. Let's create one containing the author's name, a thumbnail, a (summarized) biography and links to all social networks:

{{< code file="layouts/partials/author-info.html" download="author-info.html" >}}
{{ with .Author }}
    <h3>{{ .DisplayName }}</h3>
    <img src="{{ .Thumbnail | absURL }}" alt="{{ .DisplayName }}">
    <p>{{ .ShortBio }}</p>
    <ul>
    {{ range $network, $username := .Social }}
        <li><a href="{{ $.Author.Social.URL $network }}">{{ $network }}</a></li>
    {{ end }}
    </ul>
{{ end }}
{{< /code >}}

## Who Published What?

That question can be answered with a list of all authors and another list containing all articles that they each have written. Now we have to translate this idea into templates. The [taxonomy][] feature allows us to logically group content based on information that they have in common; e.g. a tag or a category. Well, many articles share the same author, so this should sound familiar, right?

In order to let Hugo know that we want to group content based on their author, we have to create a new taxonomy called `author` (the name corresponds to the variable in the front matter). Here is the snippet in a `config.yaml` and `config.toml`, respectively:

```
taxonomies:
    author: authors
```

```
[taxonomies]
    author = "authors"
```


### List All Authors

In the next step we can create a template to list all authors of your website. Later, the list can be accessed at `www.example.com/authors/`. Create a new template in the `layouts/taxonomy/` directory called `authors.term.html`. This template will be exclusively used for this taxonomy.

{{< code file="layouts/taxonomy/author.term.html" download="author.term.html" >}}
<ul>
{{ range $author, $v := .Data.Terms }}
    {{ $profile := $.Authors.Get $author }}
    <li>
        <a href="{{ printf "%s/%s/" $.Data.Plural $author | absURL }}">
            {{ $profile.DisplayName }} - {{ $profile.ShortBio }}
        </a>
    </li>
{{ end }}
</ul>
{{< /code >}}

`.Data.Terms` contains the identifiers of all authors and we can range over it to create a list with all author names. The `$profile` variable gives us access to the profile of the current author. This allows you to generate a nice info box with a thumbnail, a biography and social media links, like at the [end of a blog post](#linking-social-network-accounts-automatically).

### List Each Author's Publications

Last but not least, we have to create the second list that contains all publications of an author. Each list will be shown in its own page and can be accessed at `www.example.com/authors/<IDENTIFIER>`. Replace `<IDENTIFIER>` with a valid author identifier like `alice`.

The layout for this page can be defined in the template `layouts/taxonomy/author.html`.

{{< code file="layouts/taxonomy/author.html" download="author.html" >}}
{{ range .Data.Pages }}
    <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    <span>written by {{ .Author.DisplayName }}</span>
    {{ .Summary }}
{{ end }}
{{< /code >}}

The example above generates a simple list of all posts written by a single author. Inside the loop you've access to the complete set of [page variables][pagevars]. Therefore, you can add additional information about the current posts like the publishing date or the tags.

With a lot of content this list can quickly become very long. Consider to use the [pagination][] feature. It splits the list into smaller chunks and spreads them over multiple pages.

[pagevars]: /variables/page/
[pagination]: /templates/pagination/
