package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LRz00/cli-lint/common"
)

func CollectFiles(root string, lang string) ([]string, error) {
	var files []string

	// WalkDir walks the file tree rooted at root, calling fn for each file or directory in the tree, including root.
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {

		ext, ok := common.SupportedFileFormats[lang]
		if !ok {
			return fmt.Errorf("unsupported language: %s", lang)
		}

		// If there is an error accessing the path, return the error.
		if err != nil {
			return err
		}

		// If the entry is a directory, skip it.
		if d.IsDir() {
			return nil
		}

		// If the file has the lang extension, add it to the list.
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
