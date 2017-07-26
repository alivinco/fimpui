import { Component, OnInit,Inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import {MdDialog, MdDialogRef,MdSnackBar} from '@angular/material';
import {MD_DIALOG_DATA} from '@angular/material';
import { FimpService } from "app/fimp.service";
import { FimpMessage } from "app/fimp/Message";
export const BACKEND_ROOT = "http://localhost:8081"
// export const BACKEND_ROOT = ""
const msgTypeToValueTypeMap = {
  "evt.binary.report":"bool",
  "evt.binary.set":"bool",
  "evt.lvl.report":"int",
  "cmd.lvl.set":"int",
  "evt.meter.report":"float",
  "evt.sensor.report":"float",
  "evt.open.report":"bool",
  "evt.presence.report":"bool",
  "evt.alarm.report":"str_map",
  "evt.setpoint.report":"str_map",
  "cmd.setpoint.report":"str_map",
  "cmd.mode.set":"string",
  "evt.mode.report":"string",
  "evt.state.report":"string",
  "evt.lock.report":"bool_map",
  "cmd.lock.set":"bool",
  "cmd.color.set":"int_map",
  "evt.color.report":"int_map",
  "evt.scene.report":"string",

}

@Component({
  selector: 'app-flow-editor',
  templateUrl: './flow-editor.component.html',
  styleUrls: ['./flow-editor.component.css']
})
export class FlowEditorComponent implements OnInit {
  flow :Flow;
  selectedNewNodeType:string;
  constructor(private route: ActivatedRoute,private http : Http,public dialog: MdDialog) {
    this.flow = new Flow();
   }

  ngOnInit() {
    let id  = this.route.snapshot.params['id'];
    this.loadFlow(id);
  }
 
  loadFlow(id:string) {
     this.http
      .get(BACKEND_ROOT+'/fimp/flow/definition/'+id)
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         this.flow = result;
         console.dir(this.flow)
         console.log(this.flow.Name)
      });
  }
  saveFlow() {
    console.dir(this.flow)
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});
    this.http
      .post(BACKEND_ROOT+'/fimp/flow/definition/'+this.flow.Id,JSON.stringify(this.flow),  options )
      .subscribe ((result) => {
         console.log("Flow was saved");
      });

 }
 getNewNodeId():string {
   let id = 0;
   let maxId = 0;
   this.flow.Nodes.forEach(element => {
     id = parseInt(element.Id);
     if (id > maxId) {
       maxId = id ;
     }
   });
   maxId++;
   return maxId+"";
 } 
 deleteNode(node:MetaNode){
   let index: number = this.flow.Nodes.indexOf(node);
    if (index !== -1) {
        this.flow.Nodes.splice(index, 1);
    }  
 } 
 cloneNode(node:MetaNode){
   let cloneNode = new MetaNode();
   
   Object.assign(cloneNode,node);
   cloneNode.Id = this.getNewNodeId();
   this.flow.Nodes.push(cloneNode);
 } 
  
 addNode(){
    console.dir(this.selectedNewNodeType)
    let node  = new MetaNode()
    node.Id = this.getNewNodeId();
    node.Type = this.selectedNewNodeType;
    node.Address = ""
    node.Service = ""
    node.ServiceInterface = ""
    node.SuccessTransition = ""
    node.Config = null

    switch (node.Type){
      case "trigger":
        // nothing to add yet
        break;
      case "action":
        node.Config = {}; 
        node.Config["Value"] = true;
        node.Config["ValueType"] = "bool";
        break;
      case "if":
        node.Config = {}; 
        node.Config["TrueTransition"] = ""
        node.Config["FalseTransition"] = ""
        node.Config["Expression"] = [];
        let expr = {}
        expr["Operand"] = "eq";
        expr["Value"] = 100;
        expr["ValueType"] = "int";
        expr["BooleanOperator"] = "";
        node.Config["Expression"].push(expr);
        break;
      case "wait":
        node.Config = 1000;
        break;      
    }
    this.flow.Nodes.push(node) 
  }
  showSource() {
    let dialogRef = this.dialog.open(FlowSourceDialog,{
            // height: '95%',
            width: '95%',
            data:this.flow
          });
    dialogRef.afterClosed().subscribe(result => {
      if(result)
        this.flow = result;
    });      
  }

  runFlow(node:MetaNode) {
    let dialogRef = this.dialog.open(FlowRunDialog,{
            // height: '95%',
            width: '95%',
            data:node
          });
    dialogRef.afterClosed().subscribe(result => {
      // this.flow = result;
    });      
  }
}

export class MetaNode {
  Id               :string;
	Type             :string;
	Label            :string;
	SuccessTransition :string;
	TimeoutTransition :string;
	ErrorTransition   :string;
	Address           :string;
	Service           :string;
	ServiceInterface  :string;
	Config            :any;
}

export class Flow {
    Id :string ;
    Name : string ;
    Description : string ;
    Nodes : MetaNode[] ;
}

@Component({
  selector: 'flow-source-dialog',
  templateUrl: 'flow-source-dialog.html',
  styleUrls: ['flow-editor.component.css']
})
export class FlowSourceDialog {
  flowSourceText :string ;
  constructor(public dialogRef: MdDialogRef<FlowSourceDialog>,@Inject(MD_DIALOG_DATA) public data: Flow) {
    this.flowSourceText = JSON.stringify(data, null, 2)
  }
  save(){
    this.data = JSON.parse(this.flowSourceText)
    this.dialogRef.close(this.data);
    
  }
}

@Component({
  selector: 'flow-run-dialog',
  templateUrl: 'flow-run-dialog.html',
})
export class FlowRunDialog {
  constructor(public dialogRef: MdDialogRef<FlowRunDialog>,@Inject(MD_DIALOG_DATA) public data: MetaNode,private fimp:FimpService,public snackBar: MdSnackBar) {

    data.Config = {"Value":true,"ValueType":msgTypeToValueTypeMap[data.ServiceInterface]}
    // console.dir(data)
  }
  
  run(){
    let msg  = new FimpMessage(this.data.Service,this.data.ServiceInterface,this.data.Config.ValueType,this.data.Config.Value,null,null)
    this.fimp.publish(this.data.Address,msg.toString());
    let snackBarRef = this.snackBar.open('Message was sent',"",{duration:1000});
  }
}