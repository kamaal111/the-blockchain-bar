package commands

import (
	"fmt"
	"os"

	"github.com/kamaal111/the-blockchain-bar/database"
	"github.com/spf13/cobra"
)

func Balances() *cobra.Command {
	var command = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errIncorrectUsage
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	command.AddCommand(balancesListCommand)

	return command
}

var balancesListCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists all balances.",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := database.NewStateFromDisk()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer state.Close()

		fmt.Println("Accounts balances:")
		fmt.Print("__________________\n\n")
		for account, balance := range state.Balances {
			fmt.Printf("%s: %d\n", account, balance)
		}
	},
}
