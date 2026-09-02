package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultModel    = "perplexity/sonar"
	defaultProvider = "openrouter"
)

func NewArgs() CLIArgs {
	return CLIArgs{
		ModelId:  defaultModel,
		Provider: defaultProvider,
		Question: "",
	}
}


func ParseCliArgs() CLIArgs {
	args := NewArgs()
	flag.StringVar(&args.ModelId, "model", "", "model to use")
	flag.StringVar(&args.Provider, "provider", "", "provider for inference")
	flag.BoolVar(&args.Configure, "configure", false, "configure qq")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: qq [--model model] [--provider provider] [--configure] <question...>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args.Question = strings.TrimSpace(strings.Join(flag.Args(), " "))
	if args.Question == "" && !args.Configure {
		flag.Usage()
		os.Exit(1)
	}
	return args
}

