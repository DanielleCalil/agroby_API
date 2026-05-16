package handlers

import (
	"agroby_API/internal/auth"
	"agroby_API/internal/database"
	"agroby_API/internal/models"
	"agroby_API/internal/security"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var user models.Usuario
	result := database.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil || !security.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login realizado!",
		"token":   token,
		"user": gin.H{
			"id":               user.ID,
			"nome":             user.Nome,
			"email":            user.Email,
			"whatsapp":         user.Whatsapp,
			"tipo_conta":       user.TipoConta,
			"nome_propriedade": user.NomePropriedade,
			"endereco_rural":   user.EnderecoRural,
		},
	})
}

func MeHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sessão inválida"})
		return
	}

	var user models.Usuario
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":               user.ID,
			"nome":             user.Nome,
			"email":            user.Email,
			"whatsapp":         user.Whatsapp,
			"tipo_conta":       user.TipoConta,
			"nome_propriedade": user.NomePropriedade,
			"endereco_rural":   user.EnderecoRural,
		},
	})
}

func RegisterHandler(c *gin.Context) {
	var input models.UserCredentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if err := auth.CadastrarUsuario(input); err != nil {
		switch {
		case errors.Is(err, auth.ErrSenhaInvalida):
			c.JSON(http.StatusBadRequest, gin.H{"error": "A senha não atende aos requisitos de segurança."})
		case errors.Is(err, auth.ErrEmailJaCadastrado):
			c.JSON(http.StatusConflict, gin.H{"error": "Já existe uma conta com este e-mail."})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível concluir o cadastro no momento. Tente novamente."})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário cadastrado com sucesso!"})
}
