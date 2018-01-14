import { Component, OnInit ,Input } from '@angular/core';
import { MetaNode, ServiceLookupDialog } from "../flow-editor/flow-editor.component";
import { MatDialog, MatDialogRef} from '@angular/material';
import { msgTypeToValueTypeMap } from "app/things-db/mapping";
import { FlowRunDialog } from "app/flow/flow-editor/flow-editor.component"
import { Http, Response }  from '@angular/http';
import { BACKEND_ROOT } from "app/globals";

@Component({
  selector: 'app-flow-nodes',
  templateUrl: './flow-nodes.component.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class FlowNodesComponent implements OnInit {

  constructor() { }

  ngOnInit() {
  }
  
}

@Component({
  selector: 'action-node',
  templateUrl: './action-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class ActionNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  @Input() flowId:string;
  localVars:any;
  globalVars:any;
  complexValueAsString:any; //string representation of node.Config.DefaultValue.Value
  propsAsString:any;
  constructor(public dialog: MatDialog,private http : Http) { 
    this.loadContext();
   }

  ngOnInit() { 
    // backword compatability
    if (this.node.Config.VariableName=="undefined"){
      this.node["DefaultValue"] = {"Value":"","ValueType":""};
      this.node["VariableName"] = "";
      this.node["IsVariableGlobal"] = false;
    }
    try{
      this.complexValueAsString = JSON.stringify(this.node.Config.DefaultValue.Value);
    }catch (err){
      console.log("Can't stringify complex default value")
    }
    try{
      this.propsAsString = JSON.stringify(this.node.Config.Props);
    }catch (err){
      console.log("Can't stringify props ")
    }
    
  }
  updateComplexValue(){
    this.node.Config.DefaultValue.Value = JSON.parse(this.complexValueAsString)
  }
  updateProps(){
    this.node.Config.Props = JSON.parse(this.propsAsString)
  }

  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '500px',
            data:"in"
          });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
        this.nodes.forEach(element => {
           
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              if(element.Label==""||element.Label==undefined){
                element.Label =  result.serviceAlias + " at "+result.locationAlias
              }
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.DefaultValue.ValueType =  result.intfValueType

            }
        });
    });      
  }

  loadContext() {
    if (this.flowId) {
       this.http
      .get(BACKEND_ROOT+'/fimp/flow/context/'+this.flowId)
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
         this.localVars = [];
         for (var key in result){
            this.node
            this.localVars.push(result[key].Name);
         }
         
      });
    }
   
    
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
  variableSelected(event:any,config:any,isGlobal:boolean){
    // if (config.VariableName.indexOf("__global__")!=-1) {
    //   config.VariableName = config.VariableName.replace("__global__","");
    //   config.VariableIsGlobal = true;
    // }
    config.IsVariableGlobal = isGlobal;

  }

}

@Component({
  selector: 'set-variable-node',
  templateUrl: './set-variable-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class SetVariableNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
  }
}

@Component({
  selector: 'if-node',
  templateUrl: './if-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class IfNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
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
}

@Component({
  selector: 'receive-node',
  templateUrl: './receive-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class ReceiveNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }

  ngOnInit() { 
  }
  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '500px',
            data:"out"
          });
    dialogRef.afterClosed().subscribe(result => {
      console.dir(result)
      if (result)
        this.nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              if(element.Label==""||element.Label==undefined){
                element.Label =  result.serviceAlias + " at "+result.locationAlias
              }
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.ValueFilter.ValueType =  result.intfValueType
            }
        });
    });      
  }
}


/*type TimeTriggerConfig struct {
	DefaultMsg model.Variable
	Expressions []TimeExpression
	GenerateAstroTimeEvents bool
	Latitude float64
	Longitude float64
}

type TimeExpression struct {
	Name string
	Expression string   //https://godoc.org/github.com/robfig/cron#Job
	Comment string
}
 */

@Component({
  selector: 'time-trigger-node',
  templateUrl: './time-trigger-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class TimeTriggerNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
  }
}

@Component({
  selector: 'vinc-trigger-node',
  templateUrl: './vinc-trigger-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class VincTriggerNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  @Input() flowId:string;
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
  }
}


@Component({
  selector: 'loop-node',
  templateUrl: './loop-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class LoopNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
    
  }
}

@Component({
  selector: 'wait-node',
  templateUrl: './wait-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class WaitNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() { 
  }
}

@Component({
  selector: 'trigger-node',
  templateUrl: './trigger-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class TriggerNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  @Input() flowId:string;
  flowPublishService: string;
  flowPublishInterface : string;
  flowPublishAddress : string;
  constructor(public dialog: MatDialog) { }

  ngOnInit() { 
    
    // backword compatability 
    if (this.node.Config == null) {
      this.node.Config = {};
    }
    if (this.node.Config.Timeout == null) {
      this.node.Config["Timeout"] = 0;
      this.node.Config["ValueFilter"] = {"Value":"","ValueType":""}; 
    }
  }
  runFlow(node:MetaNode) {
    let dialogRef = this.dialog.open(FlowRunDialog,{
            // height: '95%',
            width: '500px',
            data:node
          });
    dialogRef.afterClosed().subscribe(result => {
      // this.flow = result;
      // this.loadContext();
    });      
  }

  onPublishServiceChange(){
    console.log(this.flowPublishService);
    var msgType = "cmd";
    try{
      if (this.flowPublishInterface.indexOf("evt.")>=0){
        msgType = "evt";
      }
    }catch(err){
      
    }
    
    this.flowPublishAddress = "pt:j1/mt:"+msgType+"/rt:dev/rn:flow/ad:1/sv:"+this.flowPublishService+"/ad:"+this.flowId;
  }
  publishFlowAsVirtualDevice(){
    this.node.ServiceInterface = this.flowPublishInterface;
    this.node.Service = this.flowPublishService;
    this.node.Address = this.flowPublishAddress;
  }
  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '500px',
            data:"out"
          });
    dialogRef.afterClosed().subscribe(result => {
      console.dir(result)
      if (result)
        this.nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              if(element.Label==""||element.Label==undefined){
                element.Label =  result.serviceAlias + " at "+result.locationAlias
              }
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.ValueFilter.ValueType =  result.intfValueType
            }
        });
    });      
  }
}
