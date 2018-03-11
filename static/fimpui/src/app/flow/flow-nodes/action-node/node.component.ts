import {MetaNode, ServiceLookupDialog} from "../../flow-editor/flow-editor.component";
import {Component, Input, OnInit} from "@angular/core";
import {MatDialog} from "@angular/material";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../../globals";

@Component({
  selector: 'action-node',
  templateUrl: './node.html',
  styleUrls: ['../flow-nodes.component.css']
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
        .get(BACKEND_ROOT+'/fimp/api/flow/context/'+this.flowId)
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
  variableSelected(event:any,config:any,isGlobal:boolean){
    // if (config.VariableName.indexOf("__global__")!=-1) {
    //   config.VariableName = config.VariableName.replace("__global__","");
    //   config.VariableIsGlobal = true;
    // }
    config.IsVariableGlobal = isGlobal;

  }

}
