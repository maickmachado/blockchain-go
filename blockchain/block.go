package blockchain

import "github.com/maickmachado/blockchain-go/entities"

type BlockChain struct {
	Blocks []*entities.Block
}

// type Block struct {
// 	//parametros que são informados
// 	Data     []byte
// 	PrevHash []byte
// 	//parametros que são obtidos
// 	Hash  []byte
// 	Nonce int
// }

//recebe alguns dados o hash anterior e retorna um pointer para o struct Block
func CreateBlock(data string, prevHash []byte) *entities.Block {
	//um novo block é criado usando o struct Block
	//block := constructor{}

	//para o campo Hash é passado um []byte vazio
	//para o campo Data nós pegamos o argumento string data e convertemos para um []byte
	//para o campo PrevHash nós colocamos o prevHash que é passado como argumento
	block := &entities.Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}

	//popula o proof of work struct com o block enviado no parametro e o target criado na função NewProof
	pow := NewProof(block)
	//executa a função Run no proof_of_work struct criado
	//retorna o hash e o nonce correto após as verificações dos requisitos
	nonce, hash := pow.Run()
	//coloca os campos obtidos no block struct
	block.Hash = hash[:]
	block.Nonce = nonce

	//com o hash criado pelo metodo DeriveHash nós retornamos um block
	return block
}

//método que permite adicionar um block ao blockchain
func (chain *BlockChain) AddBlock(data string) {
	//crio uma variável para associar ultimo block criado
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	//para criar um novo block a que será adicionado ao BlockChain
	//utilizamos a função CreateBlock e passamos o parametros data
	//e o Hash do block anterior obtido através do acesso a variável criada prevBlock
	new := CreateBlock(data, prevBlock.Hash)
	//adiciona o novo block atraves do append ao BlockChain
	chain.Blocks = append(chain.Blocks, new)
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
