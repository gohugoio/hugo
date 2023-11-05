---
title: Host on Firebase
description: You can use Firebase's free tier to host your static website; this also gives you access to Firebase's NoSQL API.
categories: [hosting and deployment]
keywords: [hosting,firebase]
menu:
  docs:
    parent: hosting-and-deployment
toc: true
---

## Assumptions

1. You have an account with [Firebase][signup]. (If you don't, you can sign up for free using your Google account.)
2. You have completed the [Quick Start] or have a completed Hugo website ready for deployment.

## Initial setup

Go to the [Firebase console][console] and create a new project (unless you already have a project). You will need to globally install `firebase-tools` (node.js):

```sh
npm install -g firebase-tools
```

Log in to Firebase (setup on your local machine) using `firebase login`, which opens a browser where you can select your account. Use `firebase logout` in case you are already logged in but to the wrong account.

```sh
firebase login
```

In the root of your Hugo project, initialize the Firebase project with the `firebase init` command:

```sh
firebase init
```

From here:

1. Choose Hosting in the feature question
2. Choose the project you just set up
3. Accept the default for your database rules file
4. Accept the default for the publish directory, which is `public`
5. Choose "No" in the question if you are deploying a single-page app

## Using Firebase & GitHub CI/CD

In new versions of Firebase, some other questions apply:

6. Set up automatic builds and deploys with GitHub?

Here you will be redirected to login in your GitHub account to get permissions. Confirm.

7. For which GitHub repository would you like to set up a GitHub workflow? (format: user/repository)

Include the repository you will use in the format above (Account/Repo)
Firebase script with retrieve credentials, create a service account you can later manage in your GitHub settings.

8. Set up the workflow to run a build script before every deploy?

Here is your opportunity to include some commands before you run the deploy.

9. Set up automatic deployment to your site's live channel when a PR is merged?

You can let in the default option (main)

After that Firebase has been set in your project with CI/CD. After that run:

```sh
hugo && firebase deploy
```

With this you will have the app initialized manually. After that you can manage and fix your GitHub workflow from: https://github.com/your-account/your-repo/actions

Don't forget to update your static pages before push!

## Manual deploy

To deploy your Hugo site, execute the `firebase deploy` command, and your site will be up in no time:

```sh
hugo && firebase deploy
```

## CI setup (other tools)

You can generate a deploy token using

```sh
firebase login:ci
```

You can also set up your CI and add the token to a private variable like `$FIREBASE_DEPLOY_TOKEN`.

{{% note %}}
This is a private secret and it should not appear in a public repository. Make sure you understand your chosen CI and that it's not visible to others.
{{% /note %}}

You can then add a step in your build to do the deployment using the token:

```sh
firebase deploy --token $FIREBASE_DEPLOY_TOKEN
```

## Reference links

* [Firebase CLI Reference](https://firebase.google.com/docs/cli/#administrative_commands)

[console]: https://console.firebase.google.com/
[Quick Start]: /getting-started/quick-start/
[signup]: https://console.firebase.google.com/
