import {MetaNode, ServiceLookupDialog} from "../../flow-editor/flow-editor.component";
import {Component, Input, OnInit} from "@angular/core";
import {MatDialog} from "@angular/material";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../../globals";
import {ContextVariable} from "../../flow-context/variable-selector.component";

@Component({
  selector: 'exec-node',
  templateUrl: './node.html',
  styleUrls: ['../flow-nodes.component.css']
})
export class ExecNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  @Input() flowId:string;
  localVars:any;
  globalVars:any;
  constructor(public dialog: MatDialog,private http : Http) {
  }

  ngOnInit() {
    this.loadDefaultConfig();

  }

  loadDefaultConfig() {
    if (this.node.Config==null) {
      this.node.Config = {
          "ExecType":"cmd", // cmd , sh-cmd , python , script
          "Command":"",
          "ScriptBody": "",
          "InputVariableName":"",
          "IsInputVariableGlobal": false,
          "OutputVariableName":"",
          "IsOutputVariableGlobal": false,
          "IsOutputJson": false,
          "IsInputJson":  false
        }
    }
  }

  inputVariableSelected(cvar:ContextVariable) {
    this.node.Config.InputVariableName = cvar.variableName;
    this.node.Config.IsInputVariableGlobal = cvar.isGlobal;
  }
  outputVariableSelected(cvar:ContextVariable) {
    this.node.Config.OutputVariableName = cvar.variableName;
    this.node.Config.IsOutputVariableGlobal = cvar.isGlobal;
  }

}
