# Life Uptime 🌐

Uma API robusta e leve desenvolvida em **Go** para monitoramento de disponibilidade e tempo de resposta (uptime) de serviços web. Projetada para verificar continuamente aplicações, APIs e sites, garantindo que você seja o primeiro a saber quando algo sai do ar.

---

## 🚀 Sobre o Projeto

O Life Uptime serve como um monitor centralizado para a saúde dos seus serviços digitais. Se você possui diversas aplicações rodando e precisa garantir que elas estão online:

- **Verificação de Endpoints HTTP/HTTPS**
- **Monitoramento de Latência e Status Codes**
- **Registro Histórico de Disponibilidade e Erros**

Esta API é o lugar ideal para cadastrar, consultar e gerenciar o status das suas aplicações de forma simples e eficiente.

---

## 🛠️ Tecnologias Utilizadas

| Tecnologia | Descrição |
|---|---|
| [Go (Golang)](https://go.dev/) | Linguagem principal — alta performance e concorrência nativa |
| [Gin Gonic](https://github.com/gin-gonic/gin) | Framework web rápido e minimalista |
| [PostgreSQL](https://www.postgresql.org/) | Banco de dados relacional robusto para armazenamento dos monitores e logs de ping |

---

## 📁 Estrutura do Projeto

A aplicação segue uma estrutura modular e limpa em Go:

```
life-uptime/
├── cmd/                 # Ponto de entrada da aplicação
├── internal/
│   ├── database/        # Configuração e conexão com o banco de dados
│   ├── handler/         # Handlers HTTP para os monitores e definição de rotas
│   ├── model/           # Modelos de dados de monitores e logs de ping
│   ├── repository/      # Operações e comunicação direta com o banco de dados
│   └── worker/          # Background workers responsáveis por realizar os pings nas URLs
├── go.mod
├── go.sum
├── Life-Uptime.postman_collection.json
└── README.md
```

---

## 📡 Integração e Rotas da API

### Monitores

Você pode cadastrar, listar e gerenciar seus monitores.

**Criar um monitor:**
`POST /api/v1/monitors/`
```json
{
  "url": "https://example.com",
  "interval": 60000000000,
  "active": true
}
```

**Listar todos os monitores:**
`GET /api/v1/monitors/`

**Consultar um monitor específico:**
`GET /api/v1/monitors/:id`

**Consultar histórico de pings (logs):**
`GET /api/v1/monitors/:id/history`

---

## ⚙️ Como Começar

### Pré-requisitos

- [Go](https://go.dev/dl/) `v1.20` ou superior
- [PostgreSQL](https://www.postgresql.org/download/)

### Instalação

1. **Clone o repositório:**

```bash
git clone https://github.com/PatrikMaltacm/life-uptime.git
cd life-uptime
```

2. **Instale as dependências:**

```bash
go mod tidy
```

3. **Configuração do Banco de Dados:**

Execute os comandos SQL abaixo no PostgreSQL para criar as tabelas e inserir alguns dados de teste:

<details>
<summary>Clique para expandir os scripts SQL de configuração</summary>

```sql
CREATE TABLE IF NOT EXISTS monitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    interval BIGINT NOT NULL,
    active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS ping_logs (
    id SERIAL PRIMARY KEY,
    monitor_id UUID NOT NULL,
    status_code INTEGER,
    latency_ms BIGINT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    error TEXT,
    CONSTRAINT fk_monitor FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
);

CREATE INDEX idx_ping_logs_monitor_id ON ping_logs(monitor_id);
CREATE INDEX idx_ping_logs_timestamp ON ping_logs(timestamp);

-- 1. Limpa os dados existentes para evitar duplicatas nos testes
TRUNCATE TABLE ping_logs, monitors RESTART IDENTITY CASCADE;

-- 2. Insere Monitores de Teste (Exemplos de 30s, 1min e 5min)
-- Nota: O intervalo está em nanosegundos (1s = 1.000.000.000)
INSERT INTO monitors (id, url, interval, active) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'https://google.com', 30000000000, true),
('550e8400-e29b-41d4-a716-446655440001', 'https://github.com', 60000000000, true),
('550e8400-e29b-41d4-a716-446655440002', 'https://site-inexistente-teste.com', 300000000000, false);

-- 3. Insere Logs de Ping para o primeiro monitor (Google)
INSERT INTO ping_logs (monitor_id, status_code, latency_ms, timestamp, error) VALUES
('550e8400-e29b-41d4-a716-446655440000', 200, 45, NOW() - INTERVAL '10 minutes', NULL),
('550e8400-e29b-41d4-a716-446655440000', 200, 42, NOW() - INTERVAL '5 minutes', NULL),
('550e8400-e29b-41d4-a716-446655440000', 200, 48, NOW(), NULL);

-- 4. Insere Logs de Ping para o segundo monitor (GitHub)
INSERT INTO ping_logs (monitor_id, status_code, latency_ms, timestamp, error) VALUES
('550e8400-e29b-41d4-a716-446655440001', 200, 120, NOW() - INTERVAL '2 minutes', NULL),
('550e8400-e29b-41d4-a716-446655440001', 503, 0, NOW(), 'Service Unavailable');

-- 5. Insere um log de erro para o monitor inativo
INSERT INTO ping_logs (monitor_id, status_code, latency_ms, timestamp, error) VALUES
('550e8400-e29b-41d4-a716-446655440002', 0, 0, NOW(), 'dial tcp: lookup site-inexistente-teste.com: no such host');
```

</details>

4. **Execute a aplicação:**

```bash
go run cmd/main.go
```
*(Nota: Ajuste o caminho `cmd/main.go` conforme necessário, dependendo de como está o arquivo principal dentro da pasta `cmd`)*

---

## 🧪 Testando a API

Uma collection do **Postman** está incluída na raiz do projeto (`Life-Uptime.postman_collection.json`) para facilitar o teste de todos os endpoints disponíveis. Basta importá-la no seu Postman e testar!

---

## 🔓 Open Source & Contribuição

Este projeto é **Open Source** e está aberto para qualquer um utilizar, modificar e melhorar. Sinta-se à vontade para:

- 🐛 Abrir **Issues** para reportar bugs ou sugerir melhorias.
- 🔀 Enviar **Pull Requests** com novas funcionalidades (ex: alertas no Discord/Telegram, frontend para os logs, etc).

---

## 📄 Licença

Distribuído sob a licença **MIT**. O uso é livre para projetos pessoais ou comerciais.

---

<p align="center">Desenvolvido por <strong>Patrik Malta</strong></p>
