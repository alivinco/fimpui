import { Component, OnInit ,Input } from '@angular/core';
import { MetaNode, ServiceLookupDialog } from "../flow-editor/flow-editor.component";
import {MdDialog, MdDialogRef} from '@angular/material';
import { msgTypeToValueTypeMap } from "app/things-db/mapping";

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
  constructor(public dialog: MdDialog) { }

  ngOnInit() { 
  }
  serviceLookupDialog(nodeId:string) {
    let dialogRef = this.dialog.open(ServiceLookupDialog,{
            width: '95%'
          });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
        this.nodes.forEach(element => {
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
              element.ServiceInterface = result.intfMsgType
              element.Address = result.intfAddress
              //element.Config.ValueFilter.ValueType =  msgTypeToValueTypeMap[element.ServiceInterface]
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