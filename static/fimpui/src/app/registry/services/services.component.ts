import {Component, ElementRef, ViewChild,OnInit,Output,EventEmitter} from '@angular/core';
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
import { ActivatedRoute } from '@angular/router';
import {ServiceInterface} from '../model';
import { BACKEND_ROOT } from "app/globals";
import {ThingIntfUiComponent} from 'app/registry/thing-intf-ui/thing-intf-ui.component'

 
@Component({
  selector: 'reg-services-main',
  templateUrl: './services-main.component.html',
  styleUrls: ['./services.component.css']
})
export class ServicesMainComponent {
  constructor() { 
    
  }

}

@Component({
  selector: 'reg-services',
  templateUrl: './services.component.html',
  styleUrls: ['./services.component.css']
})
export class ServicesComponent implements OnInit {
  displayedColumns = ['thingTech','thingAddress', 'thingAlias','serviceName','serviceAlias',
                    'intfType','intfMsgType','locationAlias','intfValueType','action'];

  dataSource: ServicesDataSource | null; 
  // usage (onSelect)="onSelected($event)">
  @Output() onSelect = new EventEmitter<ServiceInterface>();
  @ViewChild('filterThingAddr') filterThingAddr: ElementRef;
  @ViewChild('filterServiceName') filterServiceName: ElementRef;
  @ViewChild('filterInterfaceType') filterInterfaceType: ElementRef;

  constructor(private http : Http,private route: ActivatedRoute) { 
    
  }

  ngOnInit() {
    let thingId = this.route.snapshot.params['filterValue'];
    console.log("Thing id  = ",thingId);
    this.dataSource = new ServicesDataSource(this.http);
    this.dataSource.getData("","","",thingId);
    Observable.fromEvent(this.filterThingAddr.nativeElement, 'keyup')
        .debounceTime(500)
        .distinctUntilChanged()
        .subscribe(() => {
          if (!this.dataSource) { return; }
          this.dataSource.getData(this.filterThingAddr.nativeElement.value,this.filterServiceName.nativeElement.value,this.filterInterfaceType.nativeElement.value,"")
        });
    Observable.fromEvent(this.filterServiceName.nativeElement, 'keyup')
        .debounceTime(500)
        .distinctUntilChanged()
        .subscribe(() => {
          if (!this.dataSource) { return; }
          this.dataSource.getData(this.filterThingAddr.nativeElement.value,this.filterServiceName.nativeElement.value,this.filterInterfaceType.nativeElement.value,"")
        }); 
    Observable.fromEvent(this.filterInterfaceType.nativeElement, 'keyup')
        .debounceTime(500)
        .distinctUntilChanged()
        .subscribe(() => {
          if (!this.dataSource) { return; }
          this.dataSource.getData(this.filterThingAddr.nativeElement.value,this.filterServiceName.nativeElement.value,this.filterInterfaceType.nativeElement.value,"")
        });        
  }

  selectInterface(intf:ServiceInterface) {
    console.dir(intf);
    this.onSelect.emit(intf);
  }
}


export class ServicesDataSource extends DataSource<any> {
  services : ServiceInterface[] = [];
  servicesObs = new BehaviorSubject<ServiceInterface[]>([]);
  
  constructor(private http : Http) {
    super();
    
  }

  getData(thingAddr:string ,serviceName:string,interfaceType:string,thingId:string) {
    let params: URLSearchParams = new URLSearchParams();
    params.set('serviceName', serviceName);
    params.set('thingAddr', thingAddr);
    params.set('intfMsgType', interfaceType);
    if (thingId!="*") {
      params.set('thingId', thingId);
    }
    this.http
        .get(BACKEND_ROOT+'/fimp/api/registry/interfaces',{search:params})
        .map((res: Response)=>{
          let result = res.json();
          return this.mapThings(result);
        }).subscribe(result=>{
          this.servicesObs.next(result);
        });

  }
  
  connect(): Observable<ServiceInterface[]> {
    return this.servicesObs;
  }
  disconnect() {}

  mapThings(result:any):ServiceInterface[] {
    let things : ServiceInterface[] = [];
    for (var key in result){
            let thing = new ServiceInterface();
            thing.thingId = result[key].thing_id;
            thing.thingAddress = result[key].thing_address;
            thing.thingTech = result[key].thing_tech; 
            thing.thingAlias = result[key].thing_alias;
            thing.serviceId = result[key].service_id;
            thing.serviceName = result[key].service_name;
            thing.serviceAlias = result[key].service_alias;
            thing.serviceAddress = result[key].service_address;
            thing.intfType = result[key].intf_type;
            thing.intfMsgType = result[key].intf_msg_type;
            thing.intfValueType = result[key].intf_val_type;
            thing.intfAddress = result[key].intf_address;
            thing.locationId = result[key].location_id;
            thing.locationAlias = result[key].location_alias;
            thing.locationType = result[key].location_type;
            things.push(thing)
     }
     return things;     
  }
}
