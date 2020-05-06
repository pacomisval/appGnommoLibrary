import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './component/login/login.component';
import { AuthGuardService } from './services/auth-guard.service';
import { PortadaComponent } from './component/portada/portada.component';
import { AddLibroComponent } from './component/libro/add-libro/add-libro.component';
import { AddAutorComponent } from './component/autor/add-autor/add-autor.component';

const routes: Routes = [
  { path: '', component: PortadaComponent },
  { path: 'portada', component: PortadaComponent },
  { path: 'login', component: LoginComponent },
  { path: 'libro', component: AddLibroComponent },
  { path: 'autor', component: AddAutorComponent }
  //{path: '/panel', loadChildren: './panel/panel.module#PanelModule', canActivate:[AuthGuardService]},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
