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
	nSearchResults = 2
	maxAttempts = 5
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


func WebSearchFirecrawl(query string) (string, error) {
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

// type ToolDefinition struct {
// 	hldr func(params ...any) string
// 	toolDesc string
// }
//
// type ToolRegistry map[string]ToolDefinition
//
//
// type AgentQuestionAnswering struct {
// 	question string
// 	answer string
//
// 	tools ToolRegistry
// 	systemPrompt string	
// 	inferenceClient inference.InferenceClient
// }
//
//
// func NewAgentQuestionAnswering() *AgentQuestionAnswering {
// 	agent := &AgentQuestionAnswering{
// 		question: "",
// 		answer: "",
//
// 		tools: NewToolRegistry()
// 	}
// }
//
//
// func hldrWebSearchTool(query string) string {	
// 	searchAnswer, err := webSearchFirecrawl(query)
// 	if err != nil {
// 		return "ERR: Couldn't use web_search. Continue without this tool"
// 	}
// 	return searchAnswer
// }
//
// func hldrSubmitAnswer(answer string) string {
// 	return "Your answer has been submitted"
// }
//
//
// func AnswerQuestion(question string) string {
// 	contextWindow := createContextWindow()
// 	contextWindow.addSystemPrompt(systemPrompt)
// 	contextWindow.addToolDefinition(toolDef_WebSearch)
// 	contextWindow.addToolDefinition(toolDef_SubmitAnswer)
// 	contextWindow.addUserMessage(question)
//
// 	for i := range(maxAttempts) {
// 		response, err := inference(contextWindow)
// 		if err != nil {
// 			return fmt.Errorf("ERR: %w", err)
// 		}
//
// 		for _, thinking := response.Thinking {
// 			contextWindow.addThinking(thinking)
// 		}
//
// 		for _, reasoning := response.Reasoning {
// 			contextWindow.addReasoning(reasoning)
// 		}
//
//
// 		for _, toolCall := response.ToolCalls {
// 			if toolCall.Name == "web_search" {
// 				toolCallResult := hldrWebSearchTool(toolCall.Params["query"])
// 			} else if toolCall.Name == "submit_answer" {
// 				toolCallResult := hldrSubmitAnswer(toolCall.Params["answer"])
// 			} else {
// 				toolCallResult := fmt.Sprintf("ERR: Unknown tool - %v", toolCall.Name)
// 			}
// 			contextWindow.addToolCall(toolCall)
// 			contextWindow.addToolCallResult(toolCallResult)
// 		}
// 	}
//
// 	return ""
// }
