import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { timer, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ProductService } from '../../core/services/product';
import { Product } from '../../core/models/product.model';
import { UppercaseDirective } from '../../shared/directives/uppercase';
import { ModalComponent, ModalType } from '../../shared/components/modal/modal';
import { HealthCheckService } from '../../core/services/health-check';
import { ServiceStatusBannerComponent } from '../../shared/components/status-banner/status-banner';

@Component({
  selector: 'app-estoque',
  standalone: true,
  imports: [CommonModule, FormsModule, UppercaseDirective, ModalComponent, ServiceStatusBannerComponent],
  templateUrl: './estoque.html',
  styleUrl: './estoque.css'
})
export class EstoqueComponent implements OnInit, OnDestroy {
  products: Product[] = [];
  private destroy$ = new Subject<void>();
  newProduct: Product = {
    codigo: '',
    descricao: '',
    saldo: 0,
    preco_base: 0
  };

  // Controle do Modal
  isModalOpen = false;
  isAdjustmentModalOpen = false;
  modalTitle = '';
  modalMessage = '';
  modalType: ModalType = 'success';

  hasStockServiceError = false;

  selectedProduct: Product | null = null;
  adjustmentValue: number = 0;

  constructor(
    private productService: ProductService,
    private healthService: HealthCheckService
  ) {}

  ngOnInit(): void {
    console.log('Inicializando Estoque...');
    this.loadProducts();
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
        // Este componente só se importa com o estoque
        if (this.hasStockServiceError !== !health.estoque) {
          const wasOffline = this.hasStockServiceError;
          this.hasStockServiceError = !health.estoque;

          if (wasOffline && !this.hasStockServiceError) {
            console.log('[ESTOQUE] Serviço restaurado. Atualizando dados...');
            this.loadProducts(true);
          }
        }
      });
  }

  loadProducts(isSilent: boolean = false): void {
    if (!isSilent) this.hasStockServiceError = false;
    this.productService.getProducts().subscribe({
      next: (data) => {
        this.products = data || [];
        this.hasStockServiceError = false;
        this.generateNextCode();
      },
      error: (err) => {
        this.hasStockServiceError = true;
        if (!isSilent) {
          this.showModal('Erro de Resiliência', 'O serviço de ESTOQUE está offline. Não é possível cadastrar ou ajustar produtos no momento.', 'error');
        }
      }
    });
  }

  generateNextCode(): void {
    if (this.products.length === 0) {
      this.newProduct.codigo = 'PROD-001';
      return;
    }

    const codes = this.products
      .map(p => p.codigo.replace('PROD-', ''))
      .map(num => parseInt(num, 10))
      .filter(num => !isNaN(num));

    if (codes.length === 0) {
      this.newProduct.codigo = 'PROD-001';
      return;
    }

    const nextNum = Math.max(...codes) + 1;
    this.newProduct.codigo = `PROD-${nextNum.toString().padStart(3, '0')}`;
  }

  saveProduct(): void {
    if (this.hasStockServiceError) {
      this.showModal('Operação Bloqueada', 'Não é possível salvar: o serviço de ESTOQUE está indisponível.', 'error');
      return;
    }

    if (!this.newProduct.codigo || !this.newProduct.descricao || this.newProduct.saldo < 0) {
      this.showModal('Aviso', 'Preencha todos os campos obrigatórios corretamente.', 'error');
      return;
    }

    this.productService.createProduct(this.newProduct).subscribe({
      next: () => {
        this.showModal('Sucesso', 'Produto cadastrado com sucesso!', 'success');
        this.resetForm();
        this.loadProducts();
      },
      error: (err) => this.showModal('Erro', 'Não foi possível salvar o produto.', 'error')
    });
  }

  resetForm(): void {
    this.newProduct = {
      codigo: '',
      descricao: '',
      saldo: 0,
      preco_base: 0
    };
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

  openAdjustmentModal(product: Product): void {
    this.selectedProduct = product;
    this.adjustmentValue = product.saldo;
    this.isAdjustmentModalOpen = true;
  }

  closeAdjustmentModal(): void {
    this.isAdjustmentModalOpen = false;
    this.selectedProduct = null;
  }

  saveAdjustment(): void {
    if (this.hasStockServiceError) {
      this.showModal('Erro', 'O serviço de estoque está offline.', 'error');
      return;
    }
    if (!this.selectedProduct || this.adjustmentValue < 0) return;

    this.productService.adjustStock(this.selectedProduct.id!, this.adjustmentValue).subscribe({
      next: () => {
        this.showModal('Sucesso', 'Estoque ajustado com sucesso!', 'success');
        this.closeAdjustmentModal();
        this.loadProducts();
      },
      error: () => this.showModal('Erro', 'Falha ao ajustar estoque.', 'error')
    });
  }
}
