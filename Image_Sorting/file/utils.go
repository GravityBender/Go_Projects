package file

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

var extSlice = []string{".jpg", ".png", ".jpeg", ".JPG", ".PNG", ".JPEG"}

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
func ContainsFileWithExtension(folderPath string) (bool, error) {
	var found bool = false
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && slices.Contains(extSlice, filepath.Ext(path)) {
			found = true
			return filepath.SkipDir // Stop walking once we find a file with the extension
		}
		return nil
	})
	return found, err
}

/*
This function checks if the file represented by the filename is an image or not
*/
func IsImage(fileName string) bool {
	return slices.Contains(extSlice, filepath.Ext(fileName))
}

/*
This function creates a nested directory based on year and then month if not already created and returns the path of the directory
*/
func CreateDirIfNotCreated(date, folderPath string) (string, error) {
	layout := "2006:01:02 15:04:05"
	t, err := time.Parse(layout, date)
	if err != nil {
		return "", err
	}

	newDirPath := filepath.Join(folderPath, strconv.Itoa(t.Year()), t.Month().String())

	exists, err := Exists(newDirPath)
	if err != nil {
		return "", err
	} else if !exists {
		err := os.MkdirAll(newDirPath, os.ModePerm)
		if err != nil {
			return "", err
		}
		return newDirPath, nil
	} else {
		return newDirPath, nil
	}
}

/*
This function moves the file specified by the oldDir to the newDir
*/
func MoveFile(newDir, oldDir string) (bool, error) {
	err := os.Rename(oldDir, newDir)
	if err != nil {
		return false, err
	}

	return true, nil
}
