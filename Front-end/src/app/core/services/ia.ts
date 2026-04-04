import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { IAInsight } from '../models/ia.model';

@Injectable({
  providedIn: 'root'
})
export class IAService {
  private apiUrl = 'http://localhost:8082/faturas/ia/insights';

  constructor(private http: HttpClient) { }

  getInsights(): Observable<IAInsight[]> {
    return this.http.get<IAInsight[]>(this.apiUrl);
  }
}
