import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MaterialModule, 
         MdTableModule } from '@angular/material';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { ThingsComponent } from './things/things.component';
import { ServicesComponent } from './services/services.component';
import { LocationsComponent } from './locations/locations.component';
import { RegistryRoutingModule } from "./registry-routing.module";
import {CdkTableModule} from '@angular/cdk';


@NgModule({
  imports: [
    CommonModule,
    RegistryRoutingModule,
    MaterialModule,
    FormsModule,
    HttpModule,
    CdkTableModule,
    
  ],
  exports:[],

  declarations: [
    ThingsComponent,
    ServicesComponent,
    LocationsComponent
  ]
})
export class RegistryModule { }
