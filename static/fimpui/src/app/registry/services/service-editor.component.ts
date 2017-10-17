import { Component, OnInit,Inject } from '@angular/core';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import { MdDialog, MdDialogRef,MdSnackBar} from '@angular/material';
import { MD_DIALOG_DATA} from '@angular/material';
import { BACKEND_ROOT } from "app/globals";

@Component({
    selector: 'service-editor-dialog',
    templateUrl: 'service-editor-dialog.html',
  })
  export class ServiceEditorDialog {
    locationId : number;
    alias : string;
    serviceId : number;     
    constructor(public dialogRef: MdDialogRef<ServiceEditorDialog>,@Inject(MD_DIALOG_DATA) public data: any,public snackBar: MdSnackBar,private http : Http) {
          console.dir(data)
          this.serviceId = data.id
          this.alias = data.alias
          this.locationId = data.locationId
    }
    onLocationSeleted(locationId:number ) {
        console.log("Location selected = "+locationId)
    }
    save(){
      let headers = new Headers({ 'Content-Type': 'application/json' });
      let options = new RequestOptions({headers:headers});
      let request = {"id":this.serviceId,"alias":this.alias,"location_id":this.locationId}
      this.http
        .post(BACKEND_ROOT+'/fimp/api/registry/service-fields',JSON.stringify(request),  options )
        .subscribe ((result) => {
           console.log("Service fields were saved");
           this.dialogRef.close("ok");
        });
    }
    
  }
  