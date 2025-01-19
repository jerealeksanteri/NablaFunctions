package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Dockerfile string `yaml:"dockerfile"`
}

// BuildDockerImage builds a Docker image for the specified function.
// The dir parameter is the directory containing the function code.
// The language parameter is the language of the function.
// The handlerFile parameter is the name of the handler file.
// It returns the ID of the built image.
func BuildDockerImage(dir, language, handlerFile string) (string, error) {
	var dockerFileContent string

	// Load the template
	template, err := LoadTemplate(language)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %v", err)
	}

	// Generate the Dockerfilecontent
	if language == "python" {

		// Generate the Dockerfile for Python
		dockerFileContent = fmt.Sprintf(template.Dockerfile, handlerFile)

	}

	if language == "golang" {

		// Generate the Dockerfile for Golang
		dockerFileContent = template.Dockerfile

	}

	// Write the Dockerfile
	dockerFilePath := filepath.Join(dir, "Dockerfile")

	err = os.WriteFile(dockerFilePath, []byte(dockerFileContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write dockerfile: %v", err)
	}

	// Build the Docker image
	imageTag := fmt.Sprintf("NablaFunction:%s", language)
	cmd := exec.Command("docker", "build", "-t", imageTag, dir)

	// Combine the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to build docker image: %v\n%s", err, output)
	}

	// Extract the image ID
	imageId, err := ExtractImageID(string(output))
	if err != nil {
		return "", fmt.Errorf("failed to extract image ID: %v", err)
	}

	return imageId, nil

}

// ExtractImageID extracts the image ID from the output of the "docker build" command.
// The output is passed as an argument to this function.
func ExtractImageID(s string) (string, error) {

	// Split the output by newlines
	lines := strings.Split(s, "\n")

	// Iterate over the lines
	for _, line := range lines {

		// Writing the image on line sha256, so we can extract the image ID
		if strings.Contains(line, "writing image sha256") {

			// Split the line by spaces
			parts := strings.Fields(line)

			for _, part := range parts {
				if strings.HasPrefix(part, "sha256:") {
					return part, nil
				}
			}
		}
	}

	return "", fmt.Errorf("image ID not found")

}

// LoadTemplate loads the template for the specified language.
// The supported languages are "python" and "golang".
func LoadTemplate(language string) (*Template, error) {

	// Determine the path to the Dockerfile
	templateFile := fmt.Sprintf("templates/%s.yaml", language)

	// Read the Dockerfile
	data, err := os.ReadFile(templateFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read th template: %v", err)
	}

	// Parse the yaml file
	var template Template
	err = yaml.Unmarshal(data, &template)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal template: %v", err)
	}

	return &template, nil

}

func RunDockerContainer(imageId string) (string, error) {

	// Run the Docker container
	cmd := exec.Command("docker", "run", "--rm", imageId)

	// Combine the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run docker container: %v\n%s", err, output)
	}

	// Return the output
	return string(output), nil

}
