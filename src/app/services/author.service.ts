import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Author } from '../models/author';
import { Globals } from '../Global';


@Injectable({
  providedIn: 'root'
})
export class AuthorService {

  author: Author;
  id: number;

  comesAddLibro = false;
  
  constructor(private http: HttpClient) {}

  getAll() {
    return this.http.get<any>(Globals.apiUrl + '/autores');
  }

  getAutor(id: number) {
    console.log(`id ${id}`);
    const headers = new HttpHeaders();
    headers.set('Content-Type', 'application/json; charset=utf-8');
    return this.http.get<any>(Globals.apiUrl + '/autores/' + id, {
      headers
    });
    // const peticion = this.http.get<any>(apiUrl + "/autor/" + id, {
    //   headers: headers
    // });
    // peticion.subscribe(
    //   result => {
    //     console.log(result.response);
    //     respuesta = result.response;
    //   },
    //   error => {
    //     respuesta = error;
    //     console.log(error);
    //   }
    // );
  }

  postAutor(data) {
    const headers = new HttpHeaders();
    headers.set('Content-Type', 'application/json; charset=utf-8');
    console.log("Dentro de postAutor: ");
    return this.http.post<any>(Globals.apiUrl + '/autores', data, { headers });
  }

  putAutor(datosAutor) {
    const id = datosAutor.id;
    const data = {
      first_name: datosAutor.first_name,
      last_name:  datosAutor.last_name,
    };
    console.log('mandando datos');
    console.log(data);
    return this.http.put<any>(Globals.apiUrl + '/autores/' + id, data);
  }

  deleteAutor(id: number) {
    return this.http.delete<any>(Globals.apiUrl + '/autores/' + id);
  }
  // addAutorLibro() {
  //   this.comesAddLibro = true;
  // }
  modificarAuthor(datosAutor) {
    // buscar a ver si esta
    // si esta devolver id del autor
    // no esta modificar los datos
    return this.putAutor(datosAutor);
  }
}
