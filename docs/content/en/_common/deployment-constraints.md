---
_comment: Do not remove front matter.
---

Deploying a Hugo site to a web server requires more than copying files. For example, typical FTP clients do not remove orphaned files from the server, and compare files by modification time and size.

- Deployments must remove orphaned files. When you rename a page, change its permalink, or delete content, Hugo does not remove the corresponding files from your local `public` directory. Hugo also leaves a page's files in place when the page becomes a draft, expires, or has its publication date changed to a future date. To remove these files, clear the `public` directory before you build, or build with the `--cleanDestinationDir` flag or the [`cleanDestinationDir`][] configuration option enabled. Unless your deployment also removes them from the server, the old pages remain online.

- Deployments should compare files by checksum. Hugo updates the modification time of every file in the `public` directory on every build, so tools that compare modification times transfer your entire site each time. Tools that compare only file size can miss changed files, because editing text does not always change a file's size.

[`cleanDestinationDir`]: /configuration/build/#clean-destination-directory
