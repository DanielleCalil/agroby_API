package handlers

import (
	"agroby_API/internal/auth"
	"agroby_API/internal/database"
	"agroby_API/internal/models"
	"agroby_API/internal/security"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

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

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	if user.Bloqueado {
		c.JSON(http.StatusForbidden, gin.H{"error": "Usuário bloqueado. Solicite redefinição de senha."})
		return
	}

	if !security.CheckPasswordHash(input.Password, user.PasswordHash) {
		novasTentativas := user.TentativasLogin + 1

		if novasTentativas >= 3 {
			novaSenhaTemporaria, err := generateTemporaryPassword()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar tentativa de login."})
				return
			}

			hash, err := security.HashPassword(novaSenhaTemporaria)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar tentativa de login."})
				return
			}

			err = database.DB.Model(&user).Updates(map[string]interface{}{
				"tentativas_login": 0,
				"bloqueado":        true,
				"password_hash":    hash,
			}).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar tentativa de login."})
				return
			}

			c.JSON(http.StatusForbidden, gin.H{"error": "Usuário bloqueado após 3 tentativas. A senha foi resetada."})
			return
		}

		err := database.DB.Model(&user).Update("tentativas_login", novasTentativas).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar tentativa de login."})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
		return
	}

	if user.TentativasLogin > 0 {
		err := database.DB.Model(&user).Update("tentativas_login", 0).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar login."})
			return
		}
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

func generateTemporaryPassword() (string, error) {
	buffer := make([]byte, 18)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
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

func DashboardResumoHandler(c *gin.Context) {
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

	resumo := gin.H{}

	if user.TipoConta == "C" {
		var meusPedidos int64
		var produtosDisponiveis int64

		if err := database.DB.Table("pedidos").Where("id_cliente = ?", userID).Count(&meusPedidos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao carregar resumo do dashboard"})
			return
		}

		if err := database.DB.Table("produtos").Where("ativo = ?", true).Where("estoque > 0").Count(&produtosDisponiveis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao carregar resumo do dashboard"})
			return
		}

		resumo = gin.H{
			"meus_pedidos":         meusPedidos,
			"produtos_disponiveis": produtosDisponiveis,
		}
	} else {
		var safrasAtivas int64
		var produtosCadastrados int64
		var vendasRecebidas int64

		if err := database.DB.Table("safras").Where("id_produtor = ?", userID).Where("status <> ?", "finalizada").Count(&safrasAtivas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao carregar resumo do dashboard"})
			return
		}

		if err := database.DB.Table("produtos").Where("id_produtor = ?", userID).Count(&produtosCadastrados).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao carregar resumo do dashboard"})
			return
		}

		if err := database.DB.Table("pedidos").Where("id_produtor = ?", userID).Count(&vendasRecebidas).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao carregar resumo do dashboard"})
			return
		}

		resumo = gin.H{
			"safras_ativas":        safrasAtivas,
			"produtos_cadastrados": produtosCadastrados,
			"vendas_recebidas":     vendasRecebidas,
		}
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
		"resumo": resumo,
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

func ForgotPasswordHandler(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var user models.Usuario
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Resposta genérica para evitar enumeração de e-mails.
		c.JSON(http.StatusOK, gin.H{"message": "Se o e-mail existir, enviaremos instruções para redefinir a senha."})
		return
	}

	token, tokenHash, err := generateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível iniciar a recuperação de senha."})
		return
	}

	expiraEm := time.Now().Add(30 * time.Minute)
	err = database.DB.Model(&user).Updates(map[string]interface{}{
		"reset_token_hash":      tokenHash,
		"reset_token_expira_em": expiraEm,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível iniciar a recuperação de senha."})
		return
	}

	// Em produção, este token deve ser enviado por e-mail/WhatsApp, não retornado na API.
	c.JSON(http.StatusOK, gin.H{
		"message":     "Se o e-mail existir, enviaremos instruções para redefinir a senha.",
		"reset_token": token,
		"expira_em":   expiraEm,
	})
}

func ResetPasswordHandler(c *gin.Context) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Token == "" || input.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	if len(input.NewPassword) < 8 || !security.ValidarComplexidade(input.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A nova senha não atende aos requisitos de segurança."})
		return
	}

	tokenHash := hashToken(input.Token)

	var user models.Usuario
	err := database.DB.Where("reset_token_hash = ?", tokenHash).First(&user).Error
	if err != nil || user.ResetTokenExpiraEm == nil || user.ResetTokenExpiraEm.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token inválido ou expirado."})
		return
	}

	hash, err := security.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível redefinir a senha."})
		return
	}

	err = database.DB.Model(&user).Updates(map[string]interface{}{
		"password_hash":         hash,
		"bloqueado":             false,
		"tentativas_login":      0,
		"reset_token_hash":      "",
		"reset_token_expira_em": nil,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível redefinir a senha."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Senha redefinida com sucesso!"})
}

func generateResetToken() (string, string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", "", err
	}

	rawToken := base64.RawURLEncoding.EncodeToString(buffer)
	return rawToken, hashToken(rawToken), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
