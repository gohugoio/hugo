package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMakePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Foo bar  ", "Foo-bar"},
		{"Foo.Bar/foo_Bar-Foo", "Foo.Bar/foo_Bar-Foo"},
		{"fOO,bar:foo%bAR", "fOObarfoobAR"},
		{"FOo/BaR.html", "FOo/BaR.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
		{"Банковский кассир", "Банковский-кассир"},
	}

	for _, test := range tests {
		output := MakePath(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestMakePathToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  FOO bar  ", "foo-bar"},
		{"Foo.Bar/fOO_bAr-Foo", "foo.bar/foo_bar-foo"},
		{"FOO,bar:Foo%Bar", "foobarfoobar"},
		{"foo/BAR.HTML", "foo/bar.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
	}
	for _, test := range tests {
		output := MakePathToLower(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestMakeTitle(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"Make-Title", "Make Title"},
		{"MakeTitle", "MakeTitle"},
		{"make_title", "make_title"},
	}
	for i, d := range data {
		output := MakeTitle(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

//
// comment
//
func TestReplaceExtension(t *testing.T) {
	type test struct {
		input, newext, expected string
	}
	data := []test{
		// none of these work as you might expect
		{"/some/randome/path/file.xml", "html", "/some/randompath/file.html"},
		{"/banana.html", "HTML", "/banana.HTML"},
		{"./banana.html", "xml", "./banana.xml"},
		{"banana/pie/index.html", "xml", "banana/pie/index.xml"},
		{"../pies/fish/index.html", "xml", "../pies/fish/index.xml"},
		// but these all work even though they make no sense!
		{"/directory", "ext", "/directory.ext"},
		{"/directory/mydir/", "ext", "/directory/mydir/.ext"},
		{"mydir/", "ext", "mydir/.ext"},
	}

	for i, d := range data {
		output := ReplaceExtension(d.input, d.newext)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func TestDirExists(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}

	data := []test{
		{".", true},
		{"./", true},
		{"..", true},
		{"../", true},
		{"./..", true},
		{"./../", true},
		{"/tmp", true},
		{"/tmp/", true},
		{"/", true},
		{"/some-really-random-directory-name", false},
		{"/some/really/random/directory/name", false},
		{"./some-really-random-local-directory-name", false},
		{"./some/really/random/local/directory/name", false},
	}

	for i, d := range data {
		exists, _ := DirExists(d.input)
		if d.expected != exists {
			t.Errorf("Test %d failed. Expected %t got %t", i, d.expected, exists)
		}
	}
}

func TestIsDir(t *testing.T) {
	type test struct {
		input    string
		expected bool
	}
	data := []test{
		{"./", true},
		{"/", true},
		{"./this-directory-does-not-existi", false},
		{"/this-absolute-directory/does-not-exist", false},
	}

	for i, d := range data {
		exists, _ := IsDir(d.input)
		if d.expected != exists {
			t.Errorf("Test %d failed. Expected %t got %t", i, d.expected, exists)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	zeroSizedFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(zeroSizedFile)
	nonZeroSizedFile, _ := createNonZeroSizedFileInTempDir()
	defer deleteFileInTempDir(nonZeroSizedFile)
	emptyDirectory, _ := createEmptyTempDir()
	defer deleteTempDir(emptyDirectory)
	nonEmptyZeroLengthFilesDirectory, _ := createTempDirWithZeroLengthFiles()
	defer deleteTempDir(nonEmptyZeroLengthFilesDirectory)
	nonEmptyNonZeroLengthFilesDirectory, _ := createTempDirWithNonZeroLengthFiles()
	defer deleteTempDir(nonEmptyNonZeroLengthFilesDirectory)
	nonExistentFile := os.TempDir() + "/this-file-does-not-exist.txt"
	nonExistentDir := os.TempDir() + "/this/direcotry/does/not/exist/"

	fileDoesNotExist := fmt.Errorf("%q path does not exist", nonExistentFile)
	dirDoesNotExist := fmt.Errorf("%q path does not exist", nonExistentDir)

	type test struct {
		input          string
		expectedResult bool
		expectedErr    error
	}

	data := []test{
		{zeroSizedFile.Name(), true, nil},
		{nonZeroSizedFile.Name(), false, nil},
		{emptyDirectory, true, nil},
		{nonEmptyZeroLengthFilesDirectory, false, nil},
		{nonEmptyNonZeroLengthFilesDirectory, false, nil},
		{nonExistentFile, false, fileDoesNotExist},
		{nonExistentDir, false, dirDoesNotExist},
	}
	for i, d := range data {
		exists, err := IsEmpty(d.input)
		if d.expectedResult != exists {
			t.Errorf("Test %d failed. Expected result %t got %t", i, d.expectedResult, exists)
		}
		if d.expectedErr != nil {
			if d.expectedErr.Error() != err.Error() {
				t.Errorf("Test %d failed. Expected %q(%#v) got %q(%#v)", i, d.expectedErr, d.expectedErr, err, err)
			}
		} else {
			if d.expectedErr != err {
				t.Errorf("Test %d failed. Expected %q(%#v) got %q(%#v)", i, d.expectedErr, d.expectedErr, err, err)
			}
		}
	}
}

func createZeroSizedFileInTempDir() (*os.File, error) {
	filePrefix := "_path_test_"
	f, e := ioutil.TempFile("", filePrefix) // dir is os.TempDir()
	if e != nil {
		// if there was an error no file was created.
		// => no requirement to delete the file
		return nil, e
	}
	return f, nil
}

func createNonZeroSizedFileInTempDir() (*os.File, error) {
	f, err := createZeroSizedFileInTempDir()
	if err != nil {
		// no file ??
	}
	byteString := []byte("byteString")
	err = ioutil.WriteFile(f.Name(), byteString, 0644)
	if err != nil {
		// delete the file
		deleteFileInTempDir(f)
		return nil, err
	}
	return f, nil
}

func deleteFileInTempDir(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		// now what?
	}
}

func createEmptyTempDir() (string, error) {
	dirPrefix := "_dir_prefix_"
	d, e := ioutil.TempDir("", dirPrefix) // will be in os.TempDir()
	if e != nil {
		// no directory to delete - it was never created
		return "", e
	}
	return d, nil
}

func createTempDirWithZeroLengthFiles() (string, error) {
	d, dirErr := createEmptyTempDir()
	if dirErr != nil {
		//now what?
	}
	filePrefix := "_path_test_"
	_, fileErr := ioutil.TempFile(d, filePrefix) // dir is os.TempDir()
	if fileErr != nil {
		// if there was an error no file was created.
		// but we need to remove the directory to clean-up
		deleteTempDir(d)
		return "", fileErr
	}
	// the dir now has one, zero length file in it
	return d, nil

}

func createTempDirWithNonZeroLengthFiles() (string, error) {
	d, dirErr := createEmptyTempDir()
	if dirErr != nil {
		//now what?
	}
	filePrefix := "_path_test_"
	f, fileErr := ioutil.TempFile(d, filePrefix) // dir is os.TempDir()
	if fileErr != nil {
		// if there was an error no file was created.
		// but we need to remove the directory to clean-up
		deleteTempDir(d)
		return "", fileErr
	}
	byteString := []byte("byteString")
	fileErr = ioutil.WriteFile(f.Name(), byteString, 0644)
	if fileErr != nil {
		// delete the file
		deleteFileInTempDir(f)
		// also delete the directory
		deleteTempDir(d)
		return "", fileErr
	}

	// the dir now has one, zero length file in it
	return d, nil

}

func deleteTempDir(d string) {
	err := os.RemoveAll(d)
	if err != nil {
		// now what?
	}
}

func TestExists(t *testing.T) {
	zeroSizedFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(zeroSizedFile)
	nonZeroSizedFile, _ := createNonZeroSizedFileInTempDir()
	defer deleteFileInTempDir(nonZeroSizedFile)
	emptyDirectory, _ := createEmptyTempDir()
	defer deleteTempDir(emptyDirectory)
	nonExistentFile := os.TempDir() + "/this-file-does-not-exist.txt"
	nonExistentDir := os.TempDir() + "/this/direcotry/does/not/exist/"

	type test struct {
		input          string
		expectedResult bool
		expectedErr    error
	}

	data := []test{
		{zeroSizedFile.Name(), true, nil},
		{nonZeroSizedFile.Name(), true, nil},
		{emptyDirectory, true, nil},
		{nonExistentFile, false, nil},
		{nonExistentDir, false, nil},
	}
	for i, d := range data {
		exists, err := Exists(d.input)
		if d.expectedResult != exists {
			t.Errorf("Test %d failed. Expected result %t got %t", i, d.expectedResult, exists)
		}
		if d.expectedErr != err {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expectedErr, err)
		}
	}

}

// TestAbsPathify cannot be tested further because it relies on the
// viper.GetString("WorkingDir") which the test cannot know.
// viper.GetString("WorkingDir") should be passed to AbsPathify as a
// parameter.
func TestAbsPathify(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{os.TempDir(), os.TempDir()},
		{"/banana/../dir/", "/dir"},
	}

	for i, d := range data {
		expected := AbsPathify(d.input)
		if d.expected != expected {
			t.Errorf("Test %d failed. Expected %q but go %q", i, d.expected, expected)
		}
	}
}

func TestFilename(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"index.html", "index"},
		{"./index.html", "index"},
		{"/index.html", "index"},
		{"index", "index"},
		{"/tmp/index.html", "index"},
		{"./filename-no-ext", "filename-no-ext"},
		{"/filename-no-ext", "filename-no-ext"},
		{"filename-no-ext", "filename-no-ext"},
		{"directoy/", ""}, // no filename case??
		{"directory/.hidden.ext", ".hidden"},
	}

	for i, d := range data {
		output := Filename(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func TestFileAndExt(t *testing.T) {
	type test struct {
		input, expectedFile, expectedExt string
	}
	data := []test{
		{"index.html", "index", ".html"},
		{"./index.html", "index", ".html"},
		{"/index.html", "index", ".html"},
		{"index", "index", ""},
		{"/tmp/index.html", "index", ".html"},
		{"./filename-no-ext", "filename-no-ext", ""},
		{"/filename-no-ext", "filename-no-ext", ""},
		{"filename-no-ext", "filename-no-ext", ""},
		{"directoy/", "", ""}, // no filename case??
		{"directory/.hidden.ext", ".hidden", ".ext"},
	}

	for i, d := range data {
		file, ext := FileAndExt(d.input)
		if d.expectedFile != file {
			t.Errorf("Test %d failed. Expected filename %q got %q.", i, d.expectedFile, file)
		}
		if d.expectedExt != ext {
			t.Errorf("Test %d failed. Expected extension $q got %q.", i, d.expectedExt, ext)
		}
	}

}

func TestGuessSection(t *testing.T) {
	type test struct {
		input, expected string
	}

	data := []test{
		{"/", ""},
		{"", ""},
		{"/content", "/"},
		{"content/", "/"},
		{"/content/", "/"},
		{"/blog", "/blog"},
		{"/blog/", "/blog/"},
		{"blog", "blog"},
		{"content/blog", "/blog"},
		{"content/blog/", "/blog/"},
		{"/content/blog", "/blog/"},
	}

	for i, d := range data {
		expected := GuessSection(d.input)
		if d.expected != expected {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, expected)
		}
	}
}

func TestPathPrep(t *testing.T) {

}

func TestPrettifyPath(t *testing.T) {

}

func TestFindCWD(t *testing.T) {
	type test struct {
		expectedDir string
		expectedErr error
	}

	cwd, _ := os.Getwd()
	data := []test{
		{cwd, nil},
	}
	for i, d := range data {
		dir, err := FindCWD()
		if d.expectedDir != dir {
			t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedDir, dir)
		}
		if d.expectedErr != err {
			t.Error("Test %d failed. Expected %q but got %q", i, d.expectedErr, err)
		}
	}
}

func TestSafeWriteToDisk(t *testing.T) {
	emptyFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(emptyFile)
	tmpDir, _ := createEmptyTempDir()
	defer deleteTempDir(tmpDir)
	os.MkdirAll(tmpDir+"/this/dir/does/not/exist/", 0644)

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	fileExists := fmt.Errorf("%v already exists", emptyFile.Name())
	type test struct {
		filename    string
		expectedErr error
	}

	data := []test{
		{emptyFile.Name(), fileExists},
		{tmpDir + "/" + emptyFile.Name(), nil},
	}

	for i, d := range data {
		e := SafeWriteToDisk(d.filename, reader)
		if d.expectedErr != nil {
			if d.expectedErr.Error() != e.Error() {
				t.Errorf("Test %d failed. Expected error %q but got %q", i, d.expectedErr.Error(), e.Error())
			}
		} else {
			if d.expectedErr != e {
				t.Errorf("Test %d failed. Expected %q but got %q", i, d.expectedErr, e)
			}
			contents, _ := ioutil.ReadFile(d.filename)
			if randomString != string(contents) {
				t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
			}
		}
		reader.Seek(0, 0)
	}
}

func TestWriteToDisk(t *testing.T) {
	emptyFile, _ := createZeroSizedFileInTempDir()
	defer deleteFileInTempDir(emptyFile)
	tmpDir, _ := createEmptyTempDir()
	defer deleteTempDir(tmpDir)

	randomString := "This is a random string!"
	reader := strings.NewReader(randomString)

	type test struct {
		filename    string
		expectedErr error
	}

	data := []test{
		{emptyFile.Name(), nil},
		{tmpDir + "/abcd", nil},
	}

	for i, d := range data {
		//		fmt.Printf("Writing to: %s\n", d.filename)
		//		dir, _ := filepath.Split(d.filename)
		//		ospath := filepath.FromSlash(dir)

		//		fmt.Printf("dir: %q, ospath: %q\n", dir, ospath)
		e := WriteToDisk(d.filename, reader)
		//		fmt.Printf("Error from WriteToDisk: %s\n", e)
		//		f, e := os.Open(d.filename)
		//		if e != nil {
		//			fmt.Errorf("could not open file %s, error %s\n", d.filename, e)
		//		}
		//		fi, e := f.Stat()
		//		if e != nil {
		//			fmt.Errorf("Could not stat file %s, error %s\n", d.filename, e)
		//		}
		//		fmt.Printf("file size of %s is %d bytes\n", d.filename, fi.Size())
		if d.expectedErr != e {
			t.Errorf("Test %d failed. WriteToDisk Error Expected %q but got %q", i, d.expectedErr, e)
		}
		//		fmt.Printf("Reading from %s\n", d.filename)
		contents, e := ioutil.ReadFile(d.filename)
		//		fmt.Printf("Error from ReadFile: %s\n", e)
		//		fmt.Printf("Contents: %#v, %s\n", contents, string(contents))
		if e != nil {
			t.Error("Test %d failed. Could not read file %s. Reason: %s\n", i, d.filename, e)
		}
		if randomString != string(contents) {
			t.Errorf("Test %d failed. Expected contents %q but got %q", i, randomString, string(contents))
		}
		reader.Seek(0, 0)
	}
}
