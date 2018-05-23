import {Component, EventEmitter, Input, OnInit, Output} from "@angular/core";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../globals";
import {TableContextRec} from "./flow-context.component";
import {RecordEditorDialog} from "./record-editor-dialog.component";
import {MatDialog} from "@angular/material";

export class ContextVariable {

  Name : string;
  Type : string;
  Value : any;
  isGlobal : boolean;
}

@Component({
  selector: 'variable-selector',
  templateUrl: './variable-selector.html',
  // styleUrls: ['./locations.component.css']
})
export class VariableSelectorComponent implements OnInit {
    @Input() variableName : string;
    @Input() isGlobal : boolean;
    @Input() label : string;
    @Input() flowId:string;
    @Output() onSelect = new EventEmitter<ContextVariable>();
    vars:ContextVariable[];

  ngOnInit() {
    this.loadContext();
  }
  constructor(private http : Http,public dialog: MatDialog) {
  }

  loadContext() {
    this.vars = [];
    if (this.flowId) {
      this.http
        .get(BACKEND_ROOT+'/fimp/api/flow/context/'+this.flowId)
        .map(function(res: Response){
          let body = res.json();
          return body;
        }).subscribe ((result) => {
        for (var key in result){
          let v = new ContextVariable()
          v.isGlobal = false
          v.Name  = result[key].Name
          v.Type = result[key].Variable.ValueType;
          v.Value = result[key].Variable.Value;
          this.vars.push(v);
        }

      });
    }


    this.http
      .get(BACKEND_ROOT+'/fimp/api/flow/context/global')
      .map(function(res: Response){
        let body = res.json();
        return body;
      }).subscribe ((result) => {
      for (var key in result){
        let v = new ContextVariable()
        v.isGlobal = true
        v.Name  = result[key].Name;
        v.Type = result[key].Variable.ValueType;
        v.Value = result[key].Variable.Value;
        this.vars.push(v);
      }
    });
  }

  showContextVariableDialog(ctxRec:TableContextRec) {
    var ctxRec = new TableContextRec();
    if (this.isGlobal)
      ctxRec.FlowId == "global";
    else
      ctxRec.FlowId = this.flowId;
    let dialogRef = this.dialog.open(RecordEditorDialog,{
      width: '450px',
      data:ctxRec
    });
    dialogRef.afterClosed().subscribe(result => {
      if (result)
      {
        this.variableName = result.Name
        this.loadContext();
      }
    });
  }

  getVariableByName(name:string,isGlobal:boolean):ContextVariable {
    for (let v of this.vars) {
      if(v.Name == name && v.isGlobal == isGlobal) {
        return v;
      }
    }
    return null;

  }

  onSelected() {
      var event = new ContextVariable();
     // event.Name = this.variableName;
     // event.isGlobal = this.isGlobal;
     if(this.variableName=="") {
       var event = new ContextVariable();
       event.Name = "";
       this.onSelect.emit(event);
     }
     event = this.getVariableByName(this.variableName,this.isGlobal)
     if (event) {
       this.onSelect.emit(event);
     }

  }

}
