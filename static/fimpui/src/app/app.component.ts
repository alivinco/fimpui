import { Component } from '@angular/core';

@Component({
  moduleId: module.id,
  selector: 'fimp-ui',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  showHeading = true;
  
  toggleHeading() {
    this.showHeading = !this.showHeading;
  }
}
