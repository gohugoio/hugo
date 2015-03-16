package hugofs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var pathSeparator, opposite string
var baseTestPathCorrect, baseTestPathWrong string
var pcFs = new(PathCleaningOsFs)

func init() {
	// this looks weird, but there's some sanity hidden deep down!
	if runtime.GOOS == "windows" {
		pathSeparator = "\\"
		opposite = "/"
	} else {
		pathSeparator = "/"
		opposite = "\\"
	}

	tmp := filepath.Join(os.TempDir(), "pathcleaningfstest")
	os.Mkdir(tmp, 0777)
	baseTestPathCorrect = os.TempDir() + pathSeparator + "pathcleaningfstest" + pathSeparator
	baseTestPathWrong = strings.Replace(baseTestPathCorrect, pathSeparator, opposite, -1)
}

func TestClean(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect string
	}{
		{opposite + "sub1", filepath.FromSlash("/sub1")},
		{opposite + "sub1" + opposite + "sub2", filepath.FromSlash("/sub1/sub2")},
	} {
		result := clean(this.in)

		if result != this.expect {
			t.Errorf("[%d] Got %q but expected %q", i, result, this.expect)
		}
	}
}

func TestCreate(t *testing.T) {
	filename := "create.txt"
	nameWrong := baseTestPathWrong + filename
	nameCorrect := baseTestPathCorrect + filename
	f, err := pcFs.Create(nameWrong)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	f.Close()
	f, err = os.Open(nameCorrect)
	assert.Nil(t, err, fmt.Sprintf("%s", err))
	assert.NotNil(t, f)
	f.Close()
	os.Remove(nameCorrect)
}

func TestMkdir(t *testing.T) {
}

func TestMkdirAll(t *testing.T) {
}

func TestOpen(t *testing.T) {
}

func TestOpenFile(t *testing.T) {
}

func TestRemove(t *testing.T) {
}

func TestRemoveAll(t *testing.T) {
}

func TestRename(t *testing.T) {
}

func TestStat(t *testing.T) {
}

func TestChmod(t *testing.T) {
}

func TestChtimes(t *testing.T) {
}
