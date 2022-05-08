package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var tbbCommand = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCommand.AddCommand(versionCommand)

	err := tbbCommand.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
