import { Component, OnInit } from '@angular/core';
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
  constructor(private fimp: FimpService) { 
   
   };
  

  ngOnInit() {
    // this.fimp.mqtt.onConnect.subscribe((message: any) => {
    //       console.log("timeline onConnect");
           
    //  });
     this.subscribe();
     
  }
  subscribe(){
    this.fimp.getGlobalObservable().subscribe((msg) => {
      console.log("New message in timeline")
      let fimpMsg  = NewFimpMessageFromString(msg.payload.toString());
      fimpMsg.topic = msg.topic;
      this.messages.push(fimpMsg);
      console.log(this.messages.length);
    });
  }
 

}
