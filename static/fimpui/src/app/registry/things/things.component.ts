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
import {Thing} from '../model';
import { BACKEND_ROOT } from "app/globals";

@Component({
  selector: 'app-things',
  templateUrl: './things.component.html',
  styleUrls: ['./things.component.css']
})
export class ThingsComponent implements OnInit {
  displayedColumns = ['id', 'alias', 'address','manufacturerId','productId','productHash'];
  dataSource: ThingsDataSource | null;

  @ViewChild('filterAddr') filter: ElementRef;

  constructor(private http : Http) { 
    
  }

  ngOnInit() {
    this.dataSource = new ThingsDataSource(this.http);
    Observable.fromEvent(this.filter.nativeElement, 'keyup')
        .debounceTime(150)
        .distinctUntilChanged()
        .subscribe(() => {
          if (!this.dataSource) { return; }
          this.dataSource.filter = this.filter.nativeElement.value;
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
  set filter(filter: string) { this.getData() }

  constructor(private http : Http) {
    super();
    this.getData();
  }

  getData() {
    this.http
        .get(BACKEND_ROOT+'/fimp/api/registry/things')
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
            thing.alias = result[key].alias;
            thing.productId = result[key].product_id;
            thing.productHash = result[key].product_hash;
            thing.manufacturerId = result[key].manufacturer_id;
            things.push(thing)
     }
     return things;     
  }
}