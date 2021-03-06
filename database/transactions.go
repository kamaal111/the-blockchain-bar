package database

type Account string

type Transaction struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func NewAccount(value string) Account {
	return Account(value)
}

func (transaction Transaction) IsReward() bool {
	return transaction.Data == "reward"
}

func NewTransaction(from Account, to Account, value uint, data string) Transaction {
	return Transaction{from, to, value, data}
}
