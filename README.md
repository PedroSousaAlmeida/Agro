# 🌾 Agro Monitoring API

> API para monitoramento de áreas agrícolas com controle de pragas e sistema multi-tenant

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![Keycloak](https://img.shields.io/badge/Keycloak-26.5-4D4D4D?logo=keycloak)](https://www.keycloak.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)](https://redis.io/)

## 📋 Sobre o Projeto

Sistema de monitoramento agrícola multi-tenant que permite que diferentes usinas gerenciem seus monitoramentos de pragas, aplicações de herbicidas e áreas cultivadas de forma isolada e segura.

**Principais características:**
- ✅ **Multi-tenancy** com isolamento completo de dados por client (usina)
- ✅ Autenticação via **Keycloak** com OIDC/JWT
- ✅ Upload e processamento de **CSV** com dados de monitoramento
- ✅ Processamento **assíncrono** de aplicações em massa
- ✅ API **RESTful** com arquitetura limpa
- ✅ **Migrations** automáticas de banco de dados

## 🚀 Tecnologias

- **[Go 1.25](https://go.dev/)** - Linguagem de programação
- **[Chi v5](https://github.com/go-chi/chi)** - Router HTTP leve e rápido
- **[PostgreSQL 16](https://www.postgresql.org/)** - Banco de dados relacional
- **[Redis 7](https://redis.io/)** - Cache e filas
- **[Keycloak 26.5](https://www.keycloak.org/)** - Gerenciamento de identidade e autenticação
- **[gocloak](https://github.com/Nerzal/gocloak)** - Cliente Go para Keycloak Admin API
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Migrations de banco
- **[Air](https://github.com/air-verse/air)** - Live reload para desenvolvimento

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture** e **Domain-Driven Design (DDD)** com:

- **Separação de camadas**: Domain, DTO, Handler, UseCase, Repository
- **Dependency Injection**: Bootstrap centralizado
- **Repository Pattern**: Abstraindo persistência
- **Multi-tenancy**: Isolamento automático por `client_id`
- **OIDC/JWT**: Autenticação stateless

### Diagrama de Arquitetura

```
┌──────────────┐
│   Cliente    │
└──────┬───────┘
       │ HTTP/JSON
       ▼
┌──────────────────────────────────────────┐
│          Chi Router + Middlewares        │
│  (Auth, CORS, Tenancy, Logging)          │
└──────────────┬───────────────────────────┘
               │
       ┌───────┴────────┬─────────────┐
       ▼                ▼             ▼
   ┌───────┐      ┌──────────┐   ┌────────┐
   │Handler│      │ Handler  │   │Handler │
   └───┬───┘      └────┬─────┘   └───┬────┘
       │               │             │
       ▼               ▼             ▼
   ┌───────┐      ┌──────────┐   ┌────────┐
   │UseCase│      │ UseCase  │   │UseCase │
   └───┬───┘      └────┬─────┘   └───┬────┘
       │               │             │
       ▼               ▼             ▼
  ┌─────────┐    ┌──────────┐  ┌─────────┐
  │Repository│   │Repository│  │Repository│
  └─────┬────┘   └────┬─────┘  └────┬────┘
        │             │             │
        └─────────────┴─────────────┘
                      │
                      ▼
              ┌──────────────┐
              │  PostgreSQL  │
              └──────────────┘
```

## 📦 Módulos

### `clients` - Multi-Tenancy
Gerenciamento de clients (usinas) e seus usuários.
- Criação de clients com slug único
- Registro de usuários por client
- Validação de limite de usuários
- Estatísticas por client
- Integração com Keycloak (grupos, atributos)

### `monitoring`
Upload e processamento de CSVs com dados de monitoramento.
- Parse de CSV com dados agrícolas
- Validação de formato
- Criação em batch de áreas

### `area`
Gerenciamento de áreas monitoradas.
- Listagem com filtros (fazenda, praga, monitoramento)
- Busca por ID
- Gerenciamento de aplicações de herbicidas

### `jobs`
Processamento assíncrono de tarefas em massa.
- Aplicações de herbicida em batch
- Status e progresso de jobs
- Worker em background (Redis)

### `user`
Informações do usuário autenticado.
- Endpoint `/me` com claims JWT

## 🔐 Autenticação

### Keycloak OIDC

A API utiliza **Keycloak** para autenticação via **OpenID Connect (OIDC)**:

1. **Login**: POST ao Keycloak retorna `access_token` (JWT)
2. **Requests**: Incluir header `Authorization: Bearer <token>`
3. **Validação**: Middleware verifica assinatura e expiração do token
4. **Claims**: Informações do usuário extraídas do token

### Multi-Tenancy via JWT

Cada token JWT contém o claim `client_id` que identifica a usina do usuário:

```json
{
  "sub": "user-uuid",
  "email": "joao@usinasantaclara.com",
  "client_id": "uuid-usina-santa-clara",
  ...
}
```

O middleware `ExtractTenancy` extrai `client_id` e `user_id` do token e injeta no contexto. Todas as queries filtram automaticamente por `client_id`, garantindo **isolamento total** de dados.

## 🗄️ Banco de Dados

### Schema Multi-Tenant

Principais tabelas:

**`clients`** - Usinas (clients)
- `id`, `name`, `slug`, `max_users`, `active`
- `keycloak_group_id`, `metadata` (JSONB)

**`client_users`** - Usuários por client
- `id`, `client_id`, `user_id`, `email`, `role`, `active`
- UNIQUE(client_id, user_id)

**`monitoramentos`** - Uploads de CSV
- `id`, `data_upload`, `nome_arquivo`, `status`
- `client_id`, `user_id` (multi-tenancy)

**`areas_monitoramento`** - Áreas monitoradas
- `id`, `monitoramento_id`, `setor`, `cod_fazenda`, `quadra`
- `pragas_data` (JSONB), `aplicacoes` (JSONB array)
- `client_id`, `user_id` (multi-tenancy)

**`jobs`** - Tarefas assíncronas
- `id`, `type`, `status`, `payload` (JSONB)
- `progress`, `processed_items`, `total_items`
- `client_id`, `user_id` (multi-tenancy)

**`client_stats`** (VIEW) - Estatísticas agregadas

### Migrations

As migrations são versionadas e executadas com `golang-migrate`:

```bash
migrate -path migrations -database "postgresql://user:pass@localhost/agro_monitoring" up
```

**Migrations disponíveis:**
- `001` - Criar monitoramentos
- `002` - Criar areas_monitoramento
- `003` - Histórico de migrations
- `004` - Criar jobs
- `005` - Adicionar aplicacoes
- `006` - Criar clients e client_users
- `007` - Adicionar multi-tenancy (client_id, user_id)
- `008` - View client_stats

## ⚙️ Configuração

### Pré-requisitos

- **Go 1.25+** - [Instalar](https://go.dev/dl/)
- **PostgreSQL 16+** - [Instalar](https://www.postgresql.org/download/)
- **Redis 7+** - [Instalar](https://redis.io/download/)
- **Docker + Docker Compose** - [Instalar](https://docs.docker.com/get-docker/) (para Keycloak)
- **golang-migrate** (opcional) - [Instalar](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz (use `.env.example` como base):

```env
# API
PORT=8080
APP_BASE_URL=http://localhost:8080

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=agro_monitoring
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Keycloak
KEYCLOAK_URL=http://localhost:9090
KEYCLOAK_REALM=agro-realm
KEYCLOAK_CLIENT_ID=agro-api
KEYCLOAK_ADMIN_CLIENT_ID=agro-admin
KEYCLOAK_ADMIN_CLIENT_SECRET=admin-secret-change-in-prod
```

### Instalação

1. **Clone o repositório:**
```bash
git clone https://github.com/seu-usuario/agro-monitoring.git
cd agro-monitoring
```

2. **Instale as dependências:**
```bash
go mod download
```

3. **Configure o .env:**
```bash
cp .env.example .env
# Edite .env com suas configurações
```

4. **Suba o Keycloak e Redis:**
```bash
docker-compose up -d
```

5. **Crie o banco de dados:**
```bash
createdb agro_monitoring
```

6. **Execute as migrations:**
```bash
migrate -path migrations -database "postgresql://postgres:postgres@localhost/agro_monitoring?sslmode=disable" up
```

7. **Execute a aplicação:**
```bash
# Desenvolvimento (com live reload)
air

# Ou diretamente com go run
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

## 🐳 Docker

### Serviços Disponíveis

O `docker-compose.yml` fornece:

- **Keycloak** - `http://localhost:9090`
  - Admin: `admin` / `admin`
  - Realm: `agro-realm` (importado automaticamente)

- **PostgreSQL (Keycloak)** - `localhost:5433`
  - Banco exclusivo para Keycloak

- **Redis** - `localhost:6379`
  - Para cache e filas

### Comandos

```bash
# Subir todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f keycloak

# Parar serviços
docker-compose down

# Parar e remover volumes (reset completo)
docker-compose down -v
```

## 🏃 Executando

### Desenvolvimento

**Com Air (recomendado - live reload):**
```bash
air
```

**Com go run:**
```bash
go run cmd/api/main.go
```

### Produção

```bash
# Build
go build -o api cmd/api/main.go

# Executar
./api
```

### Testes

```bash
# Executar todos os testes
go test ./...

# Com coverage
go test -cover ./...

# Coverage HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📡 API Endpoints

### Públicos (sem autenticação)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check |
| POST | `/v1/register/{slug}` | Registrar usuário em um client |

### Autenticados

#### Clients
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/v1/clients/me` | Meu client |
| GET | `/v1/clients/me/stats` | Estatísticas do meu client |
| GET | `/v1/clients/me/users` | Usuários do meu client |

#### Monitoramentos
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/monitoramentos` | Upload CSV |
| GET | `/v1/monitoramentos` | Listar uploads |
| GET | `/v1/monitoramentos/{id}` | Buscar por ID |

#### Áreas
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/v1/areas` | Listar áreas |
| GET | `/v1/areas/{id}` | Buscar área por ID |
| GET | `/v1/areas/search/fazenda` | Buscar por fazenda |
| GET | `/v1/areas/search/praga` | Buscar por praga |
| POST | `/v1/areas/{id}/aplicacao` | Adicionar aplicação |

#### Jobs
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/jobs/aplicacoes` | Criar job de aplicações em massa |
| GET | `/v1/jobs/{id}` | Status do job |

#### Users
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/v1/users/me` | Claims do usuário autenticado |

### Admin (requer permissão de admin)

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/v1/admin/clients` | Criar client |
| GET | `/v1/admin/clients` | Listar clients |
| GET | `/v1/admin/clients/{id}` | Buscar client |
| GET | `/v1/admin/clients/{id}/stats` | Estatísticas do client |

## 🧪 Testes

### Postman Collection

Importe a collection em `docs/postman_collection.json` no Postman.

**Variáveis configuradas:**
- `base_url`: http://localhost:8080
- `keycloak_url`: http://localhost:9090
- `keycloak_realm`: agro-realm
- `access_token`: (preenchido automaticamente após login)

### Fluxo de Teste Completo

1. **Login** (Authentication > Login and Get Token)
   - Faz login com `testuser` / `password`
   - Salva `access_token` automaticamente

2. **Criar Client** (Admin - Clients > Create Client)
   - Cria usina com slug único
   - Salva `client_id` e `client_slug`

3. **Registrar Usuário** (Clients > Register User)
   - Registra usuário na usina
   - Cria no Keycloak + BD

4. **Login com novo usuário**
   - Trocar username/password no "Login and Get Token"
   - Token terá `client_id` da usina

5. **Testar isolamento multi-tenant**
   - Upload CSV
   - Listar monitoramentos (só do seu client)
   - Criar aplicações

## 🔑 Multi-Tenancy

### Como Funciona

1. **Criação de Client (Usina)**
   - Admin cria client com slug único (ex: `usina-santa-clara`)
   - Sistema cria grupo no Keycloak: `/clients/usina-santa-clara`
   - Gera URL de registro: `http://app.agro.com/register/usina-santa-clara`

2. **Registro de Usuário**
   - Usuário acessa URL específica do client
   - Sistema valida limite de usuários
   - Cria usuário no Keycloak
   - Adiciona ao grupo do client
   - Seta atributo `client_id` no usuário

3. **Login e Acesso**
   - Usuário faz login via Keycloak
   - Token JWT contém claim `client_id`
   - Middleware `ExtractTenancy` injeta no contexto
   - Todas queries filtram por `client_id`

### Isolamento de Dados

**Garantias:**
- ✅ Usuário só vê dados do seu client
- ✅ Impossível acessar dados de outro client
- ✅ `client_id` e `user_id` obrigatórios em novas operações
- ✅ Dados antigos (sem client_id) são nullable (retrocompatibilidade)
- ✅ Filtros automáticos em todas queries

**Implementação:**
```go
// Context helpers
clientID, _ := sharedContext.GetClientID(ctx)
userID, _ := sharedContext.GetUserID(ctx)

// Repositories filtram automaticamente
monitoramentos, _ := repo.List(ctx, clientID, page, pageSize)
```

## 📊 Estrutura do Projeto

```
agro-monitoring/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── bootstrap/
│   ├── app.go                   # Dependency Injection
│   ├── routes.go                # Configuração de rotas
│   ├── database.go              # Conexão PostgreSQL
│   └── env.go                   # Carregamento .env
├── internal/
│   ├── config/
│   │   └── env.go               # Variáveis de ambiente
│   ├── modules/
│   │   ├── clients/             # Multi-tenancy
│   │   │   ├── domain/
│   │   │   ├── dto/
│   │   │   ├── handler/
│   │   │   ├── usecase/
│   │   │   ├── repository/
│   │   │   └── service/         # Keycloak Admin API
│   │   ├── area/                # Áreas monitoradas
│   │   ├── jobs/                # Processamento assíncrono
│   │   ├── monitoring/          # Upload CSV
│   │   └── user/                # Usuário autenticado
│   ├── services/
│   │   ├── csv/                 # Parser CSV
│   │   └── queue/               # Redis Queue
│   └── shared/
│       ├── context/             # Context helpers
│       ├── errors/              # Erros globais
│       ├── middleware/          # Auth, CORS, Tenancy
│       └── response/            # Response padronizado
├── migrations/                  # Database migrations
├── docker/
│   └── keycloak/
│       └── realm-export.json    # Realm Keycloak
├── docs/
│   └── postman_collection.json  # Collection Postman
├── docker-compose.yml           # Keycloak + Redis
├── .air.toml                    # Configuração Air
├── .env.example                 # Exemplo variáveis
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🏛️ Padrão de Módulo

Cada módulo segue a mesma estrutura:

```
module/
├── domain/
│   ├── entity.go          # Entidades (structs)
│   └── repository.go      # Interface do repository
├── dto/
│   └── dto.go            # Request/Response + Converters
├── handler/
│   └── http.go           # HTTP handlers + RegisterRoutes
├── usecase/
│   └── usecase.go        # Lógica de negócio
├── repository/
│   ├── postgres.go       # Implementação PostgreSQL
│   └── inmemory.go       # Implementação em memória (testes)
└── service/              # (opcional) Serviços externos
    └── keycloak.go
```

## 🔄 Fluxo de Dados

```
HTTP Request
    ↓
[Middlewares] Auth → ExtractTenancy → RequireClient
    ↓
Handler (valida request, chama UseCase)
    ↓
UseCase (lógica de negócio, orquestra Repositories)
    ↓
Repository (acessa PostgreSQL com filtros client_id)
    ↓
Database (retorna entidades)
    ↓
DTO Converter (Domain → Response)
    ↓
HTTP Response (JSON)
```

## 🛡️ Segurança

### Autenticação e Autorização
- ✅ OIDC/JWT via Keycloak
- ✅ Verificação de assinatura digital
- ✅ Validação de expiração de tokens
- ✅ Middleware de autenticação em todas rotas `/v1/*`

### Isolamento Multi-Tenant
- ✅ `client_id` obrigatório no context
- ✅ Queries filtram automaticamente por client
- ✅ Impossível acessar dados de outro client

### Boas Práticas
- ✅ CORS configurado
- ✅ Validação de entrada
- ✅ Prepared statements (previne SQL injection)
- ✅ Erros genéricos ao usuário (sem expor internals)
- ✅ Secrets em variáveis de ambiente (.env não versionado)

## 📝 Convenções de Código

### Commits
Seguimos [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: adicionar nova funcionalidade
fix: corrigir bug
docs: atualizar documentação
refactor: refatorar código
test: adicionar testes
chore: tarefas de manutenção
```

### Go
- `go fmt` para formatação
- `go vet` para análise estática
- Nomes exportados começam com maiúscula
- Interfaces com sufixo `er` quando possível
- Erros customizados em `shared/errors`

### API
- RESTful
- Versionamento: `/v1/`
- Response padronizado: `SuccessResponse` e `ErrorResponse`
- Status HTTP semânticos

## 🤝 Contribuindo

Contribuições são bem-vindas! Siga os passos:

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'feat: adicionar nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

### Checklist PR
- [ ] Código segue convenções do projeto
- [ ] Testes adicionados/atualizados
- [ ] Documentação atualizada
- [ ] Migrations criadas (se necessário)
- [ ] Postman collection atualizado (se novos endpoints)

## 🐛 Troubleshooting

### Keycloak não inicia
```bash
# Verificar logs
docker-compose logs keycloak

# Recriar com volumes limpos
docker-compose down -v
docker-compose up -d
```

### Erro "client_id not found in token"
- Verifique se o protocol mapper `client-id-mapper` está configurado no Keycloak
- Reimporte o realm: `docker-compose down -v && docker-compose up -d`

### Migrations falham
```bash
# Verificar status
migrate -path migrations -database "postgresql://..." version

# Forçar versão
migrate -path migrations -database "postgresql://..." force <version>
```

### "Failed to create client" (500)
- Verifique permissões do service account `agro-admin` no Keycloak
- Deve ter roles: `manage-users` e `manage-groups`

## 📚 Recursos Úteis

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [Go Chi Router](https://go-chi.io/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [OIDC Spec](https://openid.net/specs/openid-connect-core-1_0.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👥 Autores

- **Pedro Sousa Almeida** - Desenvolvimento inicial

---

**Desenvolvido com ❤️ para o agronegócio brasileiro**
