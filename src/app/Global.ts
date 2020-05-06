import { Injectable } from '@angular/core';

/**
 * Guarda variables globales
 *
 * @export
 * @class Globals
 */
@Injectable()

export class Globals {


  /**
   * ruta del servidor
   *
   * @static
   * @memberof Globals
   */
  static apiUrl = `http://localhost:8000/api`;
}