import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, timer } from 'rxjs';
import { switchMap, share } from 'rxjs/operators';

export interface Progress {
  total_items: number;
  processed_items: number;
  current_url: string;
  status: string;
  last_update_time: string;
}

@Injectable({
  providedIn: 'root'
})
export class ProgressService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  // Regular progress check
  getProgress(): Observable<Progress> {
    return this.http.get<Progress>(`${this.apiUrl}/progress`);
  }

  // Polling progress every 5 seconds
  pollProgress(): Observable<Progress> {
    return timer(0, 5000).pipe(
      switchMap(() => this.getProgress()),
      share()
    );
  }

  getRecords(page: number = 1, pageSize: number = 10): Observable<any> {
    return this.http.get(`${this.apiUrl}/records`, {
      params: { page: page.toString(), pageSize: pageSize.toString() }
    });
  }
}
