package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/maickmachado/blockchain-go/entities"
)

// take the data from block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

//requirements:
// the first few bytes must contaiun 0s

//a medida que a quantidade de miners for aumentando a tendencia é aumentar a dificuldade
const Difficulty = 18

type ProofOfWork struct {
	Block *entities.Block
	//representa o requirements que é derivado do Difficulty
	Target *big.Int
}

//popula o proof of work struct com o block recebido e o target criado
func NewProof(b *entities.Block) *ProofOfWork {
	target := big.NewInt(1) //1
	// pega o numero 256 que é o numero de bytes no hash
	// 32
	// vamos usar o target para trocar o número de bytes pelo número uint(256 - Difficulty)
	// Lsh = left shift
	//So n << x is "n times 2, x times". And y >> z is "y divided by 2, z times".
	//(x << n == x*2^n )
	//For example, 1 << 5 is "1 times 2, 5 times" or 32.
	//And 32 >> 5 is "32 divided by 2, 5 times" or 1.
	//Given integer operands a and n,
	//a << n; shifts all bits in a to the left n times
	//a >> n; shifts all bits in a to the right n times
	//uint(256-Difficulty) = 2^238 = 441711766194596082395824375185729628956870974218904739530401550323154944
	target.Lsh(target, uint(256-Difficulty))
	//coloca o Block recebido no parametro e o target modificado
	pow := &ProofOfWork{b, target}
	return pow
}

//nonce é um counter
func (pow *ProofOfWork) InitData(nonce int) []byte {
	//concatena as informações contidas no Data e PrevHash e usa um []byte{} como separador
	//[][]byte{} contem vários slice of bytes, no caso abaixo o [][]byte é formado pelos slices of bytes do struct Block
	//b.Data e b.PrevHash - ambos []bytes
	//[][]byte{[]byte, []byte} = [][]byte{b.Data, b.PrevHash}
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

//fara a combinação do counter (nonce) com o dados (data) e com a dificuldade
//criará um novo bytes buffer
//Um buffer é um espaço de memória (tipicamente RAM) que armazena dados binários
//mesma coisa o decode do pacote json, só que em vez de escrever na saída 'w' ele aloca na memória
//porém aloca ele na memória codificado (encode)
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	//pega um numero e decodifica em bytes
	//binary.BigEndian indica como queremos que nossos bytes ficam organizados
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//faz um loop incrementando o hash até chega no requerimento
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	//irá preparar nossos dados (data) e transformar num hash sha256
	//depois converter esse hash em um big integer
	//depois comparar esse big integer com o big integer que está dentro do proof of work struct
	for nonce < math.MaxInt64 {
		//retorna um slice of bytes com tudo concatenado usando o nonce indicado
		data := pow.InitData(nonce)
		//pego do data e transformo em um hash do tipo sha256
		//32 "pedaços" - 32 * 8 bytes = 256
		//a função Sum256 faz o calculo do hash
		hash = sha256.Sum256(data)
		//%x	hexadecimal notation - base 16, with lower-case letters for a-f
		fmt.Printf("\r%x", hash)
		//converte o hash em um big integer
		intHash.SetBytes(hash[:])
		//compara o proof of work target com o novo big int criado o intHash
		//a comparação é pra ver se chegou no requerimento
		if intHash.Cmp(pow.Target) == -1 {
			//break porque o hash passou o valor do requisito
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

//após rodar a função do proof of work, Run
//obteremos o nonce que permiterá obter o correto hash de acordo com o target
//após isso iremos rodar novamente para provar que o hash é válido
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
