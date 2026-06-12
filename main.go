package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tickloop/qq/internal/chat"
	"github.com/tickloop/qq/internal/render"
	"github.com/tickloop/qq/internal/spinner"
)

var debug *log.Logger

const (
	defaultModel    = "perplexity/sonar"
	defaultProvider = "openrouter"
	configPath      = "/home/arya/.config/qq/qq.conf"
)

type CLIArgs struct {
	ModelId  string `json:"modelId"`
	Question string `json:"question"`
	Provider string `json:"provider"`
}

var NotConfiguredErr = errors.New("error: config file not found")

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

func NewArgs() CLIArgs {
	return CLIArgs{
		ModelId:  defaultModel,
		Provider: defaultProvider,
		Question: "",
	}
}

func parseArgs() CLIArgs {
	args := NewArgs()
	flag.StringVar(&args.ModelId, "model", "", "model to use")
	flag.StringVar(&args.Provider, "provider", "", "provider for inference")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: qq [--model model] [--provider provider] <question...>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args.Question = strings.TrimSpace(strings.Join(flag.Args(), " "))
	if args.Question == "" {
		flag.Usage()
		os.Exit(1)
	}

	dbg("model=%s", args.ModelId)
	dbg("provider=%s", args.Provider)
	dbg("argv=%v", args.Question)

	return args
}

func loadConfig() (CLIArgs, error) {
	var args CLIArgs

	// check if config exists. if not - return not configured error
	info, err := os.Stat(configPath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return args, NotConfiguredErr
	}

	// check if somehow the user managed to create a directory instead of a file
	// and this is unrecoverable
	if info.IsDir() {
		return args, fmt.Errorf("error: expect %s to be a file but found dir", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return args, err
	}
	if err := json.Unmarshal(data, &args); err != nil {
		return args, err
	}
	return args, nil
}

func saveConfig(args CLIArgs) error {
	data, err := json.Marshal(args)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("/home/arya/.config/qq", 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return err
	}
	return nil
}

func handleConfigCreation(args CLIArgs) {
	fmt.Println("Hey! Seems like we haven't configured qq yet. Let's do that right now: ")

	// get provider
	fmt.Print("pick a provider [bedrock|openrouter]: ")
	var err error

	if args.Provider, err = readLine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: couldn't read input")
		os.Exit(1)
	}

	// get modelId
	fmt.Print("model-id: ")
	if args.ModelId, err = readLine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: couldn't read input")
		os.Exit(1)
	}

	// save config
	if err = saveConfig(args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readLine() (string, error) {
	buf := bufio.NewScanner(os.Stdin)
	if !buf.Scan() {
		if err := buf.Err(); err != nil {
			return "", err
		}
	}
	return buf.Text(), nil
}

func main() {
	// first check if we have a config ready to use,
	// if not get cli args,
	// if still missing required - prompt for configuration
	configArgs, err := loadConfig()
	cliArgs := parseArgs()

	if err != nil {
		// there was no config found and no cli args either
		if errors.Is(err, NotConfiguredErr) && (cliArgs.ModelId == "" || cliArgs.Provider == "") {
			handleConfigCreation(cliArgs)
		}
	} else {
		if cliArgs.ModelId == "" {
			cliArgs.ModelId = configArgs.ModelId
		}

		if cliArgs.Provider == "" {
			cliArgs.Provider = configArgs.Provider
		}
	}

	ctx := context.Background()
	dbg("hitting %s chat completions", cliArgs.Provider)

	var spin spinner.Spinner = spinner.NewANSISpinner(os.Stderr, 100*time.Millisecond)
	spin.Start()
	defer spin.Stop()

	hldr, ok := providerConverseFnMap[cliArgs.Provider]
	if !ok {
		err := fmt.Errorf("error: unknown provider %s", cliArgs.Provider)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	answer, err := hldr(ctx, cliArgs.Question, cliArgs.ModelId)
	spin.Stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(render.Something(answer))
}
