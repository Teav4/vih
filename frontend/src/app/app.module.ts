import { NgModule, CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { RouterModule } from '@angular/router';
import { provideFluentDesignSystem, allComponents } from '@fluentui/web-components';

import { AppComponent } from './app.component';
import { DashboardComponent } from './components/dashboard/dashboard.component';

// Register FluentUI web components
provideFluentDesignSystem().register(allComponents);

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    RouterModule.forRoot([
      { path: '', component: DashboardComponent },
    ])
  ],
  providers: [],
  bootstrap: [AppComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class AppModule { }
