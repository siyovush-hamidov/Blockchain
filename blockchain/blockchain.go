package blockchain

import (
	"database/sql"
	"os"
	"time"
)

type BlockChain struct {
	DB *sql.DB
}

type Block struct {
	Nonce uint64
	Difficulty uint8
	CurrHash []byte
	PrevHash []byte
	Transaction []Transaction
	Mapping map[string]uint64 // Состояния балансов пользователей
	Miner string
	Signature []byte
	TimeStamp string
}

// DEF: Nonce - это число, которое майнер подбирает, чтобы хэш блока соответствовал заданной сложности сети.

type Transaction struct {
	RandBytes []byte
	PrevBlock []byte
	Sender string
	Receiver string
	Value uint64
	ToStorage uint64 // количество переводимых средств хранилищу
	CurrHash []byte
	Signature []byte
}
// Que: Какой максимальный размер uint в Golang?

func NewChain(filename, receiver string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	file.Close()
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}

	defer db.Close()

	_, err := db.Exec(CREATE_TABLE)
	chain := &BlockChain {
		DB: db,
	}
	genesis := &Block{
		PrevHash: []byte(GENESIS_BLOCK),
		Mapping: make(map[string]uint64),
		Miner: receiver,
		TimeStamp: time.Now().Format(time.RFC3339),
	}

	genesis.Mapping[STORAGE_CHAIN] = STORAGE_VALUE
	genesis.Mapping[receiver] = GENESIS_REWARD
	genesis.CurrHash = genesis.hash()
	chain.AddBlock(genesis) 

	return nil
}