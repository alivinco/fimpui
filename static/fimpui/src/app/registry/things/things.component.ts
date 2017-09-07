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
import { Thing } from '../model';
import { ActivatedRoute } from '@angular/router';
import { BACKEND_ROOT } from "app/globals";

@Component({
  selector: 'app-things',
  templateUrl: './things.component.html',
  styleUrls: ['./things.component.css']
})
export class ThingsComponent implements OnInit {
  displayedColumns = ['id', 'alias', 'address','manufacturerId','productId','productName','productHash','action'];
  dataSource: ThingsDataSource | null;

  @ViewChild('filterAddr') filter: ElementRef;

  constructor(private http : Http,private route: ActivatedRoute) { 
    
  }

  ngOnInit() {
    let locationId = this.route.snapshot.params['filterValue'];
    this.dataSource = new ThingsDataSource(this.http);
    if (locationId=="*"){
      locationId = "";
    }
    
    this.dataSource.getData(locationId);
    Observable.fromEvent(this.filter.nativeElement, 'keyup')
        .debounceTime(150)
        .distinctUntilChanged()
        .subscribe(() => {
          if (!this.dataSource) { return; }
          this.dataSource.filter = this.filter.nativeElement.value;
        });
  }

  deleteThing(id:string) {
    this.http
     .delete(BACKEND_ROOT+'/fimp/api/registry/thing/'+id)
     .subscribe ((result) => {
        this.dataSource.getData("");
     });
 } 
  }

  /**
 * Data source to provide what data should be rendered in the table. Note that the data source
 * can retrieve its data in any way. In this case, the data source is provided a reference
 * to a common data base, ThingsDatabase. It is not the data source's responsibility to manage
 * the underlying data. Instead, it only needs to take the data and send the table exactly what
 * should be rendered.
 */
export class ThingsDataSource extends DataSource<any> {
  _filterChange = new BehaviorSubject('');
  things : Thing[] = [];
  thingsObs = new BehaviorSubject<Thing[]>([]);
  get filter(): string { return this._filterChange.value; }
  set filter(filter: string) { this.getData("") }

  constructor(private http : Http) {
    super();
    // this.getData("");
  }

  getData(locationId:string) {
    let params: URLSearchParams = new URLSearchParams();
    params.set('locationId', locationId);
    this.http
        .get(BACKEND_ROOT+'/fimp/api/registry/things',{search:params})
        .map((res: Response)=>{
          let result = res.json();
          return this.mapThings(result);
        }).subscribe(result=>{
          this.thingsObs.next(result);
        });

  }
  
  connect(): Observable<Thing[]> {
    return this.thingsObs;
  }
  disconnect() {}

  mapThings(result:any):Thing[] {
    let things : Thing[] = [];
    for (var key in result){
            let thing = new Thing();
            thing.id = result[key].id;
            thing.address = result[key].address;
            thing.commTech = result[key].comm_tech;
            thing.alias = result[key].alias;
            thing.productId = result[key].product_id;
            thing.productName = result[key].product_name;
            thing.productHash = result[key].product_hash;
            thing.manufacturerId = result[key].manufacturer_id;
            thing.locationId = result[key].location_id;
            things.push(thing)
     }
     return things;     
  }
}