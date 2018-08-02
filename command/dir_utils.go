package command

import (
	"os"
)

func appendSlashTo(dir string) string {
	if dir[len(dir)-1:] != "/" {
		dir = dir + "/"
	}
	return dir
}

func fileOrDirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func removeDuplicates(paths []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, path := range paths {
		if _, value := keys[path]; !value {
			keys[path] = true
			list = append(list, path)
		}
	}
	return list
}
