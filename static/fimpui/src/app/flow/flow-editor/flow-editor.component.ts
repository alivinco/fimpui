import { Component, OnInit,Inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import {MatDialog, MatDialogRef,MatSnackBar} from '@angular/material';
import {MAT_DIALOG_DATA} from '@angular/material';
import { FimpService } from "app/fimp/fimp.service";
import { FimpMessage } from "app/fimp/Message";
import { msgTypeToValueTypeMap } from "app/things-db/mapping";
import { BACKEND_ROOT } from "app/globals";
import { RegistryModule} from 'app/registry/registry.module'
import { ServiceInterface } from "app/registry/model";

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
  Ui                :Ui;
}

export class Ui {
  x:number;
  y:number;
}

export class Variable {
  Value :any;
  ValueType :string;
}

@Component({
  selector: 'app-flow-editor',
  templateUrl: './flow-editor.component.html',
  styleUrls: ['./flow-editor.component.css']
})
export class FlowEditorComponent implements OnInit {
  flow :Flow;
  selectedNewNodeType:string;
  localVars:any;
  globalVars:any;
  // properties for drag-and-drop
  currentDraggableNode:any;
  dragStartPosX:number;
  dragStartPosY:number;
  currentDraggableNodeId:string;
  isDraggableLine:boolean;

  constructor(private route: ActivatedRoute,private http : Http,public dialog: MatDialog) {
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
         this.enhanceNodes();
        //  console.dir(this.flow)
        //  console.log(this.flow.Name)
         this.loadContext();
      });
  }
  loadContext() {
    this.http
      .get(BACKEND_ROOT+'/fimp/flow/context/'+this.flow.Id)
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
         this.localVars = [];
         for (var key in result){
            this.localVars.push(result[key].Name);
         }
         
      });
    
    this.http
      .get(BACKEND_ROOT+'/fimp/flow/context/global')
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
        this.globalVars = [];
        for (var key in result){
            this.globalVars.push(result[key].Name);
         }
      });  
  }  
  variableSelected(event:any,config:any){
    if (config.LeftVariableName.indexOf("__global__")!=-1) {
      config.LeftVariableName = config.LeftVariableName.replace("__global__","");
      config.LeftVariableIsGlobal = true;
    }
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
  //  let cloneNode = new MetaNode();
   
  //  Object.assign(cloneNode,node);
  //  cloneNode.Id = this.getNewNodeId();
  // temp quick dirty way how to clone nested objects;
   var cloneNode = <MetaNode>JSON.parse(JSON.stringify(node));
   this.flow.Nodes.push(cloneNode);
 }
   
 allowNodeDrop(event:any) {
   if (!this.isDraggableLine){
     console.log("allow node drop");
     event.preventDefault();
    }
  } 
  
 nodeDrop(event:any) {
   if (!this.isDraggableLine){
      console.log("node dropped");
      console.dir(event);
      var offsetX = event.clientX - this.dragStartPosX ;
      var offsetY = event.clientY - this.dragStartPosY ;
      var nodeStartPositionX = 20;
      var nodeStartPositionY = 20;

      if (this.currentDraggableNode.style.left)
        nodeStartPositionX = Number(this.currentDraggableNode.style.left.replace("px",""));
      if (this.currentDraggableNode.style.top)
        nodeStartPositionY = Number(this.currentDraggableNode.style.top.replace("px",""));
      var xPos = nodeStartPositionX + offsetX;
      var yPos = nodeStartPositionY + offsetY
      this.currentDraggableNode.style.left = (xPos)+"px";
      this.currentDraggableNode.style.top = (yPos)+"px";

      var node = this.getNodeById(this.currentDraggableNodeId)
      if(node.Ui == undefined) {
        node.Ui = new Ui();
      }
      node.Ui.x = xPos;
      node.Ui.y = yPos;

      event.preventDefault();
  }
 }
 nodeDragStart(event:any){
  this.currentDraggableNode = event.srcElement;
  this.dragStartPosX = event.clientX;
  this.dragStartPosY = event.clientY;
  if (event.srcElement.className.includes("socket")){
    console.log("Line drag start");
    this.isDraggableLine = true  
  } else {
    console.log("Node drag start");
    this.isDraggableLine = false;
    console.dir(event);
    this.currentDraggableNodeId = event.srcElement.id.replace("nodeId_","")
    console.log("active node id = "+this.currentDraggableNodeId);
  }
 }
 //////////////////////////
 allowLineDrop(event:any) {
  console.log("allow line drop");
  event.preventDefault();
 }
 lineDrop(event:any) {
    console.log("line dropped");
    console.dir(event);
    var newCord =  this.findAbsolutePosition(this.currentDraggableNode);
    this.drawCurvedLine(newCord.x,newCord.y,event.clientX,event.clientY,"black",0.5);
    event.preventDefault();
 }
 
 //////////////////////////

 drawCurvedLine(x1, y1, x2, y2, color, tension) {
  var svg = document.getElementById("flow-connections"); 
  var shape = document.createElementNS("http://www.w3.org/2000/svg","path");
  var delta = (x2-x1)*tension;
  var hx1=x1+delta;
  var hy1=y1;
  var hx2=x2-delta;
  var hy2=y2;
  var path = "M "  + x1 + " " + y1 + 
             " C " + hx1 + " " + hy1 
                   + " "  + hx2 + " " + hy2 
             + " " + x2 + " " + y2;
  shape.setAttributeNS(null, "d", path);
  shape.setAttributeNS(null, "fill", "none");
  shape.setAttributeNS(null, "stroke", color);
  svg.appendChild(shape);
}

