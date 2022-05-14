package main

import (
	"fmt"
	"os"

	"github.com/kamaal111/the-blockchain-bar/commands"
	"github.com/spf13/cobra"
)

func main() {
	var tbbCommand = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCommand.AddCommand(commands.Version)
	tbbCommand.AddCommand(commands.Balances())
	tbbCommand.AddCommand(commands.Transactions())

	err := tbbCommand.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
