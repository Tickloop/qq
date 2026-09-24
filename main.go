package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tickloop/qq/internal/app"
	"github.com/tickloop/qq/internal/config"
	"github.com/tickloop/qq/internal/utils"
)

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

	m := app.NewQAModel(args)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("ERR: %v", err)
	}
}
