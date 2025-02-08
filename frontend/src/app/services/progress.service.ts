import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProgressService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getProgress(): Observable<any> {
    return this.http.get(`${this.apiUrl}/progress`);
  }

  getRecords(page: number = 1, pageSize: number = 10): Observable<any> {
    return this.http.get(`${this.apiUrl}/records`, {
      params: { page: page.toString(), pageSize: pageSize.toString() }
    });
  }
}
