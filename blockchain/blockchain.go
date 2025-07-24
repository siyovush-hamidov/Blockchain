package blockchain

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"time"
)

const (
	CREATE_TABLE = `
	CREATE TABLE BlockChain(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Hash varchar(44) UNIQUE,
		Block TEXT
	);
	`
)

const (
	GENESIS_BLOCK = "GENESIS BLOCK"
	STORAGE_VALUE = 100
	GENESIS_REWARD = 100
	STORAGE_CHAIN = "STORAGE_CHAIN"
)

const (
	DIFFICULTY = 20
)

const (
	RAND_BYTES = 32
	START_PERCENT = 10
	STORAGE_REWARD = 1
)

const (
	TXS_LIMIT = 2
)

type BlockChain struct {
	DB *sql.DB
}

type Block struct {
	Nonce uint64
	Difficulty uint8
	CurrHash []byte
	PrevHash []byte
	Transactions []Transaction
	Mapping map[string]uint64 // Состояния балансов пользователей
	Miner string
	Signature []byte
	TimeStamp string
}

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

type User struct {
	PrivateKey *rsa.PrivateKey
}

func (chain *BlockChain) AddBlock(block *Block) {
	chain.DB.Exec("INSERT INTO BlockChain (Hash, Block) VALUES ($1, $2)",
		Base64Encode(block.CurrHash),
		SerializeBlock(Block),
	)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func SerializeBlock(block *Block) string {
	jsonData, err := json.MarshalIndent(*block, "", "\t")
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func LoadChain(filename string) *BlockChain {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil
	}
	chain := &BlockChain{
		DB: db,
	}
	return chain
}

func NewBlock(miner string, prevHash []byte) *Block{
	return &Block{
		Difficulty: DIFFICULTY,
		PrevHash: prevHash,
		Miner: miner,
		Mapping: make(map[string]uint64),
	}
}

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

func NewTransaction(user *User, lashHash []byte, to string, value *uint64) *Transaction{
	tx := &Transaction{
		RandBytes: GenerateRandomBytes(RAND_BYTES),
		PrevBlock: lastHash,
		Sender: user.Address(),
		Receiver: to,
		Value: *value,
	}
	if value > START_PERCENT {
		tx.ToStorage = STORAGE_REWARD
	}
	tx.CurrHash = tx.hash()
	tx.Signature = tx.sign(user.Private())
	return tx
}

func GenerateRandomBytes(max uint) []byte {
	var slice []byte = make([]byte, max)
	_, err := rand.Read(slice)
	if err != nil {
		return nil
	}
	return slice
}

func (user *User) Address() string {
	return StringPublic(user.Public())
}

func (user *User) Private() *rsa.PrivateKey {
	return user.PrivateKey
}

func (tx *Transaction) hash() []byte {
	return HashSum(bytes.Join(
		[][]byte{
			tx.RandBytes,
			tx.PrevBlock,
			[]byte(tx.Sender),
			[]byte(tx.Receiver),
			ToBytes(tx.Value),
			ToBytes(tx.ToStorage),
		},
		[]byte{},
	))
}

func (tx *Transaction) sign(priv *rsa.PrivateKey) []byte{
	return Sign(priv, tx.CurrHash)
}

func StringPublic(pub *rsa.PublicKey) string{
	return Base64Encode(x509.MarshalPKCS1PublicKey(pub))
}

func (user *User) Public() *rsa.PublicKey {
	return &(user.PrivateKey).PublicKey
}

func HashSum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func ToBytes(num uint64) []byte {
	var data = new(bytes.Buffer)
	err := binary.Write(data, binary.BigEndian, num)
	// Q: Как работает BigEndian и LittleEndian?
	if err != nil {
		return nil
	}
	return data.Bytes()
}

func Sign(priv *rsa.PrivateKey, data []byte) []byte {
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, data, nil)
	// Q: Что такое digest?
	if err != nil {
		return nil
	}
	return signature
}

func (block *Block) AddTransaction(chain *BlockChain, tx *Transaction) error {
	// Q: ключевое слово error
	if tx == nil {
		return errors.New("The transaction is null")
		// Q: errors.New("Some text")
	}
	if tx.Value == 0 {
		return errors.New("The transaction value is 0")
	}
	if tx.Sender != STORAGE_CHAIN && len(block.Transactions) == TXS_LIMIT {
		return errors.New("The transaction length has reached its limit")
	}
	if tx.Sender != STORAGE_CHAIN && tx.Value > START_PERCENT && tx.ToStorage != STORAGE_REWARD {
		return errors.New("Storage reward pass")
	}
	if !bytes.Equal(tx.PrevBlock, chain.LastHash()) {
		return errors.New("Previous block in the transaction is not the last hash in chain")
	}
	var balanceInChain uint64
	balanceInTX := tx.Value + tx.ToStorage
	if value, ok := block.Mapping[tx.Sender]; ok{
		// Q: Что это за запись такая? Что означает ok {} ?
		balanceInChain = value
	} else {
		balanceInChain = chain.Balance(tx.Sender, chain.Size())
	}
	if balanceInTX > balanceInChain {
		return errors.New("Insufficient funds")
	}
	block.Mapping[tx.Sender] = balanceInChain - balanceInTX
	block.addBalance(chain, tx.Receiver, tx.Value)
	block.addBalance(chain, STORAGE_CHAIN, tx.ToStorage)
	block.Transactions = append(block.Transactions, *tx)
	return nil
}

func (chain *BlockChain) Balance(address string, size uint64) uint64 {
	var (
		sblock string
		block *Block
		balance uint64
	)
	rows, err := chain.DB.Query("SELECT Block FROM BlockChain WHERE Id <= $1 ORDER BY ID DESC", size)
	if err != nil {
		return balance
	}
	defer rows.Close()
	for rows.Next() {
		// Q: Что за конструкция?
		rows.Scan(&sblock)
		block = DeserializeBlock(sblock)
		if value, ok := block.Mapping[address]; ok {
			balance = value
			break
		}
	}
	return balance
}

func (block *Block) addBalance(chain *BlockChain, receiver string, value uint64) {
	var balanceInChain uint64
	if v, ok := block.Mapping[receiver]; ok {
		balanceInChain = v
	} else {
		balanceInChain = chain.Balance(receiver, chain.Size())
	}
	block.Mapping[receiver] = balanceInChain + value
}

// Возвращает количество блоков в локальной БД
func (chain *BlockChain) Size() uint64 {
	var size uint64
	row := chain.DB.QueryRow("SELECT ID FROM BlockChain ORDER BY Id DESC")
	row.Scan(&size)
	return size
}

func DeserializeBlock(data string) *Block {
	var block Block
	err := json.Unmarshal([]byte(data), &block)
	if err != nil {
		return nil
	}
	return &block
}

