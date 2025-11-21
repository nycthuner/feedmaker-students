# Feedmaker Students

Sistema de gerenciamento de feedbacks para alunos e professores, desenvolvido como trabalho acadêmico usando Go, Gin Framework, GORM e PostgreSQL.

## 📋 Sobre o Projeto

O **Feedmaker Students** é uma API REST que permite professores fornecerem feedbacks para alunos, facilitando o acompanhamento do desempenho acadêmico. O sistema suporta criação de usuários, gerenciamento de feedbacks e consultas por aluno.

## 🛠️ Tecnologias

- **Go 1.24** - Linguagem de programação
- **Gin** - Framework web HTTP
- **GORM** - ORM para banco de dados
- **PostgreSQL 15** - Banco de dados relacional
- **Docker & Docker Compose** - Containerização
- **Air** - Hot-reload para desenvolvimento
- **godotenv** - Gerenciamento de variáveis de ambiente

## 🚀 Como Executar

### Pré-requisitos

- Docker e Docker Compose instalados
- Git

### Passos

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/anaclaraddias/feedmaker-students.git
   cd feedmaker-students
   ```

2. **Configure as variáveis de ambiente:**
   ```bash
   cp .env.example .env
   ```
   Edite o `.env` se necessário (valores padrão já funcionam com Docker Compose).

3. **Suba os containers:**
   ```bash
   docker compose up --build
   ```

4. **Acesse a API:**
   - URL base: `http://localhost:5567`
   - Health check: `http://localhost:5567/health`

## ⚙️ Variáveis de Ambiente

O serviço usa variáveis de ambiente para configurar a conexão com o banco e a porta do servidor. Você pode copiar o arquivo de exemplo:

```bash
cp .env.example .env
```

As variáveis mais importantes são:

- `DB_HOST` (ex: `db` quando usando Docker Compose)
- `DB_PORT` (ex: `5432`)
- `DB_USER` (ex: `admin`)
- `DB_PASSWORD` (ex: `admin`)
- `DB_NAME` (ex: `feedmaker_db`)
- `DB_SSLMODE` (ex: `disable`)
- `PORT` (porta HTTP do servidor, ex: `5567`)

O `docker-compose.yml` já define valores padrão compatíveis com a configuração de desenvolvimento.

## 🧑‍💻 Execução em desenvolvimento (Docker + Air)

O ambiente de desenvolvimento usa `Dockerfile.dev` que instala o `air` para hot-reload. Para iniciar a API e o banco em desenvolvimento:

```bash
docker compose up --build
```

O serviço backend fica exposto em `http://localhost:5567` (padrão definido no `docker-compose.yml`). O container backend monta o código-fonte, então alterações em `.go` reiniciam automaticamente a aplicação via `air`.

## 🗄️ Migrations / Aplicar arquivos SQL

O projeto mantém scripts SQL em `src/infra/database/`. Existem duas abordagens para aplicar o esquema:

1. Usar uma ferramenta de migrations (recomendado): instale `golang-migrate` e aponte para a pasta de migrations.

2. Aplicar os arquivos SQL diretamente no container Postgres. Exemplo (aplica `Feedback.sql`):

```bash
# copia o arquivo para o container
docker cp src/infra/database/Feedback.sql feedmaker_db:/tmp/Feedback.sql

# executa o arquivo dentro do container
docker compose exec db psql -U admin -d feedmaker_db -f /tmp/Feedback.sql
```

Se preferir, instale `psql` localmente e rode:

```bash
psql postgresql://admin:admin@localhost:5432/feedmaker_db -f src/infra/database/Feedback.sql
```

## 🧾 Observações sobre o banco (GORM)

- A aplicação usa GORM como ORM (`gorm.io/gorm`, `gorm.io/driver/postgres`).
- Alguns modelos definem `TableName()` para forçar nomes de tabela no singular (por exemplo `user` e `feedback`).

---

## 📚 Documentação da API

### Base URL
```
http://localhost:5567
```

---

### 🏥 Health Check

Verifica o status da aplicação e a conexão com o banco de dados.

**Endpoint:** `GET /health`

**Resposta de Sucesso (200):**
```json
{
  "status": "healthy"
}
```

**Resposta de Erro (500):**
```json
{
  "status": "unhealthy",
  "error": "mensagem de erro"
}
```

---

### 👤 Criar Usuário

Cria um novo usuário no sistema (aluno, professor ou coordenador).

