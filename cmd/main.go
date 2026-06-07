package main

import (
	"agroby_API/internal/auth"
	"agroby_API/internal/database"
	"agroby_API/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envPaths := []string{".env", "../.env"}
	loadedEnv := false
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			loadedEnv = true
			break
		}
	}
	if !loadedEnv {
		log.Println("Aviso: Arquivo .env não encontrado em .env ou ../.env. Usando variáveis de sistema.")
	}

	if err := auth.EnsureJWTConfigured(); err != nil {
		log.Fatal(err)
	}

	database.Connect() // Inicializa o banco

	r := gin.Default()
	r.Use(corsMiddleware())

	// Configuração de rotas que seu React vai chamar
	r.POST("/api/login", handlers.LoginHandler)
	r.POST("/api/esqueci-senha", handlers.ForgotPasswordHandler)
	r.POST("/api/resetar-senha", handlers.ResetPasswordHandler)
	r.POST("/api/cadastro", handlers.RegisterHandler)
	r.GET("/api/me", auth.AuthMiddleware(), handlers.MeHandler)
	r.GET("/api/dashboard/resumo", auth.AuthMiddleware(), handlers.DashboardResumoHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Falha ao iniciar servidor HTTP:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
