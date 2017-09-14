import { Component, OnInit , ViewChild } from '@angular/core';
import { FimpService} from 'app/fimp/fimp.service'
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
  
  sendMessage(topic:string,payload:string) {
    this.fimp.publish(topic,payload);
  } 

  ngOnInit() {
   
  }
  
  toggleExpandRow(row) {
    console.log('Toggled Expand Row!', row);
    this.table.rowDetail.toggleExpandRow(row);
  }
 

}
