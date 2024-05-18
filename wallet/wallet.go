package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/lucasacoutinho/astroboi/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func NewWallet() *Wallet {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	h2 := sha256.New()
	h2.Write(publicKey.X.Bytes())
	h2.Write(publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	chsum := digest6[:4]

	dc8 := make([]byte, 28)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	address := base58.Encode(dc8)

	return &Wallet{
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
	}
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) Address() string {
	return w.address
}

type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	receiverAddress  string
	value            float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender, receiver string, value float32) *Transaction {
	return &Transaction{
		senderPrivateKey: privateKey,
		senderPublicKey:  publicKey,
		senderAddress:    sender,
		receiverAddress:  receiver,
		value:            value,
	}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender"`
		Receiver string  `json:"receiver"`
		Value    float32 `json:"value"`
	}{
		Sender:   t.senderAddress,
		Receiver: t.receiverAddress,
		Value:    t.value,
	})
}
