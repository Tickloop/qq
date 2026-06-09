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
	"github.com/tickloop/qq/internal/render"
)

var debug *log.Logger
const defaultModel = "perplexity/sonar"
const defaultProvider = "openrouter"
var providerConverseFnMap = map[string]func(c context.Context, q, m string) (string, error){
	"openrouter": func(c context.Context, q, m string) (string, error) { return chat.OpenRouterConverse(c, q, m) },
	"bedrock": func(c context.Context, q, m string) (string, error) { return chat.AWSConverse(c, q, m) },
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

func parseArgs() map[string]string {
	/* 
		Returns a map[string]string
		keys:
			modelId: 
			question:
			provider:
	*/
	
	model := flag.String("model", defaultModel, "model to use")
	provider := flag.String("provider", defaultProvider, "provider for inference")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: qq [--model model] [--provider provider] <question...>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	question := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if question == "" {
		flag.Usage()
		os.Exit(1)
	}

	dbg("model=%s", *model)
	dbg("provider=%s", *provider)
	dbg("argv=%v", flag.Args())
	return map[string]string{
		"modelId": *model,
		"question": question,
		"provider": *provider,
	}
}


func main() {
	args := parseArgs()
	ctx := context.Background()
	dbg("hitting %s chat completions", args["provider"])

	var spin spinner.Spinner = spinner.NewANSISpinner(os.Stderr, 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()
	
	
	hldr, ok := providerConverseFnMap[args["provider"]]
	if !ok {
		err := fmt.Errorf("error: unknown provider %s", args["provider"])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	answer, err := hldr(ctx, args["question"], args["modelId"])
	spin.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(render.Something(answer))
}
