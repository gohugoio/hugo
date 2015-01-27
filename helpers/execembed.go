package helpers

import (
	"bytes"
	"os/exec"

	jww "github.com/spf13/jwalterweatherman"
)

func Exec(args ...string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var arg []string
	if len(args) == 0 {
		jww.ERROR.Print("Nothing to execute")
		return "Nothing to execute"
	}
	name := args[0]
	if len(args) > 1 {
		arg = args[1 : len(args)-1]
	}
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		jww.ERROR.Print("Error executing", err, stderr.String())
		return stderr.String()
	}
	jww.ERROR.Print("It got here")

	return out.String()
}
