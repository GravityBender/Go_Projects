package main

import (
	"Image_Sorting/file"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Welcome to the Image Sorting program!")

	if len(os.Args) < 2 {
		fmt.Println("Please provide the folder path to the image folder as a command line argument")
	}

	//	Get the folder path from the console
	folderPath := os.Args[1]
	folderPath = filepath.Clean(folderPath) //	Clean the folder path

	isValid, err := file.Exists(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else if !isValid {
		fmt.Println("Invalid folder path provided!")
		return
	}

	isEmpty, err := file.IsDirEmpty(folderPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} else if isEmpty {
		fmt.Println("folder path provided is empty!")
		return
	}

	fmt.Print("All tests passed!")

}