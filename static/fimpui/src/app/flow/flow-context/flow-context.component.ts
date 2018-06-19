import {Component, ElementRef, ViewChild,OnInit} from '@angular/core';
import {DataSource} from '@angular/cdk/collections';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {Observable} from 'rxjs/Observable';
import { Http, Response,URLSearchParams }  from '@angular/http';
import 'rxjs/add/operator/startWith';
import 'rxjs/add/observable/merge';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import 'rxjs/add/observable/fromEvent';
import { BACKEND_ROOT } from "app/globals";
import {ThingEditorDialog} from "../../registry/things/thing-editor.component";
import {Thing} from "../../registry/model";
import {MatDialog} from "@angular/material";
import {RecordEditorDialog} from "./record-editor-dialog.component";
import {TableContextRec} from "./model"


@Component({
  selector: 'flow-context',
  templateUrl: './flow-context.component.html',
  styleUrls: ['./flow-context.component.css']
})
export class FlowContextComponent implements OnInit {
  displayedColumns = ['flowId','name','description','valueType','value','updatedAt','action'];
  dataSource: FlowContextDataSource | null;
  constructor(private http : Http,public dialog: MatDialog) {
  }
  ngOnInit() {
    this.dataSource = new FlowContextDataSource(this.http);
  }

  showRecordEditorDialog(ctxRec:TableContextRec) {
    ctxRec.FlowId = "global";
    let dialogRef = this.dialog.open(RecordEditorDialog,{
      width: '450px',
      data:ctxRec
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
      {
        this.dataSource.getData()
      }
    });
  }

  showAddNewRecordEditorDialog() {
    var ctxRec = new TableContextRec();
    ctxRec.FlowId = "global";
    let dialogRef = this.dialog.open(RecordEditorDialog,{
      width: '450px',
      data:ctxRec
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
      {
        this.dataSource.getData();
      }
    });
  }

  reload() {
    this.dataSource.getData();
  }


}


export class FlowContextDataSource extends DataSource<any> {
  ctxRecordsObs = new BehaviorSubject<TableContextRec[]>([]);

  constructor(private http : Http) {
    super();
    console.log("Getting context data")
    this.getData();
  }

  getData() {
    this.http
        .get(BACKEND_ROOT+'/fimp/api/flow/context/global')
        .map((res: Response)=>{
          let result = res.json();
          return this.mapContext(result);
        }).subscribe(result=>{
          this.ctxRecordsObs.next(result);
        });

  }

  connect(): Observable<TableContextRec[]> {
    return this.ctxRecordsObs ;
  }
  disconnect() {}

  mapContext(result:any):TableContextRec[] {
    let contexts : TableContextRec[] = [];
    for (var key in result){
            let loc = new TableContextRec();
            loc.FlowId = "global";
            loc.Name = result[key].Name;
            loc.Description = result[key].Description;
            loc.UpdatedAt = result[key].UpdatedAt;
            loc.Value = result[key].Variable.Value;
            loc.ValueType = result[key].Variable.ValueType;
            contexts.push(loc)
     }
     return contexts;
  }
}
