# agroby_API

API em Go com autenticação por JWT.

## Variáveis de ambiente

- APP_PORT: porta HTTP da aplicação. Opcional, padrão 8080.
- DATABASE_URL: string de conexão do SQL Server.
- JWT_SECRET: segredo usado para assinar e validar os tokens de autenticação.

## Endpoints principais

- POST /api/cadastro: cria um usuário.
- POST /api/login: valida credenciais e retorna token JWT e dados básicos do usuário.
- GET /api/me: rota protegida. Exige header Authorization com Bearer token válido.

