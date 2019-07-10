package main

import (
	"archive/zip"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func isPathSep(r rune) bool {
	if r > 128 {
		return true
	}
	return os.IsPathSeparator(uint8(r))
}

func exists(filename string) bool {
	_, err := os.Lstat(filename)
	return !os.IsNotExist(err)
}

func zipDir(directory string) (string, error) {
	info, err := os.Lstat(directory)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", directory)
	}

	directory = strings.TrimRightFunc(directory, isPathSep)

	zipFile := directory + ".zip"
	i := 0
	for exists(zipFile) {
		i++
		zipFile = fmt.Sprintf("%s.%d.zip", directory, i)
	}

	lastPathSep := strings.LastIndexFunc(directory, isPathSep)
	zipBase := directory
	if lastPathSep != -1 {
		zipBase = directory[lastPathSep+1:]
	}

	// Get a Buffer to Write To
	outFile, err := os.Create(zipFile)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, directory, zipBase)

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return "", err
	}
	return zipFile, nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		filename := filepath.Join(basePath, file.Name())
		zipFilename := filepath.Join(baseInZip, file.Name())
		// fmt.Println(filename)
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Println(err)
			}

			f, err := w.Create(zipFilename)
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {
			addFiles(w, filename, zipFilename)
		}
	}
}
