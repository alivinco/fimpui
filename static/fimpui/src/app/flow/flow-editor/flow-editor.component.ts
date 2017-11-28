import { Component, OnInit,Inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import {MdDialog, MdDialogRef,MdSnackBar} from '@angular/material';
import {MD_DIALOG_DATA} from '@angular/material';
import { FimpService } from "app/fimp/fimp.service";
import { FimpMessage } from "app/fimp/Message";
import { msgTypeToValueTypeMap } from "app/things-db/mapping";
import { BACKEND_ROOT } from "app/globals";
import { RegistryModule} from 'app/registry/registry.module'
import { ServiceInterface } from "app/registry/model";
import * as D3NE from "d3-node-editor";

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
  
  constructor(private route: ActivatedRoute,private http : Http,public dialog: MdDialog) {
    this.flow = new Flow();
   }

  ngOnInit() {
    let id  = this.route.snapshot.params['id'];
    this.loadFlow(id);
    this.initD3Flow();
    
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
        //  console.dir(this.flow)
        //  console.log(this.flow.Name)
         this.loadContext();
      });
  }

  initD3Flow(){

    var numSocket = new D3NE.Socket("number", "Number value", "hint");
    
    var componentNum = new D3NE.Component("Number", {
       builder: function builder(node) {
          var out1 = new D3NE.Output("Number", numSocket);
          var numControl = new D3NE.Control('<input type="number">', function (el, c) {
             el.value = c.getData('num') || 1;
    
             function upd() {
                c.putData("num", parseFloat(el.value));
                editor.eventListener.trigger("change");
             }
    
             el.addEventListener("input", upd);
             el.addEventListener("mousedown", function (e) {
                e.stopPropagation();
             }); // prevent node movement when selecting text in the input field
             upd();
          });
    
          return node.addControl(numControl).addOutput(out1);
       },
       worker: function worker(node, inputs, outputs) {
          outputs[0] = node.data.num;
       }
    });

    var componentFlowNode = new D3NE.Component("FlowNode", {
      builder: function builder(node) {
         var out1 = new D3NE.Output("Number", numSocket);
         var numControl = new D3NE.Control('<button md-raised-button (click)="addNode(\'counter\')" >Counter</button>', function (el, c) {
            
         });
   
         return node.addControl(numControl).addOutput(out1);
      },
      worker: function worker(node, inputs, outputs) {
         outputs[0] = node.data.num;
      }
   });
    
    var componentAdd = new D3NE.Component("Add", {
       builder: function builder(node) {
          var inp1 = new D3NE.Input("Number", numSocket);
          var inp2 = new D3NE.Input("Number", numSocket);
          var out = new D3NE.Output("Number", numSocket);
    
          var numControl = new D3NE.Control('<input readonly type="number">', function (el, control) {
             control.setValue = function (val) {
                el.value = val;
             };
          });
    
          return node.addInput(inp1).addInput(inp2).addControl(numControl).addOutput(out);
       },
       worker: function worker(node, inputs, outputs) {
          var sum = inputs[0][0] + inputs[1][0];
          editor.nodes.find(function (n) {
             return n.id == node.id;
          }).controls[0].setValue(sum);
          outputs[0] = sum;
       }
    });
    
    var menu = new D3NE.ContextMenu({
       Values: {
          Value: componentNum,
          Action: function Action() {
             alert("ok");
          }
       },
       Add: componentAdd
    });
    
    var container = document.getElementById("nodeEditor");
    var components = [componentNum, componentAdd,componentFlowNode];
    var editor = new D3NE.NodeEditor("demo@0.1.0", container, components, menu);
    
    var nn = componentNum.newNode();
    nn.data.num = 2;
    var n1 = componentNum.builder(nn);
    var n2 = componentNum.builder(componentNum.newNode());
    var add = componentAdd.builder(componentAdd.newNode());
    var flowNode =  componentFlowNode.builder( componentFlowNode.newNode());
    
    n1.position = [80, 200];
    n2.position = [80, 400];
    add.position = [500, 240];
    flowNode.position = [200,400];
    
    editor.connect(n1.outputs[0], add.inputs[0]);
    editor.connect(n2.outputs[0], add.inputs[1]);
    
    editor.addNode(n1);
    editor.addNode(n2);
    editor.addNode(add);
    editor.addNode(flowNode);
    //  editor.selectNode(tnode);
    
    var engine = new D3NE.Engine("demo@0.1.0", components);
    
    // editor.eventListener.on("change", _asyncToGenerator(regeneratorRuntime.mark(function _callee() {
    //    return regeneratorRuntime.wrap(function _callee$(_context) {
    //       while (1) {
    //          switch (_context.prev = _context.next) {
    //             case 0:
    //                _context.next = 2;
    //                return engine.abort();
    
    //             case 2:
    //                _context.next = 4;
    //                return engine.process(editor.toJSON());
    
    //             case 4:
    //             case "end":
    //                return _context.stop();
    //          }
    //       }
    //    }, _callee, this);
    // })));
    
    editor.view.zoomAt(editor.nodes);
    editor.eventListener.trigger("change");
    editor.view.resize();


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
 addIfExpression(node:MetaNode){ 
     let rightVariable = {};
     let expr = {};
     expr["Operand"] = "eq";
     expr["LeftVariableName"] = "";
     rightVariable["Value"] = 100;
     rightVariable["ValueType"] = "int";
     expr["RightVariable"] = rightVariable
     expr["BooleanOperator"] = "";
     node.Config["Expression"].push(expr);
 } 
  
 addNode(nodeType:string){
    console.dir(this.selectedNewNodeType)
    let node  = new MetaNode()
    node.Id = this.getNewNodeId();
    node.Type = nodeType;
    node.Address = ""
    node.Service = ""
    node.ServiceInterface = ""
    node.SuccessTransition = ""
    node.Config = null

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
  
  showContextDialog() {
    let dialogRef = this.dialog.open(ContextDialog,{
            width: '95%',
            data:this.flow
          });
    dialogRef.afterClosed().subscribe(result => {
      
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
  constructor(public dialogRef: MdDialogRef<FlowSourceDialog>,@Inject(MD_DIALOG_DATA) public data: Flow) {
    this.flowSourceText = JSON.stringify(data, null, 2)
  }
  save(){
    this.data = JSON.parse(this.flowSourceText)
    this.dialogRef.close(this.data);
    
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
  constructor(public dialogRef: MdDialogRef<ContextDialog>,@Inject(MD_DIALOG_DATA) public data: Flow,private http : Http) {
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
  
  constructor(public dialogRef: MdDialogRef<FlowRunDialog>,@Inject(MD_DIALOG_DATA) public data: MetaNode,private fimp:FimpService,public snackBar: MdSnackBar) {
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
  constructor(public dialogRef: MdDialogRef<ServiceLookupDialog>,private http : Http,@Inject(MD_DIALOG_DATA) msgFlowDirectionD : string) {
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