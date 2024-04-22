package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Manifest represents the structure of the JSON manifest file
type Manifest struct {
	Repositories []struct {
		URL string `json:"url"`
	} `json:"repositories"`
}

// fetchJSONFromGitLab fetches JSON data from a GitLab repository using API token
func fetchJSONFromGitLab(url, apiToken string) ([]byte, error) {
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JSON from GitLab. Status code: %d", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	// GitLab API URL of the JSON file in the repository
	gitlabURL := "https://gitlab.com/api/v4/projects/:id/repository/files/:filename/raw"

	// Replace :id and :filename with the actual project ID and filename respectively
	gitlabURL = strings.ReplaceAll(gitlabURL, ":id", "project_id")
	gitlabURL = strings.ReplaceAll(gitlabURL, ":filename", "manifest.json")

	// Read API token from environment variable
	apiToken := os.Getenv("GITLAB_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Error: GitLab API token not found in environment variable GITLAB_API_TOKEN")
		return
	}

	// Fetch JSON data from GitLab repository
	jsonData, err := fetchJSONFromGitLab(gitlabURL, apiToken)
	if err != nil {
		fmt.Println("Error fetching JSON from GitLab:", err)
		return
	}

	// Parse JSON data
	var manifest Manifest
	if err := json.Unmarshal(jsonData, &manifest); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Process the manifest data as needed
	fmt.Println("Repositories:")
	for _, repo := range manifest.Repositories {
		fmt.Println(repo.URL)
	}
}
