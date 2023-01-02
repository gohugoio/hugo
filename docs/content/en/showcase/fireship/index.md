---
title: fireship.io
date: 2019-02-02
description: "Showcase: \"Hugo helps us create complex technical content that integrates engaging web components\""
siteURL: https://fireship.io
siteSource: https://github.com/fireship-io/fireship.io
byline: "[Jeff Delaney](https://github.com/codediodeio), Fireship.io Creator"
---

After careful consideration of JavaScript/JSX-based static site generators, it became clear that Hugo was the only tool capable of handling our project's complex demands. Not only do we have multiple content formats and taxonomies, but we often need to customize the experience at a more granular level. The problems Hugo has solved for us include:

- **Render speed.** We know from past experience that JavaScript-based static site generators become very slow when you have thousands of pages and images.
- **Feature-rich.** Our site has a long list of specialized needs and Hugo somehow manages to cover every single use case.
- **Composability.** Hugo's partial and shortcode systems empower us to write DRY and maintainable templates.
- **Simplicity.** Hugo is easy to learn (even without Go experience) and doesn't burden us with brittle dependencies.

The site is able to achieve Lighthouse performance scores of 95+, despite the fact that it is a fully interactive PWA that ships Angular and Firebase in the JS bundle. This is made possible by (1) prerendering content with Hugo and (2) lazily embedding native web components directly in the HTML and Markdown. We provide a [detailed explanation](https://youtu.be/gun8OiGtlNc) of the architecture on YouTube and can't imagine development without Hugo.
