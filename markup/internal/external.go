package internal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/markup/converter"
)

func ExternallyRenderContent(
	cfg converter.ProviderConfig,
	ctx converter.DocumentContext,
	content []byte, binaryName string, args []string,
) ([]byte, error) {
	logger := cfg.Logger

	if strings.Contains(binaryName, "/") {
		panic(fmt.Sprintf("should be no slash in %q", binaryName))
	}

	argsv := collections.StringSliceToInterfaceSlice(args)

	var out, cmderr bytes.Buffer
	argsv = append(argsv, hexec.WithStdout(&out))
	argsv = append(argsv, hexec.WithStderr(&cmderr))
	argsv = append(argsv, hexec.WithStdin(bytes.NewReader(content)))

	cmd, err := cfg.Exec.New(binaryName, argsv...)
	if err != nil {
		return nil, err
	}

	err = cmd.Run()

	// Most external helpers exit w/ non-zero exit code only if severe, i.e.
	// halting errors occurred. -> log stderr output regardless of state of err
	for _, item := range strings.Split(cmderr.String(), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			if err == nil {
				logger.Warnf("%s: %s", ctx.DocumentName, item)
			} else {
				logger.Errorf("%s: %s", ctx.DocumentName, item)
			}
		}
	}

	if err != nil {
		logger.Errorf("%s rendering %s: %v", binaryName, ctx.DocumentName, err)
	}

	return normalizeExternalHelperLineFeeds(out.Bytes()), nil
}

// Strips carriage returns from third-party / external processes (useful for Windows)
func normalizeExternalHelperLineFeeds(content []byte) []byte {
	return bytes.Replace(content, []byte("\r"), []byte(""), -1)
}

var pythonBinaryCandidates = []string{"python", "python.exe"}

func GetPythonBinaryAndExecPath() (string, string) {
	for _, p := range pythonBinaryCandidates {
		if pth := hexec.LookPath(p); pth != "" {
			return p, pth
		}
	}
	return "", ""
}
