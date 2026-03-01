# go-expert-client-server-api

## Estrutura

- `server/server.go` — API HTTP que consulta a economia.awesomeapi, grava a cotação em
  `client-server-api.db` e expõe `/cotacao`.
- `client/client.go` — cliente que consome o endpoint `/cotacao` e grava `cotacao.txt`.

## Comportamento

- Ao iniciar o `server`, o banco `client-server-api.db` será criado (se não existir) e a
  tabela `cotacao` será criada automaticamente.
- Toda chamada a `GET /cotacao` faz uma requisição externa para
  `https://economia.awesomeapi.com.br/json/last/USD-BRL`, insere um registro na tabela
  `cotacao` com os campos `id`, `code`, `bid` e `datetime` e retorna o resultado JSON.

Campos gravados na tabela `cotacao`:

- `id` — identificador único (AUTOINCREMENT)
- `code` — código da moeda (ex.: `USD`)
- `bid` — valor da cotação (string)
- `datetime` — timestamp da inserção

## Observações importantes

- O exercício solicita que o timeout de gravação no banco seja 10ms


## Como executar

1. No terminal, a partir da raiz do projeto, iniciar o servidor:

```bash
go run ./server/server.go
```

O servidor inicia na porta `8080`.

2. Em outro terminal, executar o cliente (faz a chamada a `http://localhost:8080/cotacao`):

```bash
go run ./client/client.go
```

O cliente grava ou sobrescreve o arquivo `cotacao.txt` com o conteúdo no formato:

```
Dólar: {valor}
```

## Inspecionar o banco SQLite

Abra o arquivo `client-server-api.db` com o cliente `sqlite3` ou uma extensão do VSCode.

```bash
sqlite3 client-server-api.db "SELECT id, code, bid, datetime FROM cotacao ORDER BY id DESC LIMIT 10;"
```