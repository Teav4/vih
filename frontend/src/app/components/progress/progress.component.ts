import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { ProgressService, Progress } from '../../services/progress.service';

@Component({
  selector: 'app-progress',
  template: `
    <div *ngIf="progress">
      <h3>Crawler Progress</h3>
      <p>Status: {{progress.status}}</p>
      <p>Progress: {{progress.processed_items}} / {{progress.total_items}}</p>
      <p>Current URL: {{progress.current_url}}</p>
      <p>Last Update: {{progress.last_update_time | date:'medium'}}</p>
      <mat-progress-bar
        mode="determinate"
        [value]="(progress.processed_items / progress.total_items) * 100">
      </mat-progress-bar>
    </div>
  `
})
export class ProgressComponent implements OnInit, OnDestroy {
  progress: Progress | null = null;
  private subscription: Subscription | null = null;

  constructor(private progressService: ProgressService) {}

  ngOnInit() {
    this.subscription = this.progressService.pollProgress()
      .subscribe(progress => {
        this.progress = progress;
      });
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}
