import { Component, OnInit,Inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import {MdDialog, MdDialogRef} from '@angular/material';
import {MD_DIALOG_DATA} from '@angular/material';
// export const BACKEND_ROOT = "http://localhost:8081"
export const BACKEND_ROOT = ""
@Component({
  selector: 'app-flow-editor',
  templateUrl: './flow-editor.component.html',
  styleUrls: ['./flow-editor.component.css']
})
export class FlowEditorComponent implements OnInit {
  flow :Flow;
  selectedNewNodeType:string;
  constructor(private route: ActivatedRoute,private http : Http,public dialog: MdDialog) {
    this.flow = new Flow();
   }

  ngOnInit() {
    let id  = this.route.snapshot.params['id'];
    this.loadFlow(id);
  }
 
  loadFlow(id:string) {
     this.http
      .get(BACKEND_ROOT+'/fimp/flow/definition/'+id)
      .map(function(res: Response){
        let body = res.json();
        //console.log(body.Version);
        return body;
      }).subscribe ((result) => {
         this.flow = result;
         console.dir(this.flow)
         console.log(this.flow.Name)
      });
  }
  saveFlow() {
    console.dir(this.flow)
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({headers:headers});
    this.http
      .post(BACKEND_ROOT+'/fimp/flow/definition/'+this.flow.Id,JSON.stringify(this.flow),  options )
      .subscribe ((result) => {
         console.log("Flow was saved");
      });

 }
 getNewNodeId():string {
   let id = 0;
   let maxId = 0;
   this.flow.Nodes.forEach(element => {
     id = parseInt(element.Id);
     if (id > maxId) {
       maxId = id ;
     }
   });
   maxId++;
   return maxId+"";
 } 
 deleteNode(node:MetaNode){
   let index: number = this.flow.Nodes.indexOf(node);
    if (index !== -1) {
        this.flow.Nodes.splice(index, 1);
    }  
 } 
 cloneNode(node:MetaNode){
   let cloneNode = new MetaNode();
   
   Object.assign(cloneNode,node);
   cloneNode.Id = this.getNewNodeId();
   this.flow.Nodes.push(cloneNode);
 } 
 addNode(){
    console.dir(this.selectedNewNodeType)
    let node  = new MetaNode()
    node.Id = this.getNewNodeId();
    node.Type = this.selectedNewNodeType;
    node.Address = ""
    node.Service = ""
    node.ServiceInterface = ""
    node.SuccessTransition = ""
    node.Config = null

    switch (node.Type){
      case "trigger":
        // nothing to add yet
        break;
      case "action":
        node.Config = {}; 
        node.Config["Value"] = true;
        node.Config["ValueType"] = "bool";
        break;
      case "if":
        node.Config = {}; 
        node.Config["TrueTransition"] = ""
        node.Config["FalseTransition"] = ""
        node.Config["Expression"] = [];
        let expr = {}
        expr["Operand"] = "eq";
        expr["Value"] = 100;
        expr["ValueType"] = "int";
        expr["BooleanOperator"] = "";
        node.Config["Expression"].push(expr);
        break;
      case "wait":
        node.Config = 1000;
        break;      
    }
    this.flow.Nodes.push(node) 
  }
  showSource() {
    let dialogRef = this.dialog.open(FlowSourceDialog,{
            height: '90%',
            width: '90%',
            data:JSON.stringify(this.flow, null, 2)
          });
    
  }
}

export class MetaNode {
  Id               :string;
	Type             :string;
	Label            :string;
	SuccessTransition :string;
	TimeoutTransition :string;
	ErrorTransition   :string;
	Address           :string;
	Service           :string;
	ServiceInterface  :string;
	Config            :any;
}

export class Flow {
    Id :string ;
    Name : string ;
    Description : string ;
    Nodes : MetaNode[] ;
}

@Component({
  selector: 'flow-source-dialog',
  templateUrl: 'flow-source-dialog.html',
})
export class FlowSourceDialog {
  constructor(@Inject(MD_DIALOG_DATA) public data: any) {}
}