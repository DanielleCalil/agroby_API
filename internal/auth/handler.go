package auth

import (
	"agroby_API/internal/database"
	"agroby_API/internal/models"
	"agroby_API/internal/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var input models.UserCredentials
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var user models.Usuario
	// Busca o utilizador no SQL Server pelo e-mail
	result := database.DB.Where("email = ?", input.Email).First(&user)

	// REQUISITO: Comparação com Bcrypt e Erro Genérico
	if result.Error != nil || !security.CheckPasswordHash(input.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login realizado com sucesso!"})
}
