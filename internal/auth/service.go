package auth

import (
	"agroby_API/internal/database"
	"agroby_API/internal/models"
	"agroby_API/internal/security"
	"errors"
)

func CadastrarUsuario(dados models.UserCredentials) error {
	if len(dados.Password) < 8 || !security.ValidarComplexidade(dados.Password) {
		return errors.New("a senha não atende aos requisitos de segurança")
	}

	hash, err := security.HashPassword(dados.Password)
	if err != nil {
		return err
	}

	return database.DB.Create(&models.Usuario{
		Nome:            dados.Nome,
		Email:           dados.Email,
		PasswordHash:    hash,
		Whatsapp:        dados.Whatsapp,
		TipoConta:       dados.Tipo,
		NomePropriedade: dados.NomePropriedade,
		EnderecoRural:   dados.EnderecoRural,
	}).Error
}
