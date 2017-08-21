import { RouterModule, Routes } from '@angular/router';
import { NgModule } from "@angular/core";
import {ThingsComponent} from './things/things.component'
import {ServicesComponent} from './services/services.component'
import {LocationsComponent} from './locations/locations.component'

@NgModule({
  imports: [RouterModule.forChild([
    { path: 'registry/things', component: ThingsComponent },
    { path: 'registry/services', component: ServicesComponent },
    { path: 'registry/locations', component: LocationsComponent },
  ])],
  exports: [RouterModule]
})
export class RegistryRoutingModule {}