import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class UploadService {

  urlServidor: string = "http://localhost:8000/api/upload";

  constructor(private http: HttpClient) { }

  public upload(data) {

    return this.http.post<any>(this.urlServidor, data, {
      reportProgress: true,
      observe: 'events'
    }).pipe(map((event) => {

      switch(event.type) {
        
        case HttpEventType.UploadProgress:
          const progress = Math.round(100 * event.loaded / event.total);
          console.log("valor de progress: " + progress);

          return { status: 'progress', message: progress };

        case HttpEventType.Response:
          console.log("valor de event.body: " + event.body);
          return event.body;
        
        default:
          console.log("valor de event.type: " + event.type);
          return `Unhandled evento: ${event.type}`;
      }
    })
    );
  }
}
