import { Component, OnInit } from '@angular/core';
import { Http, Response,URLSearchParams }  from '@angular/http';
import { BACKEND_ROOT } from "app/globals";

@Component({
  selector: 'flow-overview',
  templateUrl: './flow-overview.component.html',
  styleUrls: ['./flow-overview.component.css']
})
export class FlowOverviewComponent implements OnInit {
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
