package entities

type JsonBlock struct {
	InitialData string `json:"initial_data"`
}

type Block struct {
	//gorm.Model
	ID       int    `gorm:"primary_key, AUTO_INCREMENT"`
	Data     []byte `json:"data"`
	PrevHash []byte `json:"prev_hash"`
	Hash     []byte `json:"hash"`
	Nonce    int    `json:"nonce"`
}
