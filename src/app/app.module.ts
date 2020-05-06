import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './component/login/login.component';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { FormsModule, NgControl, NgForm, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from './services/auth.service';
import { ToastrModule } from 'ngx-toastr';
import { CommonModule } from '@angular/common';
//import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AuthGuardService } from './services/auth-guard.service';
import { JwtHelperService } from '@auth0/angular-jwt';
import { PortadaComponent } from './component/portada/portada.component';
import { AddLibroComponent } from './component/libro/add-libro/add-libro.component';
import { NgbModule, NgbActiveModal } from '@ng-bootstrap/ng-bootstrap';
import { BookService } from './services/book.service';
//import { Authorservice } from './services/author.service';
import { UploadService } from './services/upload.service';
import { HeaderInterceptor } from './header.interceptor';
import { AddAutorComponent } from './component/autor/add-autor/add-autor.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    PortadaComponent,
    AddLibroComponent,
    AddAutorComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    CommonModule,
    //BrowserAnimationsModule,
    ToastrModule.forRoot(),
    NgbModule
  ],
  providers: [
    AuthService,
    AuthGuardService,
    JwtHelperService,
    BookService,
    UploadService,
    NgbActiveModal,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HeaderInterceptor,
      multi: true,
    }
    //AuthorService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
