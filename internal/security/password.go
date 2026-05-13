package security

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword cumpre o requisito de Salt + Hashing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Custo 14 para robustez
	return string(bytes), err
}

// CheckPasswordHash verifica se a senha digitada bate com o hash do banco
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidarComplexidade verifica se a senha possui os requisitos mínimos de segurança
// (pelo menos uma letra maiúscula, uma minúscula, um número e um caractere especial)
func ValidarComplexidade(password string) bool {
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}