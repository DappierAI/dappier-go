package dappier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseUrl             = "https://api.dappier.com/app/datamodel"
	RealtimeDataModelId = "dm_01hpsxyfm2fwdt2zet9cg6fdxt"
	ContentType         = "application/json"
)

// AiRecommendationsRequest represents the request payload structure for the Dappier AI recommendations API
type AiRecommendationsRequest struct {
	Query          string `json:"query"`            // Natural language query or URL
	SimilarityTopK int    `json:"similarity_top_k"` // The number of articles to return (default is 9)
	Ref            string `json:"ref"`              // Domain from which to fetch recommendations (e.g., techcrunch.com)
	NumArticlesRef int    `json:"num_articles_ref"` // Guaranteed number of articles from the specified domain
}

// AiRecommendationsResult represents the response structure for the Dappier AI recommendations API
type AiRecommendationsResult struct {
	Results []struct {
		Author         string  `json:"author"`          // Author of the article
		ImageURL       string  `json:"image_url"`       // URL of the article's image
		PreviewContent string  `json:"preview_content"` // Preview content of the article
		PubDate        string  `json:"pubdate"`         // Publication date of the article
		PubDateUnix    int64   `json:"pubdate_unix"`    // Publication date in Unix format
		Score          float64 `json:"score"`           // Relevance score of the article
		Site           string  `json:"site"`            // Name of the source site
		SiteDomain     string  `json:"site_domain"`     // Domain of the source site
		Title          string  `json:"title"`           // Title of the article
		URL            string  `json:"url"`             // URL of the article
	} `json:"results"`
}

type DappierApp struct {
	APIKey  string
	Client  *http.Client // optional custom HTTP client
	baseURL string       // Optional base URL for testing or customization
}

type RealtimeSearchRequest struct {
	Query string `json:"query"`
}

type RealtimeSearchResult struct {
	Response struct {
		Results string `json:"results"`
	} `json:"response"`
}

// Option is a function that modifies the DappierApp
type Option func(*DappierApp)

// WithHTTPClient allows the user to provide a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(app *DappierApp) {
		app.Client = client
	}
}

// withBaseURL allows the user to provide a custom base URL (for testing or other purposes)
func withBaseURL(baseURL string) Option {
	return func(app *DappierApp) {
		app.baseURL = baseURL
	}
}

// NewDappierApp initializes the client with the provided API key and optional configurations
func NewDappierApp(apiKey string, opts ...Option) (*DappierApp, error) {
	if apiKey == "" {
		return nil, errors.New("API key cannot be empty")
	}

	// Initialize DappierApp with default values
	app := &DappierApp{
		APIKey:  apiKey,
		Client:  &http.Client{}, // Default client if not provided
		baseURL: fmt.Sprintf("%s/%s", BaseUrl, RealtimeDataModelId),
	}

	// Apply options if provided
	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

// RealtimeSearchAPI makes a request to the Dappier API for real-time data retrieval
func (d *DappierApp) RealtimeSearchAPI(query string) (*RealtimeSearchResult, error) {
	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	// Create the request payload
	requestData := RealtimeSearchRequest{Query: query}
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}
	url := fmt.Sprintf("%s/%s", BaseUrl, RealtimeDataModelId)
	if d.baseURL != "" {
		url = d.baseURL
	}
	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Set necessary headers
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.APIKey))

	// Use the client's HTTP client to make the request
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %d", resp.StatusCode)
	}

	// Read the response body using io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response into the expected struct
	var result []RealtimeSearchResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if results are available
	if len(result) == 0 {
		return nil, errors.New("no results found")
	}

	return &result[0], nil
}

// RecommendationsOption defines a functional option for configuring the AIRecommendations request
type RecommendationsOption func(*AiRecommendationsRequest)

// WithSimilarityTopK sets the similarity_top_k option. The number of articles to return. Default is 9.
func WithSimilarityTopK(k int) RecommendationsOption {
	return func(req *AiRecommendationsRequest) {
		req.SimilarityTopK = k
	}
}

// WithRef sets the ref option
// The domain of the site from which the recommendations should come. For example, techcrunch.com.
func WithRef(ref string) RecommendationsOption {
	return func(req *AiRecommendationsRequest) {
		req.Ref = ref
	}
}

// WithNumArticlesRef sets the num_articles_ref option
// Specifies how many articles should be guaranteed to match the domain specified in ref.
func WithNumArticlesRef(num int) RecommendationsOption {
	return func(req *AiRecommendationsRequest) {
		req.NumArticlesRef = num
	}
}

// AIRecommendations makes a request to the Dappier API for Ai recommendations using a user-provided data model ID and request payload.
// Parameters:
// - query (string): A natural language query or URL. If a URL is passed, the AI analyzes the page, creates a summary, and performs a semantic search query based on the content.
// - datamodelID (string): The data model ID to be used for the API call.
// - similarityTopK (int): The number of articles to return. Default is 9. (optional - can be set with WithSimilarityTopK)
// - ref (string): The domain of the site from which the recommendations should come. For example, techcrunch.com. (optional - can be set with WithRef)
// - numArticlesRef (int): Specifies how many articles should be guaranteed to match the domain specified in ref. (optional - can be set with WithNumArticlesRef)
func (d *DappierApp) AIRecommendations(query string, datamodelID string, opts ...RecommendationsOption) (*AiRecommendationsResult, error) {
	// Check for required parameters
	if query == "" {
		return nil, errors.New("query cannot be empty")
	}
	if datamodelID == "" {
		return nil, errors.New("datamodelID cannot be empty")
	}

	// Create the request payload
	requestData := AiRecommendationsRequest{
		Query:          query,
		SimilarityTopK: 9,  // default value for similarityTopK
		Ref:            "", // default value for ref
		NumArticlesRef: 0,  // default value for numArticlesRef
	}

	// Apply all options to the request data
	for _, opt := range opts {
		opt(&requestData)
	}

	// Marshal the request data into JSON
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Construct the request URL with the provided data model ID
	url := fmt.Sprintf("%s/%s", BaseUrl, datamodelID)
	if d.baseURL != "" {
		url = d.baseURL
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	// Set necessary headers
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.APIKey))

	// Use the client's HTTP client to make the request
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %d", resp.StatusCode)
	}

	// Read the response body using io.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the response body into the AiRecommendationsResult structure
	var result AiRecommendationsResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if results are available
	if len(result.Results) == 0 {
		return nil, errors.New("no results found")
	}

	return &result, nil
}
