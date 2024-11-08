package main

import (
	"Image_Sorting/file"
	"fmt"
	"os"
	"path/filepath"

	// "strings"

	"github.com/dsoprea/go-exif/v3"
)

func folderValidation(folderPath string) (bool, error) {
	isValid, err := file.Exists(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	} else if !isValid {
		fmt.Println("Invalid folder path provided!")
		return false, nil
	}

	isEmpty, err := file.IsDirEmpty(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	} else if isEmpty {
		fmt.Println("Folder path provided is empty!")
		return false, nil
	}

	containsImage, err := file.ContainsFileWithExtension(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return false, err
	} else if !containsImage {
		fmt.Println("Folder path provided does not contain any image!")
		return false, nil
	}

	return true, nil
}

func initiateGrouping(folderPath string) (noExifDataFoundSlice, exifParsingErrorSlice []string) {
	dir, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Error while opening folder!")
		return nil, nil
	}

	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
	}

	noExifDataFoundSlice = []string{}
	exifParsingErrorSlice = []string{}
	for _, fileDetail := range files {
		fmt.Println(fileDetail.Name())

		if !file.IsImage(fileDetail.Name()) {
			continue
		}

		// Read the EXIF data
		rawExif, err := exif.SearchFileAndExtractExif(filepath.Join(folderPath, fileDetail.Name()))
		if err != nil {
			noExifDataFoundSlice = append(noExifDataFoundSlice, fileDetail.Name())
			continue
		}

		// Parse the EXIF data
		fmt.Println("Extracting EXIF data for file: ", fileDetail.Name())
		_, err = exif.ParseExifHeader(rawExif)
		if err != nil {
			exifParsingErrorSlice = append(exifParsingErrorSlice, fileDetail.Name())
			continue
		}

		exifTagSlice, _, err := exif.GetFlatExifData(rawExif, &exif.ScanOptions{})

		if err != nil {
			return nil, nil
		}

		for _, exifTag := range exifTagSlice {
			if exifTag.TagId == 0x9003 {
				formattedDate, err := file.CreateDirIfNotCreated(exifTag.FormattedFirst, folderPath)
				if err != nil {
					fmt.Println("Error while parsing date. ", err)
					return nil, nil
				}

				_, err = file.MoveFile(filepath.Join(formattedDate, fileDetail.Name()), filepath.Join(folderPath, fileDetail.Name()))

				if err != nil {
					fmt.Println("Error while moving file! ", err)
					return nil, nil
				}
			}
		}

		//Read this value DateTimeOriginal

	}

	return noExifDataFoundSlice, exifParsingErrorSlice
}

func main() {
	fmt.Println("Welcome to the Image Sorting program!")

	if len(os.Args) < 2 {
		fmt.Println("Please provide the folder path to the image folder as a command line argument")
	}

	//	Get the folder path from the console
	folderPath := os.Args[1]
	folderPath = filepath.Clean(folderPath) //	Clean the folder path

	validationRes, err := folderValidation(folderPath)
	if err != nil || !validationRes {
		return
	}

	noExifDataFoundSlice, exifParsingErrorSlice := initiateGrouping(folderPath)

	fmt.Printf("Following files had no EXIF data associated with them: %v\n", noExifDataFoundSlice)
	fmt.Printf("Following files generated errors while having their EXIF data parsed: %v\n", exifParsingErrorSlice)

}
