module github.com/gohugoio/hugo

require (
	github.com/BurntSushi/locker v0.0.0-20171006230638-a6e239ea1c69
	github.com/BurntSushi/toml v0.3.1
	github.com/PuerkitoBio/purell v1.1.0
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38
	github.com/alecthomas/chroma v0.6.3
	github.com/alecthomas/repr v0.0.0-20181024024818-d37bc2a10ba1 // indirect
	github.com/aws/aws-sdk-go v1.16.23
	github.com/bep/debounce v1.2.0
	github.com/bep/gitmap v1.0.0
	github.com/bep/go-tocss v0.6.0
	github.com/chaseadamsio/goorgeous v1.1.0
	github.com/cpuguy83/go-md2man v1.0.8 // indirect
	github.com/disintegration/imaging v1.6.0
	github.com/dustin/go-humanize v1.0.0
	github.com/eknkc/amber v0.0.0-20171010120322-cdade1c07385
	github.com/fortytw2/leaktest v1.2.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gobwas/glob v0.2.3
	github.com/gohugoio/hugoTestHelpers/testmodBuilder/mods v0.0.0-20190513081324-4ece7d32a289
	github.com/google/go-cmp v0.2.0
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/go-immutable-radix v1.0.0
	github.com/jdkato/prose v1.1.0
	github.com/kyokomi/emoji v1.5.1
	github.com/magefile/mage v1.4.0
	github.com/markbates/inflect v1.0.0
	github.com/mattn/go-isatty v0.0.7
	github.com/miekg/mmark v1.3.6
	github.com/mitchellh/hashstructure v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/muesli/smartcrop v0.0.0-20180228075044-f6ebaa786a12
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nicksnyder/go-i18n v1.10.0
	github.com/olekukonko/tablewriter v0.0.0-20180506121414-d4647c9c7a84
	github.com/pkg/errors v0.8.1
	github.com/russross/blackfriday v1.5.2
	github.com/sanity-io/litter v1.1.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/fsync v0.0.0-20170320142552-12a01e648f05
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.3.0
	github.com/tdewolff/minify/v2 v2.3.7
	github.com/yosssi/ace v0.0.5
	gocloud.dev v0.13.0
	golang.org/x/image v0.0.0-20190321063152-3fc05d484e9f
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2
	gopkg.in/yaml.v2 v2.2.2
)

exclude github.com/chaseadamsio/goorgeous v2.0.0+incompatible

replace github.com/markbates/inflect => github.com/markbates/inflect v0.0.0-20171215194931-a12c3aec81a6

