import {MetaNode} from "../../flow-editor/flow-editor.component";
import {Component, Input, OnInit} from "@angular/core";
import {MatDialog} from "@angular/material";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../../globals";
import {ContextVariable} from "../../flow-context/variable-selector.component";

@Component({
  selector: 'transform-node',
  templateUrl: './node.html',
  styleUrls: ['../flow-nodes.component.css']
})
export class TransformNodeComponent implements OnInit {
  @Input() node :MetaNode;
  @Input() nodes:MetaNode[];
  @Input() flowId:string;
  localVars:any;
  globalVars:any;
  constructor(public dialog: MatDialog,private http : Http) { }
  ngOnInit() {
    this.loadDefaultConfig();
    this.loadContext();
  }

  loadDefaultConfig() {
    if (this.node.Config==null) {
      this.node.Config = {
        "TargetVariableName": "",
        "TargetVariableType": "",
        "IsTargetVariableGlobal": false,
        "TransformType": "calc",
        "Rtype": "var", "IsRVariableGlobal": false,
        "IsLVariableGlobal": false,
        "Operation": "add",
        "RValue": {"ValueType": "int", "Value": 0},
        "RVariableName": "",
        "LVariableName": "",
        "ValueMapping": [],
        "XPathMapping": [],
        "Template": ""
      };
    }
  }

  addValueMapping(node:MetaNode){
    let valueMap = {};
    valueMap["LValue"] = {"ValueType":"int","Value":0};
    valueMap["RValue"] = {"ValueType":"int","Value":0};
    node.Config["ValueMapping"].push(valueMap);
  }

  addXPathMapping(node:MetaNode){
    let valueMap = {
      "Path":"",
      "TargetVariableName":"",
      "TargetVariableType":"",
      "IsTargetVariableGlobal":false,
      "UpdateInputVariable":false
    };
    node.Config["XPathMapping"].push(valueMap);
  }

  targetVariableSelected(cvar:ContextVariable,vmap) {
    vmap.TargetVariableName = cvar.Name;
    vmap.TargetVariableType = cvar.Type;
    vmap.IsTargetVariableGlobal = cvar.isGlobal;
  }

  resultVariableSelected(cvar:ContextVariable) {
    this.node.Config.TargetVariableName = cvar.Name;
    this.node.Config.TargetVariableType = cvar.Type;
    this.node.Config.IsTargetVariableGlobal = cvar.isGlobal;
  }

  lVariableSelected(cvar:ContextVariable) {
    this.node.Config.LVariableName = cvar.Name;
    this.node.Config.IsLVariableGlobal = cvar.isGlobal;
  }

  loadContext() {
    if (this.flowId) {
      this.http
        .get(BACKEND_ROOT+'/fimp/api/flow/context/'+this.flowId)
        .map(function(res: Response){
          let body = res.json();
          return body;
        }).subscribe ((result) => {
        this.localVars = [];
        for (var key in result){
          this.localVars.push(result[key].Name);
        }

      });
    }


    this.http
      .get(BACKEND_ROOT+'/fimp/api/flow/context/global')
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

  deleteMappingRecord(record:any) {
    var i = this.node.Config.ValueMapping.indexOf(record);
    if(i != -1) {
      this.node.Config.ValueMapping.splice(i, 1);
    }
  }

  deleteXPathRecord(record:any) {
    var i = this.node.Config.XPathMapping.indexOf(record);
    if(i != -1) {
      this.node.Config.XPathMapping.splice(i, 1);
    }
  }

  variableSelected(event:any,config:any){



  }
}
