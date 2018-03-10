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
  selector: 'vinc-action-node',
  templateUrl: './vinc-action-node.html',
  styleUrls: ['./flow-nodes.component.css']
})
export class VincActionNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  shortcuts:any[];
  constructor(public dialog: MatDialog , private http : Http) { }
  ngOnInit() {
    this.loadShortcuts()
  }
  loadShortcuts() {
    this.http
      .get(BACKEND_ROOT+'/fimp/api/vinculum/shortcuts')
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
        this.shortcuts = result
      });
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
