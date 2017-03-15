---
title: Hosting on Firebase
linktitle: Hosting on Firebase
description: You can use Firebase's free tier to host your static website; this also gives you access to Firebase's NOSQL API.
date: 2017-03-12
publishdate: 2017-03-12
lastmod: 2017-03-15
categories: [hosting and deployment]
tags: [hosting,firebase]
authors: [Michel Racic]
weight: 20
draft: true
toc: true
aliases: []
wip: true
---

## Assumptions

- Have an account with Firebase
- Have completed the Quick Start or have a completed website ready for deployment

## Initial setup
1. Go to the [Firebase console](https://console.firebase.google.com) and create a new project (unless you already have a project and this is just an additional component to it). 
2. Install `firebase-tools` (node.js)

```sh
npm install -g firebase-tools
```

3. Login to firebase (setup on your local machine) using `firebase login` which opens a browser and you can select your account. Use `firebase logout` in case you are already logged in but to the wrong account.

```sh
firebase login
```

3. In the root of your hugo site initialize the Firebase project with `firebase init`

```sh
firebase init
```

4. Choose Hosting in the feature question
5. Choose the project you did just setup
6. Accept the default for database rules file
7. Accept the default for the publish directory which is `public`
8. Choose No in the question if it is a single-page app

## Deploy
Simply execute `firebase deploy` and your site will be up in no time.

```sh
hugo && firebase deploy
```

## CI Setup
1. Generate a deploy token using 

```sh
firebase login:ci
```

2. Setup your CI (e.g. [Wercker](/hosting-and-deployment/deployment-with-wercker))
3. Add the token to a private variable like `$FIREBASE_DEPLOY_TOKEN`
{{% note %}}
This is a private secret and it should not appear in a public repository. Make sure you understand you chosen CI and that it's not visible to others.
{{% /note %}}
4. Add a step in your build to do the deployment using the token

```sh
firebase deploy --token $FIREBASE_DEPLOY_TOKEN
```

## Reference links
* [Firebase CLI Reference](https://firebase.google.com/docs/cli/#administrative_commands)