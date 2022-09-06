package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/iotexproject/go-pkgs/hash"
)

const (
	chesumLength = 4
	version      = byte(0x00)
)

type Wallet struct {
	//elliptical curve
	//tipo de criptografia baseada em um grafico elliptical curve
	//https://www.allaboutcircuits.com/technical-articles/elliptic-curve-cryptography-in-embedded-systems/
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	fmt.Printf("pub key: %x\n", w.PublicKey)
	fmt.Printf("pub hash: %x\n", pubHash)
	fmt.Printf("address: %x\n", address)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	//outputda curva será de 256 bytes
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//converte o valor de X e Y em byte
	//... significa que pega cada item no []byte Y e concatena em X
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := hash.Hash160b(pubHash[:])

	publicRipMD2 := []byte(hasher[:])

	// hasher := ripemd160.New()
	// _, err := hasher.Write(pubHash[:])
	// if err != nil {
	// 	log.Panic(err)
	// }

	// publicRipMD2 := hasher.Sum(nil)

	return publicRipMD2
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:chesumLength]
}
