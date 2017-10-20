import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MaterialModule, 
         MdTableModule } from '@angular/material';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { ThingsComponent } from './things/things.component';
import { ServicesComponent ,ServicesMainComponent,ServiceSelectorWizardComponent } from './services/services.component';
import { LocationsComponent,LocationSelectorWizardComponent } from './locations/locations.component';
import { AdminComponent } from './admin/admin.component';
import { ThingIntfUiComponent,KeysPipe } from 'app/registry/thing-intf-ui/thing-intf-ui.component'

import { RegistryRoutingModule } from "./registry-routing.module";
import { CdkTableModule } from '@angular/cdk';
import { FimpService } from 'app/fimp/fimp.service'
import { ServiceEditorDialog }from 'app/registry/services/service-editor.component'
import { ThingEditorDialog }from 'app/registry/things/thing-editor.component'

@NgModule({
  imports: [
    CommonModule,
    RegistryRoutingModule,
    MaterialModule,
    FormsModule,
    HttpModule,
    CdkTableModule,
  ],
  exports:[ServicesComponent,ThingIntfUiComponent,KeysPipe,ServiceSelectorWizardComponent], 

  declarations: [
    ThingsComponent,
    ServicesComponent,
    ServicesMainComponent,
    LocationsComponent,
    AdminComponent,
    ThingIntfUiComponent,
    ServiceSelectorWizardComponent,
    LocationSelectorWizardComponent,
    KeysPipe,
    ServiceEditorDialog,
    ThingEditorDialog
  ],
  providers:[FimpService],
  entryComponents: [ServiceEditorDialog,ThingEditorDialog]
})
export class RegistryModule { }
