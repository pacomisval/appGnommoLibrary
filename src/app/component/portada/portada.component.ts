import { Component, OnInit, ComponentFactoryResolver } from '@angular/core';
import { Book } from '../../models/book';
import { BookService } from '../../services/book.service';
import { NgbModal, NgbActiveModal, NgbModalRef } from '@ng-bootstrap/ng-bootstrap';
import { AddLibroComponent } from '../libro/add-libro/add-libro.component';



@Component({
  selector: 'app-portada',
  templateUrl: './portada.component.html',
  styleUrls: ['./portada.component.scss']
})
export class PortadaComponent implements OnInit {

  libro: Book;
  libros: Book[] = [];

  cocina: Book[] = [];
  viajes: Book[] = [];
  negra: Book[] = [];
  juvenil: Book[] = [];
  comic: Book[] = [];
  literatura: Book[] = [];
  cifi: Book[] = [];
  historia: Book[] = [];
  misterio: Book[] = [];
  psico: Book[] = [];
  program: Book[] = [];
  contemporanea: Book[] = [];

  id: string;

  portadaModal: NgbModalRef;

  constructor(private bookService: BookService, private modalService: NgbModal) { }


  ngOnInit(): void {
    this.bookService.getAll().subscribe((results) => {
      console.log("Dentro de onInit");
      console.log(results);
      console.log(results[0].genero);
      this.libros = results;

      this.separar();
    },
    (error) => {
      console.log(error);
    });   
  }

  separar() {
    this.libros.forEach(elem => {
      switch(elem.genero) {
        case "Cocina":
          this.cocina.push(elem);
        break;
        case "Historica":
          this.historia.push(elem);
        break;
        case "Viaje":
          this.viajes.push(elem);
        break;
        case "Novela negra":
          this.negra.push(elem);
        break;
        case "Juvenil":
          this.juvenil.push(elem);
        break;
        case "Comic":
          this.comic.push(elem);
        break;
        case "Literatura":
          this.literatura.push(elem);
        break;
        case "Ciencia ficción":
          this.cifi.push(elem);
        break;
        case "Misterio":
          this.misterio.push(elem);
        break;
        case "Psicologia":
          this.psico.push(elem);
        break;
        case "Programación":
          this.program.push(elem);
        break;
        case "Contemporanea":
          this.contemporanea.push(elem);
        break;
        default: 
        console.log("No existe este genero literario: " + elem.genero);
      }

    });
  }

  libroSeleccionado(e) {
    console.log("Dentro de libroSeleccionado: " + e);
    //this.abrirPortadaModal(this.portadaModal);
    this.getLibroPorId(e);
  }

  abrirPortadaModal(portadaModal: any) {
    this.portadaModal = this.modalService.open(portadaModal, {
      ariaLabelledBy: 'modal-basic-title',
    });
  }

  cerrarPortadaModal() {
    this.modalService.dismissAll();
  }

  getLibroPorId(e) {

    this.libros.map((n) => {
      if(n.id == e) {
        this.libro = n;
      }
    });
    return this.libro;
  }

 /*  getLibroPorId(e) {
    
    this.bookService.getBookId(e).subscribe(results => {
      console.log(results);
      this.libro = results;
    },
    (error) => {
      console.log(error);
    });

    return this.libro;
  } */

}
