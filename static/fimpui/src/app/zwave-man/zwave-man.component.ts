import { Component, OnInit , OnDestroy ,Input ,ChangeDetectorRef} from '@angular/core';
import { MdDialog, MdDialogRef} from '@angular/material';
import { FimpService} from '../fimp.service';
import { Observable }    from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import {Router} from '@angular/router';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message'; 
import {
  MqttMessage,
  MqttModule,
  MqttService
}  from 'angular2-mqtt';


@Component({
  selector: 'app-zwave-man',
  templateUrl: './zwave-man.component.html',
  styleUrls: ['./zwave-man.component.css'],
  providers:[]
})
export class ZwaveManComponent implements OnInit ,OnDestroy {
  selectedOption: string; 
  nodes : any[];
  zwAdState : string;
  globalSub : Subscription;
  constructor(public dialog: MdDialog,private fimp:FimpService,private router: Router) {
  }

  ngOnInit() {
    
    this.globalSub = this.fimp.getGlobalObservable().subscribe((msg) => {
      console.log(msg.payload.toString());
      let fimpMsg = NewFimpMessageFromString(msg.payload.toString());
      if (fimpMsg.service == "zwave-ad" )
        {
        if(fimpMsg.mtype == "evt.network.all_nodes_report" )
        { 
          this.nodes = [];
          for(var key in fimpMsg.val){
            this.nodes.push({"id":key,"status":fimpMsg.val[key]}); 
          }
          localStorage.setItem("zwaveNodesList", JSON.stringify(this.nodes));
        }else if (fimpMsg.mtype == "evt.thing.exclusion_report" || fimpMsg.mtype == "evt.thing.inclusion_report"){
            console.log("Reloading nodes 2");
            this.reloadNodes();
        }else if (fimpMsg.mtype == "evt.state.report"){
            this.zwAdState = fimpMsg.val;
        }
      }
      //this.messages.push("topic:"+msg.topic," payload:"+msg.payload);
    });

    // Let's load nodes list from cache otherwise reload nodes from zwave-ad .
    if (localStorage.getItem("zwaveNodesList")==null){
      this.reloadNodes();
    }else {
      this.nodes = JSON.parse(localStorage.getItem("zwaveNodesList"));
    }
    
  }
  ngOnDestroy() {
    this.globalSub.unsubscribe();
  }

  reloadNodes(){
    let msg  = new FimpMessage("zwave-ad","cmd.network.get_all_nodes","null",null,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  resetNetwork(){
    let msg  = new FimpMessage("zwave-ad","cmd.network.reset","null",null,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  restartAdapter(){
    let msg  = new FimpMessage("zwave-ad","cmd.proc.restart","string","zwave-ad",null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    this.router.navigateByUrl("/timeline");
  }
  updateNetwork(){
    let msg  = new FimpMessage("zwave-ad","cmd.network.update","null",null,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  updateDevice(nodeId :number){
    let msg  = new FimpMessage("zwave-ad","cmd.network.node_update","int",Number(nodeId),null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  addDevice(){
    console.log("Add device")
    let msg  = new FimpMessage("zwave-ad","cmd.thing.inclusion","bool",true,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    let dialogRef = this.dialog.open(AddDeviceDialog, {
      height: '400px',
      width: '600px',
      data : "inclusion",
    });
    dialogRef.afterClosed().subscribe(result => {
      this.selectedOption = result;
    });
  }
  removeDevice(){
    console.log("Remove device ")
    let msg  = new FimpMessage("zwave-ad","cmd.thing.exclusion","bool",true,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    let dialogRef = this.dialog.open(AddDeviceDialog, {
      height: '400px',
      width: '600px',
      data:"exclusion",
    });
    dialogRef.afterClosed().subscribe(result => {
      this.selectedOption = result;
    });
  }
 

}

@Component({
  selector: 'add-device-dialog',
  templateUrl: './dialog.html',
})
export class AddDeviceDialog implements OnInit, OnDestroy  {
  private messages:string[]=[];
  globalSub : Subscription;
  constructor(public dialogRef: MdDialogRef<AddDeviceDialog>,private fimp:FimpService) {
    
    console.log("Dialog constructor Opened");
  }
  ngOnInit(){
    this.messages = [];
    this.globalSub = this.fimp.getGlobalObservable().subscribe((msg) => {
      
      let fimpMsg = NewFimpMessageFromString(msg.payload.toString());
      if (fimpMsg.service == "zwave-ad" )
        {
        if(fimpMsg.mtype == "evt.thing.inclusion_report" )
        { 
          this.messages.push("Node added :"+fimpMsg.val.address);
        } else if (fimpMsg.mtype == "evt.thing.exclusion_report" ){
          this.messages.push("Node removed :"+fimpMsg.val.address);
        }
         else if (fimpMsg.mtype == "evt.thing.inclusion_status_report" ){
          this.messages.push("New state :"+fimpMsg.val);
        } else if (fimpMsg.mtype == "evt.error.report" ){
          this.messages.push("Error : code:"+fimpMsg.val+" message:"+fimpMsg.props["msg"]);
        }
      }
      //this.messages.push("topic:"+msg.topic," payload:"+msg.payload);
    });
  }
  ngOnDestroy() {
    this.globalSub.unsubscribe();
  }
  stopInclusion(){
    let msg  = new FimpMessage("zwave-ad","cmd.thing."+this.dialogRef.config.data,"bool",false,null,null)
    
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }

}

