import { RouterModule, Routes } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import { AppComponent } from './app.component';

import { MaterialModule } from '@angular/material';
import { ZwaveManComponent , AddDeviceDialog } from './zwave-man/zwave-man.component';
import { TimelineComponent } from './timeline/timeline.component';
import { FimpService} from './fimp.service';
import { ThingsDbService} from './things-db.service';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { ThingIntfUiComponent , KeysPipe }from './thing-intf-ui/thing-intf-ui.component'

import {
  MqttMessage,
  MqttModule,
  MqttService
} from 'angular2-mqtt';
import { ThingViewComponent } from './thing-view/thing-view.component';
import { ThingsTableComponent } from './things-table/things-table.component';
import { SettingsComponent } from './settings/settings.component';


const appRoutes: Routes = [
  { path: 'settings', component: SettingsComponent },
  { path: 'zwave-man', component: ZwaveManComponent },
  { path: 'thing-view/:id', component: ThingViewComponent },
  { path: '',redirectTo:'/zwave-man',pathMatch: 'full'}
];
let mqttHost : string = "localhost";
let mqttPort : number = 8083
if (localStorage.getItem("mqttHost")!= null){
      mqttHost = localStorage.getItem("mqttHost");
}
if (localStorage.getItem("mqttPort")!= null){
      mqttPort = parseInt(localStorage.getItem("mqttPort"));
} 
console.log("Port:"+localStorage.getItem("mqttPort"));
export const MQTT_SERVICE_OPTIONS = {
  connectOnCreate: true,
  // hostname: '172.16.118.131',
  hostname:mqttHost,
  // hostname:"aleks.local",
  port: mqttPort,
  path: '/mqtt'
};

export function mqttServiceFactory() {
  console.log("Starting mqttService");
  let mqs =  new MqttService(MQTT_SERVICE_OPTIONS);
  return mqs;
}

@NgModule({
  declarations: [
    AppComponent,
    ZwaveManComponent,
    AddDeviceDialog,
    TimelineComponent,
    ThingViewComponent,
    ThingsTableComponent,
    SettingsComponent,
    ThingIntfUiComponent,
    KeysPipe
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    BrowserAnimationsModule,
    MaterialModule.forRoot(),
    MqttModule.forRoot({
      provide: MqttService,
      useFactory: mqttServiceFactory
    }),
    RouterModule.forRoot(appRoutes),
    NgxDatatableModule
    
  ],
  providers: [FimpService,ThingsDbService],
  entryComponents:[AddDeviceDialog],
  bootstrap: [AppComponent]
})
export class AppModule { }
