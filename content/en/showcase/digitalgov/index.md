---
title: Digital.gov
date: 2020-05-01
description: "Showcase: \"Guidance on building better digital services in government.\""
siteURL: https://digital.gov/
siteSource: https://github.com/gsa/digitalgov.gov
---

For over a decade, Digital.gov has provided guidance, training, and community support to the people who are responsible for delivering digital services in the U.S. government. Essentially, it is a place where people can find examples of problems being solved in government, and get links to the tools and resources they need.

Through collaboration in our communities of practice, Digital.gov is a window into the people who work in technology in government and the challenges they face making digital services stronger and more effective. [Read more about our site »](https://digital.gov/2019/12/19/a-new-digitalgov/)

Digital.gov is built using the [U.S. Web Design System](https://designsystem.digital.gov/) (USWDS) and have followed the [design principles](https://designsystem.digital.gov/maturity-model/) in building out our new site:

- **Start with real user needs**  — We used human-centered design methods to inform our product decisions (like qualitative user research), and gathered feedback from real users. We also continually test our assumptions with small experiments.
- **Earn trust**  —We recognize that trust has to be earned every time. We are including all  [required links and content](https://digital.gov/resources/required-web-content-and-links/)  on our site, clearly identifying as a government site, building with modern best practices, and using HTTPS.
- **Embrace accessibility**  —  [Accessibility](https://digital.gov/resources/intro-accessibility/)  affects everybody, and we built it into every decision. We’re continually working to conform to Section 508 requirements, use user experience best practices, and support a wide range of devices.
- **Promote continuity**  — We started from shared solutions like USWDS and  [Federalist](https://federalist.18f.gov/). We designed our site to clearly identify as a government site by including USWDS’s .gov banner, common colors and patterns, and built with modern best practices.
- **Listen**  — We actively collect user feedback and web metrics. We use the  [Digital Analytics Program](https://digital.gov/services/dap/)  (DAP) and analyze the data to discover actionable insights. We make small, incremental changes to continuously improve our website by listening to readers and learning from what we hear.

_More on the [USWDS maturity model »](https://designsystem.digital.gov/maturity-model/)_

## Open tools

We didn’t start from scratch. We built and designed the Digital.gov using many of the open-source tools and services that we develop for government here in the  [Technology Transformation Services](https://www.gsa.gov/tts/) (TTS).

Using services that make it possible to design, build, and iterate quickly are essential to modern web design and development, which is why [Federalist](https://federalist.18f.gov/) and the [U.S. Web Design System](https://designsystem.digital.gov/) are such a great combination.

**Why Hugo?** Well, with around `~3,000` files _(and growing)_ and `~9,000` built pages, we needed a site generator that could handle that volume with lightning fast speed.

Hugo was the clear option. The [Federalist](https://federalist.18f.gov/) team quickly added it to their available site generators, and we were off.

At the moment, it takes around `32 seconds` to build close to `~10,000` pages!

Take a look:

```text

                   |  EN
-------------------+-------
  Pages            | 7973
  Paginator pages  |  600
  Non-page files   |  108
  Static files     |  851
  Processed images |    0
  Aliases          | 1381
  Sitemaps         |    1
  Cleaned          |    0

Built in 32.427 seconds
```

In addition to Hugo, we are proudly using a number of other tools and services, all built by government are free to use:

- [Federalist](https://federalist.18f.gov/)
- [Search.gov](https://www.search.gov/)  — A free, hosted search platform for federal websites.
- [Cloud.gov](https://www.cloud.gov/)  — helps teams build, run, and authorize cloud-ready or legacy government systems quickly and cheaply.
- [Federal CrowdSource Mobile Testing Program](https://digital.gov/services/service_mobile-testing-program/)  — Free mobile compatibility testing by feds, for feds.
- [Digital Analytics Program](https://digital.gov/services/dap/)  (DAP) — A free analytics tool for measuring digital services in the federal government
- [Section508.gov](https://www.section508.gov/)  and  [PlainLanguage.gov](https://www.plainlanguage.gov/)  resources
- [API.data.gov](https://api.data.gov/)  — a free API management service for federal agencies
- [U.S. Digital Registry](https://digital.gov/services/u-s-digital-registry/)  — A resource for confirming the official status of government social media accounts, mobile apps, and mobile websites.

**Questions or feedback?** [Submit an issue](https://github.com/GSA/digitalgov.gov/issues) or send us an email to [digitalgov@gsa.gov](mailto:digitalgov@gsa.gov) :heart:
