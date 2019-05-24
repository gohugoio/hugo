---

title: Small Multiples 
date: 2018-02-28
description: "\"Hugo has excellent support and integration with Netlify and we were immediately blown away by how fast it was.\""
siteURL: https://smallmultiples.com.au/
byline: "[Small Multiples](https://smallmultiples.com.au/)"
draft: true
---

Previously we had built and hosted our website with SquareSpace. Although SquareSpace was adequate for quickly showcasing our work, we felt it didn’t reflect our technical capabilities and the types of products we build for our clients.

For many client applications, static front-end sites provide fast, scalable solutions that can easily connect with any back-end service, API or static data. We wanted to use the same processes and infrastructure that we use to host and deploy many of our products, so we felt that building a static site was the right solution for our website.

Netlify is a hosting and deployment service that we use for many products. Our developers really like it because it has strong integration with GitHub and it works with the build tools we use such as Yarn and Webpack. It creates a production build every time we commit to our GitHub repository. This means we can share and preview every change internally or with clients.

Application development has become increasingly complex and there is a strong motivation to simplify this process by avoiding complicated backends in favour of applications that consume static data and APIs (a JAMstack).

Libraries like React make this easy, but we also wanted something that was server rendered. This led us to look at React based tools for static site generation such as GatsbyJS. We liked GatsbyJS, but in the end, we didn’t choose it due to the lack of availability of a simple CMS driven data source.

For this, we considered Contentful. Contentful is a beautifully designed application. It’s basically a headless CMS, but it’s not specifically designed for websites and it becomes quite expensive at a commercial level. Their free tier is possibly a good option for personal sites especially with Gatsby. We also evaluated prose.io. This is a free service for editing markdown files in a GitHub repository. It works well, but it’s quite basic and didn’t provide the editing experience we were looking for.

At the same time, we started exploring Hugo. Hugo is a static site generator similar to Jekyll, but it’s written in Go. It has excellent support and integration with Netlify and we were immediately blown away by how fast it was.

We had been closely following the redevelopment of the Smashing Magazine website. We knew this was being powered by Hugo and Netlify and this showed us that Hugo could work for a large scale sites.

The deciding factor, however, was the availability of CMS options that integrate well with Hugo. Netlify has an open source project called NetlifyCMS and there are also hosted services like Forestry.io. These both provide a CMS with an editing interface for markdown files and images. There is no database, instead, changes are committed directly back into the GitHub repository.

In the end, we chose Hugo on Netlify, with Forestry as our CMS. The site is built and redeployed immediately with Netlify watching for changes to the GitHub repository.

Was this the right choice? For us, yes, but we learnt a few things along the way.

The Hugo templating language was very powerful, although also frustrating at times. The queries used to filter list pages are concise but difficult to read. Although it’s easy to get started, Hugo can have a significant learning curve as you try to do more complicated things.

Hugo has particular expectations when it comes to CMS concepts like tags, categories, RSS, related content and menus. Some parts of our initial design did not match perfectly with how these work in Hugo. It took some time to figure out how to make things work the way we wanted without losing all the benefits of structured content.

There were a few teething issues. We picked some relatively new technologies and as a result, we encountered some bugs. We were forced to find some workarounds and logged some issues with Hugo during the course of development. Most of these were fixed and features were added with releases happening frequently over the time we were working on the project. This can be exciting but also frustrating. We can see Hugo is developing in the right direction.

NetlifyCMS was also very new when we first looked at it and this is partly why we opted for Forestry. Forestry is an excellent choice for an out-of-the-box CMS and it needs very little code configuration. It provided a better editing experience for non-technical users. I would still say this is true, but it also provides fewer options for customisation when compared with NetlifyCMS.

Fortunately, the site is more portable now than it was, or would have been with a dynamic CMS like WordPress, or a fully hosted service like SquareSpace. It should be comparatively easy to swap the publishing functions from Forestry to NetlifyCMS or to change the templates. No part of the pipe-line is tightly coupled, the hosting, the CMS and the templates and the build process can all be updated independently, without changing anything else.

We have complete control over the design and mark-up produced. This means we can implement a better responsive design and have a stronger focus on accessibility and performance.

These technology choices gave us a good performance baseline. It was important to implement a site that took advantage of this. As a data visualisation agency, it can be difficult to optimise for performance with a small bundle size, while also aiming for high-quality visuals and working with large datasets. This meant we spent a lot of time optimising assets making sure there was little blocking the critical path for faster rendering and lazy-load images and videos.

The end result is a high performance site. We think this could have been achieved with GatsbyJS, Hugo or any another static site generator. However, what was important was the decision to use static infrastructure for speed, security, flexibility and hopefully a better ongoing development experience. If you are looking at choosing a static site generator or wondering whether a static is the right choice for you, we hope this has helped.
