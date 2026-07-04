package allconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

func BenchmarkLoad(b *testing.B) {
	tempDir := b.TempDir()
	configFilename := filepath.Join(tempDir, "hugo.toml")
	config := `
baseURL = "https://example.com"
defaultContentLanguage = 'en'

[module]
[[module.mounts]]
source = 'content/en'
target = 'content/en'
[module.mounts.sites.matrix]
languages = 'en'
[[module.mounts]]
source = 'content/nn'
target = 'content/nn'
[module.mounts.sites.matrix]
languages = 'nn'
[[module.mounts]]
source = 'content/no'
target = 'content/no'
[module.mounts.sites.matrix]
languages = 'no'
[[module.mounts]]
source = 'content/sv'
target = 'content/sv'
[module.mounts.sites.matrix]
languages = 'sv'
[[module.mounts]]
source = 'layouts'
target = 'layouts'

[languages]
[languages.en]
title = "English"
weight = 1
[languages.nn]
title = "Nynorsk"
weight = 2
[languages.no]
title = "Norsk"
weight = 3
[languages.sv]
title = "Svenska"
weight = 4
`
	if err := os.WriteFile(configFilename, []byte(config), 0o666); err != nil {
		b.Fatal(err)
	}
	d := ConfigSourceDescriptor{
		Fs:       afero.NewOsFs(),
		Filename: configFilename,
	}

	for b.Loop() {
		_, err := LoadConfig(d)
		if err != nil {
			b.Fatal(err)
		}
	}
}
