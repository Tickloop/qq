package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/tickloop/qq/internal/agent"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/inference"
	"github.com/tickloop/qq/internal/spinner"
	"golang.org/x/term"
)

var debug *log.Logger

var providerConverseFnMap = map[string]func(c context.Context, q, m string) (string, error){
	"openrouter": func(c context.Context, q, m string) (string, error) { return inference.OpenRouterConverse(c, q, m) },
	"bedrock":    func(c context.Context, q, m string) (string, error) { return inference.AWSConverse(c, q, m) },
}

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

func loadArgs() config.CLIArgs {
	// first check if we have a config ready to use,
	// if not get cli args,
	// if still missing required - prompt for configuration
	configArgs, err := config.LoadConfig()
	cliArgs := config.ParseArgs()

	// either the user needs to configure or the user wants to configure
	if (err != nil && errors.Is(err, config.NotConfiguredErr)) || cliArgs.Configure {
		config.HandleConfigCreation(&configArgs)
	}
	dbg("configArgs: %v", configArgs)
	dbg("cliArgs: %v", cliArgs)

	if cliArgs.ModelId == "" {
		cliArgs.ModelId = configArgs.ModelId
	}

	if cliArgs.Provider == "" {
		cliArgs.Provider = configArgs.Provider
	}
	return cliArgs
}

func styleQuestion(question string) string {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
	}

	styleQuestion := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAF0E6"))
	stylePrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#00bb00"))

	styleQuestionPrompt := lipgloss.NewStyle().Padding(1, 0)
	styleQuestionPrompt = styleQuestionPrompt.Width(w - styleQuestionPrompt.GetHorizontalFrameSize() - 1)
	
	return styleQuestionPrompt.Render(stylePrompt.Render("\033[1m›\033[0m ") + styleQuestion.Render(question))
}

func styleAnswer(answer string) string {
	out, err := glamour.Render(answer, "dark")
	if err != nil {
		// default: no markdown render
		return answer
	}
	return out
}

func addSearchResultToQuestion(question string) string {
	questionWithSearchResults := fmt.Sprintf("Question: %v\n", question)
	searchResults, err := agent.WebSearchFirecrawl(question)
	dbg("got search results: %v\n", searchResults)
	if err != nil {
		searchResults = "ERR: Search results failed to load"
	} else {
		questionWithSearchResults = fmt.Sprintf("Search Results: %v\n\nQuestion: %v\n", searchResults, question)
	}
	return questionWithSearchResults
}

func main() {
	args := loadArgs()
	dbg("model=%s", args.ModelId)
	dbg("provider=%s", args.Provider)
	dbg("question=%s", args.Question)
	dbg("configure=%v", args.Configure)

	fmt.Println(styleQuestion(args.Question))
	if args.Configure {
		fmt.Println("Configuration complete")
		os.Exit(0)
	}
	questionPrompt := addSearchResultToQuestion(args.Question)

	ctx := context.Background()
	dbg("hitting %s chat completions", args.Provider)

	var spin spinner.Spinner = spinner.NewANSISpinner(os.Stderr, 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	hldr, ok := providerConverseFnMap[args.Provider]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown provider %s\n", args.Provider)
		os.Exit(1)
	}

	answer, err := hldr(ctx, questionPrompt, args.ModelId)
	spin.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(styleAnswer(answer))
}
