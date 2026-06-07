# agroby_API

API REST em Go (Gin + GORM) com autenticação JWT para o projeto Agroby.

## Stack

- Go 1.25+
- Gin (`github.com/gin-gonic/gin`)
- GORM (`gorm.io/gorm` + `gorm.io/driver/sqlserver`)
- SQL Server
- JWT (`github.com/golang-jwt/jwt/v5`)

## Funcionalidades

- Cadastro de usuário
- Login com geração de token JWT
- Bloqueio automático após 3 tentativas de login inválidas
- Recuperação e redefinição de senha por token
- Rota de perfil autenticado (`/api/me`)
- Resumo de dashboard autenticado (`/api/dashboard/resumo`)

## Estrutura do projeto

```text
cmd/main.go                    # bootstrap HTTP, rotas e middlewares
internal/auth/                 # JWT, middleware e regras de cadastro
internal/database/             # conexão com SQL Server (GORM)
internal/handlers/             # handlers HTTP da API
internal/models/               # modelos de dados
internal/security/             # hash e validação de senha
sql/schema.sql                 # schema base do banco
```

## Pré-requisitos

- Go instalado
- SQL Server disponível
- Banco criado e acessível pela `DATABASE_URL`

## Configuração

Crie um arquivo `.env` na raiz do projeto:

```env
APP_PORT=8080
DATABASE_URL=sqlserver://usuario:senha@localhost:1433?database=agroby&encrypt=disable
JWT_SECRET=troque-por-um-segredo-forte
```

### Variáveis de ambiente

- `APP_PORT`: porta HTTP da aplicação (opcional, padrão `8080`)
- `DATABASE_URL`: string de conexão do SQL Server (obrigatória)
- `JWT_SECRET`: segredo para assinar/validar JWT (obrigatória)

## Banco de dados

O projeto faz `AutoMigrate` da entidade `usuarios` automaticamente na inicialização.

As tabelas usadas no dashboard (`safras`, `produtos`, `pedidos`, `itens_pedido`) estão definidas em `sql/schema.sql`. Execute esse script no seu SQL Server para disponibilizar todo o schema esperado pela API.

## Como executar

```bash
go mod tidy
go run ./cmd/main.go
```

Servidor padrão: `http://localhost:8080`

## Autenticação

As rotas protegidas exigem header:

```http
Authorization: Bearer <seu_token_jwt>
```

O token expira em 24 horas.

## Endpoints

### 1) Cadastro

`POST /api/cadastro`

Body:

```json
{
	"nome": "Maria Silva",
	"email": "maria@email.com",
	"whatsapp": "11999999999",
	"password": "Senha@123",
	"tipo_conta": "P",
	"nome_propriedade": "Sitio Boa Terra",
	"endereco_rural": "Zona Rural, Km 10"
}
```

Regras de senha:

- mínimo de 8 caracteres
- ao menos 1 letra maiúscula
- ao menos 1 letra minúscula
- ao menos 1 número
- ao menos 1 caractere especial

Respostas comuns:

- `201`: usuário cadastrado
- `400`: dados inválidos ou senha fraca
- `409`: e-mail já cadastrado

### 2) Login

`POST /api/login`

Body:

```json
{
	"email": "maria@email.com",
	"password": "Senha@123"
}
```

Resposta de sucesso (`200`):

```json
{
	"message": "Login realizado!",
	"token": "<jwt>",
	"user": {
		"id": 1,
		"nome": "Maria Silva",
		"email": "maria@email.com",
		"whatsapp": "11999999999",
		"tipo_conta": "P",
		"nome_propriedade": "Sitio Boa Terra",
		"endereco_rural": "Zona Rural, Km 10"
	}
}
```

Comportamento de segurança:

- após 3 tentativas inválidas, o usuário é bloqueado (`403`)
- quando bloqueia, a senha é resetada internamente para uma temporária

### 3) Esqueci minha senha

`POST /api/esqueci-senha`

Body:

```json
{
	"email": "maria@email.com"
}
```

Resposta:

- Sempre retorna mensagem genérica para evitar enumeração de e-mails
- Em ambiente atual também retorna `reset_token` e `expira_em` (somente para desenvolvimento)

Exemplo (`200`):

```json
{
	"message": "Se o e-mail existir, enviaremos instruções para redefinir a senha.",
	"reset_token": "<token>",
	"expira_em": "2026-06-07T14:30:00Z"
}
```

### 4) Redefinir senha

`POST /api/resetar-senha`

Body:

```json
{
	"token": "<reset_token>",
	"new_password": "NovaSenha@123"
}
```

Respostas comuns:

- `200`: senha redefinida com sucesso
- `400`: token inválido/expirado ou senha fora das regras

### 5) Perfil autenticado

`GET /api/me`

Requer JWT válido.

Resposta de sucesso (`200`):

```json
{
	"user": {
		"id": 1,
		"nome": "Maria Silva",
		"email": "maria@email.com",
		"whatsapp": "11999999999",
		"tipo_conta": "P",
		"nome_propriedade": "Sitio Boa Terra",
		"endereco_rural": "Zona Rural, Km 10"
	}
}
```

### 6) Resumo do dashboard

`GET /api/dashboard/resumo`

Requer JWT válido.

Retorno varia por tipo de conta:

- Cliente (`tipo_conta = "C"`): `meus_pedidos`, `produtos_disponiveis`
- Produtor (`tipo_conta != "C"`): `safras_ativas`, `produtos_cadastrados`, `vendas_recebidas`

Exemplo (`200`):

```json
{
	"user": {
		"id": 1,
		"nome": "Maria Silva",
		"email": "maria@email.com",
		"whatsapp": "11999999999",
		"tipo_conta": "P",
		"nome_propriedade": "Sitio Boa Terra",
		"endereco_rural": "Zona Rural, Km 10"
	},
	"resumo": {
		"safras_ativas": 2,
		"produtos_cadastrados": 18,
		"vendas_recebidas": 42
	}
}
```

## Códigos de status frequentes

- `200`: sucesso
- `201`: recurso criado
- `400`: requisição inválida
- `401`: não autenticado / token inválido
- `403`: acesso proibido (ex.: usuário bloqueado)
- `409`: conflito (ex.: e-mail duplicado)
- `500`: erro interno

## Observações

- CORS está liberado para qualquer origem (`*`) no estado atual.
- Em produção, o token de recuperação de senha não deve ser retornado na resposta da API.