findAbsolutePosition(htmlElement):any {
  var x = htmlElement.offsetLeft;
  var y = htmlElement.offsetTop;
  for (var el=htmlElement;el != null;el = el.offsetParent) {
         x += el.offsetLeft;
         y += el.offsetTop;
  }
  return {
      "x": x,
      "y": y
  };
}

/////////////////////////// 
 addNode(nodeType:string){
    console.dir(this.selectedNewNodeType)
    let node  = new MetaNode()
    node.Id = this.getNewNodeId();
    node.Type = nodeType;
    node.Address = "";
    node.Service = "";
    node.ServiceInterface = "";
    node.SuccessTransition = "";
    node.Config = null;
    node.Ui = new Ui();
    node.Ui.x = 70;
    node.Ui.y = 170;

    switch (node.Type){
      case "trigger":
        node.Config = {}; 
        node.Config["Timeout"] = 0;
        node.Config["ValueFilter"] = {"Value":"","ValueType":""};
        node.Config["IsValueFilterEnabled"] = false;
        break;
      case "action":
        node.Config = {"VariableName":"","IsVariableGlobal":false,"Props":{}}; 
        node.Config["DefaultValue"] = {"Value":"","ValueType":""};
        break;
      case "counter":
        node.Config = {}; 
        node.Config["StartValue"] = 0;
        node.Config["EndValue"] = 5;
        node.Config["Step"] = 1;
        node.Config["EndValueTransition"] = "";
        node.Config["SaveToVariable"] = false;

        break;  
      case "receive":
        node.Config = {}; 
        node.Config["Timeout"] = 120;
        node.Config["ValueFilter"] = {"Value":"","ValueType":""};
        node.Config["IsValueFilterEnabled"] = false;
        break;  
      case "if":
        node.Config = {}; 
        node.Config["TrueTransition"] = ""
        node.Config["FalseTransition"] = ""
        node.Config["Expression"] = [];
        let expr = {};
        let rightVariable = {};
        expr["Operand"] = "eq";
        expr["LeftVariableName"] = "";
        expr["LeftVariableIsGlobal"] = false;
        rightVariable["Value"] = 100;
        rightVariable["ValueType"] = "int";
        expr["RightVariable"] = rightVariable
        expr["BooleanOperator"] = "";
        node.Config["Expression"].push(expr);
        break;
      case "wait":
        node.Config = 1000;
        break;
      case "set_variable":
        node.Config = {}; 
        node.Config["Name"] = ""
        node.Config["UpdateGlobal"] = false
        node.Config["UpdateInputMsg"] = false
        let variable = {};
        variable["Value"] = 100;
        variable["ValueType"] = "int";
        node.Config["DefaultValue"] = variable
        break; 
      case "time_trigger":
        node.Config = {};
        node.Config["DefaultMsg"] = {"Value":"","ValueType":""};
        let expressions = [];
        expressions.push({"Name":"","Expression":"","Comments":""});
        node.Config["Expressions"] = expressions;
        node.Config["GenerateAstroTimeEvents"] = false;
        node.Config["Latitude"] = 0.0;
        node.Config["Longitude"] = 0.0;
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
  
  showNodeEditorDialog(flow:Flow,node:MetaNode) {
    let dialogRef = this.dialog.open(NodeEditorDialog,{
      // height: '95%',
      width: '95%',
      data:{"flow":flow,"node":node}
    });
    dialogRef.afterClosed().subscribe(result => {
    
    });  
  }

  showContextDialog() {
    let dialogRef = this.dialog.open(ContextDialog,{
            width: '95%',
            data:this.flow
          });
    dialogRef.afterClosed().subscribe(result => {
      
    });      
  }
  
  getNodeById(nodeId:string):MetaNode {
    console.log("GEtting node for id = "+nodeId);
    var node:MetaNode;
    this.flow.Nodes.forEach(element => {
        console.dir(element)
        if (element.Id==nodeId) {
          node = element;
          return ;
        }
    });
    return node;
  }

  enhanceNodes() {
    this.flow.Nodes.forEach(node => {
      if(node.Ui == undefined) {
        node.Ui = new Ui()
        node.Ui.x = 70;
        node.Ui.y = 170;
      }
      
  });
  }

  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '95%'
          });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
        this.flow.Nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.ValueType =  msgTypeToValueTypeMap[element.ServiceInterface]
            }
        });


    });      
  }
  
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
  constructor(public dialogRef: MatDialogRef<FlowSourceDialog>,@Inject(MAT_DIALOG_DATA) public data: Flow) {
    this.flowSourceText = JSON.stringify(data, null, 2)
  }
  save(){
    this.data = JSON.parse(this.flowSourceText)
    this.dialogRef.close(this.data);
    
  }
}

