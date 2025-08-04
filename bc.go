package main

import (
	"fmt"

	bc "github.com/siyovush-hamidov/Blockchain/blockchain"
)

const (
	DBNAME = "blockchain.db"
)

func main() {
	miner := bc.NewUser()
	bc.NewChain(DBNAME, miner.Address())
	chain := bc.LoadChain(DBNAME)
	for i := 0; i < 3; i++ {
		block := bc.NewBlock(miner.Address(), chain.LastHash())
		block.AddTransaction(chain, bc.NewTransaction(miner, chain.LastHash(), "aaa", 5))
		block.AddTransaction(chain, bc.NewTransaction(miner, chain.LastHash(), "bbb", 3))
		err := block.Accept(chain, miner, make(chan bool))
		if err != nil {
			fmt.Println("Ошибка при принятии блока:", err)
			continue
		}
		chain.AddBlock(block)
	}
	var sblock string
	rows, err := chain.DB.Query("SELECT Block FROM BlockChain")
	if err != nil {
		panic("ERROR: Query to database unsuccessful")
	}
	for rows.Next() {
		rows.Scan(&sblock)
		fmt.Println(sblock)
	}
}