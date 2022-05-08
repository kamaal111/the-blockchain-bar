package database

type Account string

type Transaction struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func (transaction Transaction) IsReward() bool {
	return transaction.Data == "reward"
}
