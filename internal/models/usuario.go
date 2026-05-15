package models

type Usuario struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	Nome            string `json:"nome"`
	Email           string `json:"email" gorm:"unique"`
	Whatsapp        string `json:"whatsapp"`
	PasswordHash    string `json:"-" gorm:"column:password_hash"` // json:"-" esconde o hash
	TipoConta       string `json:"tipo_conta"`                    // 'C' ou 'P'
	NomePropriedade string `json:"nome_propriedade"`
	EnderecoRural   string `json:"endereco_rural"`
}

type UserCredentials struct {
	Nome            string `json:"nome"`
	Email           string `json:"email"`
	Whatsapp        string `json:"whatsapp"`
	Password        string `json:"password"`
	Tipo            string `json:"tipo_conta"`
	NomePropriedade string `json:"nome_propriedade"`
	EnderecoRural   string `json:"endereco_rural"`
}
