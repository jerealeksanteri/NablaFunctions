package handlers

import (
	"fmt"
	"net/http"

	"NablaFunctions/docker"
	"NablaFunctions/utils"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
)

// imageStore is a map of image IDs to image names.
var imageStore = make(map[string]string)

// LoadHandler handles the /api/load endpoint.
// It loads a function into the server.
// The function is stored in a Docker image.
// The function is extracted from a zip file.
// The zip file contains the function code and a template file.
// The template file specifies the language of the function.
// The function is loaded into the server using the specified language.
// The function is stored in a Docker image with the specified name.
// The function is stored in the imageStore map.
// The function returns the ID of the image.
func LoadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form, with max 10 MB of memory
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("code")
	if err != nil {
		http.Error(w, "Unable to get the zip file", http.StatusBadRequest)
		return
	}
	// defer the closing of the file
	defer file.Close()

	// Create a temporary directory
	tempDir, err := utils.CreateTemporaryDirectory()
	if err != nil {
		http.Error(w, "Unable to create temporary directory", http.StatusInternalServerError)
		return
	}
	// defer the cleanup of the temporary directory
	defer utils.CleanupTemporaryDirectory(tempDir)

	// Save the zip file to the temporary directory
	zipFilePath := utils.SaveZipFile(tempDir, "function.zip", file)
	if zipFilePath == "" {
		http.Error(w, "Unable to save zip file", http.StatusInternalServerError)
		return
	}

	// Extract the zip file to the temporary directory
	extractDir, err := utils.ExtractZipFile(zipFilePath, tempDir)
	if err != nil {
		http.Error(w, "Unable to extract zip file", http.StatusInternalServerError)
		return
	}

	// Get the language of the function
	filename, language, err := utils.DetectHandlerFile(extractDir)
	if err != nil {
		http.Error(w, "Unable to get language", http.StatusBadRequest)
		return
	}

	// Build the Docker image
	imageId, err := docker.BuildDockerImage(extractDir, language, filename)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to build docker Image")
		http.Error(w, "Unable to build Docker image", http.StatusInternalServerError)
		return
	}

	// Generate uuid for the function
	functionId := uuid.New().String()

	// Store the image ID in the imageStore
	imageStore[functionId] = imageId

	// Send the function Id back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Docker Image has been build successfully for function with ID: %s and Image ID: %s", functionId, imageId)))

}

// ExecuteHandler handles the /api/execute endpoint.
// It executes a function that has been loaded into the server.
// The function is executed using the specified function ID.
// The function ID is used to retrieve the Docker image from the imageStore.
// The function is executed in a Docker container.
// The function is executed with the specified input.
// The function returns the output of the function.
func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	functionId := r.URL.Query().Get("functionId")
	if functionId == "" {
		http.Error(w, "Function ID is required", http.StatusBadRequest)
		return
	}

	// Get the image ID from the imageStore
	imageId, ok := imageStore[functionId]
	if !ok {
		http.Error(w, "Function ID not found", http.StatusNotFound)
		return
	}

	// Run the Docker container
	output, err := docker.RunDockerContainer(imageId)
	if err != nil {
		log.Error().Err(err).Msg("Unable to run Docker container")
		http.Error(w, "Unable to run Docker container", http.StatusInternalServerError)
		return
	}

	// Send the output back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Output: \n %s", output)))

}

// Logging middleware logs the request method and URI.
// It calls the next handler in the chain.
// The next handler is passed as a parameter.
// The next handler is called after logging the request.
func LoggingMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("Request Method: %s , URI: %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
