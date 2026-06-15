package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultModel    = "perplexity/sonar"
	defaultProvider = "openrouter"
)

type CLIArgs struct {
	ModelId  string
	Question string
	Provider string
}

type Config struct {
	ModelId  string `json:"ModelId"`
	Provider string `json:"Provider"`
}

var NotConfiguredErr = errors.New("error: config file not found")

func NewArgs() CLIArgs {
	return CLIArgs{
		ModelId:  defaultModel,
		Provider: defaultProvider,
		Question: "",
	}
}

func readLine() (string, error) {
	buf := bufio.NewScanner(os.Stdin)
	if !buf.Scan() {
		if err := buf.Err(); err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(buf.Text()), nil
}

func getCofnigPaths() (string, string, error) {
	base, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	configDir := filepath.Join(base, ".config/qq")
	configPath := filepath.Join(configDir, "config.json")
	return configDir, configPath, nil
}

func ParseArgs() CLIArgs {
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
	return args
}

func LoadConfig() (CLIArgs, error) {
	var args CLIArgs
	_, configPath, err := getCofnigPaths()
	if err != nil {
		return args, fmt.Errorf("error: can't load config\n%w", err)
	}

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
	config := Config{
		ModelId:  args.ModelId,
		Provider: args.Provider,
	}
	configDir, configPath, err := getCofnigPaths()
	if err != nil {
		return fmt.Errorf("error: can't save config\n%w", err)
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return err
	}
	return nil
}

func HandleConfigCreation(args *CLIArgs) {
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
	if err = saveConfig(*args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
