package file

import (
	"io"
	"os"
	"path/filepath"
	"slices"
)

/*
This function checks if the given file path exists or not
*/
func Exists(folderPath string) (bool, error) {
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}

/*
This function checks if the opened folder is empty or not
*/
func IsDirEmpty(folderPath string) (bool, error) {

	f, err := os.Open(folderPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdir(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

/*
This function checks if the provided direcotry contains atleast one file that is of type image
*/
func ContainsFileWithExtension(folderPath string, fileExt []string) (bool, error) {
	var found bool = false
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && slices.Contains(fileExt, filepath.Ext(path)) {
			found = true
			return filepath.SkipDir // Stop walking once we find a file with the extension
		}
		return nil
	})
	return found, err
}
