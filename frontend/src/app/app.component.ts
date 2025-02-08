import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <div class="app-container">
      <fluent-navigation-bar>
        <fluent-navigation-item>Crawler Dashboard</fluent-navigation-item>
      </fluent-navigation-bar>
      <main>
        <router-outlet></router-outlet>
      </main>
    </div>
  `,
  styles: [`
    .app-container {
      height: 100vh;
      display: flex;
      flex-direction: column;
    }
    main {
      padding: 20px;
      flex: 1;
    }
  `]
})
export class AppComponent {
  title = 'Crawler Dashboard';
}
