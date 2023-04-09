package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/gridiron-zone/gridiron/app"
)

func main() {
	rootCmd, _ := NewRootCmd(
		app.Name,
		app.DefaultNodeHome,
		app.ModuleBasics,
		app.NewGridiron,
		// this line is used by starport scaffolding # root/arguments
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
