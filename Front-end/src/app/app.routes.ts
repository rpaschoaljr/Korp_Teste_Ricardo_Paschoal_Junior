import { Routes } from '@angular/router';
import { EstoqueComponent } from './features/estoque/estoque';
import { FaturamentoComponent } from './features/faturamento/faturamento';

export const routes: Routes = [
  { path: 'estoque', component: EstoqueComponent },
  { path: 'faturamento', component: FaturamentoComponent },
  { path: '', redirectTo: 'estoque', pathMatch: 'full' }
];
