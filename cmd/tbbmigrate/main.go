package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kamaal111/the-blockchain-bar/database"
)

func main() {
	state, err := database.NewStateFromDisk()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	now := uint64(time.Now().Unix())
	transactions := []database.Transaction{
		database.NewTransaction("andrej", "andrej", 3, ""),
		database.NewTransaction("andrej", "andrej", 700, "reward"),
	}
	block0 := database.NewBlock(
		database.Hash{},
		now,
		transactions,
	)

	state.AddBlock(block0)
	block0hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0hash,
		uint64(time.Now().Unix()),
		[]database.Transaction{
			database.NewTransaction("andrej", "babayaga", 2000, ""),
			database.NewTransaction("andrej", "andrej", 100, "reward"),
			database.NewTransaction("babayaga", "andrej", 1, ""),
			database.NewTransaction("babayaga", "caesar", 1000, ""),
			database.NewTransaction("babayaga", "andrej", 50, ""),
			database.NewTransaction("andrej", "andrej", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}
