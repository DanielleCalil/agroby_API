package main

import (
	"agroby_API/internal/database"
	"agroby_API/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de sistema.")
	}

	database.Connect() // Inicializa o banco

	r := gin.Default()

	// Configuração de rotas que seu React vai chamar
	r.POST("/api/login", handlers.LoginHandler)
	r.POST("/api/cadastro", handlers.RegisterHandler)

	r.Run(":8080")
}