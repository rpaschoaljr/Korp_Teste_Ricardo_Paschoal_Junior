import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { IAService } from '../../../core/services/ia';
import { IAInsight } from '../../../core/models/ia.model';
import { timer, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { HealthCheckService } from '../../../core/services/health-check';

@Component({
  selector: 'app-ia-insights',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './ia-insights.html',
  styleUrl: './ia-insights.css'
})
export class IAInsightsComponent implements OnInit, OnDestroy {
  insights: IAInsight[] = [];
  isOpen = false;
  loading = false;
  hasError = false;
  
  private destroy$ = new Subject<void>();

  constructor(
    private iaService: IAService,
    private healthService: HealthCheckService
  ) {}

  ngOnInit(): void {
    this.loadInsights();
    this.startHealthMonitoring();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  toggle(): void {
    this.isOpen = !this.isOpen;
    if (this.isOpen) {
      this.loadInsights();
    }
  }

  startHealthMonitoring(): void {
    this.healthService.health$
      .pipe(takeUntil(this.destroy$))
      .subscribe(health => {
        if (this.hasError !== !health.ia) {
          const wasOffline = this.hasError;
          this.hasError = !health.ia;

          if (wasOffline && !this.hasError && this.isOpen) {
            console.log('[IA] Serviço restaurado. Atualizando insights...');
            this.loadInsights(true);
          }
        }
      });
  }

  loadInsights(isSilent: boolean = false): void {
    if (!isSilent) this.loading = true;
    this.hasError = false;

    this.iaService.getInsights().subscribe({
      next: (data) => {
        this.insights = data;
        this.loading = false;
        this.hasError = false;
      },
      error: () => {
        this.loading = false;
        this.hasError = true;
      }
    });
  }
}
