import { Component, OnInit , OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FimpService} from '../fimp.service';
import { MapFimpInclusionReportToThing } from '../things-db/integrations';
import { Thing } from '../things-db/thing-model';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message';
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-thing-view',
  templateUrl: './thing-view.component.html',
  styleUrls: ['./thing-view.component.css']
})
export class ThingViewComponent implements OnInit ,OnDestroy{
  globalSub : Subscription;
  private thing : Thing;

  rows = [
  ];
  // columns = [
  //   { prop: 'name' , name:'Service name' },
  //   { prop: 'address',width:350 },
  //   { prop: 'groups' },
  // ];

  constructor(private fimp:FimpService,private route: ActivatedRoute,) {
    this.thing = new Thing();
  }

  ngOnInit() {
    let techAdapterName  = this.route.snapshot.params['ad'];
    let id  = this.route.snapshot.params['id'];
    let serviceName = "zwave-ad";
    if (techAdapterName == "ikea"){
      serviceName = "ikea-ad";
    }
    this.getReport(techAdapterName,serviceName,id);
    this.globalSub = this.fimp.getGlobalObservable().subscribe((msg) => {
      
      let fimpMsg = NewFimpMessageFromString(msg.payload.toString());
      if (fimpMsg.service == serviceName )
        {
        if(fimpMsg.mtype == "evt.thing.inclusion_report" )
        { 
          console.log("New thing")
          this.thing = MapFimpInclusionReportToThing(fimpMsg);
          this.rows = this.thing.services;
        } 

      }else {
        console.log("Sensor report");
        for (let svc of this.thing.services){
            // console.log("Comparing "+msg.topic+" with "+ "pt:j1/mt:evt"+svc.address);
            if (msg.topic == "pt:j1/mt:evt"+svc.address) {
              // console.log("Matching service "+fimpMsg.service);
              for (let inf of svc.interfaces) {
                if ( fimpMsg.mtype == inf.msgType ) {
                  // console.log("Value updated");
                  inf.lastValue = fimpMsg
                }
              }
            }
        }
      }
    });
  }
  ngOnDestroy() {
    this.globalSub.unsubscribe();
  }
  getReport(techAdapterName:string,serviceName:string, nodeId:string){
    let msg  = new FimpMessage(serviceName,"cmd.thing.get_inclusion_report","string",nodeId,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:"+techAdapterName+"/ad:1",msg.toString());
  }

}
