package transactions

import (
	"encoding/hex"

	"github.com/maickmachado/blockchain-go/database"
	"github.com/maickmachado/blockchain-go/entities"
)

// type BlockChain struct {
// 	LastHash []byte
// 	// Database *badger.DB
// }

// type BlockChainIterator struct {
// 	CurrentHash []byte
// 	//Database    *badger.DB
// }

//transações que possuem output que não foram referenciados por outros inputs
//os inputs como informado tem o ID que faz referencia a outputs
//ou seja, outpus pode se considerado o que eu tenho disponivel para gastar
//e inputs é o que eu gastei
//logo o que não está associado a um input é o que eu tenho de unspent transactions
//revisar o fato de ser um método ded BlockChain - não vejo necessidade inicialmente
//ver como passar esse adress - ver onde essa função ta sendo usada
func FindUnspentTransactions(address string) []*entities.Transaction {
	var unspentTxs []*entities.Transaction

	var block entities.Block
	//errado não irá existir o campo address no block
	database.Instance.First(&block, address)

	spentTXOs := make(map[string][]int)

	for {
		//pega o hash antigo e coloca no atual
		//objetivo final é colocar os dados do ultimo block do banco de dados
		//vai passar por cada transação dentro do block selecionado
		for _, tx := range block.Transactions {
			//vai pegar cada ID de cada transação
			//por que eu fiz isso?
			//trasn formar o slice if bytes em string
			txID := hex.EncodeToString(tx.TransactionsRefe)
			//o label Outputs serve para o break somente sair de dentro dele e não de dentro dos for loops superiores
		Outputs:
			//para cada transação um no for loop para interagir com as outputs
			for outIdx, out := range tx.Outputs {
				//outIdx é o index varia de 0 ate x output
				//out contem os valores de tx.Outputs = (value e pubkey)
				//acessa o []int relativo ao spentTXOs[txID]
				//serve para ver se o output esta dentro do map criado
				//se estiver dentro faremos um novo for com interações dentro do map
				//populando o map com os ID das transações
				//não tem valor associado a chave, ver como vai funcionar
				//como não tem nada no map o primerio loop sempre vai se nil
				if spentTXOs[txID] != nil {
					//range chave, valor
					//txID é a chave - spentOut é o valor da chave
					//spentTXOs[txID] = []int{1, 2}
					//spentOut é um int
					//spentOut será o 1 e o 2 por exemplo
					//spentOut é o valor adicionado pelo append
					//é o valor Out do TxInput que é o index onde o TxOutput aparece na transação
					//significa que ela já foi "gasta"
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				//out contem os valores de tx.Outputs = (value e pubkey)
				//se o adress enviado como parametro for igual a pubkey retorn true
				//o que isso significa?
				//verificar se significa que como o pubkey ainda é o nome do usuario não foi enviada pra ninguem
				//ex.: como maick tem 20 reais que ainda estão no nome de maick ele pode gastar
				//verificar o porque, e se quando jogar todas as tx no unspentTxs os inputs tambem vao e se existe algum
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			//se não for uma CoinBaseTx ele verifica se nos Inputs existe transação com o address informado
			if !tx.IsCoinBase() {
				//in contem os valores de tx.Inputs = (ID, Out e Sig)
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.TxInputRefe)
						//significa que no ID tal o valor do output foi gasto
						//com isso é adicionado no map com key = ID o valor do index da transação
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
		//para caso for o primeiro bloco
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func FindUTXO(address string) []entities.TxOutput {
	var UTXOs []entities.TxOutput
	unspentTransactions := FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, *out)
			}
		}

	}
	return UTXOs
}

//pegas as transações não gastas e verifica se tem saldo
//address que queremos verificar
//amount que queremos enviar
func FindSpendebleOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.TransactionsRefe)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
