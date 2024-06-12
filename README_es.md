[bep]: https://github.com/bep
[bugs]: https://github.com/gohugoio/hugo/issues?q=is%3Aopen+is%3Aissue+label%3ABug
[contributing]: CONTRIBUTING.md
[create a proposal]: https://github.com/gohugoio/hugo/issues/new?labels=Proposal%2C+NeedsTriage&template=feature_request.md
[documentation repository]: https://github.com/gohugoio/hugoDocs
[documentation]: https://gohugo.io/documentation
[dragonfly bsd, freebsd, netbsd, and openbsd]: https://gohugo.io/installation/bsd
[forum]: https://discourse.gohugo.io
[friends]: https://github.com/gohugoio/hugo/graphs/contributors
[go]: https://go.dev/
[hugo modules]: https://gohugo.io/hugo-modules/
[installation]: https://gohugo.io/installation
[issue queue]: https://github.com/gohugoio/hugo/issues
[linux]: https://gohugo.io/installation/linux
[macos]: https://gohugo.io/installation/macos
[prebuilt binary]: https://github.com/gohugoio/hugo/releases/latest
[requesting help]: https://discourse.gohugo.io/t/requesting-help/9132
[spf13]: https://github.com/spf13
[static site generator]: https://en.wikipedia.org/wiki/Static_site_generator
[support]: https://discourse.gohugo.io
[themes]: https://themes.gohugo.io/
[twitter]: https://twitter.com/gohugoio
[website]: https://gohugo.io
[windows]: https://gohugo.io/installation/windows

<a href="https://gohugo.io/"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/static/images/hugo-logo-wide.svg?sanitize=true" alt="Hugo" width="565"></a>

Un generador de sitios estáticos rápido y flexible construido con amor por [bep], [spf13], y [friends] en [Go].

---

