# 🐶 PetCare Jobs

PetCare Jobs é um projeto de estudo em Go que simula um backend para gerenciamento de um sistema de cuidados com pets, incluindo:

* Cadastro de animais
* Recebimento de informações de veterinários
* Processamento assíncrono de tarefas usando **fila + workers**
* Estrutura inspirada em **clean architecture**

---

## 🚀 Tecnologias

* **Go (Golang)**
* **Gin** (framework HTTP)
* **PostgreSQL** (banco de dados)
* **database/sql** + `lib/pq`
* Workers e filas em memória (por enquanto)
* Variáveis de ambiente (.env)

---

## 📁 Estrutura do Projeto

```text
petcare-jobs/
├── cmd/
│   └── api/
│       └── main.go          # Entry point da aplicação
├── internal/
│   ├── animal/              # Domínio de animais (model, handler, service)
│   ├── exam/                # Domínio de exames
│   ├── job/                 # Jobs, fila e workers
│   ├── repository/          # Repositórios e conexão com banco
│   ├── http/                # Configuração de rotas HTTP
│   └── storage/             # Implementações de storage (ex: memória)
├── go.mod
├── go.sum
└── README.md
```

---

## 🧠 Conceitos Importantes

### 🔹 Domínio

Cada pasta dentro de `internal/` representa um **domínio de negócio**, por exemplo:

* `animal`
* `exam`
* `job`

O domínio contém regras de negócio e não depende de frameworks.

---

### 🔹 Workers e Fila

O projeto utiliza uma **fila em memória** e um **pool de workers** para processar tarefas assíncronas, como:

* Gerar relatórios
* Processar dados
* Simular tarefas pesadas

Fluxo simplificado:

```
API → cria Job → fila → worker → processamento
```

Isso permite que a API responda rápido, mesmo quando há tarefas demoradas.

---


## 🐘 Banco de Dados

* Banco utilizado: **PostgreSQL**
* Pode rodar localmente via Docker ou instalação local
* Ferramentas recomendadas:

  * DBeaver
  * pgAdmin

---

## ▶️ Rodando o Projeto

```bash
go run cmd/api/main.go
```

Você deverá ver logs parecidos com:

```text
PetCare Jobs API started
[worker 1] processing job
job finished
```

---

## 🌱 Branches

Este projeto utiliza branches semânticas, por exemplo:

* `chore/initial-setup`
* `feat/animal-api`
* `fix/db-connection`

---

## 📌 Status do Projeto

🚧 **Em desenvolvimento (estudo)**

Funcionalidades serão adicionadas aos poucos, sem foco comercial.

---

## ✨ Próximos Passos

* [ ] CRUD de animais
* [ ] CRUD de exames
* [ ] Persistência real de jobs
* [ ] Retry e dead-letter queue
* [ ] Dockerização
* [ ] Testes unitários

Este projeto é apenas para fins educacionais.
