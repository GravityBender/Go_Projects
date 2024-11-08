package main

import (
	"Image_Sorting/file"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dsoprea/go-exif/v3"
)

var logger *slog.Logger

func configLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	logger = slog.New(handler)
}

func folderValidation(folderPath string) (bool, error) {
	isValid, err := file.Exists(folderPath)
	if err != nil {
		logger.Error("Error while checking if given directory exists or not.", slog.String("error", err.Error()))
		return false, err
	} else if !isValid {
		logger.Error("Invalid path provided!")
		return false, nil
	}

	isEmpty, err := file.IsDirEmpty(folderPath)
	if err != nil {
		logger.Error("Error while checking if given directory is empty or not.", slog.String("error", err.Error()))
		return false, err
	} else if isEmpty {
		logger.Error("Path provided is empty!")
		return false, nil
	}

	containsImage, err := file.ContainsFileWithExtension(folderPath)
	if err != nil {
		logger.Error("Error while checking if given directory contains images or not.", slog.String("error", err.Error()))
		return false, err
	} else if !containsImage {
		logger.Error("Path provided does not contain any image!")
		return false, nil
	}

	return true, nil
}

func initiateGrouping(folderPath string) (noExifDataFoundSlice, exifParsingErrorSlice []string) {
	dir, err := os.Open(folderPath)
	if err != nil {
		logger.Error("Error while opening folder!")
		logger.Error(err.Error())
		return nil, nil
	}

	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		logger.Error("Error while opening folder!")
		logger.Error(err.Error())
		return nil, nil
	}

	noExifDataFoundSlice = []string{}
	exifParsingErrorSlice = []string{}
	for _, fileDetail := range files {
		logger.Debug("Processing file", slog.String("FileName", fileDetail.Name()))

		if !file.IsImage(fileDetail.Name()) {
			logger.Info("File is not an image", slog.String("FileName", fileDetail.Name()))
			continue
		}

		// Read the EXIF data
		rawExif, err := exif.SearchFileAndExtractExif(filepath.Join(folderPath, fileDetail.Name()))
		if err != nil {
			noExifDataFoundSlice = append(noExifDataFoundSlice, fileDetail.Name())
			continue
		}

		// Parse the EXIF data
		logger.Info("Extracting EXIF data for: ", slog.String("File", fileDetail.Name()))
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
					logger.Error("Error while parsing EXIF data!")
					logger.Error(err.Error())
					return nil, nil
				}

				_, err = file.MoveFile(filepath.Join(formattedDate, fileDetail.Name()), filepath.Join(folderPath, fileDetail.Name()))

				if err != nil {
					logger.Error("Error while parsing moving file!", slog.String("FileName", fileDetail.Name()))
					logger.Error(err.Error())
					return nil, nil
				}
			}
		}

		//Read this value DateTimeOriginal

	}

	return noExifDataFoundSlice, exifParsingErrorSlice
}

func main() {
	configLogger()
	logger.Info("Welcome to the Image Sorting program!")

	if len(os.Args) < 2 {
		logger.Warn("Please provide the folder path to the image folder as a command line argument")
	}

	//	Get the folder path from the console
	folderPath := os.Args[1]
	folderPath = filepath.Clean(folderPath) //	Clean the folder path

	validationRes, err := folderValidation(folderPath)
	if err != nil || !validationRes {
		logger.Error(err.Error())
		return
	}

	noExifDataFoundSlice, exifParsingErrorSlice := initiateGrouping(folderPath)
	if noExifDataFoundSlice == nil || exifParsingErrorSlice == nil {
		logger.Error("Some error encountered while grouping the data!")
		return
	}

	if len(noExifDataFoundSlice) != 0 {
		logger.Warn("Following files had no EXIF data associated with them.", slog.String("FileNames", strings.Join(noExifDataFoundSlice, ", ")))
	}
	if len(exifParsingErrorSlice) != 0 {
		logger.Warn("Following files had no EXIF data associated with them.", slog.String("FileNames", strings.Join(exifParsingErrorSlice, ", ")))
	}

}
