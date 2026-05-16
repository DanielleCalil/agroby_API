package auth

import (
	"agroby_API/internal/database"
	"agroby_API/internal/models"
	"agroby_API/internal/security"
	"errors"
	"strings"
)

var (
	ErrSenhaInvalida     = errors.New("a senha não atende aos requisitos de segurança")
	ErrEmailJaCadastrado = errors.New("e-mail já cadastrado")
	ErrFalhaCadastro     = errors.New("não foi possível concluir o cadastro no momento")
)

func CadastrarUsuario(dados models.UserCredentials) error {
	if len(dados.Password) < 8 || !security.ValidarComplexidade(dados.Password) {
		return ErrSenhaInvalida
	}

	hash, err := security.HashPassword(dados.Password)
	if err != nil {
		return err
	}

	err = database.DB.Create(&models.Usuario{
		Nome:            dados.Nome,
		Email:           dados.Email,
		PasswordHash:    hash,
		Whatsapp:        dados.Whatsapp,
		TipoConta:       dados.Tipo,
		NomePropriedade: dados.NomePropriedade,
		EnderecoRural:   dados.EnderecoRural,
	}).Error
	if err != nil {
		if isDuplicateEntryError(err) {
			return ErrEmailJaCadastrado
		}

		return ErrFalhaCadastro
	}

	return nil
}

func isDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "duplicate entry") ||
		strings.Contains(msg, "duplicada") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "violates unique")
}
