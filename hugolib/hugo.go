package hugolib

import (
	"html/template"
)

const Version = "0.14-DEV"

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
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		Generator:  `<meta name="generator" content="Hugo ` + Version + `" />`,
	}
}
