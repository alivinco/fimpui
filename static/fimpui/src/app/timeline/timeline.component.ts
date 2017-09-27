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
  private topic :string;
  private payload:string;

  private topicFilter:string;
  private serviceFilter:string;
  private msgTypeFilter:string;

  @ViewChild('myTable') table: any;
  constructor(private fimp: FimpService) { 
   var filter = fimp.getFilter();
   this.topicFilter = filter.topicFilter;
   this.serviceFilter = filter.serviceFilter;
   this.msgTypeFilter = filter.msgTypeFilter; 
   this.messages = this.fimp.getFilteredMessagLog();  
   };
  
  filter() {
    this.fimp.setFilter(this.topicFilter,this.serviceFilter,this.msgTypeFilter);  
  } 
  resetFilter(){
    this.topicFilter = "";
    this.serviceFilter = "";
    this.msgTypeFilter = "";
    this.fimp.setFilter(this.topicFilter,this.serviceFilter,this.msgTypeFilter);  
  }
  copyToMqttClient(topic:string ,payload:string) {
    this.topic = topic ;
    this.payload = JSON.stringify(JSON.parse(payload),null,2); 
  } 
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
