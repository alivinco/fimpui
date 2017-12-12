import {Component, ElementRef, ViewChild,OnInit,Input,Output,EventEmitter} from '@angular/core';
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
import {MatDialog, MatDialogRef,MatSnackBar} from '@angular/material';

@Component({
  selector: 'app-event-log',
  templateUrl: './event-log.component.html',
  styleUrls: ['./event-log.component.css']
})
export class EventLogComponent implements OnInit {
  displayedColumns = ['timestamp','resourceType','address','code','errSource','msg'];
  dataSource: EventLogDataSource | null;
  constructor(private http : Http,public dialog: MatDialog) { 
  }

  ngOnInit() {
    this.dataSource = new EventLogDataSource(this.http);
  }
}

export class EventLogDataSource extends DataSource<any> {
  events : any[] = [];
  eventsObs = new BehaviorSubject<any[]>([]);
  
  constructor(private http : Http) {
    super();
    this.getData();
  }

  getData() {
    let params: URLSearchParams = new URLSearchParams();
    params.set('pageSize', '100');
    params.set('page', '0');
    this.http
        .get(BACKEND_ROOT+'/fimp/api/stats/event-log',{search:params})
        .map((res: Response)=>{
          let result = res.json();
          return result;
        }).subscribe(result=>{
          this.eventsObs.next(result);
        });

  }
  
  connect(): Observable<Location[]> {
    return this.eventsObs;
  }
  disconnect() {}
  
  
}

