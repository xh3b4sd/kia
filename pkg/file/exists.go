package file

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	tilde = "~/"
)

func Exists(file string) bool {
	return exists(abs(file))
}

func abs(file string) string {
	if strings.HasPrefix(file, tilde) {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}

		return filepath.Join(u.HomeDir, strings.TrimPrefix(file, tilde))
	} else {
		abs, err := filepath.Abs(file)
		if err != nil {
			panic(err)
		}

		return abs
	}
}

func exists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}

	return true
}
