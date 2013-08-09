---
title: "Indexes"
Pubdate: "2013-07-01"
---

Hugo includes support for user defined indexes of content. In our 
terminology an index is best thought of as tags applied to content
but they can be used for far more than just tags. Other common
uses would include categories, groups, series. For the purpose of 
this document we will just use tags for our example. For a more 
complete example see [spf13.com-hugo](http://github.com/spf13/spf13.com-hugo).

## Defining Indexes for a site

Indexes must be defined in the site configuration, before they
can be used throughout the site. 

Here is an example configuration in YAML that specifies two indexes.
Notice the format is **singular key** : *plural value*. While 
we could use an inflection library to pluralize this, they currently
support only a few languages, so instead we've opted for user defined
pluralization.

**config.yaml**

    ---
    indexes:
        tag: "tags"
        topic: "topics"
    baseurl: "http://spf13.com/"
    title: "Steve Francia is spf13.com"
    ---

## Creating index templates
For each index type a template needs to be provided to render the index page.
In the case of tags, this will render the content for /tags/TAGNAME/.

The template must be called the singular name of the index and placed in 
layouts/indexes

    .
    └── layouts
        └── indexes
            └── category.html

The template will be provided Data about the index. 

### Variables

The following variables are available to the index template:

**.Title**  The title for the content. <br>
**.Date** The date the content is published on.<br>
**.Permalink** The Permanent link for this page.<br>
**.RSSLink** Link to the indexes' rss link. <br>
**.Data.Pages** The content that is assigned this index.<br>
**.Data.`singular`** The index itself.<br>

#### Example

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
        {{ range .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>

    <aside id="meta"> </aside>

    {{ template "chrome/footer.html" }}


## Assigning indexes to content

Once an index is defined at the site level, any piece of content
can be assigned to it regardless of content type or section.

Assigning content to an index is done in the front matter.
Simply create a variable with the *plural* name of the index
and assign all keys you want this content to match against. 

**Index values are case insensitive**

#### Example
    {
        "Title": "Hugo: A fast and flexible static site generator",
        "Tags": [
            "Development",
            "golang",
            "Blogging"
        ],
        "Slug": "hugo",
        "project_url": "http://github.com/spf13/hugo"
    }


## Displaying indexes within content

Within your content templates you may wish to display 
the indexes that that piece of content is assigned to.

Because we are leveraging the front matter system to 
define indexes for content, the indexes assigned to 
each content piece are located in the usual place 
(.Params.`plural`)

#### Example

    <ul id="tags">
      {{ range .Params.tags }}
        <li><a href="tags/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>

If you wish to display the list of all indexes, the index can
be retrieved from the `.Site` variable.

#### Example

    <ul id="all-tags">
      {{ range .Site.Indexes.tags }}  
        <li><a href="/tags/{{ .Name | urlize }}">{{ .Name }}</a></li>  
      {{ end }}
    </ul>

## Creating Indexes of Indexes

Hugo also supports creating pages that list your values for each 
index along with the number of content items associated with the 
index key.

This may take the form of a tag cloud or simply a list.

To have hugo create these indexes of indexes pages, simply create
a template in indexes called indexes.html

Hugo provides two different versions of the index. One alphabetically
sorted, the other sorted by most popular. It's important to recognize
that the data structure of the two is different.

#### Example indexes.html file (alphabetical)

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
       <ul>
       {{ $data := .Data }}
        {{ range $key, $value := .Data.Index }}
        <li><a href="{{ $data.Plural }}/{{ $key | urlize }}"> {{ $key }} </a> {{ len $value }} </li>
        {{ end }}
       </ul>
      </div>
    </section>

    {{ template "chrome/footer.html" }}


#### Example indexes.html file (ordered)

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
       <ul>
        {{ range $data.OrderedIndex }}
        <li><a href="{{ $data.Plural }}/{{ .Name | urlize }}"> {{ .Name }} </a> {{ .Count }} </li>
        {{ end }}
       </ul>
      </div>
    </section>

    {{ template "chrome/footer.html" }}

### Variables available to indexes of indexes pages.

**.Title**  The title for the content. <br>
**.Date** The date the content is published on.<br>
**.Permalink** The Permanent link for this page.<br>
**.RSSLink** Link to the indexes' rss link. <br>
**.Data.Singular** The singular name of the index <br>
**.Data.Plural** The plural name of the index<br>
**.Data.Index** The Alphabetical index<br>
**.Data.OrderedIndex** The popular index<br>

## Creating a menu based on indexes

Hugo can generate menus based on indexes by iterating and
nesting the index keys. This can be used to build a hierarchy
of content within your site.

To have hugo create the menu, simply create a template in chome
called menu.html, then include it using the 
`{{ template "chrome/menu.html" . }}` syntax.


#### Example menu.html file 

    <section id="menu">
      <ul>
        {{ range $indexname, $index := .Site.Indexes }}
          <li><a href="/{{ $indexname | urlize }}">{{ $indexname }}</a> 
            <ul> 
              {{ range $index }}
                <li><a href="/{{ $indexname | urlize }}/{{ .Name | urlize }}">{{ .Name }}</a></li>
              {{ end }}
            </ul>
          </li> 
        {{ end }}
      </ul>
    </section>

