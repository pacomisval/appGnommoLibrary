import { Component, OnInit } from '@angular/core';
import { AuthService } from './services/auth.service';
import { ToastrService } from 'ngx-toastr';
import { Router } from '@angular/router';
@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'appGo1 Libreria';

  constructor(private route: Router) {}

  portada() {
    this.route.navigate([""]);
  }

  libro() {
    this.route.navigate(["libro"]);
  }

  autor() {
    this.route.navigate(["autor"]);
  }

  login() {
    this.route.navigate(["login"]);
  }












}
