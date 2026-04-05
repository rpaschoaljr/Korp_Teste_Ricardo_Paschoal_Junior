import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { forkJoin, of, timer, Subject } from 'rxjs';
import { catchError, takeUntil } from 'rxjs/operators';
import { InvoiceService } from '../../core/services/invoice';
import { ProductService } from '../../core/services/product';
import { ClientService } from '../../core/services/client';
import { Invoice, InvoiceItem } from '../../core/models/invoice.model';
import { Product } from '../../core/models/product.model';
import { Client } from '../../core/models/client.model';
import { ModalComponent, ModalType } from '../../shared/components/modal/modal';
import { HealthCheckService } from '../../core/services/health-check';
import { ServiceStatusBannerComponent } from '../../shared/components/status-banner/status-banner';

@Component({
  selector: 'app-faturamento',
  standalone: true,
  imports: [CommonModule, FormsModule, ModalComponent, ServiceStatusBannerComponent],
  templateUrl: './faturamento.html',
  styleUrl: './faturamento.css'
})
export class FaturamentoComponent implements OnInit, OnDestroy {
  invoices: Invoice[] = [];
  products: Product[] = [];
  clients: Client[] = [];
  
  private destroy$ = new Subject<void>();

  selectedClientId: number | null = null;
  selectedClient: Client | null = null;
  selectedItems: InvoiceItem[] = [];
  
  currentProductId: number | null = null;
  currentQuantity: number = 1;

  isProcessing = false;
  isModalOpen = false;
  isViewModalOpen = false;
  modalTitle = '';
  modalMessage = '';
  modalType: ModalType = 'success';
  
  selectedInvoice: Invoice | null = null;
  
  hasClientError = false;
  hasProductError = false;
  hasInvoiceError = false;

  constructor(
    private invoiceService: InvoiceService,
    private productService: ProductService,
    private clientService: ClientService,
    private healthService: HealthCheckService
  ) {}

  ngOnInit(): void {
    console.log('Inicializando Faturamento...');
    this.loadData();
    this.startHealthMonitoring();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  startHealthMonitoring(): void {
    this.healthService.health$
      .pipe(takeUntil(this.destroy$))
      .subscribe(health => {
        // Se houve mudança significativa nos serviços que este componente usa
        if (this.hasInvoiceError !== !health.faturamento || 
            this.hasProductError !== !health.estoque || 
            this.hasClientError !== !health.clientes) {
          
          const wasAnyOffline = this.hasInvoiceError || this.hasProductError || this.hasClientError;
          
          this.hasInvoiceError = !health.faturamento;
          this.hasProductError = !health.estoque;
          this.hasClientError = !health.clientes;

          const isNowAllOnline = !this.hasInvoiceError && !this.hasProductError && !this.hasClientError;

          // Se estávamos offline e agora tudo voltou, recarregamos
          if (wasAnyOffline && isNowAllOnline) {
            console.log('[FATURAMENTO] Serviços restaurados. Atualizando dados...');
            this.loadData(true);
          }
        }
      });
  }

  onClientChange(): void {
    this.selectedClient = this.clients.find(c => c.id === Number(this.selectedClientId)) || null;
  }

  loadData(isSilent: boolean = false): void {
    if (!isSilent) {
      this.isProcessing = true;
      this.hasClientError = false;
      this.hasProductError = false;
      this.hasInvoiceError = false;
    }

    forkJoin({
      clients: this.clientService.getClients().pipe(catchError(() => { this.hasClientError = true; return of(null); })),
      invoices: this.invoiceService.getInvoices().pipe(catchError(() => { this.hasInvoiceError = true; return of(null); })),
      products: this.productService.getProducts().pipe(catchError(() => { this.hasProductError = true; return of(null); }))
    }).subscribe({
      next: (result) => {
        if (result.clients !== null) {
          this.clients = result.clients;
          this.hasClientError = false;
        }
        if (result.products !== null) {
          this.products = result.products;
          this.hasProductError = false;
        }
        if (result.invoices !== null) {
          this.invoices = result.invoices.map(inv => ({
            ...inv,
            cliente: this.clients.find(c => c.id === inv.cliente_id)
          }));
          this.hasInvoiceError = false;
        }
        
        if (!isSilent) this.isProcessing = false;

        if (!isSilent && (this.hasClientError || this.hasProductError || this.hasInvoiceError)) {
          const serviceList = [];
          if (this.hasClientError) serviceList.push('CLIENTES');
          if (this.hasProductError) serviceList.push('ESTOQUE');
          if (this.hasInvoiceError) serviceList.push('FATURAMENTO');
          
          this.showModal('Aviso de Indisponibilidade', `Não foi possível conectar aos serviços: ${serviceList.join(', ')}.`, 'error');
        }
      },
      error: () => {
        if (!isSilent) {
          this.isProcessing = false;
          this.showModal('Erro Crítico', 'Falha ao processar dados do sistema.', 'error');
        }
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
      this.selectedItems[existingItemIndex].quantidade = totalQuantity;
      this.selectedItems[existingItemIndex].subtotal = totalQuantity * product.preco_base;
    } else {
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
      error: (err) => {
        this.isProcessing = false;
        const backendMessage = err.error?.error || 'Falha ao abrir nota fiscal.';
        this.showModal('Erro ao Faturar', backendMessage, 'error');
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

  viewInvoice(invoice: Invoice): void {
    this.selectedInvoice = invoice;
    this.isViewModalOpen = true;
  }

  closeViewModal(): void {
    this.isViewModalOpen = false;
    this.selectedInvoice = null;
  }

  getItemStock(productId: number): number {
    const product = this.products.find(p => p.id === productId);
    return product ? product.saldo : 0;
  }

  hasStockIssue(invoice: Invoice): boolean {
    if (!invoice.itens || invoice.status === 'FECHADA') return false;
    return invoice.itens.some(item => item.quantidade > this.getItemStock(item.item_id));
  }
}
