import {MetaNode} from "../../flow-editor/flow-editor.component";
import {Component, Input, OnInit} from "@angular/core";
import {MatDialog} from "@angular/material";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../../globals";

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
        "TargetVariableName":"","IsTargetVariableGlobal":false,
        "TransformType":"calc","Rtype":"var","IsRVariableGlobal":false,
        "IsLVariableGlobal":false,
        "Operation":"add","RValue":{"ValueType":"int","Value":0},"RVariableName":"","LVariableName":"","ValueMapping":[]};
    }
  }

  addValueMapping(node:MetaNode){
    let valueMap = {};
    valueMap["LValue"] = {"ValueType":"int","Value":0};
    valueMap["RValue"] = {"ValueType":"int","Value":0};
    node.Config["ValueMapping"].push(valueMap);
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

  variableSelected(event:any,config:any){



  }
}
