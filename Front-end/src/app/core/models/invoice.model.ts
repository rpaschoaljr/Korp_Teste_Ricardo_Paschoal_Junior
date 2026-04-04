import { Client } from './client.model';
import { Product } from './product.model';

export type InvoiceStatus = 'ABERTA' | 'FECHADA';

export interface InvoiceItem {
  id?: number;
  fatura_id?: number;
  item_id: number;
  codigo_produto?: string; // Adicionado para exibição
  descricao?: string;      // Adicionado para exibição
  quantidade: number;
  preco_unitario: number;
  subtotal: number;
  product?: Product; // Campo opcional para facilitar a exibição
}

export interface Invoice {
  id?: number;
  cliente_id: number;
  status: InvoiceStatus;
  valor_total: number;
  data_criacao?: Date;
  cliente?: Client; // Campo opcional para facilitar a exibição
  itens?: InvoiceItem[];
}
