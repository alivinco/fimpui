import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CdkTableModule } from '@angular/cdk/table';
import { HttpModule } from '@angular/http';
import { EventLogComponent,EventsPerDeviceChart } from './event-log/event-log.component';
import { SystemMetricsComponent } from './system-metrics/system-metrics.component';

import { StatsRoutingModule } from "./stats-routing.module";
import { last } from 'rxjs/operator/last';
import { ChartsModule } from 'ng2-charts';
import { MatTableModule,
  MatSortModule,
  MatFormFieldModule,
  MatInputModule,
  MatPaginator,
  MatIconModule,
  MatSliderModule,
  MatListModule,
  MatDialogModule,
  MatTabsModule,
  MatExpansionModule} from '@angular/material';


@NgModule({
  imports: [
    CommonModule,
    CdkTableModule,
    MatTableModule,
    MatSortModule,
    MatInputModule,
    MatFormFieldModule, 
    MatTableModule,
    MatListModule,
    MatIconModule,
    MatSliderModule,
    MatDialogModule,
    MatExpansionModule,
    MatTabsModule,
     HttpModule,
     ChartsModule,
    StatsRoutingModule
  ],
  exports:[EventsPerDeviceChart],
  declarations: [EventLogComponent,SystemMetricsComponent,EventsPerDeviceChart]
})
export class StatsModule { }
 