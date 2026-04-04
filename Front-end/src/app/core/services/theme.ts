import { Injectable, Renderer2, RendererFactory2, Inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private renderer: Renderer2;
  private currentTheme: 'light' | 'dark' = 'light';
  private readonly STORAGE_KEY = 'korp-theme';

  constructor(
    rendererFactory: RendererFactory2,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {
    this.renderer = rendererFactory.createRenderer(null, null);
    this.initTheme();
  }

  private initTheme(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    const savedTheme = localStorage.getItem(this.STORAGE_KEY) as 'light' | 'dark' | null;
    
    if (savedTheme) {
      this.currentTheme = savedTheme;
    } else {
      // Respeita a preferência do navegador
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      this.currentTheme = prefersDark ? 'dark' : 'light';
    }

    this.applyTheme();
  }

  toggleTheme(): void {
    this.currentTheme = this.currentTheme === 'light' ? 'dark' : 'light';
    localStorage.setItem(this.STORAGE_KEY, this.currentTheme);
    this.applyTheme();
  }

  get isDark(): boolean {
    return this.currentTheme === 'dark';
  }

  private applyTheme(): void {
    if (!isPlatformBrowser(this.platformId)) return;

    const body = document.body;
    this.renderer.removeClass(body, 'light-theme');
    this.renderer.removeClass(body, 'dark-theme');
    this.renderer.addClass(body, `${this.currentTheme}-theme`);
  }
}
