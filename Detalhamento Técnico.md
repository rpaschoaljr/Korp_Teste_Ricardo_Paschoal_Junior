# Detalhamento Técnico - Sistema de Emissão de Notas Fiscais (Korp)

Este documento descreve as escolhas arquiteturais, tecnologias e padrões de implementação utilizados no desenvolvimento do desafio técnico.

## 1. Frontend (Angular)

### Ciclos de Vida do Angular Utilizados
- **ngOnInit:** Utilizado para disparar o carregamento inicial de dados e iniciar o monitoramento de saúde dos serviços.
- **ngOnDestroy:** Implementado com `Subject` e `takeUntil` para garantir o encerramento limpo de conexões SSE e assinaturas de Observables, prevenindo vazamentos de memória.

### Uso da Biblioteca RxJS
A reatividade do sistema é baseada em RxJS:
- **Observable & BehaviorSubject:** Utilizados no `HealthCheckService` para manter um estado global e reativo da saúde de todos os microsserviços.
- **forkJoin:** Utilizado para realizar chamadas paralelas otimizadas aos serviços.
- **takeUntil:** Padronizado para gerenciamento de ciclo de vida de assinaturas.
- **distinctUntilChanged:** Utilizado para evitar re-renderizações desnecessárias na interface quando o estado de saúde não sofreu alterações reais.

### Componentes e Resiliência Visual
- **ServiceStatusBanner (Componente Compartilhado):** Criamos um componente único para exibir alertas de serviço offline, garantindo consistência visual (DRY - Don't Repeat Yourself).
- **Monitoramento Ativo (Real-time):** O sistema não exige atualização de página. Ele detecta quedas e recuperações de serviços em tempo real, habilitando/desabilitando botões e campos de forma dinâmica.
- **Offline UI:** O sistema permite que o usuário continue preenchendo formulários e montando notas mesmo se a conexão cair temporariamente, bloqueando apenas o salvamento final até a restauração do serviço.

### Estilização
- **Temas Dinâmicos (Dark/Light):** Implementado via variáveis CSS globais (`:root` e `body.dark-theme`). Toda a interface, incluindo modais e alertas, adapta-se automaticamente à preferência do sistema ou usuário.

---

## 2. Backend (Golang)

### Agregador de Saúde via SSE (Server-Sent Events)
Para garantir a escalabilidade do sistema (suportando cenários de milhões de usuários), implementamos um endpoint de **SSE** no serviço de faturamento:
- **Paradigma Push:** Em vez de milhares de clientes "perguntarem" ao servidor se ele está bem (Polling), o servidor mantém uma única conexão aberta e "empurra" atualizações de status apenas quando necessário.
- **Eficiência:** Reduz drasticamente a carga no processador e no banco de dados, eliminando o custo de abrir/fechar milhares de conexões HTTP por segundo.

### Funcionalidades e Regras de Negócio
- **Ajuste Manual de Estoque:** Endpoint `PUT /produtos/ajuste` para correções de inventário com validação de integridade.
- **Transações de Banco de Dados:** Uso rigoroso de `Begin()` e `Rollback()` para garantir que a baixa no estoque e o fechamento da nota ocorram de forma atômica.
- **Códigos HTTP Semânticos:** Uso de `422 Unprocessable Entity` para erros de regra de negócio (como saldo insuficiente) e `503` para falhas de infraestrutura.

### Diferenciais Técnicos
- **Trava de Linha (FOR UPDATE):** Proteção contra condições de corrida em ambientes concorrentes.
- **IA e BI Preditivo:** Algoritmo que cruza dados de faturamento e estoque para prever rupturas e sugerir ações de compra.
- **Semente de Dados Automatizada:** Script de inicialização Docker que permite popular o sistema com massa de teste (100+ registros) via variável de ambiente `SEED_DATABASE`.

---

## 3. Resumo de Arquitetura
O sistema utiliza uma arquitetura de **Microsserviços** conteinerizados com Docker. A separação de responsabilidades permite que o serviço de faturamento funcione mesmo que o de IA esteja offline, e que o frontend se recupere automaticamente de falhas de rede, seguindo as melhores práticas de sistemas distribuídos e resilientes.
