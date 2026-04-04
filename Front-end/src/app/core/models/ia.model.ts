export interface IAInsight {
  id: string;
  tipo: 'ESTOQUE' | 'VENDA';
  mensagem: string;
  prioridade: 'ALTA' | 'MEDIA' | 'BAIXA';
  acao_sugerida: string;
}
