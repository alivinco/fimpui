import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CdkTableModule } from '@angular/cdk/table';
import { HttpModule } from '@angular/http';
import { EventLogComponent } from './event-log/event-log.component';
import { StatsRoutingModule } from "./stats-routing.module";
import { last } from 'rxjs/operator/last';
import { MatTableModule,
  MatFormFieldModule,
  MatInputModule,
  MatPaginator,
  MatIconModule,
  MatSliderModule,
  MatListModule,
  MatDialogModule,
  MatExpansionModule} from '@angular/material';


@NgModule({
  imports: [
    CommonModule,
    CdkTableModule,
    MatTableModule,
    MatInputModule,
    MatFormFieldModule, 
    MatTableModule,
    MatListModule,
    MatIconModule,
    MatSliderModule,
    MatDialogModule,
    MatExpansionModule,
     HttpModule,
    StatsRoutingModule
  ],
  declarations: [EventLogComponent]
})
export class StatsModule { }
 