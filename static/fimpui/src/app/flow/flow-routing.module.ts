import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import {FlowEditorComponent} from './flow-editor/flow-editor.component'
import {FlowOverviewComponent} from './flow-overview/flow-overview.component'

@NgModule({
  imports: [RouterModule.forChild([
    { path: 'flow/overview', component: FlowOverviewComponent },
    { path: 'flow/editor/:id', component: FlowEditorComponent },
  ])],
  exports: [RouterModule]
})
export class FlowRoutingModule {}