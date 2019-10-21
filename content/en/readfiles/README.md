# readdirs Directory for Reusable Content

Files in this directory are:

1. Used in *more than one place* within the Hugo docs
2. Used in Examples of readdir (i.e. in local file templates)

These files are called using the [`readfile` shortcode (source)](../layouts/readfile.html).

You can call this shortcode in the docs as follows:


<code>{</code><code>{</code>% readfile file="/path/to/file.txt" markdown="true" %<code>}</code><code>}</code>


`markdown="true"` is optional (default = `"false"`) and parses the string through the Blackfriday renderer.
