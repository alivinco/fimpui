import {Component, ElementRef, ViewChild,OnInit} from '@angular/core';
import {DataSource} from '@angular/cdk';
import {BehaviorSubject} from 'rxjs/BehaviorSubject';
import {Observable} from 'rxjs/Observable';
import { Http, Response,URLSearchParams }  from '@angular/http';
import 'rxjs/add/operator/startWith';
import 'rxjs/add/observable/merge';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import 'rxjs/add/observable/fromEvent';
import {Location} from '../model';
import { BACKEND_ROOT } from "app/globals";

@Component({
  selector: 'app-locations',
  templateUrl: './locations.component.html',
  styleUrls: ['./locations.component.css']
})
export class LocationsComponent implements OnInit {
displayedColumns = ['id','type','alias','address','geo','action'];

// displayedColumns = ['thingAddress', 'thingAlias',
// 'serviceName','serviceAlias','intfMsgType'];
  dataSource: LocationsDataSource | null;

  constructor(private http : Http) { 
    
  }

  ngOnInit() {
    this.dataSource = new LocationsDataSource(this.http);
   
  }
}


export class LocationsDataSource extends DataSource<any> {
  locations : Location[] = [];
  locationsObs = new BehaviorSubject<Location[]>([]);
  
  constructor(private http : Http) {
    super();
    this.getData();
  }

  getData() {
    let params: URLSearchParams = new URLSearchParams();
    this.http
        .get(BACKEND_ROOT+'/fimp/api/registry/locations',{search:params})
        .map((res: Response)=>{
          let result = res.json();
          return this.mapThings(result);
        }).subscribe(result=>{
          this.locationsObs.next(result);
        });

  }
  
  connect(): Observable<Location[]> {
    return this.locationsObs;
  }
  disconnect() {}

  mapThings(result:any):Location[] {
    let locations : Location[] = [];
    for (var key in result){
            let loc = new Location();
            loc.id = result[key].id;
            loc.type = result[key].type;
            loc.alias = result[key].alias; 
            loc.address = result[key].address; 
            loc.long = result[key].long; 
            loc.lat = result[key].lat; 
            locations.push(loc)
     }
     return locations;     
  }
}
