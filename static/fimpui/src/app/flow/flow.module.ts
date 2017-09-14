import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlowOverviewComponent } from './flow-overview/flow-overview.component';
import { FlowContextComponent } from './flow-context/flow-context.component';
import { FlowEditorComponent, FlowSourceDialog, FlowRunDialog, ServiceLookupDialog,ContextDialog } from './flow-editor/flow-editor.component';
import { FlowNodesComponent ,ActionNodeComponent,SetVariableNodeComponent,TimeTriggerNodeComponent } from './flow-nodes/flow-nodes.component';
import { TriggerNodeComponent ,CounterNodeComponent} from './flow-nodes/flow-nodes.component';
import { ReceiveNodeComponent } from './flow-nodes/flow-nodes.component';
import { FlowRoutingModule } from "app/flow/flow-routing.module";
import { MaterialModule } from '@angular/material';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RegistryModule} from 'app/registry/registry.module'
import { VariableElementComponent} from 'app/flow/flow-nodes/ui-elements/ui-elements.component'
import { CdkTableModule } from '@angular/cdk';


@NgModule({
  imports: [
    CommonModule,
    FlowRoutingModule,
    MaterialModule,
    FormsModule,
    HttpModule,
    RegistryModule,
    CdkTableModule
  ],
  declarations: [
     FlowOverviewComponent,
     FlowContextComponent,
     FlowEditorComponent,
     FlowSourceDialog,
     FlowRunDialog,
     FlowNodesComponent,
     ActionNodeComponent,
     TriggerNodeComponent,
     ReceiveNodeComponent,
     SetVariableNodeComponent,
     ServiceLookupDialog, 
     ContextDialog,
     VariableElementComponent,
     TimeTriggerNodeComponent,
     CounterNodeComponent
  ],
  entryComponents: [FlowSourceDialog,FlowRunDialog,ServiceLookupDialog,ContextDialog]  
})
export class FlowModule { }
