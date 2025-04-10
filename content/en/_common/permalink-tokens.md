---
_comment: Do not remove front matter.
---

`:year`
: The 4-digit year as defined in the front matter `date` field.

`:month`
: The 2-digit month as defined in the front matter `date` field.

`:monthname`
: The name of the month as defined in the front matter `date` field.

`:day`
: The 2-digit day as defined in the front matter `date` field.

`:weekday`
: The 1-digit day of the week as defined in the front matter `date` field  (Sunday = `0`).

`:weekdayname`
: The name of the day of the week as defined in the front matter `date` field.

`:yearday`
: The 1- to 3-digit day of the year as defined in the front matter `date` field.

`:section`
: The content's section.

`:sections`
: The content's sections hierarchy. You can use a selection of the sections using _slice syntax_: `:sections[1:]` includes all but the first, `:sections[:last]` includes all but the last, `:sections[last]` includes only the last, `:sections[1:2]` includes section 2 and 3. Note that this slice access will not throw any out-of-bounds errors, so you don't have to be exact.

`:title`
: The `title` as defined in front matter, else the automatic title. Hugo generates titles automatically for section, taxonomy, and term pages that are not backed by a file.

`:slug`
: The `slug` as defined in front matter, else the `title` as defined in front matter, else the automatic title. Hugo generates titles automatically for section, taxonomy, and term pages that are not backed by a file.

`:filename`
: The content's file name without extension, applicable to the `page` page kind.

  {{< deprecated-in v0.144.0 >}}
  The `:filename` token has been deprecated. Use `:contentbasename` instead.
  {{< /deprecated-in >}}

`:slugorfilename`
: The `slug` as defined in front matter, else the content's file name without extension, applicable to the `page` page kind.

  {{< deprecated-in v0.144.0 >}}
  The `:slugorfilename` token has been deprecated. Use `:slugorcontentbasename` instead.
  {{< /deprecated-in >}}

`:contentbasename`
: {{< new-in 0.144.0 />}}
: The [content base name].

[content base name]: /methods/page/file/#contentbasename

`:slugorcontentbasename`
: {{< new-in 0.144.0 />}}
: The `slug` as defined in front matter, else the [content base name].

For time-related values, you can also use the layout string components defined in Go's [time package]. For example:

[time package]: https://pkg.go.dev/time#pkg-constants

{{< code-toggle file=hugo >}}
permalinks:
  posts: /:06/:1/:2/:title/
{{< /code-toggle >}}

[content base name]: /methods/page/file/#contentbasename
