package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Snapshot [32]byte

type State struct {
	Balances              map[Account]uint
	transactionMemoryPool []Transaction

	databaseFile *os.File
	snapshot     Snapshot
}

func (state *State) doSnapshot() (Snapshot, error) {
	// Re-read the whole file from the first byte
	_, err := state.databaseFile.Seek(0, 0)
	if err != nil {
		return Snapshot{}, err
	}

	transactionData, err := ioutil.ReadAll(state.databaseFile)
	if err != nil {
		return Snapshot{}, err
	}

	snapshot := sha256.Sum256(transactionData)
	state.snapshot = snapshot

	return snapshot, err
}

func (state *State) Add(transaction Transaction) error {
	err := state.apply(transaction)
	if err != nil {
		return err
	}

	state.transactionMemoryPool = append(state.transactionMemoryPool, transaction)

	return nil
}

func (state *State) Persist() (Snapshot, error) {
	// Make a copy of mempool because the s.txMempool will be modified
	// in the loop below
	memoryPool := make([]Transaction, len(state.transactionMemoryPool))
	copy(memoryPool, state.transactionMemoryPool)

	var snapshot Snapshot
	for _, transaction := range memoryPool {
		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return snapshot, err
		}

		log.Printf("Persisting new transaction in to disk:\n\t%s\n", transactionJSON)
		_, err = state.databaseFile.Write(append(transactionJSON, '\n'))
		if err != nil {
			return snapshot, err
		}

		snapshot, err = state.doSnapshot()
		if err != nil {
			return snapshot, err
		}

		log.Printf("New database snapshot: %x\n", snapshot)
		// Remove the transaction written to a file from the mempool
		state.transactionMemoryPool = state.transactionMemoryPool[1:]
	}

	return snapshot, nil
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

	state := &State{balances, make([]Transaction, 0), transactionsDatabaseFile, Snapshot{}}

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
