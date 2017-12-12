import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FlowOverviewComponent } from './flow-overview/flow-overview.component';
import { FlowContextComponent } from './flow-context/flow-context.component';
import { FlowEditorComponent, FlowSourceDialog, FlowRunDialog, ServiceLookupDialog,ContextDialog } from './flow-editor/flow-editor.component';
import { FlowNodesComponent ,ActionNodeComponent,SetVariableNodeComponent,TimeTriggerNodeComponent } from './flow-nodes/flow-nodes.component';
import { TriggerNodeComponent ,CounterNodeComponent} from './flow-nodes/flow-nodes.component';
import { ReceiveNodeComponent } from './flow-nodes/flow-nodes.component';
import { FlowRoutingModule } from "app/flow/flow-routing.module";
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { RegistryModule} from 'app/registry/registry.module'
import { VariableElementComponent} from 'app/flow/flow-nodes/ui-elements/ui-elements.component'
import { CdkTableModule } from '@angular/cdk/table';
import { MatTableModule,
  MatFormFieldModule,
  MatInputModule,
  MatButtonModule,
  MatChipsModule,
  MatIconModule,
  MatSliderModule,
  MatCheckboxModule,
  MatListModule,
  MatSelectModule, 
  MatOptionModule,
  MatDialogModule,
  MatCardModule,
  MatSidenavModule,
  MatRadioModule,
  MatExpansionModule,
  MatTabsModule,
  MatCheckbox} from '@angular/material';


@NgModule({
  imports: [
    CommonModule,
    FlowRoutingModule,
    MatInputModule,
    MatButtonModule,
    MatFormFieldModule, 
    MatTableModule,
    MatChipsModule,
    MatOptionModule,
    MatSelectModule,
    MatListModule,
    MatIconModule,
    MatSliderModule,
    MatCheckboxModule,
    MatDialogModule,
    MatCardModule,
    MatSidenavModule,
    MatRadioModule,
    MatExpansionModule,
    FormsModule,
    HttpModule,
    RegistryModule,
    MatTabsModule,
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