**Endpoint:** `POST /user`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Ana Clara",
  "username": "anaclara",
  "password": "senha123",
  "type": "STUDENT"
}
```

**Campos:**
- `name` (string, obrigatório): Nome completo do usuário
- `username` (string, obrigatório): Nome de usuário único
- `password` (string, obrigatório): Senha do usuário
- `type` (string, obrigatório): Tipo de usuário
  - Valores aceitos: `"STUDENT"`, `"TEACHER"`, `"COORDINATOR"`

**Resposta de Sucesso (201):**
```json
{
  "id": 1,
  "name": "Ana Clara",
  "username": "anaclara",
  "type": "STUDENT"
}
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid request body"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```

---

### 📝 Criar Feedback

Cria um novo feedback de um aluno para um professor.

**Endpoint:** `POST /feedback`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "score": 8,
  "body": "Ótima aula na ultima sexta!",
  "student_id": 2,
  "teacher_id": 1
}
```

**Campos:**
- `score` (integer, obrigatório): Nota do feedback (0-10)
- `body` (string, obrigatório): Texto descritivo do feedback
- `student_id` (integer, obrigatório): ID do aluno que recebe o feedback
- `teacher_id` (integer, obrigatório): ID do professor que dá o feedback

**Resposta de Sucesso (201):**
```json
{
  "id": 1,
  "score": 8,
  "body": "Ótimo desempenho na apresentação do trabalho!",
  "student_id": 2,
  "teacher_id": 1,
}
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid request body"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```

---

### 📋 Listar Feedbacks de um Aluno

Retorna todos os feedbacks criados por um aluno específico.

**Endpoint:** `GET /student/:id/feedbacks`

**Parâmetros de Rota:**
- `id` (integer, obrigatório): ID do aluno

**Exemplo de Requisição:**
```bash
GET /student/2/feedbacks
```

**Resposta de Sucesso (200):**
```json
[
  {
    "id": 1,
    "score": 8,
    "body": "Aula passada foi massa!",
    "student_id": 2,
    "teacher_id": 1,
  },
  {
    "id": 2,
    "score": 9,
    "body": "Estou gostando bastante das novas atividades.",
    "student_id": 2,
    "teacher_id": 3,
  }
]
```

**Resposta quando não há feedbacks (200):**
```json
[]
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid student id"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```

---

### 🔍 Buscar Feedback por ID

Retorna os detalhes de um feedback específico.

**Endpoint:** `GET /feedback/:id`

**Parâmetros de Rota:**
- `id` (integer, obrigatório): ID do feedback

**Exemplo de Requisição:**
```bash
GET /feedback/1
```

**Resposta de Sucesso (200):**
```json
{
  "id": 1,
  "score": 8,
  "body": "As aulas estão muito legais!",
  "student_id": 2,
  "teacher_id": 1,
}
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid feedback id"
}
```

**Resposta de Erro (404):**
```json
{
  "error": "feedback not found"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```

---

### 🔐 Login

Autentica um usuário no sistema.

**Endpoint:** `POST /login`

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "username": "anaclara",
  "password": "senha123"
}
```

**Campos:**
- `username` (string, obrigatório): Nome de usuário
- `password` (string, obrigatório): Senha do usuário

**Resposta de Sucesso (200):**
```json
{
  "id": 1,
  "name": "Ana Clara",
  "username": "anaclara",
  "type": "STUDENT"
}
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid request body"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```

---

### 📊 Listar Feedbacks de um Professor

Retorna todos os feedbacks recebidos por um professor específico.

**Endpoint:** `GET /teacher/:id/feedbacks`

**Parâmetros de Rota:**
- `id` (integer, obrigatório): ID do professor

**Exemplo de Requisição:**
```bash
GET /teacher/1/feedbacks
```

**Resposta de Sucesso (200):**
```json
[
  {
    "id": 1,
    "score": 8,
    "body": "Ótimas explicações durante a aula!",
    "student_id": 2,
    "teacher_id": 1
  },
  {
    "id": 3,
    "score": 9,
    "body": "Material didático muito bom.",
    "student_id": 5,
    "teacher_id": 1
  }
]
```

**Resposta quando não há feedbacks (200):**
```json
[]
```

**Resposta de Erro (400):**
```json
{
  "error": "invalid teacher id"
}
```

**Resposta de Erro (500):**
```json
{
  "error": "mensagem de erro"
}
```
