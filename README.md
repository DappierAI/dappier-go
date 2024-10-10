# dappier-go

A Go client library for interacting with Dappier's API's.

## Overview

`dappier-go` provides a straightforward way to interact with Dappier's API's, which allows for real-time data search on the internet and other datamodels from the marketplace. The library is designed to be easy to use and integrate into existing Go projects.

## Installation

To install the package, run:

```bash
go get github.com/DappierAI/dappier-go

```

## Real Time Search API Usage

Real time search API - allows you to search the internet in real-time and get the most relevant search results.
More information on the API can be found [here](https://docs.dappier.com).

Here’s a basic example of how to use the library:

```go
package main

import (
	"fmt"
	"log"

	"github.com/dappier/dappier-go"
)

func main() {
	// Initialize Dappier client
	client, err := dappier.NewDappierApp("your-api-key")
	if err != nil {
		log.Fatalf("Failed to initialize Dappier client: %v", err)
	}

	// Make real-time search API request
	result, err := client.RealtimeSearchAPI("when is election in USA")
	if err != nil {
		log.Fatalf("Failed to get search result: %v", err)
	}

	// Print the search result
	fmt.Println("Search Result:", result.Response.Results)
}
```

## AI Recommendations API Usage

AI Recommendations API - allows you to get AI-driven recommendations for a selected datamodel from [marketplace](https://marketplace.dappier.com).

Here’s a basic example of how to use the library:

```go
package main

import (
	"fmt"
	"log"

	"github.com/dappier/dappier-go"
)

func main() {
	// Initialize Dappier client
	client, err := dappier.NewDappierApp("your-api-key")
	if err != nil {
		log.Fatalf("Failed to initialize Dappier client: %v", err)
	}

	// Example usage of DappierRagAPI
	// Parameters:
	// - query (string): A natural language query or a URL. If a URL is passed, 
	//                      the AI analyzes the page and performs a semantic search query.
	// - similarityTopK (int): Number of articles to return (default is 9).
	// - ref (string): The domain of the site to fetch recommendations from (e.g., techcrunch.com). optional
	// - numArticlesRef (int): Number of guaranteed articles from the specified domain (if ref is provided). optional
	// - datamodelID (string): The Data Model ID for the API request.
	
	recommendations, err := client.AIRecommendations("latest tech news", "dm_02hr75e8ate6adr15hjrf3ikol")
	if err != nil {
		log.Fatalf("Failed to get AI recommendations: %v", err)
	}

	// Print the results
	for _, result := range recommendations.Results {
		fmt.Printf("Title: %s\nAuthor: %s\nSite: %s\nURL: %s\n\n", result.Title, result.Author, result.Site, result.URL)
	}
	
	// With optional params
	recommendations, err = client.AIRecommendations(
		"latest tech news",
		"dm_02hr75e8ate6adr15hjrf3ikol",
		dappier.WithSimilarityTopK(5), // Set custom similarity_top_k
		dappier.WithRef("techcrunch.com"), // Set custom ref
		dappier.WithNumArticlesRef(5),
	)
	if err != nil {
		log.Fatalf("Failed to get AI recommendations: %v", err)
	}

	// Print the results
	for _, result := range recommendations.Results {
		fmt.Printf("Title: %s\nAuthor: %s\nSite: %s\nURL: %s\n\n", result.Title, result.Author, result.Site, result.URL)
	}
}

```

## Parameters

### `query` (string):
- A natural language query or URL.
- If a URL is passed, the AI analyzes the page, creates a summary, and performs a semantic search query based on the content.

### `similarity_top_k` (integer):
- The number of articles to return. Default is 9.

### `ref` (string):
- The domain of the site from which the recommendations should come.
- Example: `techcrunch.com`.

### `num_articles_ref` (integer):
- Specifies how many articles should be guaranteed to match the domain specified in `ref`.
- Use this to ensure a set number of articles from the desired domain appear in the results.

### `search_algorithm` (string):
- Options: `"most_recent"` or `"semantic"`.
    - `"semantic"` (default): Contextual matching of the query to retrieve articles.
    - `"most_recent"`: Retrieves articles sorted by the most recent publication date.

## Setup

1. Replace `/datamodel/[DATAMODEL_ID]` with the appropriate Data Model ID from your Data Model API for the respective AI Agent.
2. Replace the Bearer token in the Authorization header with your Dappier API key. This can be retrieved from the user profile settings.


## Search API Usage

The Search API allows you to query with different data_model_ids, enabling flexible data search across different models.

Here’s how to use the Search API:

```go
package main

import (
  "fmt"
  "log"
  "github.com/dappier/dappier-go"
)

func main() {
  // Initialize Dappier client
  client, err := dappier.NewDappierApp("your-api-key")
  if err != nil {
    log.Fatalf("Failed to initialize Dappier client: %v", err)
  }

  // Example: Search with default options
  searchResult, err := client.Search("Latest Microsoft News", "dm_01htjq2njgecvah7ncepm8v87y")
  if err != nil {
    log.Fatalf("Failed to get search results: %v", err)
  }

  // Print the results
  for _, result := range searchResult.Response.Results {
    fmt.Printf("Search Result - %v\n", result)
  }

  // Example: Search with custom options
  searchCustomResult, err := client.Search(
    "Latest Microsoft News",
    "dm_01htjq2njgecvah7ncepm8v87y",
    dappier.WithSearchSimilarityTopK(6),    // Custom similarity_top_k
    dappier.WithSearchRef("familyproof.com"), // Custom ref
    dappier.WithSearchNumArticlesRef(3),    // Ensure 3 articles from the domain
  )
  if err != nil {
    log.Fatalf("Failed to get search results: %v", err)
  }

  // Print the custom results
  for _, result := range searchCustomResult.Response.Results {
    fmt.Printf("Search Result (Custom) - %v\n", result)
  }
}


```

## Parameters

### `query` (string):
- A natural language query or URL.
- If a URL is passed, the AI analyzes the page, creates a summary, and performs a semantic search query based on the content.

### `similarity_top_k` (integer):
- The number of articles to return. Default is 9.

### `ref` (string):
- The domain of the site from which the recommendations should come.
- Example: `techcrunch.com`.

### `num_articles_ref` (integer):
- Specifies how many articles should be guaranteed to match the domain specified in `ref`.
- Use this to ensure a set number of articles from the desired domain appear in the results.

### `search_algorithm` (string):
- Options: `"most_recent"` or `"semantic"`.
  - `"semantic"` (default): Contextual matching of the query to retrieve articles.
  - `"most_recent"`: Retrieves articles sorted by the most recent publication date.


