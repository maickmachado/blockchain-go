package transactions

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/maickmachado/blockchain-go/blockchain/entities"
)

//primeira transação do block
//temos somente um input e somente um output
//contem um reward - verificar se o 100 do txout é como se fosse uma mineração
func CoinBaseTx(to, data string) *entities.Transaction {
	//verificar o porque desse if
	//talvez seja pra criação de blocks
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	//input não faz referência a output pois não existe outputs
	//como o ID não faz referencia a nenhum output, referencia a apenas um output vazio
	txin := entities.TxInput{
		TxInputRefe: []byte{},
		//-1 pois não faz referencia a nenhum output (entender isso aqui pq -1)
		//se o out é um indexz o -1 ta fora do range de index 0...n
		Out: -1,
		Sig: data,
	}

	txout := entities.TxOutput{
		//reward
		Value:  100,
		PubKey: to,
	}

	tx := entities.Transaction{
		//TransactionsRefe: nil,
		Inputs:  []*entities.TxInput{&txin},
		Outputs: []*entities.TxOutput{&txout},
	}
	tx.SetID()
	tx.SetInputID()
	tx.SetOutputID()

	// txf := entities.Transaction{
	// 	ID:      tx.ID,
	// 	Inputs:  []entities.TxInput{txin},
	// 	Outputs: []entities.TxOutput{txout},
	// }
	return &tx
}

func NewTransaction(from, to string, amount int) *entities.Transaction {
	var inputs []*entities.TxInput
	var outputs []*entities.TxOutput

	acc, validOutputs := FindSpendebleOutputs(from, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			//cria os inputs para os fundos outputs que serão usado
			input := entities.TxInput{
				TxInputRefe: txID,
				Out:         out,
				Sig:         from,
			}
			inputs = append(inputs, &input)
		}
	}
	//primeira transação que mostra que o usuário esta mandando amount x
	outputs = append(outputs, &entities.TxOutput{
		Value:  amount,
		PubKey: to,
	})
	//segunda transação é criada para mostrar o tanto que sobrou na conta do usuário
	if acc > amount {
		outputs = append(outputs, &entities.TxOutput{
			Value:  acc - amount,
			PubKey: from,
		})
	}

	tx := entities.Transaction{
		TransactionsRefe: nil,
		Inputs:           inputs,
		Outputs:          outputs,
	}
	tx.SetID()
	tx.SetInputID()
	tx.SetOutputID()

	return &tx
}
