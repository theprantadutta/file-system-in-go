package main

import (
	"bufio"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing_bubble_tea/logger"
	"time"
)

var totalFiles int
var currentFileIndex int

func main() {
	startTime := time.Now()

	reader := bufio.NewReader(os.Stdin)

	sourceFolder, _ := getInput("Enter source directory: ", reader)
	destinationFolder, _ := getInput("Enter destination directory: ", reader)

	numberOfFiles, countErr := countFiles(sourceFolder)
	totalFiles = numberOfFiles
	currentFileIndex = 1

	if countErr != nil {
		logger.Error("Something Went Wrong When Counting Files")
	}

	logger.Info("Number of Files ", numberOfFiles)

	err := copyDir(sourceFolder, destinationFolder)
	if err != nil {
		logger.Error("Error:", err)
	} else {
		logger.Success("Files and folders copied successfully!")
	}

	elapsedTime := time.Since(startTime).Seconds()
	logger.Info(fmt.Sprintf("Process took %.2f seconds", elapsedTime))
}

func getInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			logger.Error("Failed to Close the File: %v", err)
		}
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			logger.Error("Failed to Close the File: %v", err)
		}
	}(destinationFile)

	logger.Info("Copying Files ", sourceFile.Name(), "...")

	// Get file size for progress bar
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	description := fmt.Sprintf("[cyan][%d/%d][reset] Copying file...", currentFileIndex, totalFiles)

	// Use DefaultBytes to create the progress bar
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Copy the file with progress updates
	_, err = io.Copy(io.MultiWriter(destinationFile, bar), sourceFile)
	if err != nil {
		return err
	}

	fmt.Println()
	currentFileIndex++

	return nil
}

func copyDir(src, dst string) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dst, fileInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destinationPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func countFiles(src string) (int, error) {
	var fileCount int

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		fileCount++
		return nil
	})

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}
