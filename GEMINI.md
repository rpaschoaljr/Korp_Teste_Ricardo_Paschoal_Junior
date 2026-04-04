# Contexto do Projeto: Sistema de Emissão de Notas Fiscais (Korp)

## Objetivo Principal
Objetivo
Desenvolver uma aplicação em Angular, conforme os requisitos descritos abaixo, e
apresentar os resultados em formato de vídeo, demonstrando:
● As telas desenvolvidas;
● As funcionalidades implementadas;
● Um detalhamento técnico da solução.
No detalhamento técnico, informar:
● Quais ciclos de vida do Angular foram utilizados;
● Se foi feito uso da biblioteca RxJS e, em caso afirmativo, como;
● Quais outras bibliotecas foram utilizadas e para qual finalidade;
● Para componentes visuais, quais bibliotecas foram utilizadas;
● Como foi realizado o gerenciamento de dependências no Golang (se
aplicável);
● Quais frameworks foram utilizados no Golang ou C#;
● Como foram tratados os erros e exceções no backend;
● Caso a implementação utilize C#, indicar se foi utilizado LINQ e de que forma.
Escopo
1. Funcionalidades a serem desenvolvidas
Cadastro de Produtos
Campos obrigatórios:
● Código
● Descrição (nome do produto)
● Saldo (quantidade disponível em estoque)
Resultado esperado: permitir que um produto seja previamente cadastrado
para posterior utilização em notas fiscais.
Cadastro de Notas Fiscais
Campos obrigatórios:
● Numeração sequencial
● Status: Aberta ou Fechada
● Inclusão de múltiplos produtos com respectivas quantidades
Resultado esperado: permitir a criação de uma nota fiscal com numeração sequencial e
status inicial Aberta.
Impressão de Notas Fiscais
● Botão de impressão visível e intuitivo em tela.
Resultado esperado:
● Ao clicar no botão, exibir indicador de processamento;
● Após finalização, atualizar o status da nota para Fechada;
● Não permitir a impressão de notas com status diferente
de Aberta;
● Atualizar o saldo dos produtos conforme a quantidade
utilizada na nota.
○ Exemplo: saldo anterior = 10; nota utiliza 2 unidades → novo saldo =
8.
Requisitos obrigatórios
1. Arquitetura de Microsserviços:
Estruturar o sistema com no mínimo dois microsserviços:
● Serviço de Estoque – controle de produtos e saldos;
● Serviço de Faturamento – gestão de notas fiscais.
2. Tratamento de Falhas:
Implementar um cenário em que um dos microsserviços falha.
O sistema deve ser capaz de se recuperar da falha e fornecer
feedback apropriado ao usuário sobre o erro.
3. Conexão Real com banco de dados:
É esperado que os cadastros sejam persistidos fisicamente em um banco
de dados de sua escolha.
Requisitos opcionais
O candidato poderá, a seu critério, implementar também:
a. Tratamento de Concorrência:
Cenário: produto com saldo 1 sendo utilizado simultaneamente por duas notas.
b. Uso de Inteligência Artificial:
Implementar alguma funcionalidade do sistema que utilize IA.
c. Implementação de Idempotência:
Garantir que operações repetidas não causem efeitos colaterais indesejados.
Orientações para entrega
1. O projeto deverá ser disponibilizado em um repositório público do
GitHub com o nome: Korp_Teste_SeuNome
2. Após finalizar o desenvolvimento, o candidato deve enviar por e-mail, em
até 7 dias corridos após o recebimento deste desafio, o material concluído.
(Envie para julia.canever@korp.com.br)
3. O envio deve conter obrigatoriamente:
● Link público do repositório GitHub Korp_Teste_SeuNome
● Link para vídeo de apresentação das telas e funcionalidades implementadas
(Google Drive, One Drive ou alguma outra nuvem da sua preferência)
● Detalhamento técnico conforme itens descritos no documento de especificação
deste teste

