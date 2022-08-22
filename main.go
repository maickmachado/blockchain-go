package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

//método que permite criar um novo HASH baseado no HASH anterior e no dados (DATA).
func (b *Block) DeriveHash() {
	//concatena as informações contidas no Data e PrevHash e usa o []byte como separador
	//[][]byte contem vários slice of bytes, no caso abaixo o [][]byte é formado pelos slices of bytes do struct Block
	//b.Data e b.PrevHash - ambos []bytes
	//[][]byte{[]byte, []byte} = [][]byte{b.Data, b.PrevHash}
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	//o atual hash é criado usando a bliclioteca sha256
	hash := sha256.Sum256(info)
	//com o hash criado é possível coloca-lo no struct
	//seleciona todos os itens do hash e coloca no Hash do struct
	b.Hash = hash[:]
}

//recebe alguns dados o hash anterior e retorna um pointer para o struct Block
func CreateBlock(data string, prevHash []byte) *Block {
	//um novo block é criado usando o struct Block
	//block := constructor{}
	//para o campo Hash é passado um []byte vazio
	//para o campo Data nós pegamos o argumento string data e convertemos para um []byte
	//para o campo PrevHash nós colocamos o prevHash que é passado como argumento
	block := &Block{[]byte{}, []byte(data), prevHash}
	//pegamos nosso objeto criado block e acessamos o método DeriveHash
	block.DeriveHash()
	//com o hash criado pelo metodo DeriveHash nós retornamos um block
	return block
}

//método que permite adicionar um block ao blockchain
func (chain *BlockChain) AddBlock(data string) {
	//crio uma variável para associar ultimo block criado
	prevBlock := chain.blocks[len(chain.blocks)-1]
	//para criar um novo block a que será adicionado ao BlockChain
	//utilizamos a função CreateBlock e passamos o parametros data
	//e o Hash do block anterior obtido através do acesso a variável criada prevBlock
	new := CreateBlock(data, prevBlock.Hash)
	//adiciona o novo block atraves do append ao BlockChain
	chain.blocks = append(chain.blocks, new)
}

//função responsável pela criação do primeiro block
// func Genesis() *Block {
// 	//o {} do []byte é para informar que está vazio
// 	//é diferente de simplesmente passar []byte
// 	return CreateBlock("Genesis", []byte{})
// }

// func InitBlockChain() *BlockChain {
// 	return &BlockChain{[]*Block{Genesis()}}
// }

func main() {
	//criação do primeiro block
	firstBlock := CreateBlock("Genesis", []byte{})
	//colocando o primeiro block no BlockChain
	chain := BlockChain{[]*Block{firstBlock}}
	// chain := InitBlockChain()

	chain.AddBlock("First block after Genesis")
	chain.AddBlock("Second block after Genesis")
	chain.AddBlock("Third block after Genesis")

	for _, value := range chain.blocks {
		//fmt.Printf("Previous Hash: %x\n", value.PrevHash)
		fmt.Printf("Data in block: %s\n", value.Data)
		fmt.Printf("Hash: %x\n", value.Hash)
	}

}
