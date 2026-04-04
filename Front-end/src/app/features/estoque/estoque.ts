import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ProductService } from '../../core/services/product';
import { Product } from '../../core/models/product.model';
import { UppercaseDirective } from '../../shared/directives/uppercase';
import { ModalComponent, ModalType } from '../../shared/components/modal/modal';

@Component({
  selector: 'app-estoque',
  standalone: true,
  imports: [CommonModule, FormsModule, UppercaseDirective, ModalComponent],
  templateUrl: './estoque.html',
  styleUrl: './estoque.css'
})
export class EstoqueComponent implements OnInit {
  products: Product[] = [];
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
  
  selectedProduct: Product | null = null;
  adjustmentValue: number = 0;

  constructor(private productService: ProductService) {}

  ngOnInit(): void {
    console.log('Inicializando Estoque...');
    this.loadProducts();
  }

  loadProducts(): void {
    this.productService.getProducts().subscribe({
      next: (data) => {
        this.products = data || [];
        this.generateNextCode();
      },
      error: (err) => this.showModal('Erro', 'Falha ao carregar produtos. O serviço de estoque está online?', 'error')
    });
  }

  generateNextCode(): void {
    if (this.products.length === 0) {
      this.newProduct.codigo = 'PROD-001';
      return;
    }

    // Extrai os números dos códigos atuais (formato PROD-XXX)
    const codes = this.products
      .map(p => p.codigo.replace('PROD-', ''))
      .map(num => parseInt(num, 10))
      .filter(num => !isNaN(num));

    if (codes.length === 0) {
      this.newProduct.codigo = 'PROD-001';
      return;
    }

    // Pega o maior número e incrementa
    const nextNum = Math.max(...codes) + 1;
    this.newProduct.codigo = `PROD-${nextNum.toString().padStart(3, '0')}`;
  }

  saveProduct(): void {
    if (!this.newProduct.codigo || !this.newProduct.descricao || this.newProduct.saldo < 0) {
      this.showModal('Aviso', 'Preencha todos os campos obrigatórios corretamente.', 'error');
      return;
    }

    this.productService.createProduct(this.newProduct).subscribe({
      next: () => {
        this.showModal('Sucesso', 'Produto cadastrado com sucesso!', 'success');
        this.resetForm();
        this.loadProducts(); // Recarrega e gera o próximo código
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
