import {Component, EventEmitter, Input, OnInit, Output} from "@angular/core";
import {Http, Response} from "@angular/http";
import {BACKEND_ROOT} from "../../globals";
import {TableContextRec} from "./flow-context.component";
import {RecordEditorDialog} from "./record-editor-dialog.component";
import {MatDialog} from "@angular/material";

export class ContextVariable {

  variableName : string;
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
    localVars:any;
    globalVars:any;

  ngOnInit() {
    this.loadContext();
  }
  constructor(private http : Http,public dialog: MatDialog) {
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

  showContextVariableDialog(ctxRec:TableContextRec) {
    var ctxRec = new TableContextRec();
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

  onSelected() {
     var event = new ContextVariable();
     event.variableName = this.variableName;
     event.isGlobal = this.isGlobal;
     this.onSelect.emit(event);
  }

}
