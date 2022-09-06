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

	var block entities.Block
	database.Instance.Preload(clause.Associations).Preload("Transactions." + clause.Associations).Last(&block)

	transactions := transactions.NewTransaction(inputData.From, inputData.To, inputData.Amount)
	transactions.TransactionHash = block.Hash[:]

	database.Instance.Create(&transactions)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
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
		Hash:     []byte{},
		PrevHash: []byte(""),
		Nonce:    0,
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
	nonce, hash := pow.Run(firstTransaction)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newBlock)
}

func GetAllData(w http.ResponseWriter, r *http.Request) {

	var block []entities.Block

	//database.Instance.Model(&entities.Block{}).Preload("Transaction").Find(&block)
	database.Instance.Preload(clause.Associations).Preload("Transactions." + clause.Associations).Find(&block)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(block)
}

// func checkIfExistsData(blockPrevHash []byte) bool {
// 	var block entities.Block
// 	//passar como parâmetro o endereço de memória da variável e segundo parametro o ID
// 	//usa GORM para achar o primeiro item
// 	prevHash := database.Instance.Last(&block, blockPrevHash)
// 	//se não tiver dados o GORM retorna o valor 0
// 	if errors.Is(prevHash.Error, gorm.ErrRecordNotFound) {
// 		return false
// 	}
// 	return true
// 	// if product.ID == 0 {
// 	// 	return false
// 	// }
// 	// return true
// }

// func GetPokemonById(w http.ResponseWriter, r *http.Request) {
// 	//utiliza o ID passado na URL host/api/products/{id}
// 	//usa o MUX para extrair o ID da URL recebida no 'r' e associar a variável criada productId
// 	pokemonId := mux.Vars(r)["id"]
// 	if !checkIfExistsData(pokemonId) {
// 		//faz onde, o que
// 		//no w escreve a frase em formato json (?)
// 		json.NewEncoder(w).Encode("Pokemon não encontrado!")
// 		return
// 	}
// 	var pokemon entities.Block
// 	//utiliza o endereço na memória do product para procurar o primeiro item com o productId
// 	database.Instance.Preload(clause.Associations).First(&pokemon, pokemonId)
// 	//database.Instance.First(&pokemon, pokemonId)
// 	w.Header().Set("Content-Type", "application/json")
// 	//codifica product e envia para w
// 	json.NewEncoder(w).Encode(pokemon)
// }

// func GetPokemonByName(w http.ResponseWriter, r *http.Request) {
// 	pokemonName := mux.Vars(r)["name"]
// 	if !checkIfPokemonExists(pokemonName) {
// 		//faz onde, o que
// 		//no w escreve a frase em formato json (?)
// 		json.NewEncoder(w).Encode("Pokemon não encontrado!")
// 		return
// 	}
// }

// func UpdateDataPokemon(w http.ResponseWriter, r *http.Request) {
// 	//utiliza o ID passado na URL host/api/products/{id}
// 	//usa o MUX para extrair o ID da URL recebida no 'r' e associar a variável criada productId
// 	pokemonId := mux.Vars(r)["id"]
// 	if !checkIfPokemonExists(pokemonId) {
// 		//faz onde, o que
// 		//no w escreve a frase em formato json (?)
// 		json.NewEncoder(w).Encode("Pokemon não encontrado!")
// 		return
// 	}
// 	var pokemon entities.NamesDataBase
// 	//pega o primeiro item no banco de dados com o determinado ID
// 	database.Instance.First(&pokemon, pokemonId)
// 	//decodifica o dado recebido em 'r' no tipo product
// 	json.NewDecoder(r.Body).Decode(&pokemon)
// 	//usa o GORM para salvar no banco de dados o tipo decodificado
// 	database.Instance.Save(&pokemon)
// 	w.Header().Set("Content-Type", "application/json")
// 	//codifico o product e envio para 'w'
// 	json.NewEncoder(w).Encode(pokemon)
// }

// func DeleteDataPokemon(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	//utiliza o ID passado na URL host/api/products/{id}
// 	//usa o MUX para extrair o ID da URL recebida no 'r' e associar a variável criada productId
// 	pokemonId := mux.Vars(r)["id"]
// 	if !checkIfPokemonExists(pokemonId) {
// 		w.WriteHeader(http.StatusNotFound)
// 		//faz onde, o que
// 		//no w escreve a frase em formato json (?)
// 		json.NewEncoder(w).Encode("Pokemon não encontrado!")
// 		return
// 	}
// 	var pokemon entities.NamesDataBase
// 	//GORM acessa o banco de dados e deleta o product
// 	database.Instance.Delete(&pokemon, pokemonId)
// 	json.NewEncoder(w).Encode("Pokemon deletado com sucesso!")
// }

// //será chamado de CreateBlockchain
// //será necessário separar as funções
// func Teste(w http.ResponseWriter, r *http.Request) {
// 	//define o estilo que os dados serão mostrados em w
// 	//w.Header().Set("Content-Type", "application/json")
// 	//cria uma variável do tipo Product do pacote entities
// 	var block entities.JsonBlock
// 	json.NewDecoder(r.Body).Decode(&block)
// 	//ver se os dados que estão no firstBlock impacta no find se não usar ele em vez da var oquesera
// 	//com * retorna nil

// 	// var result int64
// 	// database.Instance.Table("blocks").Count(&result)

// 	// database.Instance.Table("blocks").Last()

// 	var lastBlock entities.Block

// 	newBlock := &entities.Block{
// 		Hash:         []byte{},
// 		Transactions: []entities.Transaction{},
// 		PrevHash:     []byte(""),
// 		Nonce:        0,
// 	}
// 	//VERIFICAR SE PRECISA ESPECIFICAR OU QUAL ELE CONSIDERA O ULTIMO SE ELE DA UM ID AUTOMATICO
// 	//com um mesmo data ele puxa
// 	//tentar incluir um time pra puxar o ultimo com base no time
// 	//"CreatedAt": "2022-08-24T17:13:58.611-03:00", igual gom.Model
// 	//usar len para pegar o ultimo item
// 	blockIsData := database.Instance.Last(&lastBlock)

// 	//para quando não existe dados no banco de dados
// 	if errors.Is(blockIsData.Error, gorm.ErrRecordNotFound) {

// 		genesisData := "First Transaction from Genesis"

// 		firstTransaction := transactions.CoinBaseTx(block.Address, genesisData)
// 		//da pra usar for range no lugar de append
// 		newBlock.Transactions = append(newBlock.Transactions, *firstTransaction)
// 		pow := blockchain.NewProof(newBlock)
// 		nonce, hash := pow.Run()
// 		newBlock.Hash = hash[:]
// 		newBlock.Nonce = nonce

// 		database.Instance.Create(&newBlock)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(newBlock)

// 		return
// 	}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(newBlock)

// 		return
// 	}
// 	transactions := transations.CoinBaseTx(block.Address, lock.InitialData)
// 	newBlock.Transactions = append(newBlock.Transactions,*transactions)

// 	pow := blockchain.NewProof(newBlock)
// 	nonce, hash := pow.Run()
// 	newBlock.Hash = hash[:]
// 	newBlock.Nonce = nonce

// 	database.Instance.Creae(&newBlock)

// 	w.Header().Set("Content-Type", "appication/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(neBlock)
// }
