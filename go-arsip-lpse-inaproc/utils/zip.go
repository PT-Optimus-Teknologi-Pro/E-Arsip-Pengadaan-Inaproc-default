package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// createZip creates a single zip file with the given name and list of source files.
func CreateZip(filesToZip []string, filename string) (string, error) {
	outFile, err := os.CreateTemp("", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	
	for _, fileName := range filesToZip {
		err := AddFileToZip(zipWriter, fileName)
		if err != nil {
			zipWriter.Close()
			return "",fmt.Errorf("adding file %s to temp file: %w", fileName, err)
		}
	}
	
	// Explicitly close to ensure central directory is written
	zipWriter.Close()

	return outFile.Name(), nil
}

// addFileToZip opens a source file and copies its contents into the zip writer.
func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a zip file header with the file's name in the archive
	filename = filepath.ToSlash(filepath.Clean(filename))
	// get only filename , not fullpath
	filename = filepath.Base(filename)
	zipFile, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	// Copy the file content into the zip archive
	_, err = io.Copy(zipFile, file)
	return err
}
