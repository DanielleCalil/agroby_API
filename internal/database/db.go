package database

import (
    "log"
    "os"

	"agroby_API/internal/models"
    "gorm.io/driver/sqlserver" 
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    var err error
    
    // Formato da DSN para SQL Server:
    // sqlserver://username:password@localhost:1433?database=dbname
    dsn := os.Getenv("DATABASE_URL")
    
    // Aqui fazemos a troca para sqlserver.Open
    DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    
    if err != nil {
        log.Fatal("Falha ao conectar no banco de dados SQL Server:", err)
    }

	// Cria a tabela 'usuarios' automaticamente baseada no modelo
	DB.AutoMigrate(&models.Usuario{})

    log.Println("Conexão com SQL Server estabelecida com sucesso!")
}

func FindByEmail(email string) (*models.Usuario, error) {
	var user models.Usuario
	result := DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func CreateUser(user *models.Usuario) error {
	result := DB.Create(user)
	return result.Error
}