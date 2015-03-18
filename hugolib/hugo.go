package hugolib

import (
	"fmt"
	"github.com/spf13/hugo/helpers"
	"html/template"
)

var (
	CommitHash string
	BuildDate  string
)

var hugoInfo *HugoInfo

// HugoInfo contains information about the current Hugo environment
type HugoInfo struct {
	Version    string
	Generator  template.HTML
	CommitHash string
	BuildDate  string
}

func init() {
	hugoInfo = &HugoInfo{
		Version:    helpers.HugoVersion(),
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		Generator:  template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, helpers.HugoVersion())),
	}
}
