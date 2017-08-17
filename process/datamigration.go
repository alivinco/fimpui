package process

import (
	"github.com/alivinco/fimpui/integr/fhcore"
	"github.com/alivinco/fimpui/registry"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

func LoadVinculumDeviceInfoToStore(thingRegistryStore *registry.ThingRegistryStore, vincClient *fhcore.VinculumClient) error {

	commTechMap := map[string]string{"zwave-ad": "zw", "ikea": "ikea"}
	vincToServiceNameMap := map[string]string{
		"power":"out_bin_switch",
		"dimValue":"out_lvl_switch",
		"batteryPercentage":"battery",
		"illuminance":"sensor_lumin",
		"presence":"sensor_presence",
		"temperature":"sensor_temp",
		"targetTemperature":"thermostat",
		"openState":"sensor_contact",

	}

	msg, err := vincClient.GetMessage([]string{"device","room"})
	if err != nil {
		log.Errorf("Vinculum client error :",err)
		return err
	}

	rooms := msg.Msg.Data.Param.Room
	for i := range rooms {
		location,err := thingRegistryStore.GetLocationByIntegrationId(strconv.Itoa(rooms[i].ID))
		if err != nil {
			newLocation := registry.Location{}
			newLocation.IntegrationId = strconv.Itoa(rooms[i].ID)
			if rooms[i].Client.Name == "" {
				newLocation.Alias = rooms[i].Type
			}else {
				newLocation.Alias = rooms[i].Client.Name
			}
			newLocation.Type = "room"
			thingRegistryStore.UpsertLocation(&newLocation)
			log.Infof("<VincMigration> Location %s was added. New ID = %d ",newLocation.Alias,newLocation.ID)
		} else {
			if rooms[i].Client.Name == "" {
				location.Alias = rooms[i].Type
			}else {
				location.Alias = rooms[i].Client.Name
			}
			thingRegistryStore.UpsertLocation(location)
		}
	}

	devices := msg.Msg.Data.Param.Device
	for i := range devices {
		if devices[i].Fimp.Address != "" {
			tech := commTechMap[devices[i].Fimp.Adapter]
			thing ,err := thingRegistryStore.GetThingByAddress(tech,devices[i].Fimp.Address)
			if err != nil {
				log.Infof("Device %s not found in registry. Generate inclusion report first",devices[i].Client.Name)
				//newThing := registry.Thing{}
				//newThing.Address = devices[i].Fimp.Address
				//newThing.CommTechnology = tech
				//newThing.Alias = devices[i].Client.Name
				//newThing.IntegrationId = strconv.Itoa(devices[i].ID)
				//thingRegistryStore.UpsertThing(newThing)
			}else {
				thing.Alias = devices[i].Client.Name
			    for si := range thing.Services {
					for _,group := range thing.Services[si].Groups {
						if group == devices[i].Fimp.Group {
							for k,_ := range devices[i].Param {
								if thing.Services[si].Name == vincToServiceNameMap[k] {
									thing.Services[si].IntegrationId = strconv.Itoa(devices[i].ID)
									thing.Services[si].Alias = devices[i].Client.Name
									location,_ := thingRegistryStore.GetLocationByIntegrationId(strconv.Itoa(devices[i].Room))
									if location != nil {
										thing.Services[si].LocationId = location.ID
										thing.LocationId = location.ID
									}
								}
							}
						}
					}
				}
				thingRegistryStore.UpsertThing(thing)

			}

		}
	}



	return nil
}

