import { Component, OnInit,Inject } from '@angular/core';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import { MatDialog, MatDialogRef,MatSnackBar} from '@angular/material';
import { MAT_DIALOG_DATA} from '@angular/material';
import { BACKEND_ROOT } from "app/globals";
import {Variable} from "../flow-editor/flow-editor.component";
import {TableContextRec} from "./flow-context.component";

class ContextRecord{
  Name : string;
  Description : string;
  UpdatedAt : string;
  Variable : Variable;
}

@Component({
    selector: 'record-editor-dialog',
    templateUrl: 'record-editor-dialog.html',
    styleUrls: ['./flow-context.component.css']
  })
  export class RecordEditorDialog {
    ctxRec : TableContextRec;
    constructor(public dialogRef: MatDialogRef<RecordEditorDialog>, @Inject(MAT_DIALOG_DATA) public data: TableContextRec, public snackBar: MatSnackBar, private http : Http) {
          this.ctxRec = data;
          console.dir(data)
    }

    save(){
      let headers = new Headers({ 'Content-Type': 'application/json' });
      let options = new RequestOptions({headers:headers});
      let request = new ContextRecord();
      request.Name = this.ctxRec.Name;
      request.Variable = new Variable();
      request.Variable.Value = this.ctxRec.Value;
      request.Variable.ValueType = this.ctxRec.ValueType;
      request.Description = this.ctxRec.Description;

      this.http
        .post(BACKEND_ROOT+'/fimp/api/flow/context/record',JSON.stringify(request),  options )
        .subscribe ((result) => {
           console.log("Context record was saved");
           this.dialogRef.close("ok");
        });
    }

    delete() {
    this.http
      .delete(BACKEND_ROOT+'/fimp/api/flow/context/record/'+this.ctxRec.Name)
      .subscribe ((result) => {
        this.dialogRef.close("ok");
      });
    }

}
