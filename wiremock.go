package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Mapping represents a WireMock mapping
type Mapping struct {
	ID      string            `json:"id"`
	Request RequestDefinition `json:"request"`
	// Add other fields as needed
}

// RequestDefinition represents the request part of a WireMock mapping
type RequestDefinition struct {
	URL string `json:"url"`
	// Add other fields as needed
}

// Manifest represents the structure of the JSON manifest file
type Manifest struct {
	Repositories []struct {
		URL string `json:"url"`
	} `json:"repositories"`
}

// fetchMappings fetches mappings from a repository URL
func fetchMappings(url string) ([]Mapping, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch mappings from %s. Status code: %d", url, response.StatusCode)
	}

	var mappings []Mapping
	if err := json.NewDecoder(response.Body).Decode(&mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

func main() {
	// Read repository URLs from JSON manifest file
	manifestFile := "manifest.json"
	manifestData, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		fmt.Println("Error reading manifest file:", err)
		return
	}

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		fmt.Println("Error parsing JSON manifest:", err)
		return
	}

	// URL of WireMock's admin API to delete all mappings
	adminURL := "http://localhost:8080/__admin/mappings"

	// Delete all existing mappings
	request, err := http.NewRequest(http.MethodDelete, adminURL, nil)
	if err != nil {
		fmt.Println("Error creating delete request:", err)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error deleting mappings:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Failed to delete mappings. Status code: %d\n", response.StatusCode)
		return
	}

	fmt.Println("All mappings deleted successfully.")

	// Import mappings from each repository
	for _, repo := range manifest.Repositories {
		mappings, err := fetchMappings(repo.URL)
		if err != nil {
			fmt.Printf("Error fetching mappings from %s: %v\n", repo.URL, err)
			continue
		}

		// Import fetched mappings to WireMock
		// Uncomment the following lines to import mappings
		/*
			importURL := "http://localhost:8080/__admin/mappings/import"
			err = importMappings(importURL, mappings)
			if err != nil {
				fmt.Printf("Error importing mappings from %s: %v\n", repo.URL, err)
				continue
			}
		*/
		fmt.Printf("Mappings from %s imported successfully.\n", repo.URL)
	}
}
