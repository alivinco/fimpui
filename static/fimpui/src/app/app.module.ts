import { RouterModule, Routes } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule,APP_INITIALIZER } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';
import { BrowserAnimationsModule} from '@angular/platform-browser/animations';
import { AppComponent } from './app.component';
import { MatTableModule,
  MatFormFieldModule,
  MatButtonModule,
  MatInputModule,
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
  MatMenuModule,
  MatToolbarModule,
  MatProgressBarModule,
  MatTabsModule,
  MatPaginatorModule,
  MatCheckbox} from '@angular/material';
import { CdkTableModule } from '@angular/cdk/table';
import { ZwaveManComponent , AddDeviceDialog,RemoveDeviceDialog } from './zwave-man/zwave-man.component';
import { TemplateEditorDialog } from './zwave-man/zwave-man.component';
import { IkeaManComponent } from './ikea-man/ikea-man.component';
import { ZigbeeManComponent } from './zigbee-man/zigbee-man.component';
import { SystemsManComponent } from './systems-man/systems-man.component';
import { TimelineComponent,MsgDetailsDialog } from './timeline/timeline.component';
import { ReportComponent } from './report/report.component';
import { FlightRecorderComponent } from './flight-recorder/flight-recorder.component';
import { FimpService} from './fimp/fimp.service';
import { ThingsDbService } from './things-db.service';
import { ConfigsService } from './configs.service';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
// import { ThingIntfUiComponent , KeysPipe }from './thing-intf-ui/thing-intf-ui.component'
import 'hammerjs';
import {
  MqttMessage,
  MqttModule,
  MqttService
} from 'angular2-mqtt';
import { ThingViewComponent } from './thing-view/thing-view.component';
import { ThingsTableComponent } from './things-table/things-table.component';
import { SettingsComponent } from './settings/settings.component';
import { FlowModule} from './flow/flow.module'
import { StatsModule} from './stats/stats.module'
import { RegistryModule} from './registry/registry.module'

const appRoutes: Routes = [
  { path: 'settings', component: SettingsComponent },
  { path: 'zwave-man', component: ZwaveManComponent },
  { path: 'ikea-man', component: IkeaManComponent },
  { path: 'zigbee-man', component: ZigbeeManComponent },
  { path: 'systems-man', component: SystemsManComponent },
  { path: 'timeline', component: TimelineComponent },
  { path: 'report', component: ReportComponent },
  { path: 'flight-recorder', component: FlightRecorderComponent },
  { path: 'thing-view/:ad/:id', component: ThingViewComponent },
  { path: 'thing-view-addr/:id', component: ThingViewComponent },
  { path: '',redirectTo:'/zwave-man',pathMatch: 'full'}
];
let mqttHost : string = window.location.hostname;
let mqttPort : number = Number(window.location.port);
if (localStorage.getItem("mqttHost")!= null){
      mqttHost = localStorage.getItem("mqttHost");
}else {
  localStorage.setItem("mqttHost",mqttHost);
}
if (localStorage.getItem("mqttPort")!= null){
      mqttPort = parseInt(localStorage.getItem("mqttPort"));
} else {
  localStorage.setItem("mqttPort",String(mqttPort));
}
console.log("Port:"+localStorage.getItem("mqttPort"));
export const MQTT_SERVICE_OPTIONS = {
  connectOnCreate: false
  // hostname:mqttHost,
  // port: mqttPort,
  // path: '/mqtt',
  // username:"5Qm19y",
  // password:"66ldpVL19cab"
};

export function mqttServiceFactory() {
  console.log("Starting mqttService");
  let mqs =  new MqttService(MQTT_SERVICE_OPTIONS);
  return mqs;
}

export function startupServiceFactory(startupService: ConfigsService): Function {
  return () => startupService.load();
}

@NgModule({
  declarations: [
    AppComponent,
    ZwaveManComponent,
    IkeaManComponent,
    ZigbeeManComponent,
    SystemsManComponent,
    AddDeviceDialog,
    RemoveDeviceDialog,
    TimelineComponent,
    ThingViewComponent,
    ThingsTableComponent,
    SettingsComponent,
    ReportComponent,
    FlightRecorderComponent,
    TemplateEditorDialog,
    MsgDetailsDialog
    // ThingIntfUiComponent,
    // KeysPipe,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    HttpClientModule,
    BrowserAnimationsModule,
    // MaterialModule,
    MqttModule.forRoot({
      provide: MqttService,
      useFactory: mqttServiceFactory
    }),
    MatButtonModule,
    MatPaginatorModule,
    MatInputModule,
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
    MatMenuModule,
    MatExpansionModule,
    MatToolbarModule,
    MatProgressBarModule,
    MatTabsModule,
    RouterModule.forRoot(appRoutes),
    NgxDatatableModule,
    FlowModule,
    StatsModule,
    RegistryModule,
    CdkTableModule

  ],
  providers: [FimpService,ThingsDbService,ConfigsService,{
    // Provider for APP_INITIALIZER
    provide: APP_INITIALIZER,
    useFactory: startupServiceFactory,
    deps: [ConfigsService],
    multi: true
}],
  entryComponents:[AddDeviceDialog,RemoveDeviceDialog,TemplateEditorDialog,MsgDetailsDialog],
  bootstrap: [AppComponent]
})
export class AppModule { }
