import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { forkJoin, of } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { InvoiceService } from '../../core/services/invoice';
import { ProductService } from '../../core/services/product';
import { ClientService } from '../../core/services/client';
import { Invoice, InvoiceItem } from '../../core/models/invoice.model';
import { Product } from '../../core/models/product.model';
import { Client } from '../../core/models/client.model';
import { ModalComponent, ModalType } from '../../shared/components/modal/modal';

@Component({
  selector: 'app-faturamento',
  standalone: true,
  imports: [CommonModule, FormsModule, ModalComponent],
  templateUrl: './faturamento.html',
  styleUrl: './faturamento.css'
})
export class FaturamentoComponent implements OnInit {
  invoices: Invoice[] = [];
  products: Product[] = [];
  clients: Client[] = [];

  selectedClientId: number | null = null;
  selectedClient: Client | null = null;
  selectedItems: InvoiceItem[] = [];
  
  currentProductId: number | null = null;
  currentQuantity: number = 1;

  isProcessing = false;
  isModalOpen = false;
  modalTitle = '';
  modalMessage = '';
  modalType: ModalType = 'success';
  
  hasClientError = false;
  hasProductError = false;
  hasInvoiceError = false;

  constructor(
    private invoiceService: InvoiceService,
    private productService: ProductService,
    private clientService: ClientService
  ) {}

  ngOnInit(): void {
    console.log('Inicializando Faturamento...');
    this.loadData();
  }

  onClientChange(): void {
    this.selectedClient = this.clients.find(c => c.id === Number(this.selectedClientId)) || null;
  }

  loadData(): void {
    this.isProcessing = true;
    this.hasClientError = false;
    this.hasProductError = false;
    this.hasInvoiceError = false;

    forkJoin({
      clients: this.clientService.getClients().pipe(catchError(() => { this.hasClientError = true; return of([]); })),
      invoices: this.invoiceService.getInvoices().pipe(catchError(() => { this.hasInvoiceError = true; return of([]); })),
      products: this.productService.getProducts().pipe(catchError(() => { this.hasProductError = true; return of([]); }))
    }).subscribe({
      next: (result) => {
        this.clients = result.clients || [];
        this.products = result.products || [];
        this.invoices = (result.invoices || []).map(inv => ({
          ...inv,
          cliente: this.clients.find(c => c.id === inv.cliente_id)
        }));
        this.isProcessing = false;

        if (this.hasClientError || this.hasProductError || this.hasInvoiceError) {
          const serviceList = [];
          if (this.hasClientError) serviceList.push('CLIENTES');
          if (this.hasProductError) serviceList.push('ESTOQUE');
          if (this.hasInvoiceError) serviceList.push('FATURAMENTO');
          
          this.showModal('Aviso de Indisponibilidade', `Não foi possível conectar aos serviços: ${serviceList.join(', ')}.`, 'error');
        }
      },
      error: () => {
        this.isProcessing = false;
        this.showModal('Erro Crítico', 'Falha ao processar dados do sistema.', 'error');
      }
    });
  }

  addItem(): void {
    if (!this.currentProductId || this.currentQuantity <= 0) return;
    const product = this.products.find(p => p.id === Number(this.currentProductId));
    if (!product) return;

    const existingItemIndex = this.selectedItems.findIndex(item => item.item_id === product.id);
    if (existingItemIndex > -1) {
      const totalQuantity = this.selectedItems[existingItemIndex].quantidade + this.currentQuantity;
      if (totalQuantity > product.saldo) {
        this.showModal('Aviso', `Saldo insuficiente para acumular ${product.descricao}.`, 'error');
        return;
      }
      this.selectedItems[existingItemIndex].quantidade = totalQuantity;
      this.selectedItems[existingItemIndex].subtotal = totalQuantity * product.preco_base;
    } else {
      if (this.currentQuantity > product.saldo) {
        this.showModal('Aviso', `Saldo insuficiente para ${product.descricao}.`, 'error');
        return;
      }
      this.selectedItems.push({
        item_id: product.id!,
        quantidade: this.currentQuantity,
        preco_unitario: product.preco_base,
        subtotal: this.currentQuantity * product.preco_base,
        product: product
      });
    }
    this.currentProductId = null;
    this.currentQuantity = 1;
  }

  calculateTotal(): number {
    return this.selectedItems.reduce((acc, item) => acc + item.subtotal, 0);
  }

  saveInvoice(): void {
    if (this.hasInvoiceError) {
      this.showModal('Erro', 'Serviço de faturamento offline.', 'error');
      return;
    }
    if (!this.selectedClientId || this.selectedItems.length === 0) {
      this.showModal('Aviso', 'Selecione cliente e itens.', 'error');
      return;
    }

    const newInvoice: Invoice = {
      cliente_id: Number(this.selectedClientId),
      status: 'ABERTA',
      valor_total: this.calculateTotal(),
      itens: this.selectedItems
    };

    this.isProcessing = true;
    this.invoiceService.createInvoice(newInvoice).subscribe({
      next: () => {
        this.isProcessing = false;
        this.showModal('Sucesso', 'Nota Fiscal aberta com sucesso!', 'success');
        this.resetForm();
        this.loadData();
      },
      error: () => {
        this.isProcessing = false;
        this.showModal('Erro', 'Falha ao abrir nota fiscal.', 'error');
      }
    });
  }

  printInvoice(invoiceId: number): void {
    console.log('[DEBUG-FRONT] Solicitando impressao para Fatura ID:', invoiceId);
    this.isProcessing = true;
    
    this.invoiceService.printInvoice(invoiceId).subscribe({
      next: (blob) => {
        console.log('[DEBUG-FRONT] PDF recebido. Tamanho:', blob.size, 'bytes');
        const url = window.URL.createObjectURL(blob);
        window.open(url, '_blank');
        this.isProcessing = false;
        this.showModal('Sucesso', 'Nota Fiscal impressa e fechada com sucesso!', 'success');
        this.loadData();
      },
      error: async (err) => {
        console.error('[DEBUG-FRONT] Erro na resposta de impressao:', err);
        this.isProcessing = false;
        if (err.error instanceof Blob) {
          const text = await err.error.text();
          console.error('[DEBUG-FRONT] Conteudo do erro (Blob):', text);
          const errorData = JSON.parse(text);
          this.showModal('Erro de Resiliência', errorData.error, 'error');
        } else {
          this.showModal('Erro', 'Falha ao processar impressão.', 'error');
        }
      }
    });
  }

  resetForm(): void {
    this.selectedClientId = null;
    this.selectedClient = null;
    this.selectedItems = [];
    this.currentProductId = null;
    this.currentQuantity = 1;
  }

  showModal(title: string, message: string, type: ModalType): void {
    this.modalTitle = title;
    this.modalMessage = message;
    this.modalType = type;
    this.isModalOpen = true;
  }

  closeModal(): void {
    this.isModalOpen = false;
  }
}
