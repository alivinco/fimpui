package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ThingRegistryStore struct {
	thingRegistryStoreFile string
	thingRegistry          ThingRegistry
}

func NewThingRegistryStore(storeFile string) *ThingRegistryStore {
	store := ThingRegistryStore{thingRegistryStoreFile: storeFile}
	store.LoadThingRegistry()
	return &store
}

func (st *ThingRegistryStore) LoadThingRegistry() error {
	if _, err := os.Stat(st.thingRegistryStoreFile); os.IsNotExist(err) {
		st.SaveThingRegistry()
	}
	file, err := ioutil.ReadFile(st.thingRegistryStoreFile)
	if err != nil {
		fmt.Println("Can't open DB file.")
		return err
	}
	reg := ThingRegistry{}
	err = json.Unmarshal(file, &reg)
	if err != nil {
		fmt.Println("Can't unmarshel DB file.")
		return err
	}
	st.thingRegistry = reg
	return nil

}

func (st *ThingRegistryStore) SaveThingRegistry() error {

	data, err := json.Marshal(st.thingRegistry)
	if err != nil {
		return err
	}
	ioutil.WriteFile(st.thingRegistryStoreFile, data, 0644)
	return nil
}

func (st *ThingRegistryStore) getNewId() ID {
	var maxId ID
	if len(st.thingRegistry.Things) == 0 {
		return 1
	}
	for i := range st.thingRegistry.Things {
		if st.thingRegistry.Things[i].Id > maxId {
			maxId = st.thingRegistry.Things[i].Id
		}
	}
	return maxId + 1
}

func (st *ThingRegistryStore) getNewLocationId() ID {
	var maxId ID
	if len(st.thingRegistry.Locations) == 0 {
		return 1
	}
	for i := range st.thingRegistry.Locations {
		if st.thingRegistry.Locations[i].Id > maxId {
			maxId = st.thingRegistry.Locations[i].Id
		}
	}
	return maxId + 1
}

func (st *ThingRegistryStore) GetThingById(Id ID) (*Thing, error) {
	for i := range st.thingRegistry.Things {
		if st.thingRegistry.Things[i].Id == Id {
			return &st.thingRegistry.Things[i], nil
		}
	}
	return nil, nil
}

func (st *ThingRegistryStore) GetLocationById(Id ID) *Location {
	for i := range st.thingRegistry.Locations {
		if st.thingRegistry.Locations[i].Id == Id {
			return &st.thingRegistry.Locations[i]
		}
	}
	return nil
}

func (st *ThingRegistryStore) GetAllThings() []Thing {
	st.LoadThingRegistry()
	return st.thingRegistry.Things
}

func (st *ThingRegistryStore) GetAllServices() []Service {
	var result []Service
	for i := range st.thingRegistry.Things {
		result = append(result,st.thingRegistry.Things[i].Services...)
	}
	return result
}

func (st *ThingRegistryStore) GetAllLocations() []Location {
	return st.thingRegistry.Locations
}

func (st *ThingRegistryStore) GetServicesByName(name string ) []Service {
	var result []Service
	for i := range st.thingRegistry.Things {
		for i2 := range st.thingRegistry.Things[i].Services {
			if st.thingRegistry.Things[i].Services[i2].Name == name {
				result = append(result,st.thingRegistry.Things[i].Services[i2])
			}
		}
	}
	return result
}

func (st *ThingRegistryStore) GetThingByAddress(technology string, address string) (*Thing, error) {
	for i := range st.thingRegistry.Things {
		if st.thingRegistry.Things[i].Address == address && st.thingRegistry.Things[i].CommTechnology == technology {
			return &st.thingRegistry.Things[i], nil
		}
	}
	return nil, nil
}

func (st *ThingRegistryStore) GetThingByIntegrationId(id string) *Thing {
	for i := range st.thingRegistry.Things {
		if st.thingRegistry.Things[i].IntegrationId == id {
			return &st.thingRegistry.Things[i]
		}
	}
	return nil
}

func (st *ThingRegistryStore) GetLocationByIntegrationId(id string) *Location {
	for i := range st.thingRegistry.Locations {
		if st.thingRegistry.Locations[i].IntegrationId == id {
			return &st.thingRegistry.Locations[i]
		}
	}
	return nil
}

func (st *ThingRegistryStore) UpsertThing(thing Thing) error {
	exThing, err := st.GetThingByAddress(thing.CommTechnology, thing.Address)
	if err != nil {
		return err
	}
	if exThing == nil {
		thing.Id = st.getNewId()
		st.thingRegistry.Things = append(st.thingRegistry.Things, thing)
	} else {
		thing.Id = exThing.Id
		*exThing = thing
	}

	st.SaveThingRegistry()
	return nil
}

func (st *ThingRegistryStore) UpsertLocation(location Location) error {
	exLocation := st.GetLocationById(location.Id)

	if exLocation == nil {
		location.Id = st.getNewId()
		st.thingRegistry.Locations = append(st.thingRegistry.Locations, location)
	} else {
		location.Id = location.Id
		*exLocation = location
	}
	st.SaveThingRegistry()
	return nil
}

func (st *ThingRegistryStore) ClearAll() error {
	st.thingRegistry.Things = st.thingRegistry.Things[:0]
	st.thingRegistry.Locations = st.thingRegistry.Locations[:0]
	return st.SaveThingRegistry()
}
