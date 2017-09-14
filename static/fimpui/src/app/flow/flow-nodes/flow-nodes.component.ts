import { Component, OnInit ,Input } from '@angular/core';
import { MetaNode, ServiceLookupDialog } from "../flow-editor/flow-editor.component";
import { MdDialog, MdDialogRef} from '@angular/material';
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
  constructor(public dialog: MdDialog,private http : Http) { 
    this.loadContext();
   }

  ngOnInit() { 
    // backword compatability
    if (this.node.Config.VariableName=="undefined"){
      this.node["DefaultValue"] = {"Value":"","ValueType":""};
      this.node["VariableName"] = "";
      this.node["IsVariableGlobal"] = false;
    }
  

  }
  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '95%'
          });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
        this.nodes.forEach(element => {
           
            if (element.Id==nodeId) {
              // console.dir(result);
              element.Service = result.serviceName
              element.Label = element.Label + result.serviceAlias + " at "+result.loactionAlias
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              // element.Config.DefaultValue.ValueType =  msgTypeToValueTypeMap[element.ServiceInterface]
              element.Config.DefaultValue.ValueType =  result.intfValueType;

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
  variableSelected(event:any,config:any){
    if (config.VariableName.indexOf("__global__")!=-1) {
      config.VariableName = config.LeftVariableName.replace("__global__","");
      config.VariableIsGlobal = true;
    }
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
  constructor(public dialog: MdDialog) { }

  ngOnInit() { 
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
  constructor(public dialog: MdDialog) { }

  ngOnInit() { 
  }
  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '95%'
          });
    dialogRef.afterClosed().subscribe(result => {
      console.dir(result)
      if (result)
        this.nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              element.Label = element.Label + result.serviceAlias + " at "+result.loactionAlias
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.ValueFilter.ValueType =  msgTypeToValueTypeMap[element.ServiceInterface]
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
  constructor(public dialog: MdDialog) { }
  ngOnInit() { 
  }
}

@Component({
  selector: 'counter-node',
  templateUrl: './counter-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class CounterNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  constructor(public dialog: MdDialog) { }
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
  constructor(public dialog: MdDialog) { }

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
            width: '95%',
            data:node
          });
    dialogRef.afterClosed().subscribe(result => {
      // this.flow = result;
      // this.loadContext();
    });      
  }


  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '95%'
          });
    dialogRef.afterClosed().subscribe(result => {
      console.dir(result)
      if (result)
        this.nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.serviceName
              element.Label = element.Label + result.serviceAlias + " at "+result.loactionAlias
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              element.Config.ValueFilter.ValueType =  result.intfValueType
            }
        });
    });      
  }
}
