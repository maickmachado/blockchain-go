package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"

	"github.com/maickmachado/blockchain-go/blockchain/entities"
	"github.com/maickmachado/blockchain-go/blockchain/proof"
	"github.com/maickmachado/blockchain-go/blockchain/transactions"
	"github.com/maickmachado/blockchain-go/database"
)

// type BlockChain struct {
// 	Blocks []*entities.Block
// }

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	var inputData entities.InputTransaction
	json.NewDecoder(r.Body).Decode(&inputData)

	var block []entities.Block
	database.Instance.Preload(clause.Associations).Preload("Transactions." + clause.Associations).Find(&block)

	lastCounter := len(block)

	fmt.Println("last block in new transaction:", block)

	newBlock := entities.Block{
		CounterBlock: lastCounter + 1,
		Hash:         []byte{},
		PrevHash:     block[lastCounter-1].Hash,
		Nonce:        0,
	}

	pow := proof.NewProof(&newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	transactions := transactions.NewTransaction(inputData.From, inputData.To, inputData.Amount)
	//ok
	fmt.Println("new transactions in new transaction:", transactions)

	database.Instance.Create(&newBlock)

	var transactionDatabase entities.Transaction
	transactionDatabase.TransactionsRefe = transactions.TransactionsRefe
	transactionDatabase.TransactionHash = hash[:]

	database.Instance.Create(&transactionDatabase)

	//criação do banco de dados output
	var outputDatabase entities.TxOutput
	outputDatabase.TxOutputID = transactionDatabase.TransactionsRefe
	for _, value := range transactions.Outputs {
		outputDatabase.Value = value.Value
		outputDatabase.PubKey = value.PubKey
		database.Instance.Create(&outputDatabase)
	}
	//criação do banco de dados input
	var inputDatabase entities.TxInput
	inputDatabase.TxInputID = transactionDatabase.TransactionsRefe
	for _, value := range transactions.Inputs {
		inputDatabase.TxInputRefe = value.TxInputRefe
		inputDatabase.Out = value.Out
		inputDatabase.Sig = value.Sig
		database.Instance.Create(&inputDatabase)
	}

	fmt.Println("newblock:", newBlock)
	fmt.Println("output:", outputDatabase)
	fmt.Println("input:", inputDatabase)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactionDatabase)
}

func GetBalanceAddress(w http.ResponseWriter, r *http.Request) {
	address := mux.Vars(r)["address"]

	balance := 0
	UTXOs := transactions.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(fmt.Sprintf("Balance of %s: %d\n", address, balance))
}

func CreateGenesisBlock(w http.ResponseWriter, r *http.Request) {
	//define o estilo que os dados serão mostrados em w
	//w.Header().Set("Content-Type", "application/json")
	//cria uma variável do tipo Product do pacote entities
	var inputData entities.GenesisData
	json.NewDecoder(r.Body).Decode(&inputData)
	//ver se os dados que estão no firstBlock impacta no find se não usar ele em vez da var oquesera
	//com * retorna nil
	// var result int64
	// database.Instance.Table("blocks").Count(&result)
	// database.Instance.Table("blocks").Last()
	newBlock := entities.Block{
		CounterBlock: 1,
		Hash:         []byte{},
		PrevHash:     []byte(""),
		Nonce:        0,
	}
	//VERIFICAR SE PRECISA ESPECIFICAR OU QUAL ELE CONSIDERA O ULTIMO SE ELE DA UM ID AUTOMATICO
	//com um mesmo data ele puxa
	//tentar incluir um time pra puxar o ultimo com base no time
	//"CreatedAt": "2022-08-24T17:13:58.611-03:00", igual gom.Model
	//usar len para pegar o ultimo item
	genesisData := "First Transaction from Genesis"

	firstTransaction := transactions.CoinBaseTx(inputData.Address, genesisData)
	//da pra usar for range no lugar de append
	//newBlock.Transactions = append(newBlock.Transactions, firstTransaction)

	pow := proof.NewProof(&newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	database.Instance.Create(&newBlock)

	var transactionDatabase entities.Transaction
	transactionDatabase.TransactionsRefe = firstTransaction.TransactionsRefe

	transactionDatabase.TransactionHash = hash[:]
	// fmt.Println(transactionDatabase.TransactionHash)
	database.Instance.Create(&transactionDatabase)

	//criação do banco de dados output
	var outputDatabase entities.TxOutput
	outputDatabase.TxOutputID = transactionDatabase.TransactionsRefe
	for _, value := range firstTransaction.Outputs {
		outputDatabase.Value = value.Value
		outputDatabase.PubKey = value.PubKey
		database.Instance.Create(&outputDatabase)
	}
	//criação do banco de dados input
	var inputDatabase entities.TxInput
	inputDatabase.TxInputID = transactionDatabase.TransactionsRefe
	for _, value := range firstTransaction.Inputs {
		inputDatabase.TxInputRefe = value.TxInputRefe
		inputDatabase.Out = value.Out
		inputDatabase.Sig = value.Sig
		database.Instance.Create(&inputDatabase)
	}

	// database.Instance.Create(&firstTransaction.Outputs)
	// database.Instance.Create(&firstTransaction.Inputs)
	fmt.Println("newblock:", newBlock)
	fmt.Println("output:", outputDatabase)
	fmt.Println("input:", inputDatabase)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newBlock)
}

func GetAllData(w http.ResponseWriter, r *http.Request) {

	var block []entities.Block

	//database.Instance.Model(&entities.Block{}).Preload("Transaction").Find(&block)
	database.Instance.Preload(clause.Associations).Preload("Transactions." + clause.Associations).Find(&block)
	fmt.Println(block)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(block)
}
