package hugolib

import (
	"os"
	"strings"
)

func fileExt(path string) (file, ext string) {
	if strings.Contains(path, ".") {
		i := len(path) - 1
		for path[i] != '.' {
			i--
		}
		return path[:i], path[i+1:]
	}
	return path, ""
}

func replaceExtension(path string, newExt string) string {
	f, _ := fileExt(path)
	return f + "." + newExt
}

// Check if Exists && is Directory
func dirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
