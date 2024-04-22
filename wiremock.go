package main

import (
	"bytes"
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

// fetchMappings fetches mappings from WireMock
func fetchMappings(url string) ([]Mapping, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch mappings. Status code: %d", response.StatusCode)
	}

	var mappings []Mapping
	if err := json.NewDecoder(response.Body).Decode(&mappings); err != nil {
		return nil, err
	}

	return mappings, nil
}

// updateMappings updates mappings in WireMock
func updateMappings(url string, mappings []Mapping) error {
	mappingsJSON, err := json.Marshal(mappings)
	if err != nil {
		return err
	}

	response, err := http.Put(url, "application/json", bytes.NewReader(mappingsJSON))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update mappings. Status code: %d", response.StatusCode)
	}

	fmt.Println("Mappings updated successfully.")
	return nil
}

func main() {
	// URL of WireMock's mappings endpoint
	wiremockURL := "http://localhost:8080/__admin/mappings"

	// Fetch existing mappings
	existingMappings, err := fetchMappings(wiremockURL)
	if err != nil {
		fmt.Println("Error fetching mappings:", err)
		return
	}

	// Assuming you have other arrays of mappings to merge
	// mappings1 := fetchMappings(someOtherURL)
	// mappings2 := fetchMappings(yetAnotherURL)

	// Merge mappings (assuming mappings1 and mappings2 are fetched from somewhere)
	// mergedMappings := append(existingMappings, mappings1...)
	// mergedMappings = append(mergedMappings, mappings2...)

	// Update WireMock with merged mappings
	// err = updateMappings(wiremockURL, mergedMappings)
	err = updateMappings(wiremockURL, existingMappings)
	if err != nil {
		fmt.Println("Error updating mappings:", err)
		return
	}
}
