---
title: Hugo Deploy
linktitle: Hugo Deploy
description: You can upload your site to GCS, S3, or Azure using the Hugo CLI.
date: 2019-05-30
publishdate: 2019-05-30
lastmod: 2019-05-30
categories: [hosting and deployment]
keywords: [s3,gcs,azure,hosting,deployment]
authors: [Robert van Gent]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 2
weight: 2
sections_weight: 2
draft: false
aliases: []
toc: true
---

You can use the "hugo deploy" command to upload your site directly to a Google Cloud Storage (GCS) bucket, an AWS S3 bucket, and/or an Azure Storage bucket.

## Assumptions

* You have completed the [Quick Start][] or have a Hugo website you are ready to deploy and share with the world.
* You have an account with the service provider ([Google Cloud](https://cloud.google.com/), [AWS](https://aws.amazon.com), or [Azure](https://azure.microsoft.com)) that you want to deploy to.
* You have authenticated.
  * Google Cloud: [Install the CLI](https://cloud.google.com/sdk) and run [`gcloud auth login`](https://cloud.google.com/sdk/gcloud/reference/auth/login).
  * AWS: [Install the CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) and run [`aws configure`](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).
  * Azure: [Install the CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) and run [`az login`](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
  * NOTE: Each service supports alternatives for authentication, including using environment variables. See [here](https://gocloud.dev/howto/blob/#services) for more details.

## Create a bucket to deploy to

Create a storage bucket to deploy your site to. If you want your site to be
public, be sure to configure the bucket to be publicly readable.

### Google Cloud Storage (GCS)

Follow the [GCS instructions for how to create a bucket](https://cloud.google.com/storage/docs/creating-buckets).

### AWS S3

Follow the [AWS instructions for how to create a bucket](https://docs.aws.amazon.com/AmazonS3/latest/gsg/CreatingABucket.html).

### Azure Storage

Follow the [Azure instructions for how to create a bucket](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-portal).

## Configure the deployment

In the configuration file for your site, add a `[deployment]` section with one
or more `[[deployment.targets]]` section, one for each deployment target. Here's
a detailed example:

```
[deployment]
# By default, files are uploaded in an arbitrary order.
# Files that match the regular expressions in the "Order" list
# will be uploaded first, in the listed order.
order = [".jpg$", ".gif$"]


[[deployment.targets]]
# An arbitrary name for this target.
name = "mydeployment"
# The Go Cloud Development Kit URL to deploy to. Examples:
# GCS; see https://gocloud.dev/howto/blob/#gcs
# URL = "gs://<Bucket Name>"

# S3; see https://gocloud.dev/howto/blob/#s3
# For S3-compatible endpoints, see https://gocloud.dev/howto/blob/#s3-compatible
# URL = "s3://<Bucket Name>?region=<AWS region>"

# Azure Blob Storage; see https://gocloud.dev/howto/blob/#azure
# URL = "azblob://$web"

# You can use a "prefix=" query parameter to target a subfolder of the bucket:
# URL = "gs://<Bucket Name>?prefix=a/subfolder/"

# If you are using a CloudFront CDN, deploy will invalidate the cache as needed.
cloudFrontDistributionID = <ID>


# ... add more [[deployment.targets]] sections ...


# [[deployment.matchers]] configure behavior for files that match the Pattern.
# Samples:

[[deployment.matchers]]
#  Cache static assets for 20 years.
pattern = "^.+\\.(js|css|svg|ttf)$"
cacheControl = "max-age=630720000, no-transform, public"
gzip = true

[[deployment.matchers]]
pattern = "^.+\\.(png|jpg)$"
cacheControl = "max-age=630720000, no-transform, public"
gzip = false

[[deployment.matchers]]
pattern = "^.+\\.(html|xml|json)$"
gzip = true
```

## Deploy

To deploy to a target:
```
hugo deploy [--target=<target name>, defaults to first target]
```

Hugo will identify and apply any local changes that need to be reflected to the
remote target. You can use `--dryRun` to see the changes without applying them,
or `--confirm` to be prompted before making changes.

See `hugo help deploy` for more command-line options.

[Quick Start]: /getting-started/quick-start/
[Google Cloud]: [https://cloud.google.com]
[AWS]: [https://aws.amazon.com]
[Azure]: [https://azure.microsoft.com]

