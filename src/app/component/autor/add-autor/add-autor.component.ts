import { Component, OnInit, ViewChild, TemplateRef } from '@angular/core';
import { Router } from '@angular/router';
import { AuthorService } from 'src/app/services/author.service';
import { FormGroup, FormBuilder, Validators, FormControl } from '@angular/forms';
import { NgForm } from '@angular/forms';
import { Location } from '@angular/common';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { THIS_EXPR } from '@angular/compiler/src/output/output_ast';
/**
 * Componente para añadir Autor
 *
 * @export
 * @class AddautorComponent
 * @implements {OnInit}
 */

@Component({
  selector: 'app-add-autor',
  templateUrl: './add-autor.component.html',
  styleUrls: ['./add-autor.component.scss']
})
export class AddAutorComponent implements OnInit {
 /**
   * View child Ventana Modal con un mensaje
   */
  @ViewChild('modalInformation', { static: false })
  modalInformation: TemplateRef<any>;

  /**
   * Formulario de Autor
   */
  authorForm: FormGroup;
  /**
   * Formulario emitido
   *
   * @memberof AddautorComponent
   */
  submittedAuthor = false;
  /**
   * Inicializa autor
   *
   * @memberof AddautorComponent
   */
  autor = {
    first_name: '',
    last_name: '',
    nacionalidad: '',
    fechaNacimiento: ''
  };
  /**
   * Mensaje en ventana modal
   *
   * @type {string}
   * @memberof AddautorComponent
   */
  information: string;
  /**
   * Verificacion del formulario
   *
   * @type {boolean}
   * @memberof AddautorComponent
   */
  invalidated: boolean;


  /**
   * Creates an instance of AddautorComponent.
   * @param {Router} router Para enrutar
   * @param {AuthorService} authorService Servicio para Author
   * @param {NgbModal} modalService Servivio para ventanas Modales
   * @memberof AddautorComponent
   */
  constructor(
    private router: Router,
    private authorService: AuthorService,
    private modalService: NgbModal,
    private formBuilder: FormBuilder
  ) {}

  ngOnInit() {
    /* this.authorForm = this.formBuilder.group({
      first_name: ['', Validators.required],
      last_name: ['', Validators.required],
      nacionalidad: ['', Validators.required],
      fechaNacimiento: ['', Validators.required],
    }); */

    /* this.authorForm = new FormGroup({
      firstName: new FormControl()
   }); */
  }
  /**
   * AbreAbreviatura de autorForm.controls
   *
   * @memberof AddautorComponent
   */
  get afc() {
    return this.authorForm.controls;
  }
  /**
   * Guarda Autor en BD
   *
   * @memberof AddautorComponent
   */

  comprobacionFinal(){
    console.log("entraaa rafa");
    console.log(this.autor)
    var res=true;
    var letras : RegExp = /^[A-Za-z\s]{2,50}$/g;
    var fecha: RegExp = / ^ \ d { 4 } \ - \ d { 1 , 2 } \ - \ d { 1 , 2 } $ /g;
    
    console.log(this.autor.first_name.length);
    console.log(this.autor.first_name);
    console.log(this.autor.last_name.length);
    console.log(this.autor.last_name);
    console.log(this.autor.nacionalidad.length);
    console.log(this.autor.nacionalidad);
    console.log(this.autor.fechaNacimiento);


    if(this.autor.first_name.length > 50) {
      this.information = "-Has superado el límite de carácteres máximos en el campo nombre";
      res=false; 
    }
    
    else if(letras.test(this.autor.first_name)) {
      this.information = "-En el campo nombre solo se permiten letras";
      res=false;
    }

    if(this.autor.last_name.length > 50) {
      this.information = "-Has superado el límite de carácteres máximos en el campo apellido";
      res=false;
    }
    else if(letras.test(this.autor.last_name)){
      console.log("apellidos: " + this.autor.last_name);
      this.information = "-En el campo apellido solo se permiten letras";
      res=false;
    }

    if(this.autor.nacionalidad.length > 50) {
      this.information = "-Has superado el límite de carácteres máximos en el campo nombre";
      res=false; 
    }
    else if(letras.test(this.autor.nacionalidad)) {
      this.information = "-En el campo nombre solo se permiten letras";
      res=false;
    }

    if(!this.autor.fechaNacimiento) {
      this.information = "Introduce una fecha"
      res = false;
    }
    else if(fecha.test(this.autor.fechaNacimiento)) {
      this.information = "El formato de fecha es YYYY-MM-DD";
      res = false;
    }
    

    if(!res){
      this.openInformationWindows();
      
    }
    
    return res;
   }

  Guardar() {
    //if(this.comprobacionFinal()){
     // console.log(this.authorForm.controls);
      this.submittedAuthor = true;
      // if (this.authorForm.invalid) {
      //   console.log("autorForm invalid: " + this.authorForm.invalid);
      //   return;
      // }
      console.log(this.authorService.comesAddLibro);
    
      const data = {
        id: "",
        first_name: this.autor.first_name,
        last_name: this.autor.last_name,
        nacionalidad: this.autor.nacionalidad,
        fechaNacimiento: this.autor.fechaNacimiento
      };

      // controlamos que no este repetido.
      // repetido=> Damos como agregado

      // no repetido =>agregamos
      this.authorService.postAutor(data).subscribe(
        (results) => {
          this.information = 'Autor añadido';
          this.openInformationWindows();
          this.backRoute();

        },
        (error) => {
          this.information = 'Autor no añadido';
          this.openInformationWindows();

          this.authorService.comesAddLibro = true;
          this.backRoute();
          //this.router.navigate(['/']);
        }
      );
     
  }

  /**
   * Enruta en segun el valor de comesAddLibro
   */
  backRoute() {
    if (this.authorService.comesAddLibro) {
      this.authorService.comesAddLibro = false;
      this.router.navigate(['autor']);
      window.location.reload();
    } else {
      this.router.navigate(['portada']);
    }
  }

  // checkForm() {
  //  // if (this.autor.first_name != '' || this.autor.last_name != '')

  //     //hacer busqueda sql select * where first= and last = ... result
  //      //   = null -> añadir
  //     //   != null -> informar autor existe
  // }

  /**
   * Abre Ventana Modal informativa
   *
   * @memberof ListarComponent
   */
  openInformationWindows() {
    this.modalService.open(this.modalInformation);
  }

}
