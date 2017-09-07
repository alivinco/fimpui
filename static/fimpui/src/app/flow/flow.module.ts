import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlowOverviewComponent } from './flow-overview/flow-overview.component';
import { FlowEditorComponent, FlowSourceDialog, FlowRunDialog, ServiceLookupDialog,ContextDialog } from './flow-editor/flow-editor.component';
import { FlowNodesComponent ,ActionNodeComponent,SetVariableNodeComponent,TimeTriggerNodeComponent } from './flow-nodes/flow-nodes.component';
import { ReceiveNodeComponent } from './flow-nodes/flow-nodes.component';
import { FlowRoutingModule } from "app/flow/flow-routing.module";
import { MaterialModule } from '@angular/material';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RegistryModule} from 'app/registry/registry.module'
import { VariableElementComponent} from 'app/flow/flow-nodes/ui-elements/ui-elements.component'


@NgModule({
  imports: [
    CommonModule,
    FlowRoutingModule,
    MaterialModule,
    FormsModule,
    HttpModule,
    RegistryModule
  ],
  declarations: [
     FlowOverviewComponent,
     FlowEditorComponent,
     FlowSourceDialog,
     FlowRunDialog,
     FlowNodesComponent,
     ActionNodeComponent,
     ReceiveNodeComponent,
     SetVariableNodeComponent,
     ServiceLookupDialog, 
     ContextDialog,
     VariableElementComponent,
     TimeTriggerNodeComponent,
  ],
  entryComponents: [FlowSourceDialog,FlowRunDialog,ServiceLookupDialog,ContextDialog]  
})
export class FlowModule { }
