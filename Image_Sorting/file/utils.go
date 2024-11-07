package file

import (
	"io"
	"os"
)

/*
* This function checks if the given file path exists or not
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
*	This function checks if the opened folder is empty or not
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
