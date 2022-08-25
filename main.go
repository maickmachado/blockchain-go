package main

import (
	"github.com/maickmachado/blockchain-go/database"
	"github.com/maickmachado/blockchain-go/routes"
)

func main() {
	database.Connect("root:abc@tcp(127.0.0.1:3306)/bitcoin")
	database.Migrate()
	routes.HandleRequest()
	//criação do primeiro block
	//firstBlock := blockchain.CreateBlock("Genesis", []byte{})
	//colocando o primeiro block no BlockChain
	//chain := &blockchain.BlockChain{Blocks: []*blockchain.Block{firstBlock}}
	// chain := InitBlockChain()

	// chain.AddBlock("First block after Genesis")
	// chain.AddBlock("Second block after Genesis")
	// chain.AddBlock("Third block after Genesis")

	// for _, value := range chain.Blocks {
	// 	//fmt.Printf("Previous Hash: %x\n", value.PrevHash)
	// 	fmt.Printf("Data in block: %s\n", value.Data)
	// 	fmt.Printf("Hash: %x\n", value.Hash)

	// 	pow := blockchain.NewProof(value)
	// 	//converte o bool em um tipo string
	// 	fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
	// 	fmt.Println()
	// }
}
