import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import {EventLogComponent} from './event-log/event-log.component'


@NgModule({
  imports: [RouterModule.forChild([
    { path: 'stats/event-log', component: EventLogComponent },
  ])],
  exports: [RouterModule]
})
export class StatsRoutingModule {}