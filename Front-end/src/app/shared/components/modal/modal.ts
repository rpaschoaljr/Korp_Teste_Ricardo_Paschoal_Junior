import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';

export type ModalType = 'success' | 'error' | 'confirm';

@Component({
  selector: 'app-modal',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './modal.html',
  styleUrl: './modal.css'
})
export class ModalComponent {
  @Input() title: string = '';
  @Input() message: string = '';
  @Input() type: ModalType = 'success';
  @Input() isOpen: boolean = false;

  @Output() closeEvent = new EventEmitter<void>();
  @Output() confirmEvent = new EventEmitter<void>();

  onClose(): void {
    this.closeEvent.emit();
  }

  onConfirm(): void {
    this.confirmEvent.emit();
  }
}
