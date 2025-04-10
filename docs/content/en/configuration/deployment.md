---
title: Configure deployment
linkTitle: Deployment
description: Configure deployments to Amazon S3, Azure Blob Storage, or Google Cloud Storage.
categories: []
keywords: []
---

> [!note]
> This configuration is only relevant when running `hugo deploy`. See&nbsp;[details](/host-and-deploy/deploy-with-hugo-deploy/).

## Top-level options

These settings control the overall behavior of the deployment process. This is the default configuration:

{{< code-toggle file=hugo config=deployment />}}

confirm
: (`bool`) Whether to prompt for confirmation before deploying. Default is `false`.

dryRun
: (`bool`) Whether to simulate the deployment without any remote changes. Default is `false`.

force
: (`bool`) Whether to re-upload all files. Default is `false`.

invalidateCDN
: (`bool`) Whether to invalidate the CDN cache listed in the deployment target. Default is `true`.

maxDeletes
: (`int`) The maximum number of files to delete, or `-1` to disable. Default is `256`.

matchers
: (`[]*Matcher`) A slice of [matchers](#matchers-1).

order
: (`[]string`) An ordered slice of [regular expressions](g) that determines upload priority (left to right). Files not matching any expression are uploaded last in an arbitrary order.

target
: (`string`) The target deployment [`name`](#name). Defaults to the first target.

targets
: (`[]*Target`) A slice of [targets](#targets-1).

workers
: (`int`) The number of concurrent workers to use when uploading files. Default is `10`.

## Targets

A target represents a deployment target such as "staging" or "production".

cloudFrontDistributionID
: (`string`) The CloudFront Distribution ID, applicable if you are using the Amazon Web Services CloudFront CDN. Hugo will invalidate the CDN when deploying this target.

exclude
: (`string`) A [glob](g) pattern matching files to exclude when deploying to this target. Local files failing the include/exclude filters are not uploaded, and remote files failing these filters are not deleted.

googleCloudCDNOrigin
: (`string`) The Google Cloud project and CDN origin to invalidate when deploying this target, specified as `<project>/<origin>`.

include
: (`string`) A [glob](g) pattern matching files to include when deploying to this target. Local files failing the include/exclude filters are not uploaded, and remote files failing these filters are not deleted.

name
: (`string`) An arbitrary name for this target.

stripIndexHTML
: (`bool`) Whether to map files named `<dir>/index.html` to `<dir>` on the remote (except for the root `index.html`). This is useful for key-value cloud storage (e.g., Amazon S3, Google Cloud Storage, Azure Blob Storage) to align canonical URLs with object keys. Default is `false`.

url
: (`string`) The [destination URL](#destination-urls) for deployment.

## Matchers

A Matcher represents a configuration to be applied to files whose paths match
the specified pattern.

cacheControl
: (`string`) The caching attributes to use when serving the blob. See&nbsp;[details][cacheControl].

contentEncoding
: (`string`) The encoding used for the blob's content, if any. See&nbsp;[details][contentEncoding].

contentType
: (`string`) The media type of the blob being written. See&nbsp;[details][contentType].

force
: (`bool`) Whether matching files should be re-uploaded. Useful when other route-determined metadata (e.g., `contentType`) has changed. Default is `false`.

gzip
: (`bool`) Whether the file should be gzipped before upload. If so, the `ContentEncoding` field will automatically be set to `gzip`. Default is `false`.

pattern
: (`string`) A [regular expression](g) used to match paths. Paths are converted to use forward slashes (`/`) before matching.

[cacheControl]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
[contentEncoding]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
[contentType]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type

## Destination URLs

Service|URL example
:--|:--
Amazon Simple Storage Service (S3)|`s3://my-bucket?region=us-west-1`
Azure Blob Storage|`azblob://my-container`
Google Cloud Storage (GCS)|`gs://my-bucket`

With Google Cloud Storage you can target a subdirectory:

```text
gs://my-bucket?prefix=a/subdirectory
```

You can also to deploy to storage servers compatible with Amazon S3 such as:

- [Ceph]
- [MinIO]
- [SeaweedFS]

[Ceph]: https://ceph.com/
[Minio]: https://www.minio.io/
[SeaweedFS]: https://github.com/chrislusf/seaweedfs

For example, the `url` for a MinIO deployment target might resemble this:

```text
s3://my-bucket?endpoint=https://my.minio.instance&awssdk=v2&use_path_style=true&disable_https=false
```

## Example

{{< code-toggle file=hugo >}}
[deployment]
  order = ['.jpg$', '.gif$']
  [[deployment.matchers]]
    cacheControl = 'max-age=31536000, no-transform, public'
    gzip = true
    pattern = '^.+\.(js|css|svg|ttf)$'
  [[deployment.matchers]]
    cacheControl = 'max-age=31536000, no-transform, public'
    gzip = false
    pattern = '^.+\.(png|jpg)$'
  [[deployment.matchers]]
    contentType = 'application/xml'
    gzip = true
    pattern = '^sitemap\.xml$'
  [[deployment.matchers]]
    gzip = true
    pattern = '^.+\.(html|xml|json)$'
  [[deployment.targets]]
    url = 's3://my_production_bucket?region=us-west-1'
    cloudFrontDistributionID = 'E1234567890ABCDEF0'
    exclude = '**.{heic,psd}'
    name = 'production'
  [[deployment.targets]]
    url = 's3://my_staging_bucket?region=us-west-1'
    exclude = '**.{heic,psd}'
    name = 'staging'
{{< /code-toggle >}}
