import { Component, OnInit ,Input ,Pipe, PipeTransform } from '@angular/core';
import { Interface } from "../things-db/thing-model";
import { FimpService } from '../fimp.service';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message';

@Component({
  selector: 'thing-intf-ui',
  templateUrl: './thing-intf-ui.component.html',
  styleUrls: ['./thing-intf-ui.component.css']
})
export class ThingIntfUiComponent implements OnInit {
  @Input() intf: Interface;
  @Input() addr: string;
  @Input() service :string;
  
  constructor(private fimp:FimpService) { 
    // this.fimp.getGlobalObservable().subscribe((msg) => {
    // let fimpMsg = NewFimpMessageFromString(msg.payload.toString()); 
    // });
   }

  ngOnInit() {
  }
  cmdBinarySet(val:boolean){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdConfigSet(name:string,value:string){
    let val = {};
    val[name] = value;
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdLvlSet(level:number,duration:number){
    var props = new Map<string,string>() ;
    props["duration"] = String(duration);
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,level,props,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdSetpointSet(setpointType:string,temp:string){
    let val = {};
    val["type"] = setpointType;
    val["temp"] = temp;
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdSetpointReportGet(name:string){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,name,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdModeSet(mode:string){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,mode,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdLevelStart(direction:string,duration:number){
    var val = direction;
    var props = new Map<string,string>() ;
    props["duration"] = String(duration);
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,props,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  } 
  cmdLevelStop(direction:string){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,null,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  } 

  cmdGroupSet(group:string,member:string){
    let val = {};
    val["group"] = group;
    val["members"] = [member];
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdGroupReportGet(group:string){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,group,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
  cmdConfigReportGet(name:string){
    let val = [name];
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,val,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }

  cmdGetReportNull(){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,null,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }

  cmdMeterReportGet(unit:string){
    let msg  = new FimpMessage(this.service,this.intf.msgType,this.intf.valueType,unit,null,null)
    this.fimp.publish("pt:j1/mt:cmd"+this.addr,msg.toString());
  }
}
@Pipe({name: 'keys'})
export class KeysPipe implements PipeTransform {
  transform(value, args:string[]) : any {
    if (!value) {
      return value;
    } 

    let keys = [];
    for (let key in value) {
      keys.push({key: key, value: value[key]});
    } 
    return keys;
  } 
} 