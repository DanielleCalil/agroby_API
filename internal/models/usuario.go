package models

import "time"

type Usuario struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`
	Nome               string     `json:"nome"`
	Email              string     `json:"email" gorm:"unique"`
	Whatsapp           string     `json:"whatsapp"`
	PasswordHash       string     `json:"-" gorm:"column:password_hash"` // json:"-" esconde o hash
	TentativasLogin    int        `json:"-" gorm:"column:tentativas_login;default:0"`
	Bloqueado          bool       `json:"-" gorm:"column:bloqueado;default:false"`
	ResetTokenHash     string     `json:"-" gorm:"column:reset_token_hash"`
	ResetTokenExpiraEm *time.Time `json:"-" gorm:"column:reset_token_expira_em"`
	TipoConta          string     `json:"tipo_conta"` // 'C' ou 'P'
	NomePropriedade    string     `json:"nome_propriedade"`
	EnderecoRural      string     `json:"endereco_rural"`
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
