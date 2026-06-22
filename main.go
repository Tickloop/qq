package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tickloop/qq/internal/chat"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/render"
	"github.com/tickloop/qq/internal/spinner"
)

var debug *log.Logger

var providerConverseFnMap = map[string]func(c context.Context, q, m string) (string, error){
	"openrouter": func(c context.Context, q, m string) (string, error) { return chat.OpenRouterConverse(c, q, m) },
	"bedrock":    func(c context.Context, q, m string) (string, error) { return chat.AWSConverse(c, q, m) },
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

func main() {
	args := loadArgs()
	dbg("model=%s", args.ModelId)
	dbg("provider=%s", args.Provider)
	dbg("question=%s", args.Question)
	dbg("configure=%v", args.Configure)

	if(args.Configure) {
		fmt.Println("Configuration complete")
		os.Exit(0)
	}

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

	answer, err := hldr(ctx, args.Question, args.ModelId)
	spin.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(render.RenderMarkdown(answer))
}