## Modo de agir
-Nunca deve gerar o codigo antes de explicar a logica que irá aplicar
-Sempre que for criar algo deve esperar uma resposta de confirmação
-Sempre deve roda e testar os codigo e modificações antes de implementar

## Stack Tecnológica
- **Frontend:** Angular
- **Backend:** Golang
- **Arquitetura:** Microsserviços (mínimo dois: Serviço de Estoque e Serviço de Faturamento)
- **Banco de Dados:** Persistência física real em PostgreSQL
- **Comunicação:** HTTP, utilize as mensagens do http para saber se deu certo ou não e ter uma resposta mais imediata

## Regras de Negócio e Funcionalidades (Obrigatórias)
1. **Cadastro de Produtos:** Campos obrigatórios: Código, Descrição, Saldo (quantidade)
2. **Cadastro de Notas Fiscais:** Campos obrigatórios: Numeração sequencial, Status (Aberta/Fechada), Itens (múltiplos produtos com quantidades)
3. **Impressão de Nota:**
   - Botão visível com indicador de processamento
   - Só permitir se o status for "Aberta"
   - Após "imprimir", mudar o status para "Fechada"
   - **Crucial:** Atualizar automaticamente o saldo dos produtos no estoque ao fechar a nota.

## Diretrizes Técnicas
- **Tratamento de Erros:** Implementar resiliência para falha de microsserviços (ex: simular queda do serviço de estoque) com feedback claro ao usuário no frontend
- **Backend (Go):** Explicar como os erros e exceções são tratados, como foi feito o gerenciamento de dependências (`go mod`) e quais frameworks foram usados
- **Frontend (Angular):** Usar RxJS para reatividade e detalhar os ciclos de vida utilizados (OnInit, OnDestroy, etc.), além de listar as bibliotecas visuais
- **Testes Automáticos:** Ao gerar código em Go ou TS, sempre criar a estrutura de testes automáticos para que o desenvolvedor possa validar a lógica localmente.

## Diferenciais Opcionais
- **Tratamento de Concorrência:** Impedir que duas notas usem o mesmo saldo simultaneamente
- **Uso de Inteligência Artificial:** Implementar alguma funcionalidade que utilize IA
- **Idempotência:** Garantir que operações repetidas (ex: múltiplos cliques no botão imprimir) não causem efeitos colaterais indesejados

## Regras de Entrega
- **Repositório:** Código no GitHub público com o nome `Korp_Teste_Paschoal`
- **Apresentação:** Gravar e enviar link de um vídeo demonstrando as telas, as funcionalidades e o detalhamento técnico
- **Prazo e Contato:** Enviar o material concluído em até 7 dias corridos para `julia.canever@korp.com.br`

