---
title: Host on Azure Static Web Apps
description: Deploy Hugo to Azure Static Web Apps and automate the whole process with Github Action Workflow
categories: [hosting and deployment]
keywords: [azure,git,deployment,hosting]
menu:
  docs:
    parent: hosting-and-deployment
    weight: 10
weight: 10
toc: true
---

[Azure Static Web Apps] is a service that automatically builds and deploys full stack web apps to Azure from a Git repository, using [GitHub Actions] or [Azure DevOps].

_The following documentation covers how to use GitHub Actions for the deployment. If you are using Azure DevOps, follow the Microsoft documentation._

## Assumptions

1. You have Git 2.8 or greater [installed on your machine][installgit].
2. You have a GitHub account. [Signing up][ghsignup] for GitHub is free.
3. You have an Azure account. You can sign up for a [Free Trial][azuretrial].
4. You have a ready-to-publish Hugo website or have at least completed the [Quick Start].

## Deploy Hugo to Azure Static Web Apps

1. Navigate to the [Azure Portal][azureportal]
2. Click **Create a Resource**
3. Search for **Static Web Apps**
4. Click **Static Web Apps**
5. Click **Create**

![Create in Azure Portal](/images/hosting-and-deployment/hosting-on-azure/create-in-portal.png)

6. For **Subscription**, accept the subscription that is listed or select a new one from the drop-down list.
7. In _Resource group_, select **New**. In _New resource group name_, enter **hugo-static-app** and select **OK**.
8. Next, a name for your app in the **Name** box. Valid characters include `a-z`, `A-Z`, `0-9` and `-`.
9. For _Region_, select an available region close to you.
10. For _SKU_, select **Free**.

![Basic app details](/images/hosting-and-deployment/hosting-on-azure/basic-app-details.png)

11. Click the **Sign in with GitHub** button.
12. Select the **Organization** under which your repo exists.
13. Select the Hugo app you wish to deploy as the _Repository_ .
14. For the _Branch_ select the branch you want to deploy (eg: **main**).
15. Select **Hugo** under the _Build Presets_, which will populate the configuration files with the standard Hugo build options
  * **App Location** is the path in the Git repo where Hugo's config file is
  * **Api Location** is the path where the Serverless API is (or left blank if there is no API)
  * **Artifact Location** is the path where Hugo publishes to
16. Click **Review + Create** to review the details and then **Create** to start the creation of the Azure Static Web Apps and create the GitHub Action workflow for deployment.

A GitHub Action workflow will immediately start a build using Hugo and deployment to Azure. The website can be accessed via the URL shown on the _Overview_ page of the Azure Static Web Apps resource in Azure.

## Using A Custom Hugo Version

When you create a Static Web App, a [workflow file][swaconfig] is generated which contains the deployment settings for the site. You can configure a specific Hugo version in the workflow file by providing a value for `HUGO_VERSION` in the `env` section of the `Azure/static-web-apps-deploy` GitHub Action.

```yaml
jobs:
  build_and_deploy_job:
    if: github.event_name == 'push' || (github.event_name == 'pull_request' && github.event.action != 'closed')
    runs-on: ubuntu-latest
    name: Build and Deploy Job
    steps:
      - uses: actions/checkout@v3
        with:
          submodules: true
      - name: Build And Deploy
        id: builddeploy
        uses: Azure/static-web-apps-deploy@v1
        with:
          azure_static_web_apps_api_token: ${{ secrets.AZURE_STATIC_WEB_APPS_API_TOKEN }}
          repo_token: ${{ secrets.GITHUB_TOKEN }}
          action: "upload"
          app_location: "/" # App source code path
          api_location: "api" # Api source code path - optional
          output_location: "public" # Built app content directory - optional
        env:
          HUGO_VERSION: 0.100.2
```

## Use a Custom Domain

Azure Static Web Apps supports custom domains as a CNAME or APEX domain mapping. You can configure the custom domains via the Azure Portal. Refer to the [official documentation for custom domains][domains] for more information.

[Azure Static Web Apps]: https://docs.microsoft.com/azure/static-web-apps/?WT.mc_id=javascript-26008-aapowell
[GitHub Actions]: https://docs.github.com/en/actions
[Azure DevOps]: https://docs.microsoft.com/azure/static-web-apps/publish-devops?WT.mc_id=javascript-26008-aapowell
[ghsignup]: https://github.com/join
[installgit]: https://git-scm.com/downloads
[azuretrial]: https://azure.microsoft.com/free/?WT.mc_id=javascript-26008-aapowell
[azureportal]: https://portal.azure.com/
[swaconfig]: https://docs.microsoft.com/azure/static-web-apps/github-actions-workflow?WT.mc_id=javascript-26008-aapowell
[domains]: https://docs.microsoft.com/azure/static-web-apps/custom-domain?WT.mc_id=javascript-26008-aapowell
[Quick Start]: /getting-started/quick-start/
