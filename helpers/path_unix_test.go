// +build !windows

package helpers

import (
	"github.com/spf13/viper"
	"testing"
)

func TestPlatformAbsPathify(t *testing.T) {
	type test struct {
		inPath, workingDir, expected string
	}
	data := []test{
		{"/banana/../dir/", "/work", "/dir"},
	}

	for i, d := range data {
		// todo see comment in AbsPathify
		viper.Set("WorkingDir", d.workingDir)

		expected := AbsPathify(d.inPath)
		if d.expected != expected {
			t.Errorf("Test %d failed. Expected %q but got %q", i, d.expected, expected)
		}
	}
}