## Comandos de Execução
- **Rodar Backend (Estoque):** `cd estoque && go run main.go`
- **Rodar Backend (Faturamento):** `cd faturamento && go run main.go`
- **Rodar Frontend:** `ng serve`
- **Rodar Testes:** `ng test` (Angular) e `go test ./...` (Golang)

 1. Arquitetura do docker-compose.yml
  O arquivo definirá 4 serviços distintos, permitindo que você controle cada um separadamente:
   * db (PostgreSQL): O banco de dados central.
       * Persistência: Usarei um Docker Volume nomeado (ex: pg_data:/var/lib/postgresql/data). Isso garante que, mesmo que você destrua o contêiner do banco, os dados continuam salvos no seu
         computador.
   * estoque (Golang): O microserviço de estoque.
   * faturamento (Golang): O microserviço de faturamento.
   * frontend (Angular): A aplicação visual.

  2. Como "Matar" e Voltar com os Microserviços
  Como os serviços estarão separados no docker-compose.yml, você terá o controle exato que pediu usando comandos nativos do Docker no terminal.
  Por exemplo, para testar como o Front-end reage se o Estoque cair, você fará:
   * Para matar o serviço: docker compose stop estoque
   * Para voltar com ele: docker compose start estoque
  (Isso afeta apenas o Estoque; o banco, o faturamento e o front continuarão rodando normalmente).

  3. Lógica dos Scripts de Banco de Dados
  Para atender ao seu pedido de dois scripts distintos:

   * Script 1: Inicialização Normal (Estrutura)
       * Criarei um arquivo 01-init-schema.sql.
       * Ele conterá comandos como CREATE TABLE IF NOT EXISTS... para as tabelas de itens e faturas.
       * Como roda: O PostgreSQL do Docker é programado para rodar automaticamente qualquer script colocado na pasta /docker-entrypoint-initdb.d/ na primeira vez que o contêiner é criado.
         Assim, a estrutura sempre nasce pronta.

   * Script 2: Popular Dados (Rodado apenas uma vez, se vazio)
       * Criarei um arquivo 02-seed-data.sql.
       * A lógica de "SE estiver vazio" será tratada dentro do próprio SQL. Usarei comandos como INSERT INTO itens (...) ON CONFLICT DO NOTHING; ou farei uma verificação condicional (IF NOT
         EXISTS (SELECT 1 FROM itens) THEN INSERT...).
       * Dessa forma, mesmo que o script tente rodar novamente ou você o execute manualmente, ele será inteligente o suficiente para não duplicar dados se o banco já tiver conteúdo.
       * Ele também será colocado na pasta de inicialização para rodar automaticamente na primeira vez, mas pela sua natureza condicional, não causará problemas nos re-starts subsequentes.
         (Ou, se preferir controle total manual, podemos deixá-lo fora da pasta de auto-init para você rodar via comando quando desejar).
  1. Tabela clientes
  Permanece igual, guardando os dados de contato e os documentos para a nota fiscal.

   1 CREATE TABLE clientes (
   2     id SERIAL PRIMARY KEY,
   3     nome VARCHAR(255) NOT NULL,
   4     telefone VARCHAR(20),
   5     endereco TEXT,
   6     cpf VARCHAR(14),
   7     cnpj VARCHAR(18)
   8 );

  2. Tabela itens (Seu Cadastro de Estoque)
  Atualizada com o Código único, o Saldo e o Preço Base.

   1 CREATE TABLE itens (
   2     id SERIAL PRIMARY KEY,
   3     codigo VARCHAR(50) UNIQUE NOT NULL, -- Identificador Textual (ex: PROD-001)
   4     descricao VARCHAR(255) NOT NULL,    -- Descrição do produto
   5     saldo INT NOT NULL DEFAULT 0,       -- Quantidade disponível
   6     preco_base DECIMAL(10, 2) NOT NULL  -- O preço padrão de catálogo
   7 );

  3. Tabela faturas (Cabeçalho da Nota)
  A "capa" da nota fiscal, agora simplificada pois os itens foram movidos para a tabela filha, mas mantendo o valor total pronto para leitura rápida.

   1 CREATE TABLE faturas (
   2     id SERIAL PRIMARY KEY,                  -- Número Sequencial da Fatura
   3     cliente_id INT REFERENCES clientes(id), -- A quem a fatura pertence
   4     status VARCHAR(20) DEFAULT 'ABERTA',    -- 'ABERTA' (editável) ou 'FECHADA' (impressa/concluída)
   5     valor_total DECIMAL(10, 2) NOT NULL DEFAULT 0.00, -- Soma de todos os sub-totais dos itens
   6     data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   7 );

  4. Nova Tabela itens_fatura (Corpo da Nota)
  A tabela intermediária de relação (1 Fatura -> N Itens).

   1 CREATE TABLE itens_fatura (
   2     id SERIAL PRIMARY KEY,
   3     fatura_id INT REFERENCES faturas(id) ON DELETE CASCADE, -- Se a fatura for deletada, os itens somem
   4     item_id INT REFERENCES itens(id),                       -- Referência ao produto no estoque
   5     quantidade INT NOT NULL,                                -- Quantidade vendida NESTA nota
   6     preco_unitario DECIMAL(10, 2) NOT NULL,                 -- Preço que foi cobrado no momento da venda
   7     subtotal DECIMAL(10, 2) NOT NULL                        -- quantidade * preco_unitario
   8 );

