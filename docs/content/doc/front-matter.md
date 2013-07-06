---
title: "Front Matter"
Pubdate: "2013-07-01"
...

The front matter is one of the features that gives Hugo it's strength. It enables
you to include the meta data of the content right with it. Hugo supports a few 
different formats. The main format supported is JSON. Here is an example:

    {
        "Title": "spf13-vim 3.0 release and new website",
        "Description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
        "Tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ],
        "Pubdate": "2012-04-06",
        "Categories": [ "Development", "VIM" ],
        "Slug": "spf13-vim-3-0-release-and-new-website"
    }

### Variables
There are a few predefined variables that Hugo is aware of and utilizes. The user can also create
any variable they want to. These will be placed into the `.Params` variable available to the templates.

#### Required

**Title**  The title for the content. <br>
**Description** The description for the content.<br>
**Pubdate** The date the content will be sorted by.<br>
**Indexes** These will use the field name of the plural form of the index (see tags and categories above)

#### Optional

**Draft** If true the content will not be rendered unless `hugo` is called with -d<br>
**Type** The type of the content (will be derived from the directory automatically if unset).<br>
**Slug** The token to appear in the tail of the url.<br>
  *or*<br>
**Url** The full path to the content from the web root.<br>
*If neither is present the filename will be used.*

