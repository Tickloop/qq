package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tickloop/qq/internal/app"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/inference"
	"github.com/tickloop/qq/internal/utils"
)

func checkProviderSupported(args config.CLIArgs) {
	_, ok := inference.ProviderConverseFnMap[args.Provider]
	if !ok {
		fmt.Println("ERR: Unknown provider")
		os.Exit(1)
	}
}

func main() {
	args := config.LoadArgs()
	utils.Dbg(
		"model=%s provider=%s question=%s configure=%v",
		args.ModelId,
		args.Provider,
		args.Question,
		args.Configure,
	)

	if args.Configure {
		fmt.Println("Configuration complete")
		os.Exit(0)
	}

	checkProviderSupported(args)
	m := app.NewQAModel(args)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("ERR: %v", err)
	}
}
