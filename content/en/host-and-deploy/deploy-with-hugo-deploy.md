---
title: Deploy with hugo
description: Deploy your site with the hugo CLI.
categories: []
keywords: []
aliases: [/hosting-and-deployment/hugo-deploy/]
---

Use the `hugo deploy` command to deploy your site Amazon S3, Azure Blob Storage, or Google Cloud Storage.

> [!note]
> This feature requires the Hugo extended/deploy edition. See the [installation] section for details.

## Assumptions

1. You have completed the [Quick Start] or have a Hugo website you are ready to deploy and share with the world.
1. You have an account with the service provider ([AWS], [Azure], or [Google Cloud]) that you want to deploy to.
1. You have authenticated.
    - AWS: [Install the CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) and run [`aws configure`](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html).
    - Azure: [Install the CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) and run [`az login`](https://docs.microsoft.com/en-us/cli/azure/authenticate-azure-cli).
    - Google Cloud: [Install the CLI](https://cloud.google.com/sdk) and run [`gcloud auth login`](https://cloud.google.com/sdk/gcloud/reference/auth/login).

    Each service supports various authentication methods, including environment variables. See&nbsp;[details](https://gocloud.dev/howto/blob/#services).

1. You have created a bucket to deploy to. If you want your site to be
  public, be sure to configure the bucket to be publicly readable as a static website.
    - AWS: [create a bucket](https://docs.aws.amazon.com/AmazonS3/latest/gsg/CreatingABucket.html) and [host a static website](https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html)
    - Azure: [create a storage container](https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-portal) and [host a static website](https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blob-static-website)

    - Google Cloud: [create a bucket](https://cloud.google.com/storage/docs/creating-buckets) and [host a static website](https://cloud.google.com/storage/docs/hosting-static-website)

## Configuration

Create a deployment target in your [site configuration]. The only required parameters are [`name`] and [`url`]:

{{< code-toggle file=hugo >}}
[deployment]
  [[deployment.targets]]
    name = 'production'
    url = 's3://my_bucket?region=us-west-1'
{{< /code-toggle >}}

## Deploy

To deploy to a target:

```bash
hugo deploy [--target=<target name>]
```

This command syncs the contents of your local `public` directory (the default publish directory) with the destination bucket. If no target is specified, Hugo deploys to the first configured target.

For more command-line options, see `hugo help deploy` or the [CLI documentation].

### File list creation

`hugo deploy` creates local and remote file lists by traversing the local publish directory and the remote bucket. Inclusion and exclusion are determined by the deployment target's [configuration]:

- `include`: All files are skipped by default except those that match the pattern.
- `exclude`: Files matching the pattern are skipped.

> [!note]
> During local file list creation, Hugo skips `.DS_Store` files and hidden directories (those starting with a period, like `.git`), except for the [`.well-known`] directory, which is traversed if present.

### File list comparison

Hugo compares the local and remote file lists to identify necessary changes. It first compares file names. If both exist, it compares sizes and MD5 checksums. Any difference triggers a re-upload, and remote files not present locally are deleted.

> [!note]
> Excluded remote files (due to `include`/`exclude` configuration) won't be deleted.

The `--force` flag forces all files to be re-uploaded, even if Hugo detects no local/remote differences.

The `--confirm` or `--dryRun` flags cause Hugo to display the detected differences and then pause or stop.

### Synchronization

Hugo applies the changes to the remote bucket: uploading missing or changed files and deleting remote files not present locally. Uploaded file headers are configured remotely based on the matchers configuration.

> [!note]
> To prevent accidental data loss, Hugo will not delete more than 256 remote files by default. Use the `--maxDeletes` flag to override this limit.

## Advanced configuration

See [configure deployment](/configuration/deployment/).

[`.well-known`]: https://en.wikipedia.org/wiki/Well-known_URI
[`name`]: /configuration/deployment/#name
[`url`]: /configuration/deployment/#url
[AWS]: https://aws.amazon.com
[Azure]: https://azure.microsoft.com
[CLI documentation]: /commands/hugo_deploy/
[configuration]: /configuration/deployment/#targets-1
[Google Cloud]: https://cloud.google.com/
[installation]: /installation/
[Quick Start]: /getting-started/quick-start/
[site configuration]: /configuration/deployment/
