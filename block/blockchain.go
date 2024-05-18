package block

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lucasacoutinho/astroboi/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "CHAIN"
	MINING_REWARD     = 1.0
)

type Transaction struct {
	sender   string
	receiver string
	value    float32
}

func NewTransaction(sender, receiver string, value float32) *Transaction {
	return &Transaction{
		sender:   sender,
		receiver: receiver,
		value:    value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender   %s\n", t.sender)
	fmt.Printf("receiver %s\n", t.receiver)
	fmt.Printf("value    %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender"`
		Receiver string  `json:"receiver"`
		Value    float32 `json:"value"`
	}{
		Sender:   t.sender,
		Receiver: t.receiver,
		Value:    t.value,
	})
}

func (t *Transaction) Hash() [32]byte {
	m, _ := json.Marshal(t)
	return sha256.Sum256([]byte(m))
}

type Block struct {
	nonce        int
	prevHash     [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(nonce int, prevHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		nonce:        nonce,
		prevHash:     prevHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PrevHash     string         `json:"prev_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PrevHash:     fmt.Sprintf("%x", b.prevHash),
		Timestamp:    b.timestamp,
		Transactions: b.transactions,
	})
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) Print() {
	fmt.Printf("Nonce: 			%d\n", b.nonce)
	fmt.Printf("PrevHash: 		%x\n", b.prevHash)
	fmt.Printf("Timestamp: 		%d\n", b.timestamp)

	for _, t := range b.transactions {
		t.Print()
	}
}

type BlockChain struct {
	pool    []*Transaction
	chain   []*Block
	address string
	port    uint16
}

func NewBlockChain(address string, port uint16) *BlockChain {
	block := &Block{}

	chain := &BlockChain{address: address, port: port}
	chain.CreateBlock(0, block.Hash())
	return chain
}

func (bc *BlockChain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *BlockChain) CreateBlock(nonce int, prevHash [32]byte) *Block {
	block := NewBlock(nonce, prevHash, bc.pool)
	bc.chain = append(bc.chain, block)
	bc.pool = []*Transaction{}
	return block
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddTransaction(
	sender, receiver string,
	value float32,
	senderPublicKey *ecdsa.PublicKey,
	s *utils.Signature) bool {
	t := NewTransaction(sender, receiver, value)

	if sender == MINING_SENDER {
		bc.pool = append(bc.pool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		// if bc.CalculateTotalAmount(sender) < value {
		// 	return false
		// }

		bc.pool = append(bc.pool, t)
		return true
	}

	return false
}

func (bc *BlockChain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *BlockChain) CopyPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.pool {
		transactions = append(transactions, NewTransaction(t.sender, t.receiver, t.value))
	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, prevHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guess := Block{
		nonce:        nonce,
		prevHash:     prevHash,
		timestamp:    0,
		transactions: transactions,
	}

	guessHash := fmt.Sprintf("%x", guess.Hash())

	return guessHash[:difficulty] == zeros
}

func (bc *BlockChain) ProofOfWork() int {
	transactions := bc.CopyPool()
	prevHash := bc.LastBlock().Hash()
	nonce := 0

	for !bc.ValidProof(nonce, prevHash, transactions, 3) {
		nonce += 1
	}

	return nonce
}

func (bc *BlockChain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.address, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	prevHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, prevHash)
	log.Println("action=mining, status=success")

	return true
}

func (bc *BlockChain) CalculateTotalAmount(address string) float32 {
	var total float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if address == t.receiver {
				total += value
			}
			if address == t.sender {
				total -= value
			}
		}
	}

	return total
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}
