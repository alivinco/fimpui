import { Component, OnInit ,Input } from '@angular/core';
import { MetaNode } from "app/flow-editor/flow-editor.component";

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
  constructor() { }

  ngOnInit() {
  }

}