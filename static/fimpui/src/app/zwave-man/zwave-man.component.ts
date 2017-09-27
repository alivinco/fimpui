import { Component, OnInit , OnDestroy ,Input ,ChangeDetectorRef,Inject} from '@angular/core';
import { MdDialog, MdDialogRef,MD_DIALOG_DATA} from '@angular/material';
import { FimpService} from 'app/fimp/fimp.service';
import { Observable }    from 'rxjs/Observable';
import { Subscription } from 'rxjs/Subscription';
import {Router} from '@angular/router';
import { FimpMessage ,NewFimpMessageFromString } from '../fimp/Message'; 
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import { BACKEND_ROOT } from "app/globals";
import {MdSnackBar} from '@angular/material';
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
  localTemplates : string[];
  localTemplatesCache : string[];
  constructor(public dialog: MdDialog,private fimp:FimpService,private router: Router,private http : Http) {
  }

  ngOnInit() {
    this.showProgress(false);
    this.loadLocalTemplates();
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

  loadLocalTemplates () {
    ///fimp/api/products/list-local-templates?type=cache
    this.http.get(BACKEND_ROOT+'/fimp/api/zwave/products/list-local-templates')
    .map(function(res: Response){
      let body = res.json();
      return body;
    }).subscribe ((result) => {
         this.localTemplates = result     
    });
    this.http.get(BACKEND_ROOT+'/fimp/api/zwave/products/list-local-templates?type=cache')
    .map(function(res: Response){
      let body = res.json();
      return body;
    }).subscribe ((result) => {
         this.localTemplatesCache = result     
    });
  }
  downloadTemplatesFromCloud(){
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});
    this.http
    .post(BACKEND_ROOT+'/fimp/api/zwave/products/download-from-cloud',  options )
    .subscribe ((result) => {
       console.log("Flow was saved");
    });
  }
  uploadCacheToCloud() {
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});
    this.http
    .post(BACKEND_ROOT+'/fimp/api/zwave/products/upload-to-cloud',  options )
    .subscribe ((result) => {
       console.log("Flow was saved");
    });
  }
  
  openTemplateEditor(templateName:string,templateType :string ) {
    let dialogRef = this.dialog.open(TemplateEditorDialog,{
            // height: '95%',
            width: '95%',
            data:{"name":templateName,"type":templateType} 
          });
    dialogRef.afterClosed().subscribe(result => {
              this.loadLocalTemplates();
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

@Component({
  selector: 'template-editor-dialog',
  templateUrl: './template-editor-dialog.html',
})
export class TemplateEditorDialog implements OnInit, OnDestroy  {
  template : any;
  templateStr : string;
  templateName :string;
  templateType :string;
  constructor(public dialogRef: MdDialogRef<TemplateEditorDialog>,@Inject(MD_DIALOG_DATA) public data: any,private http : Http) {
    this.templateName = data["name"];
    this.templateType = data["type"]
    this.template = {};
    this.template["auto_configs"] = {"assoc":[],"configs":[]};
    this.template["dev_custom"] = {"service_grouping":[],"service_descriptor":[],"basic_mapping":[]}
    this.template["docs_ref"] = ""
    console.log("Dialog constructor Opened");
  }

  ngOnInit(){
    this.loadTemplate();
  }

  loadTemplate(){
    this.http.get(BACKEND_ROOT+'/fimp/api/zwave/products/template?name='+this.templateName+'&type='+this.templateType)
    .map(function(res: Response){
      let body = res.json();
      return body;
    }).subscribe ((result) => {
         this.template = result;     
         if(this.template.auto_configs == undefined) {
           this.template["auto_configs"] = {"assoc":[],"configs":[]}
         }
         if(this.template.dev_custom == undefined) {
           this.template["dev_custom"] = {"service_grouping":[],"service_descriptor":[],"basic_mapping":[]}
         }
         if(this.template.comment == undefined){
           this.template["comment"]=""
         }
         if(this.template.wakeup_interval == undefined){
           this.template.wakeup_interval = this.template.wkup_intv;
         } 
         if( this.template["docs_ref"] == undefined){
          this.template["docs_ref"] = "";
         }
         // Converting json object into string, needed for editor 
         this.template.dev_custom.service_descriptor.forEach(element => {
           element.descriptor = JSON.stringify(element.descriptor, null, 2);
         });
        //  this.templateStr = JSON.stringify(result, null, 2);
    });
  }
  addNewAssoc() {
      this.template.auto_configs.assoc.push({"group":1,"node":1,"comment":""})
  }
  deleteAssoc(assoc:any) {
    var i = this.template.auto_configs.assoc.indexOf(assoc);
    if(i != -1) {
      this.template.auto_configs.assoc.splice(i, 1);
    }
  }
  addNewConfig() {
    this.template.auto_configs.configs.push({"key":1,"value":1,"size":1,"comment":""})
  }
  deleteConfig(configObj:any) {
    var i = this.template.auto_configs.configs.indexOf(configObj);
    if(i != -1) {
      this.template.auto_configs.configs.splice(i, 1);
    }
  }
  addNewServiceGrouping() {
    this.template.dev_custom.service_grouping.push({"endp":1,"service":"sensor_temp","group":"ch_0","comment":""})
  }

  deleteServiceGrouping(serviceGrp:any) {
    var i = this.template.dev_custom.service_grouping.indexOf(serviceGrp);
    if(i != -1) {
      this.template.dev_custom.service_grouping.splice(i, 1);
    }
  }
  addNewServiceDescriptor() {
    this.template.dev_custom.service_descriptor.push({"endp":0,"operation":"add","descriptor":"","comment":""});
  }
  deleteServiceDescriptor(serviceDescriptor:any) {
    var i = this.template.dev_custom.service_descriptor.indexOf(serviceDescriptor);
    if(i != -1) {
      this.template.dev_custom.service_descriptor.splice(i, 1);
    }
  }
  addNewBasicMapping() {
    this.template.dev_custom.basic_mapping.push({"endp":0,"basic_value":0,"serv":"","msg_type":"","fimp_value":{},"map_range":false,"min":0,"max":100,"comment":"" });
  }
  deleteBasicMapping(basicMapping:any) {
    var i = this.template.dev_custom.basic_mapping.indexOf(basicMapping);
    if(i != -1) {
      this.template.dev_custom.basic_mapping.splice(i, 1);
    }
  }

  templateOperation(opName:string) {
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});
    this.http
      .post(BACKEND_ROOT+'/fimp/api/zwave/products/template-op/'+opName+'/'+this.templateName,null,  options )
      .subscribe ((result) => {
         console.log("Operation executed");
         this.dialogRef.close();
         
      });
  }

  deleteTemplate() {
    this.http
    .delete(BACKEND_ROOT+'/fimp/api/zwave/products/template/'+this.templateType+'/'+this.templateName)
    .subscribe ((result) => {
      console.log("Template deleted");
      this.dialogRef.close();
    });
  }

  prepareTemplate(){
    // Converting descriptor back from string to object 
    this.template.dev_custom.service_descriptor.forEach(element => {
      element.descriptor = JSON.parse(element.descriptor);
    });
  }

  showSource() {
     this.prepareTemplate();
     this.templateStr = JSON.stringify(this.template, null, 2);
     this.template.dev_custom.service_descriptor.forEach(element => {
        element.descriptor = JSON.stringify(element.descriptor,null,2);
     });
  }
  saveSource() {
    this.template = JSON.parse(this.templateStr);
    this.template.dev_custom.service_descriptor.forEach(element => {
      element.descriptor = JSON.stringify(element.descriptor,null,2);
   });
  }


  saveTemplate(){
    this.prepareTemplate();
    console.dir(this.template)
     let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});


    this.http
      .post(BACKEND_ROOT+'/fimp/api/zwave/products/template/'+this.templateType+'/'+this.templateName,JSON.stringify(this.template),  options )
      .subscribe ((result) => {
         console.log("Template is saved");
         
      });
  }

  ngOnDestroy() {
    
  }
  

}

