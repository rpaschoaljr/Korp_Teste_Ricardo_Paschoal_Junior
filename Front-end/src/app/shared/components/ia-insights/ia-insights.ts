import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { IAService } from '../../../core/services/ia';
import { IAInsight } from '../../../core/models/ia.model';

@Component({
  selector: 'app-ia-insights',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './ia-insights.html',
  styleUrl: './ia-insights.css'
})
export class IAInsightsComponent implements OnInit {
  insights: IAInsight[] = [];
  isOpen = false;
  loading = false;

  constructor(private iaService: IAService) {}

  ngOnInit(): void {
    this.loadInsights();
  }

  toggle(): void {
    this.isOpen = !this.isOpen;
    if (this.isOpen) {
      this.loadInsights();
    }
  }

  loadInsights(): void {
    this.loading = true;
    this.iaService.getInsights().subscribe({
      next: (data) => {
        this.insights = data;
        this.loading = false;
      },
      error: () => {
        this.loading = false;
      }
    });
  }
}
