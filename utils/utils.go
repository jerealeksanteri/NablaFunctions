package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/rs/zerolog/log"
)

// CreateTemporaryDirectory creates a temporary directory.
func CreateTemporaryDirectory() (string, error) {
	return os.MkdirTemp("", "zip-extract")
}

// CleanupTemporaryDirectory removes a temporary directory.
// The dir parameter is the directory to remove.
func CleanupTemporaryDirectory(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Info().Msgf("Failed to remove temporary directory: %v\n", err)

	}
}

// SaveZipFile saves a zip file to a temporary directory.
// It returns the path to the saved file.
// The tempDir parameter is the directory where the zip file is saved.
// The filename parameter is the name of the zip file.
// The file parameter is the zip file to save.
func SaveZipFile(tempDir, filename string, file io.Reader) string {
	filePath := filepath.Join(tempDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		log.Info().Msgf("Failed to create file: %v\n", err)
		return ""
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Info().Msgf("Failed to save zip file: %v\n", err)
		return ""
	}

	return filePath
}

// ExtractZipFile extracts a zip file to a temporary directory.
// It returns the directory where the files were extracted.
// The zip file is deleted after extraction.
// The tempDir parameter is the directory where the zip file is saved.
// The zipFile parameter is the path to the zip file.
// The extractDir parameter is the directory where the files are extracted.
func ExtractZipFile(zipFile, tempDir string) (string, error) {

	// Create a directory to extract the zip file
	extractDir := filepath.Join(tempDir, "extracted")
	err := os.Mkdir(extractDir, 0755)
	if err != nil {
		log.Info().Msgf("Failed to create extract directory: %v\n", err)
		return "", err
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Info().Msgf("Failed to open zip file: %v\n", err)
		return "", err
	}

	// Close the zip file when done
	defer reader.Close()

	// Extract the files
	for _, file := range reader.File {
		// Open the file inside the zip
		path := filepath.Join(extractDir, file.Name)

		// Create the directory if it doesn't exist
		if file.FileInfo().IsDir() {
			err1 := os.MkdirAll(path, file.Mode())
			if err1 != nil {
				log.Info().Msgf("Failed to create directory: %v\n", err1)
				return "", err1
			}
			continue
		}

		// Create the file
		outfile, err1 := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err1 != nil {
			log.Error().Msgf("Failed to open file: %v", err1)
			return "", err1
		}

		// Close the file when done
		defer outfile.Close()

		// Open the file inside the zip
		zipFile, err1 := file.Open()
		if err1 != nil {
			log.Info().Msgf("Failed to open zip file: %v\n", err1)
			return "", err1
		}

		// Close the zip file when done
		defer zipFile.Close()

		// Copy the file
		_, err1 = io.Copy(outfile, zipFile)
		if err1 != nil {
			log.Info().Msgf("Failed to copy file: %v\n", err1)
			return "", err1
		}

	}

	// Return the directory where the files were extracted
	return extractDir, nil
}

// DetectHandlerFile detects the handler file in a directory.
// It returns the name of the handler file and the language.
// The supported languages are "python" and "golang".
func DetectHandlerFile(dir string) (string, string, error) {

	// Read the directory
	files, err := os.ReadDir(dir)

	// Check for errors
	if err != nil {
		log.Info().Msgf("Failed to read directory: %v\n", err)
		return "", "", err
	}

	// Find the handler file
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Switch statement to check the file extension
		switch filepath.Ext(file.Name()) {
		case ".py":
			return file.Name(), "python", nil
		case ".go":
			return file.Name(), "golang", nil
		}

	}

	// Return an error if no handler file was found
	return "", "", fmt.Errorf("No handler file found")
}
