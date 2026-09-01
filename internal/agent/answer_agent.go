package agent

import (
	"context"
	"time"
	"fmt"
	"strings"

	firecrawl "github.com/firecrawl/firecrawl/apps/go-sdk"
	"github.com/firecrawl/firecrawl/apps/go-sdk/option"
)

const (
	nSearchResults = 3
)

func conertFirecrawlJSONToText(results []map[string]any) string {
	var output strings.Builder 
	for i, result := range results {
		fmt.Fprintf(
			&output,
			"Result: %d\nTitle: %v\nURL: %v\nSummary: %v\n\n",
			i + 1,
			result["title"],
			result["url"],
			result["summary"],
		)
	}
	return strings.TrimSpace(output.String())
}


func webSearchFirecrawl(query string) (string, error) {
	client, err := firecrawl.NewClient(
		option.WithMaxRetries(1),
		option.WithBackoffFactor(0.5),
		option.WithTimeout(10 * time.Second),
	)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	results, err := client.Search(ctx, query, &firecrawl.SearchOptions{
		Limit: firecrawl.Int(nSearchResults),
		ScrapeOptions: &firecrawl.ScrapeOptions{
			Formats: []string{"summary"},
		},
	})
	if err != nil {
		return "", err
	}
	return conertFirecrawlJSONToText(results.Web), nil
}


