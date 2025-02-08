import { Component, OnInit } from '@angular/core';
import { ProgressService } from '../../services/progress.service';

@Component({
  selector: 'app-dashboard',
  template: `
    <div class="dashboard-container">
      <fluent-card>
        <h2>Crawling Progress</h2>
        <fluent-progress-bar
          [value]="progress"
          [max]="100">
        </fluent-progress-bar>
        <div class="metrics">
          <div>
            <h3>Pages Crawled</h3>
            <p>{{crawledPages}} / {{totalPages}}</p>
          </div>
          <div>
            <h3>Status</h3>
            <p>{{status}}</p>
          </div>
        </div>
      </fluent-card>

      <fluent-card>
        <h2>Latest Records</h2>
        <fluent-data-grid [items]="records">
          <!-- Add columns based on your data structure -->
        </fluent-data-grid>
      </fluent-card>
    </div>
  `,
  styles: [`
    .dashboard-container {
      display: grid;
      gap: 20px;
      max-width: 1200px;
      margin: 0 auto;
    }
    .metrics {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 20px;
      margin-top: 20px;
    }
  `]
})
export class DashboardComponent implements OnInit {
  progress = 0;
  crawledPages = 0;
  totalPages = 0;
  status = 'Initializing...';
  records: any[] = [];

  constructor(private progressService: ProgressService) {}

  ngOnInit() {
    this.startProgressPolling();
  }

  private startProgressPolling() {
    setInterval(() => {
      this.progressService.getProgress().subscribe(
        (data: any) => {
          this.crawledPages = data.crawledPages;
          this.totalPages = data.totalPages;
          this.status = data.status;
          this.progress = (this.crawledPages / this.totalPages) * 100;
        }
      );
    }, 5000);
  }
}
