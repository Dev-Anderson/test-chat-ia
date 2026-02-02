# MeuChatIA – Backend

Backend do **MeuChatIA**, uma plataforma de atendimento automatizado com suporte a múltiplas empresas, usuários (admin e user) e integração futura com serviços de IA.

O foco do projeto é **arquitetura limpa**, **escala** e **manutenibilidade**, usando Go como linguagem principal.

---

## 🧠 Visão Geral da Arquitetura

O backend segue uma arquitetura modular inspirada em **Clean Architecture**, separando claramente responsabilidades:

```
backend/
├── cmd/
│   └── api/          # Ponto de entrada da aplicação
│
├── internal/
│   ├── config/       # Configurações e variáveis de ambiente
│   ├── database/     # Conexão com PostgreSQL
│   ├── http/
│   │   ├── handlers/ # Controllers / Handlers HTTP
│   │   ├── middleware/ # Autenticação e autorização
│   │   └── routes/  # Definição de rotas
│   ├── models/       # Entidades de domínio (User, etc.)
│   └── services/     # Regras de negócio (JWT, senha, etc.)
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── .env.example
└── README.md
```

### 📌 Princípios adotados

* Separação de responsabilidades
* Código desacoplado
* Fácil evolução para microserviços (ex: IA em Python)
* Pronto para ambientes de produção

---

## 🔐 Autenticação

* Login via **email e senha**
* Senhas armazenadas com **bcrypt**
* Autenticação via **JWT**
* Controle de acesso por **roles**:

  * `ADMIN`
  * `USER`

---

## 🧱 Stack Utilizada

### Backend

* **Golang 1.22**
* **Gin** – framework HTTP
* **GORM** – ORM
* **JWT** – autenticação
* **bcrypt** – hash de senha

### Infraestrutura

* **Docker**
* **Docker Compose**
* **PostgreSQL 15**

### Ferramentas

* **DBeaver** (acesso ao banco)
* **Git**

---

## 🐳 Como rodar o projeto

### 1️⃣ Criar o arquivo `.env`

```env
APP_NAME=MeuChatIA
APP_PORT=8000
ENV=development

POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=meuchat
POSTGRES_HOST=db
POSTGRES_PORT=5432

JWT_SECRET=supersecretjwtkey
```

---

### 2️⃣ Subir os containers

```bash
docker-compose up --build
```

A API ficará disponível em:

```
http://localhost:8000
```

---

## 🗄️ Acessar o banco pelo DBeaver

| Campo    | Valor       |
| -------- | ----------- |
| Host     | `localhost` |
| Porta    | `5434`      |
| Database | `meuchat`   |
| Usuário  | `postgres`  |
| Senha    | `postgres`  |

---

## 🔮 Próximos Passos

* Migrations com `golang-migrate`
* Cadastro de usuários via API
* Multi-empresa
* Refresh Token
* Integração com serviço de IA (Python)

---

## 👨‍💻 Autor

Projeto em desenvolvimento para estudos e evolução de uma plataforma de atendimento automatizado com IA.
