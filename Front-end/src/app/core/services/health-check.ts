import { Injectable, OnDestroy, NgZone } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { distinctUntilChanged } from 'rxjs/operators';

export interface SystemHealth {
  estoque: boolean;
  faturamento: boolean;
  clientes: boolean;
  ia: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class HealthCheckService implements OnDestroy {
  private eventSource?: EventSource;
  
  private healthState = new BehaviorSubject<SystemHealth>({
    estoque: true,
    faturamento: true,
    clientes: true,
    ia: true
  });

  public health$ = this.healthState.asObservable().pipe(
    distinctUntilChanged((a, b) => JSON.stringify(a) === JSON.stringify(b))
  );

  constructor(private ngZone: NgZone) {
    this.connectToHealthStream();
  }

  private connectToHealthStream() {
    console.log('[HEALTH] Conectando ao canal de eventos de saúde (SSE)...');
    
    this.eventSource = new EventSource('http://localhost:8082/health-stream');

    this.eventSource.onmessage = (event) => {
      // NgZone.run garante que o Angular perceba a mudança de dados vinda de fora (SSE)
      this.ngZone.run(() => {
        try {
          const health: SystemHealth = JSON.parse(event.data);
          this.healthState.next(health);
        } catch (e) {
          console.error('[HEALTH] Erro ao processar dados de saúde:', e);
        }
      });
    };

    this.eventSource.onerror = (error) => {
      this.ngZone.run(() => {
        console.warn('[HEALTH] Falha na conexão SSE. O serviço de faturamento pode estar offline.');
        // Se o canal cair, assumimos que o faturamento (que hospeda o SSE) caiu
        this.healthState.next({
          estoque: false,
          faturamento: false,
          clientes: false,
          ia: false
        });
      });
    };
  }

  ngOnDestroy() {
    if (this.eventSource) {
      this.eventSource.close();
    }
  }
}
