import { Component, OnInit , OnDestroy ,Input ,ChangeDetectorRef,Inject} from '@angular/core';
import { MdDialog, MdDialogRef,MD_DIALOG_DATA} from '@angular/material';
import { FimpService} from 'app/fimp/fimp.service';
import { Observable }    from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import {Router} from '@angular/router';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message'; 
import { Http, Response,URLSearchParams }  from '@angular/http';
import { BACKEND_ROOT } from "app/globals";
import {
  MqttMessage,
  MqttModule,
  MqttService
}  from 'angular2-mqtt';


@Component({
  selector: 'app-zwave-man',
  templateUrl: './zwave-man.component.html',
  styleUrls: ['./zwave-man.component.css']
})
export class ZwaveManComponent implements OnInit ,OnDestroy {
  selectedOption: string; 
  nodes : any[];
  zwAdState : string;
  errorMsg : string;
  globalSub : Subscription;
  progressBarMode : string ;
  constructor(public dialog: MdDialog,private fimp:FimpService,private router: Router,private http : Http) {
  }

  ngOnInit() {
    this.showProgress(false);
    this.globalSub = this.fimp.getGlobalObservable().subscribe((msg) => {
      console.log(msg.payload.toString());
      let fimpMsg = NewFimpMessageFromString(msg.payload.toString());
      if (fimpMsg.service == "zwave-ad" )
        {
        if(fimpMsg.mtype == "evt.network.all_nodes_report" )
        { 
          this.nodes = fimpMsg.val;
          //this.loadThingsFromRegistry()

          // for(var key in fimpMsg.val){
          //   this.nodes.push({"id":key,"status":fimpMsg.val[key]}); 
          // }
          this.showProgress(false);
          localStorage.setItem("zwaveNodesList", JSON.stringify(this.nodes));
        }else if (fimpMsg.mtype == "evt.thing.exclusion_report" || fimpMsg.mtype == "evt.thing.inclusion_report"){
            console.log("Reloading nodes 2");
            this.reloadNodes();
        }else if (fimpMsg.mtype == "evt.state.report"){
            this.zwAdState = fimpMsg.val;
            if (fimpMsg.val == "NET_UPDATED" || fimpMsg.val == "RUNNING") {
              this.showProgress(false);
            }else if (fimpMsg.val == "STARTING" || fimpMsg.val == "TERMINATED") {
              this.showProgress(true);
            }
        }else if (fimpMsg.mtype == "evt.error.report") {
            this.errorMsg = fimpMsg.props["msg"];
        }else if (fimpMsg.mtype == "evt.network.update_report") {
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
  showProgress(start:boolean){
    if (start){
      this.progressBarMode = "indeterminate";
    }else {
      this.progressBarMode = "determinate";
    }
  } 

  loadThingsFromRegistry() {
     this.http
      .get(BACKEND_ROOT+'/fimp/api/registry/interfaces')
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
        //  console.log(result.report_log_files);
         for(let node of this.nodes) {
           for (let thing of result) {
              // change node.id to node.address
               if (node.address == thing.thing_address) {
                  node["alias"] = thing.location_alias + thing.thing_alias
               }
           }
         }
         localStorage.setItem("zwaveNodesList", JSON.stringify(this.nodes));         
      });
  }
 
  reloadNodes(){
    let msg  = new FimpMessage("zwave-ad","cmd.network.get_all_nodes","null",null,null,null)
    this.showProgress(true);
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
    this.showProgress(true);
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  updateDevice(nodeId :number){
    let msg  = new FimpMessage("zwave-ad","cmd.network.node_update","int",Number(nodeId),null,null)
    this.showProgress(true);
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
  }
  deleteFailedDevice(nodeId :number){
     let val = {"address":String(nodeId),"stop":""}
    let msg  = new FimpMessage("zwave-ad","cmd.thing.delete","str_map",val,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    let dialogRef = this.dialog.open(AddDeviceDialog, {
      height: '400px',
      width: '600px',
      data : "exclusion",
    });
  }
  replaceDevice(nodeId :number){
    let val = {"address":String(nodeId),"stop":""}
    let msg  = new FimpMessage("zwave-ad","cmd.thing.replace","str_map",val,null,null)
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    let dialogRef = this.dialog.open(AddDeviceDialog, {
      height: '400px',
      width: '600px',
      data : "inclusion",
    });
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
  constructor(public dialogRef: MdDialogRef<AddDeviceDialog>,private fimp:FimpService,@Inject(MD_DIALOG_DATA) public data: any) {
    
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
          this.messages.push("Product name :"+fimpMsg.val.product_name);
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
    let msg  = new FimpMessage("zwave-ad","cmd.thing."+this.data,"bool",false,null,null)
    
    this.fimp.publish("pt:j1/mt:cmd/rt:ad/rn:zw/ad:1",msg.toString());
    this.dialogRef.close();
  }

}

