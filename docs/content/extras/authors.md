---
date: 2016-08-22T14:27:50+02:00
lastmod: 2016-08-22
menu:
  main:
    parent: extras
title: Authors
weight: 20
next: /extras/builders
prev: /extras/analytics
---

For larger websites it's not unusual to have multiple publishing content creators. Hugo tries to provide a standardized approach to organize relations between content and their authors. This can be achieved with author profiles.

## Author profiles

For each author you can create a profile that will contain metadata of him or her. Those profiles have to be saved under `data/_authors/`. The filename of the profile will later be used as an identifier. This way Hugo can associate content with one or multiple authors through their identifiers. An author's profile can be defined in the JSON, YAML or TOML format.

### Profile example

Let's suppose Alice Allison is a blogger. A simple unique identifier would be `alice`. Now, we have to create a file called `alice.toml` in the `data/_authors/` directory. The following example is the standardized template written in TOML:

```toml
# file: data/_authors/alice.toml

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
```

All variables are optional but it's advised to fill all important ones (e.g. names and biography) because themes can vary in their usage.

You can store files for the `thumbnail` and `image` attributes in the `static` folder. Then add the path to the photos relative to `static`, e.g. `/static/path/to/thumbnail.jpg`.

Weight allows you to define the order of an author in the `.Authors` list, that can be accessed on nodes or via `.Site.Authors`.

The `social` section contains all the links to the social network accounts of an author. Hugo is able to generate the account links for the most popular social networks automatically. This way, you only have to enter your username. You can find a list of all supported social networks [here](#linking-social-network-accounts-automatically). All other variables, like `website` in the example above remain untouched.

The `params` section can contain arbitrary data much like the same-named section in the config file. What it contains is up to you.

## Associating content through identifiers

Earlier it was mentioned that content can be associated with an author through their corresponding identifer. In our case blogger Alice has the identifer `alice`. In the frontmatter of a content file you can create a list of identifiers and assign it to the `authors` variable.

```yaml
---
title: Why Hugo is so awesome
date: 2016-08-22T14:27:50+02:00
authors: ["alice"]
---

Nothing to read here. Move along...
```

If multiple authors work on this blog post then append their identifiers to the `authors` list in the frontmatter as well.

## Working with templates

After a successful setup it's time to give some credit to the authors by showing them on the website. Within the templates Hugo provides a list of the author's profiles if they are listed in the `authors` variable within the frontmatter.

The list is accessible via the `.Authors` template variable. Printing all authors of a the blog post is straight forward:

```
{{ range .Authors }}
    {{ .DisplayName }}
{{ end }}

# output: Alice Allison
```

Even if there are co-authors you may only want to show the main author. For this case you can use the `.Author` template variable **(note the singular form)**. The template variable contains the profile of the author that is first listed with his identifier in the frontmatter.

> **Note:** you can find a list of all template variables to access the profile information [here]({{< relref "templates/variables.md#author-variables" >}})

### Linking social network accounts automatically

As aforementioned, Hugo is able to generate links to profiles of the most popular social networks. The following social networks with their corrersponding identifiers are supported:  `github`, `facebook`, `twitter`, `googleplus`, `pinterest`, `instagram`, `youtube` and `linkedin`.

This is can be done with the `.Social.URL` function. Its only parameter is the name of the social network as they are defined in the profile (e.g. `facebook`, `googleplus`). Custom variables like `website` remain as they are.

Most articles feature a small section with information about the author at the end. Let's create one containing the author's name, a thumbnail, a (summarized) biography and links to all social networks:

```html
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
```

## Who published what?

That question can be answered with a list of all authors and another list containing all articles that they each have written. Now we have to translate this idea into templates. The [taxonomy]({{< relref "taxonomies/overview.md" >}}) feature allows us to logically group content based on an information that they have in common, e.g. a tag or a category. Well, many articles share the same author, so this should sound familiar, right?

In order to let Hugo know that we want to group content based on their author, we have to create a new taxonomy called `author` (the name corresponds to the variable in the frontmatter). Open your config file and add the following information:

```toml
# file: config.toml

[taxonomies]
    author = "authors"
```

### Listing all authors

In the next step we can create a template to list all authors of your website. Later, the list can be accessed at `www.example.com/authors/`. Create a new template in the `layouts/taxonomy/` directory called `authors.term.html`. This template will be exclusively used for this taxonomy.

```html
<!-- file: layouts/taxonomy/author.term.html -->

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
```

`.Data.Terms` contains the identifiers of all authors and we can range over it to create a list with all author names. The `$profile` variable gives us access to the profile of the current author. This allows you to generate a nice info box with a thumbnail, a biography and social media links, like at the [end of a blog post](#linking-social-network-accounts-automatically).

### Listing each author's publications

Last but not least we have to create the second list that contains all publications of an author. Each list will be shown in its own page and can be accessed at `www.example.com/authors/<IDENTIFIER>`. Replace `<IDENTIFIER>` with a valid author identifier like `alice`.

The layout for this page can be defined in the template `layouts/taxonomy/author.html`.

```html
<!-- file: layouts/taxonomy/author.html -->

{{ range .Data.Pages }}
    <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    <span>written by {{ .Author.DisplayName }}</span>

    {{ .Summary }}
{{ end }}
```

The example above generates a simple list of all posts written by a single author. Inside the loop you've access to the complete set of [page variables]({{< relref "templates/variables.md#page-variables" >}}). Therefore, you can add additional information about the current posts like the publishing date or the tags.

With a lot of content this list can quickly become very long. Consider to use the [pagination]({{< relref "extras/pagination.md" >}}) feature. It splits the list into smaller chunks and spreads them over multiple pages.
