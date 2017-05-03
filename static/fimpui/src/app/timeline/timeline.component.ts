import { Component, OnInit , ViewChild } from '@angular/core';
import { FimpService} from '../fimp.service'
import { FimpMessage,NewFimpMessageFromString } from '../fimp/Message'; 

@Component({
  selector: 'timeline',
  //providers:[FimpService],
  templateUrl: './timeline.component.html',
  styleUrls: ['./timeline.component.css']
})
export class TimelineComponent implements OnInit {
  private messages:FimpMessage[]=[];
  @ViewChild('myTable') table: any;
  constructor(private fimp: FimpService) { 
   this.messages = fimp.getMessagLog();
   };
  

  ngOnInit() {
    // this.fimp.mqtt.onConnect.subscribe((message: any) => {
    //       console.log("timeline onConnect");
           
    //  });
    //  this.subscribe();
     
  }
  // subscribe(){
  //   this.fimp.getGlobalObservable().subscribe((msg) => {
  //     console.log("New message in timeline")
  //     let fimpMsg  = NewFimpMessageFromString(msg.payload.toString());
  //     fimpMsg.topic = msg.topic;
  //     fimpMsg.raw = msg.payload.toString();
  //     fimpMsg.localTs =  Date.now();
  //     this.messages.push(fimpMsg);
  //     console.log(this.messages.length);
  //   });
  // }

  toggleExpandRow(row) {
    console.log('Toggled Expand Row!', row);
    this.table.rowDetail.toggleExpandRow(row);
  }
 

}
