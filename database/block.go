package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func (hash Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(hash[:])), nil
}

func (hash *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(hash[:], data)
	return err
}

func NewBlock(parent Hash, time uint64, transcations []Transaction) Block {
	header := BlockHeader{parent, time}
	return Block{header, transcations}
}

func (block *Block) Hash() (Hash, error) {
	json, err := json.Marshal(block)
	if err != nil {
		return Hash{}, err
	}

	hash := sha256.Sum256(json)
	return hash, nil
}
