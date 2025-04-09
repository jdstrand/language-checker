package walker

import (
	"os"
	"path/filepath"
)

// Walk is a helper function that will automatically skip the `.git` directory.
func Walk(root string, walkFn func(path string, typ os.DirEntry) error) error {
	return filepath.WalkDir(root, func(path string, typ os.DirEntry, err error) error {
		path = filepath.Clean(path)

		if typ.IsDir() && isDotGit(path) {
			return filepath.SkipDir
		}

		return walkFn(path, typ)
	})
}

func isDotGit(path string) bool {
	return filepath.Base(path) == ".git"
}
