package dappier

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test using a custom HTTP client
func TestNewDappierApp_WithCustomClient(t *testing.T) {
	// Custom HTTP client
	customClient := &http.Client{}

	// Initialize with custom client
	client, err := NewDappierApp("mock-api-key", WithHTTPClient(customClient))
	if err != nil {
		t.Fatalf("Failed to initialize Dappier client: %v", err)
	}

	if client.Client != customClient {
		t.Errorf("Expected custom client, but got different client")
	}
}

// Test using the default HTTP client
func TestNewDappierApp_DefaultClient(t *testing.T) {
	// Initialize without custom client (default should be set)
	client, err := NewDappierApp("mock-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize Dappier client: %v", err)
	}

	if client.Client == nil {
		t.Error("Expected default HTTP client, but got nil")
	}
}

// Test RealtimeSearchAPI with a mock server
func TestRealtimeSearchAPI(t *testing.T) {
	// Mock server for testing
	mockResponse := `[{"response": {"results": "Election Day is November 5, 2024"}}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	url := fmt.Sprintf("%s/%s", server.URL, RealtimeDataModelId)
	client, err := NewDappierApp("mock-api-key", withBaseURL(url))
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test valid query
	result, err := client.RealtimeSearchAPI("when is election in USA")
	if err != nil {
		t.Fatalf("RealtimeSearchAPI failed: %v", err)
	}

	expectedResult := "Election Day is November 5, 2024"
	if !strings.Contains(result.Response.Results, expectedResult) {
		t.Errorf("Expected result to contain: %s, got: %s", expectedResult, result.Response.Results)
	}

	// Test empty query (should return an error)
	_, err = client.RealtimeSearchAPI("")
	if err == nil {
		t.Error("Expected error for empty query, got none")
	}
}

// Test HTTP request errors
func TestRealtimeSearchAPI_RequestError(t *testing.T) {
	// Invalid server URL to simulate error
	client, err := NewDappierApp("mock-api-key", withBaseURL("http://invalid-url"))
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	_, err = client.RealtimeSearchAPI("when is election in USA")
	if err == nil {
		t.Error("Expected error due to invalid server URL, got none")
	}
}

// Test handling of non-200 status code
func TestRealtimeSearchAPI_Non200StatusCode(t *testing.T) {
	// Mock server that returns 500 status code
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	url := fmt.Sprintf("%s/%s", server.URL, RealtimeDataModelId)
	client, err := NewDappierApp("mock-api-key", withBaseURL(url))

	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	_, err = client.RealtimeSearchAPI("when is election in USA")
	if err == nil {
		t.Error("Expected error for non-200 status code, got none")
	}

	if !strings.Contains(err.Error(), "received non-OK response status: 500") {
		t.Errorf("Expected error message to contain status code 500, got: %v", err)
	}
}

// Test handling of invalid JSON response
func TestRealtimeSearchAPI_InvalidJSONResponse(t *testing.T) {
	// Mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	url := fmt.Sprintf("%s/%s", server.URL, RealtimeDataModelId)
	client, err := NewDappierApp("mock-api-key", withBaseURL(url))
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	_, err = client.RealtimeSearchAPI("when is election in USA")
	if err == nil {
		t.Error("Expected error for invalid JSON response, got none")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal response") {
		t.Errorf("Expected error message to indicate JSON unmarshaling issue, got: %v", err)
	}
}

// Test handling of empty response
func TestRealtimeSearchAPI_EmptyResponse(t *testing.T) {
	// Mock server that returns empty response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`)) // Empty response
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	url := fmt.Sprintf("%s/%s", server.URL, RealtimeDataModelId)
	client, err := NewDappierApp("mock-api-key", withBaseURL(url))

	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	_, err = client.RealtimeSearchAPI("when is election in USA")
	if err == nil {
		t.Error("Expected error for empty response, got none")
	}

	if !strings.Contains(err.Error(), "no results found") {
		t.Errorf("Expected error message to indicate no results found, got: %v", err)
	}
}

// TestAIRecommendations_DefaultValues tests the AIRecommendations method with default values (no optional params)
func TestAIRecommendations_DefaultValues(t *testing.T) {
	// Mock server for testing
	mockResponse := `{
		"results": [
			{
				"author": "Test Author",
				"image_url": "https://example.com/test.jpg",
				"preview_content": "Test content",
				"pubdate": "Tue, 10 Sep 2024 18:24:27 +0000",
				"pubdate_unix": 1725992667,
				"score": 0.75,
				"site": "Test Site",
				"site_domain": "testsite.com",
				"title": "Test Title",
				"url": "https://example.com/test"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	client, err := NewDappierApp("mock-api-key", withBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test AIRecommendations with default values (no optional params)
	result, err := client.AIRecommendations("latest tech news", "dm_02hr75e8ate6adr15hjrf3ikol")
	if err != nil {
		t.Fatalf("AIRecommendations failed: %v", err)
	}

	expectedTitle := "Test Title"
	if result.Results[0].Title != expectedTitle {
		t.Errorf("Expected title to be: %s, got: %s", expectedTitle, result.Results[0].Title)
	}
}

// TestAIRecommendations_WithOptionalParams tests the AIRecommendations method with custom similarityTopK and ref
func TestAIRecommendations_WithOptionalParams(t *testing.T) {
	// Mock server for testing
	mockResponse := `{
		"results": [
			{
				"author": "Test Author",
				"image_url": "https://example.com/test.jpg",
				"preview_content": "Test content",
				"pubdate": "Tue, 10 Sep 2024 18:24:27 +0000",
				"pubdate_unix": 1725992667,
				"score": 0.75,
				"site": "Test Site",
				"site_domain": "testsite.com",
				"title": "Test Title",
				"url": "https://example.com/test"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Use the mock server's URL instead of the real API URL
	client, err := NewDappierApp("mock-api-key", withBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test AIRecommendations with optional params
	result, err := client.AIRecommendations(
		"latest tech news",
		"dm_02hr75e8ate6adr15hjrf3ikol",
		WithSimilarityTopK(5),     // Custom similarityTopK
		WithRef("techcrunch.com"), // Custom ref
	)
	if err != nil {
		t.Fatalf("AIRecommendations failed: %v", err)
	}

	expectedTitle := "Test Title"
	if result.Results[0].Title != expectedTitle {
		t.Errorf("Expected title to be: %s, got: %s", expectedTitle, result.Results[0].Title)
	}
}

// TestAIRecommendations_EmptyQuery tests error handling for empty query
func TestAIRecommendations_EmptyQuery(t *testing.T) {
	client, err := NewDappierApp("mock-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test with empty query
	_, err = client.AIRecommendations("", "dm_02hr75e8ate6adr15hjrf3ikol")
	if err == nil {
		t.Error("Expected error for empty query, got none")
	}

	if !strings.Contains(err.Error(), "query cannot be empty") {
		t.Errorf("Expected 'query cannot be empty' error, got: %v", err)
	}
}

// TestAIRecommendations_EmptyDataModelID tests error handling for empty data model ID
func TestAIRecommendations_EmptyDataModelID(t *testing.T) {
	client, err := NewDappierApp("mock-api-key")
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}

	// Test with empty data model ID
	_, err = client.AIRecommendations("latest tech news", "")
	if err == nil {
		t.Error("Expected error for empty data model ID, got none")
	}

	if !strings.Contains(err.Error(), "datamodelID cannot be empty") {
		t.Errorf("Expected 'datamodelID cannot be empty' error, got: %v", err)
	}
}