@Component({
  selector: 'node-editor-dialog',
  templateUrl: 'node-editor-dialog.html',
  styleUrls: ['flow-editor.component.css']
})
export class NodeEditorDialog {
  flow :Flow;
  node :MetaNode;
  constructor(public dialogRef: MatDialogRef<NodeEditorDialog>,@Inject(MAT_DIALOG_DATA) public data:any) {
    this.flow = data.flow;
    this.node = data.node;
   }
  
}

@Component({
  selector: 'context-dialog',
  templateUrl: 'context-dialog.html',
  styleUrls: ['flow-editor.component.css']
})
export class ContextDialog {
  localContext :string ;
  globalContext : string;
  constructor(public dialogRef: MatDialogRef<ContextDialog>,@Inject(MAT_DIALOG_DATA) public data: Flow,private http : Http) {
     this.http
      .get(BACKEND_ROOT+'/fimp/flow/context/'+data.Id)
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
         this.localContext = JSON.stringify(result, null, 2);
      });
    
    this.http
      .get(BACKEND_ROOT+'/fimp/flow/context/global')
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
         this.globalContext = JSON.stringify(result, null, 2);
      });
      
      
    // this.localContext = JSON.stringify(data, null, 2)
    // this.globalContext = JSON.stringify(data, null, 2)
  }
  
}

@Component({
  selector: 'flow-run-dialog',
  templateUrl: 'flow-run-dialog.html',
})
export class FlowRunDialog {
  value : any;
  valueType : string ;
  actionData : MetaNode;
  
  constructor(public dialogRef: MatDialogRef<FlowRunDialog>,@Inject(MAT_DIALOG_DATA) public data: MetaNode,private fimp:FimpService,public snackBar: MatSnackBar) {
    // data.Config. = {"Value":true,"ValueType":msgTypeToValueTypeMap[data.ServiceInterface]}
    // this.valueType = msgTypeToValueTypeMap[data.ServiceInterface];
    // this.actionData = new MetaNode();
    // this.actionData.Address = data.Address;
    // this.actionData.Id = data.Id;
    // this.actionData.Service = data.Service;
    // this.actionData.ServiceInterface = data.ServiceInterface;
    // this.actionData.Type = data.Type;
    // this.actionData.Config = {"Value":true,"ValueType":data.Config.ValueFilter.ValueType}
    this.valueType = data.Config.ValueFilter.ValueType
  
  }
  
  run(){
    let msg  = new FimpMessage(this.data.Service,this.data.ServiceInterface,this.valueType,this.value,null,null)
    this.fimp.publish(this.data.Address,msg.toString());
    let snackBarRef = this.snackBar.open('Message was sent',"",{duration:1000});
  }
}

@Component({
  selector: 'service-lookup-dialog',
  templateUrl: 'service-lookup-dialog.html',
  styleUrls: ['./flow-editor.component.css']
})
export class ServiceLookupDialog  implements OnInit {
  interfaces :any;
  msgFlowDirectionD = "";
  constructor(public dialogRef: MatDialogRef<ServiceLookupDialog>,private http : Http,@Inject(MAT_DIALOG_DATA) msgFlowDirectionD : string) {
    console.log("Msg flow direction:"+msgFlowDirectionD);
    this.msgFlowDirectionD = msgFlowDirectionD
  }
  ngOnInit() {
    console.log("ng on init Msg flow  direction:"+this.msgFlowDirectionD);  
  }
 
  onSelected(intf :ServiceInterface){
    console.dir(intf);
    this.dialogRef.close(intf);

  }
}