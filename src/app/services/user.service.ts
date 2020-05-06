import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { User } from '../models/user';
import { Globals } from '../Global';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  currentUserType;

  constructor(private http: HttpClient) { }

  createUser(data) {
    console.log("Usuario: " + data.nombre + " insertado con exito");
    return this.http.post<any>(Globals.apiUrl + '/usuarios', data);
  }

  userAdmin() {
    console.log('userAdmin');
    console.log(this.currentUserType);
    return this.currentUserType == 'admin' ? true : false;
  }



}
