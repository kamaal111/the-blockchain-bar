package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Balances              map[Account]uint
	transactionMemoryPool []Transaction

	databaseFile *os.File
}

func (state *State) Add(transaction Transaction) error {
	err := state.apply(transaction)
	if err != nil {
		return err
	}

	state.transactionMemoryPool = append(state.transactionMemoryPool, transaction)

	return nil
}

func (state *State) Persist() error {
	// Make a copy of mempool because the s.txMempool will be modified
	// in the loop below
	memoryPool := make([]Transaction, len(state.transactionMemoryPool))
	copy(memoryPool, state.transactionMemoryPool)

	for _, transaction := range memoryPool {
		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return err
		}

		_, err = state.databaseFile.Write(append(transactionJSON, '\n'))
		if err != nil {
			return err
		}

		// Remove the transaction written to a file from the mempool
		state.transactionMemoryPool = state.transactionMemoryPool[1:]
	}

	return nil
}

func (state *State) Close() {
	state.databaseFile.Close()
}

func NewStateFromDisk() (*State, error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	databaseDirectory := filepath.Join(currentWorkingDirectory, "database")

	genesisFilepath := filepath.Join(databaseDirectory, "genesis.json")
	genesisContent, err := loadGenesis(genesisFilepath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range genesisContent.Balances {
		balances[account] = balance
	}

	transactionsDatabaseFilepath := filepath.Join(databaseDirectory, "transactions.db")
	transactionsDatabaseFile, err := os.OpenFile(transactionsDatabaseFilepath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	state := &State{balances, make([]Transaction, 0), transactionsDatabaseFile}

	transactionsDatabaseFileScanner := bufio.NewScanner(transactionsDatabaseFile)
	// Iterate over each the transaction database file's line
	for transactionsDatabaseFileScanner.Scan() {
		err = transactionsDatabaseFileScanner.Err()
		if err != nil {
			return nil, err
		}

		// Convert JSON encoded Transaction into an object (struct)
		var transaction Transaction
		json.Unmarshal(transactionsDatabaseFileScanner.Bytes(), &transaction)

		// Rebuild the state (user balances),
		// as a series of events
		err = state.apply(transaction)
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (state *State) apply(transaction Transaction) error {
	if transaction.IsReward() {
		state.Balances[transaction.To] += transaction.Value
		return nil
	}

	if transaction.Value > state.Balances[transaction.From] {
		return fmt.Errorf("insufficient balance")
	}

	state.Balances[transaction.From] -= transaction.Value
	state.Balances[transaction.To] += transaction.Value

	return nil
}
