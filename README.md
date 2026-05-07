# Life Uptime

O Life Uptime é uma aplicação de monitoramento de código aberto (open-source) desenvolvida em Go. Ele permite que você monitore a integridade e a disponibilidade de suas aplicações web, APIs e serviços, realizando pings em intervalos definidos e registrando os códigos de status HTTP, latência e quaisquer possíveis erros.

## Funcionalidades

- **Monitoramento de URLs:** Verifique continuamente a disponibilidade de qualquer endpoint HTTP/HTTPS.
- **Intervalos Customizados:** Defina intervalos específicos de checagem para cada monitor.
- **Histórico de Ping:** Mantenha um registro do histórico de latência, códigos de status e tempo de inatividade (downtime).
- **API RESTful:** Gerencie seus monitores facilmente através de um conjunto completo de endpoints de API.
- **Open Source:** Sinta-se à vontade para usar, modificar e contribuir!

## Tecnologias Utilizadas

- **Go (Golang):** Backend rápido, compilado e eficiente.
- **Gin Framework:** Framework web HTTP de alta performance.
- **PostgreSQL:** Banco de dados relacional confiável para armazenar os monitores e logs de ping.

## Como Começar

### Pré-requisitos

- [Go](https://golang.org/doc/install) (1.20+)
- [PostgreSQL](https://www.postgresql.org/download/)

### Instalação

1. Clone o repositório:
   ```bash
   git clone https://github.com/PatrikMaltacm/life-uptime.git
   cd life-uptime
   ```

2. Instale as dependências:
   ```bash
   go mod download
   ```

### Configuração do Banco de Dados

Execute os comandos SQL abaixo em seu banco de dados PostgreSQL para criar as tabelas necessárias e inserir alguns dados de teste:

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

## Testando a API

Uma collection do Postman está incluída na pasta raiz do projeto (`Life-Uptime.postman_collection.json`) para facilitar o teste de todos os endpoints disponíveis da API. Basta importá-la no Postman para começar rapidamente!

## Como Contribuir

Contribuições são bem-vindas! Se você encontrar algum bug ou quiser propor uma nova funcionalidade, por favor abra uma *issue* ou envie um *pull request*. Este projeto é orgulhosamente Open Source.
