# Sistema de Emissão de Notas Fiscais (Korp)

Este projeto é uma solução de faturamento e estoque baseada em microsserviços, desenvolvida como desafio técnico para a Korp. A aplicação utiliza **Angular** no Frontend, **Go (Golang)** no Backend e **PostgreSQL** para persistência física.

## 🚀 Como Rodar o Projeto (Docker)

O projeto está totalmente conteinerizado. Certifique-se de ter o Docker e o Docker Compose instalados.

### 1. Subir a aplicação completa
Para iniciar todos os microsserviços, banco de dados e o frontend:
```bash
docker compose up -d
```
Acesse em: `http://localhost:4200`

### 2. Popular o banco com dados de teste (Opcional)
Implementamos uma semente automática que gera 100 produtos, 100 faturas e 10 clientes. Para usá-la:
1. Altere no `docker-compose.yml` a variável `SEED_DATABASE` para `"true"`.
2. Derrube os volumes e suba novamente:
```bash
docker compose down -v
docker compose up -d
```

---

## 🛠️ Operações de Resiliência (Teste de Falhas)

Como solicitado nos requisitos, o sistema é capaz de lidar com a queda de microsserviços. Você pode simular falhas usando os comandos abaixo:

### Matar um serviço específico (ex: Estoque)
```bash
docker compose stop estoque_backend
```
*O frontend exibirá um aviso de que o serviço de estoque está offline e desabilitará funções dependentes.*

### Subir um serviço específico
```bash
docker compose start estoque_backend
```

### Reiniciar um serviço após alteração de código
```bash
docker compose up --build -d faturamento_backend
```

---

## 🧪 Como Rodar os Testes

Os testes validam as regras de negócio críticas, como a proibição de fechar notas sem estoque.

### Testes do Backend (Go)
Existem suítes de teste para cada microsserviço na pasta `/Testes/Teste Back-end`.
```bash
# Exemplo: Testando Faturamento
cd Testes/Teste Back-end/faturamento
go test -v faturamento_test.go
```

### Testes do Frontend (Angular)
```bash
cd Front-end
npm test
```

---

## 🏗️ Arquitetura e Funcionalidades

### Microsserviços
- **Serviço de Estoque (Porta 8081):** Cadastro de produtos e controle de saldos.
- **Serviço de Faturamento (Porta 8082):** Gestão de notas fiscais e integração.
- **Serviço de Clientes (Porta 8083):** Cadastro de clientes para faturamento.
- **Serviço de Impressão (Porta 8084):** Geração dinâmica de PDFs.

### Principais Funcionalidades
1. **Cadastro e Ajuste de Estoque:** Além do cadastro, permite o ajuste manual de saldo via modal.
2. **Fluxo de Nota Fiscal:** Abertura de nota -> Inclusão de Itens (mesmo sem estoque) -> Impressão (validação final de estoque).
3. **Destaque de Ruptura:** Notas com itens faltando estoque são destacadas em vermelho na listagem principal.
4. **Visualização Detalhada:** Botão "VER ITENS" permite conferir o conteúdo da nota antes de fechar.
5. **IA Insights:** Painel preditivo que avisa quando produtos estão prestes a acabar com base no histórico de vendas.

---

## 🛑 Encerrar o Projeto
Para parar todos os serviços:
```bash
docker compose down
```
Para apagar também os dados salvos no banco:
```bash
docker compose down -v
```
