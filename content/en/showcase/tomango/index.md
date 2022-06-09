---

title: Tomango

date: 2018-05-04

description: "Showcase: \"Tomango site relaunch: Building our JAMstack site\""

siteURL: https://www.tomango.co.uk

siteSource: https://github.com/trys/tomango-2018

byline: "[Trys Mudford](https://www.trysmudford.com), Lead Developer, Tomango"

---

Hugo is our static site generator (SSG) of choice. It's **really quick**. After using it on a number of [client projects](/showcase/hartwell-insurance/), it became clear that our new site _had_ to be built with Hugo.

The big benefit of an SSG is how it moves all the heavy lifting to the build time.

For example in WordPress, all the category pages are created at runtime, generating a lot of database queries. In Hugo, the paginated category pages are created at build time - so all the computational complexity is done once, and doesn't impact the user at all.

Similarly, instead of running a live, or even a heavily cached Instagram feed that checked for new photos on page load, we used IFTTT to flip the feature to work performantly. I've [written about it](https://www.trysmudford.com/blog/making-the-static-dynamic-instagram-importer/) in detail on my blog but in essence: IFTTT sends a webhook to a Netlify Cloud Function every time a photo is uploaded. The function scrapes the photo and commits it to our GitHub repo which triggers a Hugo build on Netlify, deploying the site immediately!

Shortcodes allow copy editors to continue using WordPress-esque features, Markdown keeps our developers happy, and our users don't have any of the database overheads. It's win-win!

---

This is an extract from our [technical launch post](https://www.tomango.co.uk/thinks/tomango-progressive-web-app/).
