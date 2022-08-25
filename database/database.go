package database

import (
	"log"

	"github.com/maickmachado/blockchain-go/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//definição de uma variável para o banco de dados e uma para o erro
//o fato de ser criada fora de uma função pode ser acessada por diversas
var Instance *gorm.DB
var err error

//função que faz conecção com o banco de dados utilizando o GORM
//após acessado, a variável Instance vai conseguir executar operações no banco de dados
func Connect(connectionString string) {
	Instance, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

//função que assegura que a tabela do banco de dados será igual ao struct criado
func Migrate() {
	blockDataBase := &entities.Block{}
	Instance.AutoMigrate(blockDataBase)
	log.Println("Database Migration Completed...")
}
