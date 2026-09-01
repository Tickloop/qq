package agent

// Testing internal implemenations of tools for the agent
import (
	"testing"
)

// Functional test
func TestFirecrawlWebSearch(t *testing.T) {
	query := "What is a test?"
	_, err := webSearchFirecrawl(query)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log(responses)
}

// Unit
func TestConvertFirecrawlJSONToText(t *testing.T) {
	tests := []struct {
		name   string
		input  []map[string]any
		output string
	}{
		{
			name: "usual response",
			input: []map[string]any{
				{"title": "title 1", "url": "https://text1.com", "summary": "First summary"},
				{"title": "title 2", "url": "https://text2.com", "summary": "Second summary"},
			},
			output: "Result: 1\nTitle: title 1\nURL: https://text1.com\nSummary: First summary\n\n" +
				"Result: 2\nTitle: title 2\nURL: https://text2.com\nSummary: Second summary",
		},
		{
			name:   "empty results",
			input:  []map[string]any{},
			output: "",
		},
		{
			name: "different response format",
			input: []map[string]any{
				{"headline": "title 1", "link": "https://text1.com", "description": "First description"},
			},
			output: "Result: 1\nTitle: <nil>\nURL: <nil>\nSummary: <nil>",
		},
		{
			name: "one of many results is missing a field",
			input: []map[string]any{
				{"title": "title 1", "url": "https://text1.com", "summary": "First summary"},
				{"title": "title 2", "url": "https://text2.com"},
			},
			output: "Result: 1\nTitle: title 1\nURL: https://text1.com\nSummary: First summary\n\n" +
				"Result: 2\nTitle: title 2\nURL: https://text2.com\nSummary: <nil>",
		},
		{
			name:   "nil results",
			input:  nil,
			output: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := conertFirecrawlJSONToText(test.input)
			if got != test.output {
				t.Errorf("conertFirecrawlJSONToText() = %q, want %q", got, test.output)
			}
		})
	}
}
