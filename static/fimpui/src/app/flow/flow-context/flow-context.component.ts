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

export class Context {
  FlowId : string ;
  Name :string;
  Description:string;
  UpdatedAt:string;
  Value:any;
  ValueType:string;
}

@Component({
  selector: 'flow-context',
  templateUrl: './flow-context.component.html',
  styleUrls: ['./flow-context.component.css']
})
export class FlowContextComponent implements OnInit {
  displayedColumns = ['flowId','name','description','valueType','value','updatedAt','action'];
  dataSource: FlowContextDataSource | null;
  constructor(private http : Http) { 
  }
  ngOnInit() {
    this.dataSource = new FlowContextDataSource(this.http);
  }
}


export class FlowContextDataSource extends DataSource<any> {
  locations : Location[] = [];
  locationsObs = new BehaviorSubject<Context[]>([]);
  
  constructor(private http : Http) {
    super();
    console.log("Getting context data")
    this.getData();
  }

  getData() {
    this.http
        .get(BACKEND_ROOT+'/fimp/flow/context/global')
        .map((res: Response)=>{
          let result = res.json();
          return this.mapContext(result);
        }).subscribe(result=>{
          this.locationsObs.next(result);
        });

  }
  
  connect(): Observable<Context[]> {
    return this.locationsObs;
  }
  disconnect() {}

  mapContext(result:any):Context[] {
    let locations : Context[] = [];
    for (var key in result){
            let loc = new Context();
            loc.FlowId = "global";
            loc.Name = result[key].Name;
            loc.Description = result[key].Description; 
            loc.UpdatedAt = result[key].UpdatedAt;
            loc.Value = result[key].Variable.Value; 
            loc.ValueType = result[key].Variable.ValueType; 
            locations.push(loc)
     }
     return locations;     
  }
}
