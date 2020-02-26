---

title: Aether
date: 2020-02-26
description: "Showcase: \"Hugo is not just a static site generator for us, it's the core framework at the heart of our entire web front-end.\""

# The URL to the site on the internet.
siteURL: https://getaether.net

# Add credit to the article author. Leave blank or remove if not needed/wanted.
byline: "[Burak Nehbit](https://twitter.com/nehbit), Maintainer, Aether"

---

To say that this website, our main online presence, needed to do a lot would be an understatement.

Our site is home to both *Aether* and *Aether Pro*, our **knowledgebase for each product**, a **server for static assets that we use in our emails**, the **interactive sign-up flows**, **payments client**, **downloads provider**, and even a **mechanism for delivering auto-update notifications** to our native clients. We are using a single Hugo site for all these — it's not a static site generator for us, it's the core framework at the heart of our *entire* web front-end.

Not only that, this had to work with one developer crunched for time who spends most of his time working on two separate apps across 3 desktop platforms — someone whose main job is very far from building static websites. We only had scraps of time to design and build this Hugo site, make it performant and scalable, and Hugo did a phenomenal job delivering on that promise.

The last piece is, funnily enough, moving our blog to Hugo, which it is not as of now. This was an inherited mistake we are currently rectifying. Soon, our entire web footprint will be living in Hugo.

### Structure

Our website is built in such a way that there is a separate Vue.js instance for each of the contexts since we are no using JS-based single-page navigation. We use Hugo for navigation and to build most pages. For the pages we need to make interactive, we use Vue.js to build individual, self-contained single-page Javascript apps. One such example is our sign-up flow at [aether.app](https://aether.app), an individual Vue app living within a Hugo page, with its own JS-based navigation.

This is a relatively complex setup, and somewhat out of the ordinary. Yet, even with this custom setup, using Hugo was painless.

### Tools

**CMS**: Hugo

**Theme**: Custom-designed

**Hosting**: Netlify, pushed to production via `git push`.

**Javascript runtime**: Vue.js


