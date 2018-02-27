import {MetaNode} from "../../flow-editor/flow-editor.component";
import {Component, Input, OnInit} from "@angular/core";
import {MatDialog} from "@angular/material";

@Component({
  selector: 'rest-action-node',
  templateUrl: './node.html',
  styleUrls: ['../flow-nodes.component.css']
})
export class RestActionNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  shortcuts:any[];
  constructor(public dialog: MatDialog) { }
  ngOnInit() {

  }
}
