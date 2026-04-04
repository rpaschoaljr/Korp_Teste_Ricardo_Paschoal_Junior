import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-service-status-banner',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="status-banner" *ngIf="isOffline">
      <span class="icon">⚠</span>
      <span class="message">{{ message }}</span>
    </div>
  `,
  styles: [`
    .status-banner {
      background-color: var(--warning-bg);
      color: var(--warning-text);
      border: 1px solid var(--border-color);
      padding: 1rem;
      border-radius: 4px;
      margin: 1rem 0;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      font-weight: bold;
      font-size: 0.9rem;
      box-shadow: 0 2px 4px rgba(0,0,0,0.05);
      animation: fadeIn 0.3s ease-in-out;
    }

    .icon {
      font-size: 1.2rem;
    }

    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(-10px); }
      to { opacity: 1; transform: translateY(0); }
    }
  `]
})
export class ServiceStatusBannerComponent {
  @Input() isOffline: boolean = false;
  @Input() message: string = 'SERVIÇO OFFLINE. AGUARDE A RECONEXÃO.';
}
