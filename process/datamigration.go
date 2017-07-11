package process

import (
	"github.com/alivinco/fimpui/integr/fhcore"
	"github.com/alivinco/fimpui/registry"
)

func LoadVinculumDeviceInfoToStore(thingRegistryStore *registry.ThingRegistryStore, vincClient *fhcore.VinculumClient) error {

	commTechMap := map[string]string{"zwave-ad": "zw", "ikea": "ikea"}

	msg, err := vincClient.GetMessage([]string{"device"})
	if err != nil {
		return err
	}
	devices := msg.Msg.Data.Param.Device
	for i := range devices {
		if devices[i].Fimp.Address != "" {
			thing := registry.Thing{}
			thing.Address = devices[i].Fimp.Address
			thing.CommTechnology = commTechMap[devices[i].Fimp.Adapter]
			thing.Alias = devices[i].Client.Name
			thingRegistryStore.UpsertThing(thing)
		}
	}
	return nil
}
