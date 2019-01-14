---
title: Deployment with Rsync
linktitle: Deployment with Rsync
description: If you have access to your web host with SSH, you can use a simple rsync one-liner to incrementally deploy your entire Hugo website.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [hosting and deployment]
keywords: [rsync,deployment]
authors: [Adrien Poupin]
menu:
  docs:
    parent: "hosting-and-deployment"
    weight: 70
weight: 70
sections_weight: 70
draft: false
aliases: [/tutorials/deployment-with-rsync/]
toc: true
notesforauthors:
---

## Assumptions

* Access to your web host with SSH
* A functional static website built with Hugo

The spoiler is that you can deploy your entire website with a command that looks like the following:

```
hugo && rsync -avz --delete public/ www-data@ftp.topologix.fr:~/www/
```

As you will see, we put it in a shell script file, which makes building and deployment as easy as executing `./deploy`.

## Install SSH Key

If it is not done yet, we will make an automated way to SSH to your server. If you have already installed an SSH key, switch to the next section.

First, install the ssh client. On Debian/Ubuntu/derivates, use the following command:

{{< code file="install-openssh.sh" >}}
sudo apt-get install openssh-client
{{< /code >}}

Then generate your ssh key by entering the following commands:

```
~$ cd && mkdir .ssh & cd .ssh
~/.ssh/$ ssh-keygen -t rsa -q -C "For SSH" -f rsa_id
~/.ssh/$ cat >> config <<EOF
Host HOST
     Hostname HOST
     Port 22
     User USER
     IdentityFile ~/.ssh/rsa_id
EOF
```

Don't forget to replace the `HOST` and `USER` values with your own ones. Then copy your ssh public key to the remote server:

```
~/.ssh/$ ssh-copy-id -i rsa_id.pub USER@HOST.com
```

Now you can easily connect to the remote server:

```
~$ ssh user@host
Enter passphrase for key '/home/mylogin/.ssh/rsa_id':
```

And you've done it!

## Shell Script

We will put the first command in a script at the root of your Hugo tree:

```
~/websites/topologix.fr$ editor deploy
```

Here you put the following content. Replace the `USER`, `HOST`, and `DIR` values with your own:

```
#!/bin/sh
USER=my-user
HOST=my-server.com
DIR=my/directory/to/topologix.fr/   # might sometimes be empty!

hugo && rsync -avz --delete public/ ${USER}@${HOST}:~/${DIR}

exit 0
```

Note that `DIR` is the relative path from the remote user's home. If you have to specify a full path (for instance `/var/www/mysite/`) you must change `~/${DIR}` to `${DIR}` inside the command line. For most cases you should not have to.

Save and close, and make the `deploy` file executable:

```
~/websites/topologix.fr$ chmod +x deploy
```

Now you only have to enter the following command to deploy and update your website:

```
~/websites/topologix.fr$ ./deploy
Started building sites ...
Built site for language en:
0 draft content
0 future content
0 expired content
5 pages created
0 non-page files copied
0 paginator pages created
0 tags created
0 categories created
total in 56 ms
sending incremental file list
404.html
index.html
index.xml
sitemap.xml
cours-versailles/index.html
exercices/index.html
exercices/index.xml
exercices/barycentre-et-carres-des-distances/index.html
posts/
posts/index.html
sujets/index.html
sujets/index.xml
sujets/2016-09_supelec-jp/index.html
tarifs-contact/index.html

sent 9,550 bytes  received 1,708 bytes  7,505.33 bytes/sec
total size is 966,557  speedup is 85.86
```
