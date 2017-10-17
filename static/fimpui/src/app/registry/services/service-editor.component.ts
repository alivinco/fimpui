import { Component, OnInit,Inject } from '@angular/core';
import { Http, Response,URLSearchParams,RequestOptions,Headers }  from '@angular/http';
import {MdDialog, MdDialogRef,MdSnackBar} from '@angular/material';
import {MD_DIALOG_DATA} from '@angular/material';

@Component({
    selector: 'service-editor-dialog',
    templateUrl: 'service-editor-dialog.html',
  })
  export class ServiceEditorDialog {
    value : any;
    valueType : string ;
    service : any;
    
    constructor(public dialogRef: MdDialogRef<ServiceEditorDialog>,@Inject(MD_DIALOG_DATA) public data: any,public snackBar: MdSnackBar) {
          
    }
    
    
  }
  