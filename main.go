package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tickloop/qq/internal/chat"
	"github.com/tickloop/qq/internal/spinner"
)

const defaultModel = "perplexity/sonar"
var debug *log.Logger

func init() {
	if os.Getenv("QQ_DEBUG") != "" {
		debug = log.New(os.Stderr, "[qq] ", log.Ltime)
	}
}

func dbg(format string, args ...any) {
	if debug != nil {
		debug.Printf(format, args...)
	}
}

func main() {
	model := flag.String("m", defaultModel, "model to use")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: qq [-m model] <question...>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	dbg("model=%s args=%v", *model, flag.Args())

	question := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if question == "" {
		flag.Usage()
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: OPENROUTER_API_KEY is not set")
		os.Exit(1)
	}

	var provider chat.Provider = chat.NewHTTPClient("https://openrouter.ai/api/v1", apiKey, nil)

	req := chat.CompletionRequest{
		Model: *model,
		Messages: []chat.Message{
			{Role: "user", Content: question},
		},
	}
	dbg("request: %+v", req)

	ctx := context.Background()
	dbg("hitting OpenRouter chat completions")

	var spin spinner.Spinner = spinner.NewANSISpinner(os.Stderr, 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	resp, err := provider.Complete(ctx, req)
	spin.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dbg("raw response: %+v", resp)

	if len(resp.Choices) == 0 {
		fmt.Fprintln(os.Stderr, "error: no response from model")
		os.Exit(1)
	}

	text := strings.TrimSpace(resp.Choices[0].Message.Content)
	if text == "" {
		fmt.Fprintln(os.Stderr, "error: empty response from model")
		os.Exit(1)
	}

	dbg("parsed message: %s", text)
	dbg("displaying final answer")
	fmt.Println(text)
}
