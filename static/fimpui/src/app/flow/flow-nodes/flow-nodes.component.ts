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
              element.Service = result.service_name
              element.ServiceInterface = result.intf_msg_type
              element.Address = result.intf_address
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
      if (result)
        this.nodes.forEach(element => {
            if (element.Id==nodeId) {
              element.Service = result.service_name
              element.ServiceInterface = result.intf_msg_type
              element.Address = result.intf_address
              element.Config.ValueType =  msgTypeToValueTypeMap[element.ServiceInterface]
            }
        });
    });      
  }
}