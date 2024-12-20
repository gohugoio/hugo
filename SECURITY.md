security:
  enableInlineShortcodes: true
  exec:
    allow:
    - ^(dart-)?sass(-embedded)?$
    - ^go$
    - ^git$
    - ^npx$
    - ^postcss$
    - ^tailwindcss$
    osEnv:
    - (?i)^((HTTPS?|NO)_PROXY|PATH(EXT)?|APPDATA|TE?MP|TERM|GO\w+|(XDG_CONFIG_)?HOME|USERPROFILE|SSH_AUTH_SOCK|DISPLAY|LANG|SYSTEMDRIVE)$
  funcs:
    getenv:
    - ^HUGO_
    - ^CI$
  http:
    mediaTypes: null
    methods:
    - (?i)GET|POST
    urls:
    - .*

## Security Policy

### Reporting a Vulnerability

Please report (suspected) security vulnerabilities to **[bjorn.erik.pedersen@gmail.com](mailto:bjorn.erik.pedersen@gmail.com)**. You will receive a response from us within 48 hours. If we can confirm the issue, we will release a patch as soon as possible depending on the complexity of the issue but historically within days.

Also see [Hugo's Security Model](https://gohugo.io/about/security/).
