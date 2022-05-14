package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Balances              map[Account]uint
	transactionMemoryPool []Transaction

	blockFile       *os.File
	latestBlockHash Hash
}

func (state *State) AddTransaction(transaction Transaction) error {
	err := state.apply(transaction)
	if err != nil {
		return err
	}

	state.transactionMemoryPool = append(state.transactionMemoryPool, transaction)

	return nil
}

func (state *State) AddBlock(block Block) error {
	for _, transaction := range block.Transactions {
		err := state.AddTransaction(transaction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (state *State) Persist() (Hash, error) {
	now := uint64(time.Now().Unix())
	block := NewBlock(
		state.latestBlockHash,
		now,
		state.transactionMemoryPool,
	)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFS := BlockFS{blockHash, block}
	blockJSON, err := json.Marshal(blockFS)
	if err != nil {
		return Hash{}, nil
	}

	log.Printf("Persisting new Block in to disk:\n\t%s\n", blockJSON)
	_, err = state.blockFile.Write(append(blockJSON, '\n'))
	if err != nil {
		return Hash{}, err
	}

	state.latestBlockHash = blockHash
	state.transactionMemoryPool = []Transaction{}

	return blockHash, nil
}

func (state *State) Close() {
	state.blockFile.Close()
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

	blockFilepath := filepath.Join(databaseDirectory, "block.db")
	blockFile, err := os.OpenFile(blockFilepath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	state := &State{balances, make([]Transaction, 0), blockFile, Hash{}}

	scanner := bufio.NewScanner(blockFile)
	for scanner.Scan() {
		err = scanner.Err()
		if err != nil {
			return nil, err
		}

		var blockFS BlockFS
		err = json.Unmarshal(scanner.Bytes(), &blockFS)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockFS.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFS.Key
	}

	return state, nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
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

func (state *State) applyBlock(block Block) error {
	for _, transaction := range block.Transactions {
		err := state.apply(transaction)
		if err != nil {
			return err
		}
	}

	return nil
}
