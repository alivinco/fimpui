import { Component, OnInit } from '@angular/core';
import { Http, Response,URLSearchParams }  from '@angular/http';

// export const BACKEND_ROOT = "http://localhost:8081"
export const BACKEND_ROOT = ""

@Component({
  selector: 'app-flow',
  templateUrl: './flow.component.html',
  styleUrls: ['./flow.component.css']
})
export class FlowComponent implements OnInit {
  flows : any[];
  constructor(private http : Http) {  }

  ngOnInit() {
    this.loadListOfFlows()
  }
  loadListOfFlows() {
     this.http
      .get(BACKEND_ROOT+'/fimp/flow/list')
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         this.flows = result;
      });
  } 
  deleteFlow(id:string) {
     this.http
      .delete(BACKEND_ROOT+'/fimp/flow/definition/'+id)
      .subscribe ((result) => {
         this.loadListOfFlows()
      });
  } 
}
