import { Directive, HostListener, ElementRef, Renderer2, Optional } from '@angular/core';
import { NgControl } from '@angular/forms';

@Directive({
  selector: '[appUppercase]',
  standalone: true
})
export class UppercaseDirective {
  constructor(
    private el: ElementRef,
    private renderer: Renderer2,
    @Optional() private ngControl: NgControl
  ) {}

  @HostListener('input', ['$event'])
  onInput(event: any): void {
    const input = event.target as HTMLInputElement;
    const originalValue = input.value;
    
    // Converte para MAIÚSCULO e remove ACENTOS
    const normalizedValue = originalValue
      .toUpperCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '');

    if (originalValue !== normalizedValue) {
      // Atualiza o valor visual no input
      this.renderer.setProperty(this.el.nativeElement, 'value', normalizedValue);
      
      // Atualiza o valor no model do Angular (FormControl/ngModel)
      if (this.ngControl && this.ngControl.control) {
        this.ngControl.control.setValue(normalizedValue, { emitEvent: false });
      }
    }
  }
}
