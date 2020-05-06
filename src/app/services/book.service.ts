import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Globals } from '../Global';

@Injectable({
  providedIn: 'root'
})
export class BookService {

  books: any;

  constructor(private http: HttpClient) {}

  getAll() {
    console.log(Globals.apiUrl);
    return this.http.get<any>(Globals.apiUrl + '/libros/all');
  }
  getBookId(id: number) {
    console.log(`id ${id}`);
    return this.http.get<any>(Globals.apiUrl + '/libros/' + id);
  }

  getBookFromAutor(id: number) {
    return this.http.get<any>(Globals.apiUrl + '/libros/autor/' + id);
  }

  createBook(data) {
    console.log(data);
    return this.http.post<any>(Globals.apiUrl + '/libros', data);
  }

  updateBook(data) {
    console.log('entra en updateBook');
    console.log(data);
    return this.http.put<any>(Globals.apiUrl + '/libros/' + data.id, data);
  }

  deleteBook(id) {
    return this.http.delete<any>(Globals.apiUrl + '/libros/' + id);
  }
}
