# Contexto do Front-end: Sistema de Emissão de Notas Fiscais (Korp)

Este documento é a base para o desenvolvimento do Front-end. Siga rigorosamente as instruções abaixo.

## 🚀 Stack e Workflow
- **Framework:** Angular (Versão mais recente: 21.x.x ou superior).
- **Regra de Ouro:** NUNCA crie arquivos de componentes, serviços, diretivas ou módulos manualmente. Utilize obrigatoriamente os comandos do Angular CLI:
  - `ng generate component features/nome-componente`
  - `ng generate service core/services/nome-servico`
  - `ng generate directive shared/directives/nome-diretiva`
- **Ambiente:** Execução via Docker (Node 20+).

## 🏗️ Estrutura de Pastas
- `app/core/`: Serviços globais, interceptors e modelos de dados.
- `app/shared/`: Componentes reutilizáveis (modais, tabelas, autocomplete) e diretivas.
- `app/features/`: Telas principais (Estoque, Faturamento).

## 🛡️ Segurança e Padronização (A Regra do "FOGAO")
- **Normalização:** Todo texto (Descrição, Nome) deve ser convertido para **MAIÚSCULO** e ter **ACENTOS REMOVIDOS** em tempo real via diretiva `[appUppercase]`.
- **Sanitização:** Bloquear caracteres especiais maliciosos e tags HTML nos inputs.
- **Validação:** Impedir saldos negativos e garantir campos obrigatórios antes de habilitar botões de envio.
- **Documentos:** CPF, CNPJ e Telefone devem ser limpos (apenas números) antes do envio para a API.

## 💎 Interface e UX (Modais e Temas)
- **Alertas:** PROIBIDO o uso de `alert()` ou `confirm()` nativos. Use componentes de Modal próprios para:
  - Erros técnicos (exibir mensagem real do microsserviço).
  - Confirmação de cancelamento (se houver dados preenchidos).
  - Confirmação de impressão.
- **Temas:** Suporte nativo a Light e Dark Mode, respeitando a preferência do navegador e permitindo troca manual com persistência em `localStorage`.
- **Resiliência:** Se a API falhar ao salvar, o formulário NÃO deve ser resetado. Os dados devem permanecer na tela para nova tentativa.

## 🧪 Estratégia de Testes (Vitest)
- **Unitários:** Devem cobrir a lógica de normalização ("fogão" -> "FOGAO"), cálculos de subtotal/total e bloqueios de segurança.
- **Mocking:** Usar `vi.spyOn` para simular as respostas dos microsserviços (Sucesso, Erro 500, Saldo Insuficiente).
- **Integração:** Testar a leitura real dos dados vindos dos microsserviços rodando no Docker (Ports 8081, 8082, 8083).

---
**IMPORTANTE:** Antes de qualquer `build`, verifique se o `tsconfig.json` e o `angular.json` estão apontando para as versões corretas de TypeScript e Builders suportados.
