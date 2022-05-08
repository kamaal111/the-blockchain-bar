package commands

import (
	"fmt"
	"os"

	"github.com/kamaal111/the-blockchain-bar/database"
	"github.com/spf13/cobra"
)

const (
	FLAG_FROM  = "from"
	FLAG_TO    = "to"
	FLAG_VALUE = "value"
	FLAG_DATA  = "data"
)

func Transactions() *cobra.Command {
	var command = &cobra.Command{
		Use:   "tx",
		Short: "Interact with transactions (add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return errIncorrectUsage
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	command.AddCommand(transactionsAddCommand())

	return command
}

func transactionsAddCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "add",
		Short: "Adds new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(FLAG_FROM)
			to, _ := cmd.Flags().GetString(FLAG_TO)
			value, _ := cmd.Flags().GetUint(FLAG_VALUE)
			data, _ := cmd.Flags().GetString(FLAG_DATA)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			transaction := database.NewTransaction(fromAcc, toAcc, value, data)

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// defer means, at the end of this function execution,
			// execute the following statement (close DB file with all TXs)
			defer state.Close()

			// Add the TX to an in-memory array (pool)
			err = state.Add(transaction)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			// Flush the memory pool transactions to disk
			err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("Transaction successfully added to the ledger.")
		},
	}

	command.Flags().String(FLAG_FROM, "", "From what account to send tokens")
	command.MarkFlagRequired(FLAG_FROM)

	command.Flags().String(FLAG_TO, "", "To what account to send tokens")
	command.MarkFlagRequired(FLAG_TO)

	command.Flags().Uint(FLAG_VALUE, 0, "How many tokens to send")
	command.MarkFlagRequired(FLAG_VALUE)

	command.Flags().String(FLAG_DATA, "", "Possible values: 'reward'")

	return command
}
