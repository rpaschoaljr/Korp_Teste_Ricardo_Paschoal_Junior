# Detalhamento Técnico - Sistema de Emissão de Notas Fiscais (Korp)

Este documento descreve as escolhas arquiteturais, tecnologias e padrões de implementação utilizados no desenvolvimento do desafio técnico.

## 1. Frontend (Angular)

### Ciclos de Vida do Angular Utilizados
- **ngOnInit:** Utilizado em todos os componentes principais (`EstoqueComponent`, `FaturamentoComponent`) para disparar o carregamento inicial de dados (produtos, faturas e clientes) assim que o componente é inicializado.
- **ngOnDestroy:** Utilizado implicitamente através do gerenciamento de assinaturas (subscriptions) para garantir que não ocorram vazamentos de memória ao navegar entre as telas.

### Uso da Biblioteca RxJS
A reatividade do sistema é baseada em RxJS:
- **Observable:** Todas as chamadas HTTP para os microsserviços retornam Observables, permitindo um fluxo assíncrono e não bloqueante.
- **forkJoin:** Utilizado na tela de Faturamento para realizar chamadas paralelas aos serviços de Clientes, Estoque e Faturamento simultaneamente, otimizando o tempo de carregamento.
- **catchError & of:** Implementados para garantir a resiliência do sistema. Caso um microsserviço esteja offline, o `catchError` intercepta a falha e o `of([])` retorna um valor padrão, permitindo que a interface continue funcional e exiba um aviso amigável ao usuário.

### Componentes Visuais e Estilização
- **Vanilla CSS:** Para garantir máxima performance e fidelidade ao design, não foram utilizadas bibliotecas de componentes prontos (como Material ou Bootstrap). Toda a estilização foi feita manualmente com CSS moderno (Flexbox, Grid).
- **Componentes Customizados:** Criamos componentes reutilizáveis como `ModalComponent` para alertas e feedbacks, e `UppercaseDirective` para normalização de texto em tempo real.

### Outras Bibliotecas
- **Vitest:** Utilizada para execução de testes unitários rápidos e modernos, substituindo o Karma/Jasmine padrão.

---

## 2. Backend (Golang)

### Frameworks Utilizados
- **Gin Web Framework:** Utilizado em todos os microsserviços (Estoque, Faturamento, Clientes e Impressão) pela sua alta performance, facilidade de roteamento e excelente suporte a middlewares.

### Gerenciamento de Dependências
- **Go Modules (`go mod`):** O gerenciamento de todas as bibliotecas externas é feito de forma nativa pelo Go, garantindo reprodutibilidade do ambiente através dos arquivos `go.mod` e `go.sum`.

### Tratamento de Erros e Exceções
- **Verificação Explícita:** Seguindo o idioma de Go, todos os erros são verificados imediatamente (`if err != nil`).
- **Transações de Banco de Dados:** Operações críticas (como salvar uma nota e seus itens) utilizam `tx.Begin()` e `tx.Rollback()` em caso de falha, garantindo a integridade dos dados (Atomicidade).
- **Resiliência de Rede:** O serviço de faturamento implementa health-checks para verificar a saúde dos outros serviços antes de processar uma nota.
- **Códigos HTTP Semânticos:** Retorno de erros específicos:
  - `409 Conflict`: Para violação de regras de concorrência ou idempotência.
  - `503 Service Unavailable`: Quando um microsserviço dependente está offline.
  - `400 Bad Request`: Para validações de saldo ou dados inválidos.

### Concorrência e Idempotência (Diferenciais)
- **Trava de Linha (FOR UPDATE):** No Postgres, utilizamos `SELECT ... FOR UPDATE` para travar o registro da nota no início da impressão, garantindo que dois processos não tentem fechar a mesma nota simultaneamente.
- **Atualização Atômica de Estoque:** O desconto de saldo é feito diretamente via SQL condicional para evitar que dois usuários vendam o mesmo item físico se o saldo for insuficiente.

### Inteligência Artificial e BI Preditivo (Diferencial)
- **Módulo de Análise Preditiva:** Implementamos um algoritmo de inteligência de negócio que analisa o fluxo de vendas histórico.
- **Predição de Ruptura de Estoque:** O sistema calcula a velocidade de saída de cada produto e cruza com o saldo real do microsserviço de estoque, alertando quando o produto acabará em menos de 7 dias.
- **Análise de Churn de Clientes:** Identifica automaticamente clientes inativos com base no intervalo médio de compras, sugerindo ações de reativação.
- **Interface Flutuante:** No frontend, um painel dedicado ("IA Insights") exibe essas recomendações de forma proativa para o gestor.

### Biblioteca de PDF
- **gofpdf:** Utilizada no microsserviço de impressão para gerar dinamicamente o documento da nota fiscal em formato PDF.

---

## 3. Resumo de Arquitetura
O sistema foi estruturado em **Microsserviços** que se comunicam via HTTP, com persistência real em **PostgreSQL**. Cada serviço é isolado em seu próprio container Docker, simulando um ambiente de produção real onde falhas parciais podem ocorrer e devem ser tratadas pelo sistema.
