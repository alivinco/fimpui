package registry

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type ThingRegistryStore struct {
	thingRegistryStoreFile string
	thingRegistry ThingRegistry
}

func NewThingRegistryStore(storeFile string) *ThingRegistryStore {
	store := ThingRegistryStore{thingRegistryStoreFile:storeFile}
	store.LoadThingRegistry()
	return &store
}

func (st * ThingRegistryStore) LoadThingRegistry() (error) {
	if _, err := os.Stat(st.thingRegistryStoreFile); os.IsNotExist(err) {
		st.SaveThingRegistry()
	}
	file , err := ioutil.ReadFile(st.thingRegistryStoreFile)
	if err != nil {
		fmt.Println("Can't open DB file.")
		return err
	}
	reg := ThingRegistry{}
	err = json.Unmarshal(file,&reg)
	if err != nil {
		fmt.Println("Can't unmarshel DB file.")
		return err
	}
	st.thingRegistry = reg
	return nil

}

func (st * ThingRegistryStore) SaveThingRegistry() error {

	data ,err := json.Marshal(st.thingRegistry)
	if err != nil {
		return err
	}
	ioutil.WriteFile(st.thingRegistryStoreFile,data,0644)
	return nil
}

func (st * ThingRegistryStore) getNewId() ID {
	var maxId ID
	if len(st.thingRegistry) == 0 {
		return 1
	}
	for i := range st.thingRegistry {
		if st.thingRegistry[i].Id > maxId {
			maxId = st.thingRegistry[i].Id
		}
	}
	return maxId+1
}


func (st * ThingRegistryStore) GetThingById(Id ID) (*Thing,error) {
	for i := range st.thingRegistry {
		if st.thingRegistry[i].Id == Id {
			return &st.thingRegistry[i],nil
		}
	}
	return nil,nil
}

func (st * ThingRegistryStore) GetAllThings()(ThingRegistry){
	st.LoadThingRegistry()
	return st.thingRegistry
}

func (st * ThingRegistryStore) GetThingByAddress(technoloy string , address string) (*Thing,error) {
	for i := range st.thingRegistry {
		if st.thingRegistry[i].Address == address && st.thingRegistry[i].CommTechnology == technoloy {
			return &st.thingRegistry[i],nil
		}
	}
	return nil,nil
}

func (st * ThingRegistryStore) UpsertThing(thing Thing) error {
	exThing ,err := st.GetThingByAddress(thing.CommTechnology,thing.Address)
	if err != nil {
		return err
	}
	if exThing == nil {
		thing.Id = st.getNewId()
		st.thingRegistry = append(st.thingRegistry, thing)
	}else {
		thing.Id = exThing.Id
		*exThing = thing
 	}

	st.SaveThingRegistry()
	return nil
}

func (st * ThingRegistryStore) ClearAll() error {
	st.thingRegistry = st.thingRegistry[:0]
	return st.SaveThingRegistry();
}
