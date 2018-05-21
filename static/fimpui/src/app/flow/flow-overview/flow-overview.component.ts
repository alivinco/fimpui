import { Component, OnInit } from '@angular/core';
import { Http, Response,URLSearchParams }  from '@angular/http';
import { BACKEND_ROOT } from "app/globals";
import { DatePipe } from '@angular/common';
import {FlowLogDialog} from "../flow-editor/flow-editor.component";
import {MatDialog} from "@angular/material";

@Component({
  selector: 'flow-overview',
  templateUrl: './flow-overview.component.html',
  styleUrls: ['./flow-overview.component.css']
})
export class FlowOverviewComponent implements OnInit {
  flows : any[];
  constructor(private http : Http,public dialog: MatDialog) {  }

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
  showLog() {
    let dialogRef = this.dialog.open(FlowLogDialog,{
      // height: '95%',
      width: '95%',
      data:{flowId:"",mode:"all_flows"}
    });
    dialogRef.afterClosed().subscribe(result => {

    });
  }
}
