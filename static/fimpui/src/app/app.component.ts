import { Component } from '@angular/core';
import { Http, Response }  from '@angular/http';
import { Observable } from 'rxjs/Rx';
@Component({
  moduleId: module.id,
  selector: 'fimp-ui',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  showHeading = true;
  private version :String;
  constructor (private http : Http){
    this.loadSystemInfo();
  }
  toggleHeading() {
    this.showHeading = !this.showHeading;
  }

  loadSystemInfo() {
     console.log("Loading system info")
     
     this.http
      .get('/fimp/system-info')
      .map(function(res: Response){
        let body = res.json();
        console.log(body.Version);
        return body;
      }).subscribe (function(result){
         console.log(result.Version);
         this.version = result.Version;
         
      });
  }


}