[![GoDoc](https://godoc.org/github.com/gohugoio/hugo?status.svg)](https://godoc.org/github.com/gohugoio/hugo)
[![Tests on Linux, MacOS and Windows](https://github.com/gohugoio/hugo/workflows/Test/badge.svg)](https://github.com/gohugoio/hugo/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gohugoio/hugo)](https://goreportcard.com/report/github.com/gohugoio/hugo)

[Website] | [Installation] | [Documentation] | [Support] | [Contributing] | <a rel="me" href="https://fosstodon.org/@gohugoio">Mastodon</a>

## Resumen

Hugo es un [generador de sitios estáticos] escrito en [Go], optimizado para la velocidad y diseñado para la flexibilidad. Con su avanzado sistema de plantillas y rápidos pipelines de activos, Hugo renderiza un sitio completo en segundos, a menudo menos.

Debido a su marco flexible, soporte multilingüe y poderoso sistema de taxonomía, Hugo es ampliamente utilizado para crear:

- Sitios corporativos, gubernamentales, sin fines de lucro, educativos, de noticias, de eventos y de proyectos
- Sitios de documentación
- Portafolios de imágenes
- Páginas de aterrizaje
- Blogs empresariales, profesionales y personales
- Currículums y CVs

Usa el servidor web integrado de Hugo durante el desarrollo para ver instantáneamente los cambios en el contenido, la estructura, el comportamiento y la presentación. Luego despliega el sitio en tu anfitrión, o sube los cambios a tu proveedor de Git para construcciones y despliegues automatizados.

Los rápidos pipelines de activos de Hugo incluyen:

- Empaquetado de CSS &ndash; transpile (Sass), eliminación de código muerto, minificación, mapas de fuente, hashing SRI e integración con PostCSS
- Empaquetado de JavaScript &ndash; transpile (TypeScript, JSX), eliminación de código muerto, minificación, mapas de fuente y hashing SRI
- Procesamiento de imágenes &ndash; convertir, redimensionar, recortar, rotar, ajustar colores, aplicar filtros, superponer texto e imágenes y extraer datos EXIF

Y con [Hugo Modules], puedes compartir contenido, activos, datos, traducciones, temas, plantillas y configuración con otros proyectos a través de repositorios Git públicos o privados.

## Patrocinadores

<p>&nbsp;</p>
<p float="left">
  <a href="https://www.linode.com/?utm_campaign=hugosponsor&utm_medium=banner&utm_source=hugogithub" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/assets/images/sponsors/linode-logo_standard_light_medium.png" width="200" alt="Linode"></a>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="https://cloudcannon.com/hugo-cms/?utm_campaign=HugoSponsorship&utm_source=sponsor&utm_content=gohugo" target="_blank"><img src="https://raw.githubusercontent.com/gohugoio/gohugoioTheme/master/assets/images/sponsors/cloudcannon-blue.svg" width="220" alt="CloudCannon"></a>
<p>&nbsp;</p>

## Instalación

Instala Hugo desde un [binario precompilado], gestor de paquetes o repositorio de paquetes. Por favor, consulta las instrucciones de instalación para tu sistema operativo:

- [macOS]
- [Linux]
- [Windows]
- [DragonFly BSD, FreeBSD, NetBSD y OpenBSD]

## Construir desde el código fuente

Hugo está disponible en dos ediciones: estándar y extendida. Con la edición extendida puedes:

- Codificar al formato WebP al procesar imágenes. Puedes decodificar imágenes WebP con cualquiera de las ediciones.
- Transpilar Sass a CSS usando el transpiler LibSass integrado. La edición extendida no es requerida para usar el transpiler Dart Sass.

Requisitos previos para construir Hugo desde el código fuente:

- Edición estándar: Go 1.20 o posterior
- Edición extendida: Go 1.20 o posterior, y GCC

Construir la edición estándar:

```text
go install github.com/gohugoio/hugo@latest
```

Construir la edición extendida:
```text
CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
```

## Documentación

La [documentación] de Hugo incluye instrucciones de instalación, una guía rápida, explicaciones conceptuales, información de referencia y ejemplos.

Por favor, envía problemas de documentación y pull requests al [repositorio de documentación].

## Soporte

Por favor no uses la cola de problemas para preguntas o solución de problemas. A menos que estés seguro de que tu problema es un defecto del software, usa el [foro].

El [foro] de Hugo es una comunidad activa de usuarios y desarrolladores que responden preguntas, comparten conocimiento y proporcionan ejemplos. Una búsqueda rápida de más de 20,000 temas a menudo responderá tu pregunta. Por favor asegúrate de leer sobre [solicitar ayuda] antes de hacer tu primera pregunta.

## Contribuir

Puedes contribuir al proyecto Hugo:

- Respondiendo preguntas en el [foro]
- Mejorando la [documentación]
- Monitoreando la [cola de problemas]
- Creando o mejorando [temas]
- Eliminando [errores]

Por favor, envía problemas de documentación y pull requests al [repositorio de documentación].

Si tienes una idea para una mejora o una nueva característica, crea un nuevo tema en el [foro] en la categoría "Feature". Esto te ayudará a:

- Determinar si la capacidad ya existe
- Medir el interés
- Refinar el concepto
Si hay suficiente interés, [crea una propuesta]. No envíes un pull request hasta que el líder del proyecto acepte la propuesta.

Para una guía completa sobre cómo contribuir a Hugo, consulta la Guía de Contribución.

Dependencias
Hugo se apoya en grandes bibliotecas de código abierto. Ejecuta hugo env --logLevel info para mostrar una lista de dependencias.

<details>
<summary>Ver dependencias actuales</summary>
cloud.google.com/go/compute/metadata="v0.2.3"
cloud.google.com/go/iam="v1.1.3"
cloud.google.com/go/storage="v1.31.0"
cloud.google.com/go="v0.110.8"
github.com/Azure/azure-sdk-for-go/sdk/azcore="v1.7.0"
github.com/Azure/azure-sdk-for-go/sdk/azidentity="v1.3.0"
github.com/Azure/azure-sdk-for-go/sdk/internal="v1.3.0"
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob="v1.1.0"
github.com/Azure/go-autorest/autorest/to="v0.4.0"
github.com/AzureAD/microsoft-authentication-library-for-go="v1.0.0"
github.com/BurntSushi/locker="v0.0.0-20171006230638-a6e239ea1c69"
github.com/PuerkitoBio/purell="v1.1.1"
github.com/PuerkitoBio/urlesc="v0.0.0-20170810143723-de5bf2ad4578"
github.com/alecthomas/chroma/v2="v2.11.1"
github.com/armon/go-radix="v1.0.0"
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream="v1.4.11"
github.com/aws/aws-sdk-go-v2/config="v1.18.32"
github.com/aws/aws-sdk-go-v2/credentials="v1.13.31"
github.com/aws/aws-sdk-go-v2/feature/ec2/imds="v1.13.7"
github.com/aws/aws-sdk-go-v2/feature/s3/manager="v1.11.76"
github.com/aws/aws-sdk-go-v2/internal/configsources="v1.1.37"
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2="v2.4.31"
github.com/aws/aws-sdk-go-v2/internal/ini="v1.3.38"
github.com/aws/aws-sdk-go-v2/internal/v4a="v1.1.0"
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding="v1.9.12"
github.com/aws/aws-sdk-go-v2/service/internal/checksum="v1.1.32"
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url="v1.9.31"
github.com/aws/aws-sdk-go-v2/service/internal/s3shared="v1.15.0"
github.com/aws/aws-sdk-go-v2/service/s3="v1.38.1"
github.com/aws/aws-sdk-go-v2/service/sso="v1.13.1"
github.com/aws/aws-sdk-go-v2/service/ssooidc="v1.15.1"
github.com/aws/aws-sdk-go-v2/service/sts="v1.21.1"
github.com/aws/aws-sdk-go-v2="v1.20.0"
github.com/aws/aws-sdk-go="v1.48.2"
github.com/aws/smithy-go="v1.14.0"
github.com/bep/clocks="v0.5.0"
github.com/bep/debounce="v1.2.0"
github.com/bep/gitmap="v1.1.2"
github.com/bep/goat="v0.5.0"
github.com/bep/godartsass/v2="v2.0.0"
github.com/bep/godartsass="v1.2.0"
github.com/bep/golibsass="v1.1.1"
github.com/bep/gowebp="v0.3.0"
github.com/bep/lazycache="v0.2.0"
github.com/bep/logg="v0.3.0"
github.com/bep/mclib="v1.20400.20402"
github.com/bep/overlayfs="v0.6.0"
github.com/bep/simplecobra="v0.3.2"
github.com/bep/tmc="v0.5.1"
github.com/clbanning/mxj/v2="v2.7.0"
github.com/cli/safeexec="v1.0.1"
github.com/cpuguy83/go-md2man/v2="v2.0.2"
github.com/disintegration/gift="v1.2.1"
github.com/dlclark/regexp2="v1.10.0"
github.com/dustin/go-humanize="v1.0.1"
github.com/evanw/esbuild="v0.19.7"
github.com/fatih/color="v1.16.0"
github.com/frankban/quicktest="v1.14.6"
github.com/fsnotify/fsnotify="v1.7.0"
github.com/getkin/kin-openapi="v0.120.0"
github.com/ghodss/yaml="v1.0.0"
github.com/go-openapi/jsonpointer="v0.19.6"
github.com/go-openapi/swag="v0.22.4"
github.com/gobuffalo/flect="v1.0.2"
github.com/gobwas/glob="v0.2.3"
github.com/gohugoio/go-i18n/v2="v2.1.3-0.20230805085216-e63c13218d0e"
github.com/gohugoio/locales="v0.14.0"
github.com/gohugoio/localescompressed="v1.0.1"
github.com/golang-jwt/jwt/v4="v4.5.0"
github.com/golang/groupcache="v0.0.0-20210331224755-41bb18bfe9da"
github.com/golang/protobuf="v1.5.3"
github.com/google/go-cmp="v0.6.0"
github.com/google/s2a-go="v0.1.7"
github.com/google/uuid="v1.4.0"
github.com/google/wire="v0.5.0"
github.com/googleapis/enterprise-certificate-proxy="v0.3.2"
github.com/googleapis/gax-go/v2="v2.12.0"
github.com/gorilla/websocket="v1.5.1"
github.com/hairyhenderson/go-codeowners="v0.4.0"
github.com/hashicorp/golang-lru/v2="v2.0.1"
github.com/invopop/yaml="v0.2.0"
github.com/jdkato/prose="v1.2.1"
github.com/jmespath/go-jmespath="v0.4.0"
github.com/josharian/intern="v1.0.0"
github.com/kr/pretty="v0.3.1"
github.com/kr/text="v0.2.0"
github.com/kylelemons/godebug="v1.1.0"
github.com/kyokomi/emoji/v2="v2.2.12"
github.com/mailru/easyjson="v0.7.7"
github.com/marekm4/color-extractor="v1.2.1"
github.com/mattn/go-colorable="v0.1.13"
github.com/mattn/go-isatty="v0.0.20"
github.com/mattn/go-runewidth="v0.0.9"
github.com/mitchellh/hashstructure="v1.1.0"
github.com/mitchellh/mapstructure="v1.5.0"
github.com/mohae/deepcopy="v0.0.0-20170929034955-c48cc78d4826"
github.com/muesli/smartcrop="v0.3.0"
github.com/niklasfasching/go-org="v1.7.0"
github.com/olekukonko/tablewriter="v0.0.5"
github.com/pelletier/go-toml/v2="v2.1.0"
github.com/perimeterx/marshmallow="v1.1.5"
github.com/pkg/browser="v0.0.0-20210911075715-681adbf594b8"
github.com/pkg/errors="v0.9.1"
github.com/rogpeppe/go-internal="v1.11.0"
github.com/russross/blackfriday/v2="v2.1.0"
github.com/rwcarlsen/goexif="v0.0.0-20190401172101-9e8deecbddbd"
github.com/sanity-io/litter="v1.5.5"
github.com/sass/dart-sass/compiler="1.63.2"
github.com/sass/dart-sass/implementation="1.63.2"
github.com/sass/dart-sass/protocol="2.0.0"
github.com/sass/libsass="3.6.5"
github.com/spf13/afero="v1.10.0"
github.com/spf13/cast="v1.5.1"
github.com/spf13/cobra="v1.7.0"
github.com/spf13/fsync="v0.9.0"
github.com/spf13/pflag="v1.0.5"
github.com/tdewolff/minify/v2="v2.20.7"
github.com/tdewolff/parse/v2="v2.7.5"
github.com/webmproject/libwebp="v1.3.2"
github.com/yuin/goldmark-emoji="v1.0.2"
github.com/yuin/goldmark="v1.6.0"
go.opencensus.io="v0.24.0"
go.uber.org/atomic="v1.11.0"
go.uber.org/automaxprocs="v1.5.3"
gocloud.dev="v0.34.0"
golang.org/x/crypto="v0.15.0"
golang.org/x/exp="v0.0.0-20221031165847-c99f073a8326"
golang.org/x/image="v0.13.0"
golang.org/x/mod="v0.14.0"
golang.org/x/net="v0.18.0"
golang.org/x/oauth2




