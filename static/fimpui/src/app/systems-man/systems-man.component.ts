import { Component, OnInit } from '@angular/core';
import { FimpService} from 'app/fimp/fimp.service';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message'; 

@Component({
  selector: 'app-systems-man',
  templateUrl: './systems-man.component.html',
  styleUrls: ['./systems-man.component.css']
})
export class SystemsManComponent implements OnInit {

  constructor(private fimp:FimpService) { }

  ngOnInit() {
  }

  public connect(service:string,securityKey:string,address:string) {
     let val = {"address":address,"security_key":securityKey};
     let msg  = new FimpMessage(service,"cmd.system.connect","str_map",val,null,null);
     this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:"+service+"/ad:1",msg.toString());
  }
  public disconnect(service:string) {
    let msg  = new FimpMessage(service,"cmd.system.disconnect","null",null,null,null);
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:"+service+"/ad:1",msg.toString());
 }

}
