package hugolib

const Version = "0.13-DEV"

var (
	CommitHash string
	BuildDate  string
)

// HugoInfo contains information about the current Hugo environment
type HugoInfo struct {
	Version    string
	Generator  string
	CommitHash string
	BuildDate  string
}

func newHugoInfo() *HugoInfo {
	return &HugoInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		Generator:  `<meta name="generator" content="Hugo ` + Version + `" />`,
	}
}
