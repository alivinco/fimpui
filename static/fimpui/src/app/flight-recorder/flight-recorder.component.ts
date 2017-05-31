import { Component, OnInit } from '@angular/core';
import { Http, Response,URLSearchParams }  from '@angular/http';
@Component({
  selector: 'app-flight-recorder',
  templateUrl: './flight-recorder.component.html',
  styleUrls: ['./flight-recorder.component.css']
})
export class FlightRecorderComponent implements OnInit {
  private reportLogFiles:string[]=[];
  private reportLogMaxSize:number = 0;
  private hostAlias:string = localStorage.getItem("hostAlias") ;
  constructor(private http : Http) { }

  ngOnInit() {
    this.loadSystemConfigs()
  }
  show() {
    console.dir(this.reportLogFiles)
  }
  uploadLogSnapshot(hostAlias) {
    localStorage.setItem("hostAlias",hostAlias)
    let params = new URLSearchParams();
    params.set('hostAlias', hostAlias);
     this.http
      .get('/fimp/fr/upload-log-snapshot',{search:params})
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         this.reportLogFiles = result;
        
      });
  }
  loadSystemConfigs() {
     console.log("Loading system info")
     
     this.http
      .get('/fimp/configs')
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         console.log(result.report_log_files);
         this.reportLogFiles = result.report_log_files;
         this.reportLogMaxSize = result.report_log_size_limit;
         
      });
  }
}
