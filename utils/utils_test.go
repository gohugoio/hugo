package utils

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	jww "github.com/spf13/jwalterweatherman"
)

type testData struct {
	logLevel        string
	logError        string
	logStr          []string
	logFileExpected bool
}

func TestCutUsageMessage(t *testing.T) {
	tests := []struct {
		message    string
		cutMessage string
	}{
		{"", ""},
		{" Usage of hugo: \n  -b, --baseUrl=...", ""},
		{"Some error Usage of hugo: \n", "Some error"},
		{"Usage of hugo: \n -b --baseU", ""},
		{"CRITICAL error for usage of hugo ", "CRITICAL error for usage of hugo"},
		{"Invalid short flag a in -abcde", "Invalid short flag a in -abcde"},
	}

	for _, test := range tests {
		message := cutUsageMessage(test.message)
		if message != test.cutMessage {
			t.Errorf("Expected %#v, got %#v", test.cutMessage, message)
		}
	}
}

func TestCheckErr(t *testing.T) {
	tests := []testData{
		{"ERROR", "first test case", []string{""}, true},
		{"ERROR", "second test case", []string{"banana", "man"}, true},
		{"ERROR", "third test case", []string{"multi-word string"}, true},
		{"ERROR", "fourth test case", []string{"multiple", "multi-word strings"}, true},
		{"CRITICAL", "Oops no array of strings", []string{}, true},
	}
	for _, test := range tests {
		filename := setup(t)
		defer teardown(t, filename)
		CheckErr(errors.New(test.logError), test.logStr...) // converts the array of strings in test.logStr to a varadic - cool!
		checkLogFile(t, filename, &test)
	}
}

func TestDoStopOnErr(t *testing.T) {
	tests := []struct {
		message    string
		cutMessage string
		t          testData
	}{
		{"", "", testData{"", "", []string{}, false}},
		{" Usage of hugo: \n  -b, --baseUrl=...", "", testData{"", "", []string{}, false}},
		{"Some error Usage of hugo: \n", "Some error", testData{"CRITICAL", "Some error", []string{}, true}},
		// sould get the same output if we pass any array of strings and not via the error
		{"Some error Usage of hugo: \n", "Some error", testData{"CRITICAL", "", []string{"Some error"}, true}},
		{"Usage of hugo: \n -b --baseU", "", testData{"", "", []string{""}, false}},
		{"CRITICAL error for usage of hugo ", "CRITICAL error for usage of hugo", testData{"CRITICAL", "CRITICAL error for usage of hugo", []string{""}, false}},
		{"CRITICAL error for usage of hugo ", "CRITICAL error for usage of hugo", testData{"CRITICAL", "", []string{"CRITICAL error for usage of hugo"}, true}},
		{"Invalid short flag a in -abcde", "Invalid short flag a in -abcde", testData{"CRITICAL", "Invalid short flag a in -abcde", []string{""}, false}},
		{"Invalid short flag a in -abcde", "Invalid short flag a in -abcde", testData{"CRITICAL", "", []string{"Invalid short flag a in -abcde"}, true}},
	}

	for _, test := range tests {
		filename := setup(t)
		defer teardown(t, filename)
		doStopOnErr(errors.New(test.t.logError), test.t.logStr...) // converts the array of strings in test.logStr to a varadic - cool!
		checkLogFile(t, filename, &test.t)
	}

}

func checkLogFile(t *testing.T, filename string, test *testData) {
	if err := jww.LogHandle.(*os.File).Close(); err != nil {
		t.Errorf("Error: Could not close file \"f\". Error: %v\n", err)
	}
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not open the log file \"%s\". Failed with %v\n", filename, err)
	}
	// does the test expect to have a log file. If so, it must also have contents.
	if !logFileIsExpectedAndValid(t, filename, test, &contents) {
		return
	}
	r := bytes.NewReader(contents)
	scanner := bufio.NewScanner(r)
	errorMessageMatches := false
	for scanner.Scan() {
		line := scanner.Text()
		// lines in the log file are of the form:
		// <log level>: yyyy/mm/dd <string|error message>
		// we pase this format left to right in three sections
		checkForExpectedLogLevelOrFail(t, line, test)
		errorMessageMatches = checkForExpectedErrorMsg(t, line, test)
		// There was no match against the error message. So see if it matchs one of the error strings
		if !errorMessageMatches {
			checkForExpectedStingOrFail(t, line, test)
		}
	}
	if err = scanner.Err(); err != nil {
		t.Fatalf("Could not scan the next token in the log file. Failed with: %v\n", err)
	}
}

func logFileIsExpectedAndValid(t *testing.T, filename string, test *testData, contents *[]byte) bool {
	if test.logFileExpected {
		// yup, so then the file cannot be empty.
		if len(*contents) == 0 {
			t.Fatalf("Unexpected empty log file! Filename:\"%s\"\n", filename)
		}
		return true
	}
	// we don't expect a log file for this test so bail here.
	return false
}

func checkForExpectedLogLevelOrFail(t *testing.T, line string, test *testData) {
	regexpErrorLabel := "^" + test.logLevel
	validErrorLevel := regexp.MustCompile(regexpErrorLabel)
	if !validErrorLevel.MatchString(line) {
		// can't find the expected start of line string. So fail
		t.Fatalf("Did not find the expected log level \"%s\" at the start of the line \"%s\"\n", test.logLevel, line)
	}
}

func checkForExpectedErrorMsg(t *testing.T, line string, test *testData) bool {
	regexpValidErrorMsg := test.logError + "$"
	validErrorMsg := regexp.MustCompile(regexpValidErrorMsg)
	return validErrorMsg.MatchString(line)
}

func checkForExpectedStingOrFail(t *testing.T, line string, test *testData) {
	for _, s := range test.logStr {
		regexpstr := s + "$"
		validLineEnd := regexp.MustCompile(regexpstr)
		if validLineEnd.MatchString(line) {
			return
		}
	}
	// if we reach here there was no match.
	// Note: It's not possibe for this to be called with test.logStr as an empty array.
	// The proceeding call to checkForExpectedErrorMsg
	// in checkLogFile guarentees this. i.e. checkForExpecedErrorMsg will return true in this case.
	t.Fatalf("Did not find any of the strings \"%v\" in \"%s\"\n", test.logStr, line)
}

func setup(t *testing.T) string {
	// first set the logger
	// we can't use jww.UseLogTempFile for this, becase we need the file name
	// so we can delete the file in teardown function.
	// We should really fix jww.UseLogTempFile so we can access the temp file, or
	// better yet provide a "DeleteTempLogFile" function
	const logfilename = "utils_test_"
	f, err := ioutil.TempFile(os.TempDir(), logfilename)
	if err != nil {
		t.Errorf("Error: Could not create temporary file for the logger. Error: %#v\n", err)
	}
	jww.SetStdoutThreshold(jww.LevelFatal)
	// jww.SetLogFile generates the "Logging to .... " line on stdout.
	// Maybe we should update jww to remove the fmt.PrintF calls?
	jww.SetLogFile(f.Name())
	return f.Name()
}

func teardown(t *testing.T, f string) {
	if err := os.Remove(f); err != nil {
		t.Errorf("Error: Could not remove file \"f\". Error: %v\n", err)
	}
}
