package entities

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type GenesisData struct {
	//	InitialData string `json:"initial_data"`
	Address string `json:"address"`
}

// type BlockChain struct {
// 	Blocks []*Block
// }

//cada block pode ter quantas transactions for necessario
//e tem que ter pelo menos uma
type Block struct {
	//testar e ver se o é necessário
	//ID           int            `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Hash         []byte         `json:"hash" gorm:"primaryKey;size:255"`
	Nonce        int            `json:"nonce"`
	Transactions []*Transaction `json:"transactions" gorm:"foreignKey:TransactionHash;references:Hash"`
	PrevHash     []byte         `json:"prev_hash" gorm:"size:255"`
}

type Transaction struct {
	//ID               int         `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	TransactionHash  []byte      `json:"transaction_hash" gorm:"size:255"`
	TransactionsRefe []byte      `json:"transactions_refe" gorm:"primaryKey;size:255"`
	Outputs          []*TxOutput `json:"outputs" gorm:"foreignKey:TxOutputID;references:TransactionsRefe"`
	Inputs           []*TxInput  `json:"inputs" gorm:"foreignKey:TxInputID;references:TransactionsRefe"`
}

//colocar as foreing key nos campos abaixo, talvez sera necessario criar dois database igual no app do pokemon
type TxOutput struct {
	TxOutputID []byte `json:"txoutput_id" gorm:"size:255"`
	//contem a quantidade de 'notas'
	//exemplo 3 notas de 2 reais 5 notas de 50 reais
	//temos que criar outputs para cada 'nota'
	//são indivisíveis, não tem como rasgar a nota de 10 reais para pagar uma conta de 5 reais
	Value int `json:"value"`
	//no bitcoin trata-se de uma criptografia complexa
	//nesse app iremos colocar apenas o nome da conta do usuário
	//necessario para unlock as tokens do value
	PubKey string `json:"pubkey"`
}

//são referências de outputs anteriores
type TxInput struct {
	TxInputID []byte `json:"txinput_id" gorm:"size:255"`
	//referencia a transação que o outputs está contido
	TxInputRefe []byte `json:"txinput_refe" gorm:"size:255"`
	//index de onde o output aparece na transação
	Out int `json:"out"`
	//mesmo propósito do PubKey
	//será o mesmo valor da PubKey
	Sig string `json:"sig"`
}

func (tx *Transaction) IsCoinBase() bool {
	//verifica que a transação é uma coinbasetx
	//len(tx.Inputs) == 1 -- verifica se o comprimento do input é igual a 1 pois o coinbasetx tem apenas 1 input
	//len(tx.Inputs[0].ID) == 0 -- verifica se o ID do primeiro input é igual a 0
	//tx.Inputs[0].Out == -1 -- verifica se o Out do tx input é -1 como foi definido na coinbasetx
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxInputRefe) == 0 && tx.Inputs[0].Out == -1
}

//pelo que entendi quer dizer que o valor associado ao adress foi "utilizado"
func (in *TxInput) CanUnlock(adress string) bool {
	return in.Sig == adress
}

//pelo que entendi quer dizer que o valor associado ao adress não foi "utilizado"
func (out *TxOutput) CanBeUnlocked(adress string) bool {
	return out.PubKey == adress
}

//cria um hash baseado nos bytes que representam a transção
//método de Transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	//usar a biblioteca do JSON para ver se serve
	//json.NewEncoder(w).Encode(newBlock)
	//pelo que entendi a biblioteca gob usa qualquer dado, a JSON somente json
	gob.NewEncoder(&encoded).Encode(tx)
	hash = sha256.Sum256(encoded.Bytes())
	tx.TransactionsRefe = hash[:]
}

func (tx *Transaction) SetInputID() {
	for _, value := range tx.Inputs {
		value.TxInputID = tx.TransactionsRefe
	}
}

func (tx *Transaction) SetOutputID() {
	for _, value := range tx.Outputs {
		value.TxOutputID = tx.TransactionsRefe
	}
}

//hash que representa todas as transações combinadas
func (b *Block) HashTransaction(t *Transaction) []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, value := range b.Transactions {
		txHashes = append(txHashes, value.TransactionsRefe)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
