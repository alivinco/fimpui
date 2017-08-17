package registry

import (
	"strings"
	"github.com/asdine/storm"
	gobcodec "github.com/asdine/storm/codec/gob"
	log "github.com/Sirupsen/logrus"
	"github.com/asdine/storm/q"
	"encoding/gob"
)

type ThingRegistryStore struct {
	thingRegistryStoreFile string
	thingRegistry          ThingRegistry
	db *storm.DB
}



func NewThingRegistryStore(storeFile string) *ThingRegistryStore {
	store := ThingRegistryStore{thingRegistryStoreFile: storeFile}
	store.Connect()
	return &store
}

func (st *ThingRegistryStore) Connect() error {
	var err error
	gob.Register([]interface {}{})
	st.db, err = storm.Open(st.thingRegistryStoreFile, storm.Codec(gobcodec.Codec))
	if err != nil {
		log.Error("<Reg> Can't open DB file . Error : ",err)
		return err
	}

	err = st.db.Init(&Thing{})
	if err != nil {
		log.Error("<Reg> Can't Init Things . Error : ",err)
		return err
	}

	err = st.db.Init(&Location{})
	if err != nil {
		log.Error("<Reg> Can't Init Things . Error : ",err)
		return err
	}

	return nil

}

func (st *ThingRegistryStore)Disconnect() {
	st.db.Close()
}

func (st *ThingRegistryStore) GetThingById(Id ID) (*Thing, error) {
	var thing Thing
	err := st.db.One("ID", Id, &thing)
	return &thing, err
}

func (st *ThingRegistryStore) GetLocationById(Id ID) (*Location,error) {
	var location Location
	err := st.db.One("ID", Id, &location)
	return &location, err
}

func (st *ThingRegistryStore) GetAllThings() ([]Thing,error) {
	var things []Thing
	err := st.db.All(&things)
	return things , err
}

func (st *ThingRegistryStore) GetAllServices() ([]Service,error) {
	things,err := st.GetAllThings()
	if err != nil {
		return nil,err
	}
	var result []Service
	for i := range things {
		result = append(result,st.thingRegistry.Things[i].Services...)
	}
	return result,nil
}

func (st *ThingRegistryStore) GetAllLocations() ([]Location,error) {
	var locations []Location
	err := st.db.All(&locations)
	return locations , err
}

//func (st *ThingRegistryStore) GetServicesByName(name string ) []Service {
//
//
//	var result []Service
//	for i := range st.thingRegistry.Things {
//		for i2 := range st.thingRegistry.Things[i].Services {
//			if st.thingRegistry.Things[i].Services[i2].Name == name {
//				result = append(result,st.thingRegistry.Things[i].Services[i2])
//			}
//		}
//	}
//	return result
//}

func (st *ThingRegistryStore) GetThingByAddress(technology string, address string) (*Thing, error) {
	var thing Thing
	err := st.db.Select(q.And(q.Eq("Address",address),q.Eq("CommTechnology",technology))).First(&thing)
	return &thing,err
}

func (st *ThingRegistryStore) GetThingByIntegrationId(id string) (*Thing,error) {
	var thing Thing
	err := st.db.Select(q.Eq("IntegrationId",id)).Find(&thing)
	return &thing,err
}

func (st *ThingRegistryStore) GetLocationByIntegrationId(id string) (*Location,error) {
	var location Location
	err := st.db.Select(q.Eq("IntegrationId",id)).Find(&location)
	return &location,err
}

func (st *ThingRegistryStore) GetFlatInterfaces() ([]InterfaceFlatView,error) {
	var result []InterfaceFlatView
	things, err  := st.GetAllThings()
	if err != nil {
		return nil,err
	}
	for thi := range things {
		for si := range things[thi].Services {
			for inti := range things[thi].Services[si].Interfaces{
				flatIntf := InterfaceFlatView{}
				flatIntf.ThingId = things[thi].ID
				flatIntf.ThingAddress = things[thi].Address
				flatIntf.ThingAlias = things[thi].Alias
				flatIntf.ServiceId = things[thi].Services[si].ID
				flatIntf.ServiceName = things[thi].Services[si].Name
				flatIntf.ServiceAlias = things[thi].Services[si].Alias
				flatIntf.ServiceAddress = things[thi].Services[si].Address
				flatIntf.InterfaceType = things[thi].Services[si].Interfaces[inti].Type
				flatIntf.InterfaceMsgType = things[thi].Services[si].Interfaces[inti].MsgType
				//pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:21_0
				prefix := "pt:j1/mt:evt"
				if strings.Contains(prefix+things[thi].Services[si].Interfaces[inti].MsgType,"cmd"){
					prefix = "pt:j1/mt:cmd"
				}
				flatIntf.InterfaceAddress = prefix+things[thi].Services[si].Address
				location,_ := st.GetLocationById(things[thi].Services[si].LocationId)
				if location != nil {
					flatIntf.LocationId = location.ID
					flatIntf.LocationAlias = location.Alias
					flatIntf.LocationType = location.Type

				}
				location,_ = st.GetLocationById(things[thi].LocationId)
				if location != nil {
					if location.Alias != flatIntf.LocationAlias{
						flatIntf.LocationAlias = location.Alias +" "+flatIntf.LocationAlias
					}
					if flatIntf.LocationType == "" {
						flatIntf.LocationType = location.Type
					}
				}

				result = append(result,flatIntf)
			}
		}
	}
	return result,nil
}

func (st *ThingRegistryStore) UpsertThing(thing *Thing) (ID, error) {
	var  err error
	if thing.ID == 0 {
		err = st.db.Save(thing)
	}else {
		err = st.db.Update(thing)
	}

	if err != nil {
		log.Error("<Reg> Can't save thing . Error :",err)
		return 0,err
	}else {
		log.Debug("<Reg> Thing saved ")
	}

	return thing.ID,nil
}

func (st *ThingRegistryStore) UpsertLocation(location *Location) (ID,error) {
	var  err error
	if location.ID == 0 {
		err = st.db.Save(location)
	}else {
		err = st.db.Update(location)
	}

	if err != nil {
		log.Error("Can't save location . Error :",err)
		return 0,err
	}else {
		log.Debug("Location saved ")
	}

	return location.ID,nil
}

func (st *ThingRegistryStore) DeleteThing(id ID) error {
	thing,err := st.GetThingById(id)
	if err != nil {
		return err
	}
	st.db.DeleteStruct(&thing)
	return nil
}

func (st *ThingRegistryStore) DeleteLocation(id ID) error {
	location,err := st.GetLocationById(id)
	if err != nil {
		return err
	}
	st.db.DeleteStruct(&location)
	return nil
}

func (st *ThingRegistryStore) ClearAll() error {
	//st.thingRegistry.Things = st.thingRegistry.Things[:0]
	//st.thingRegistry.Locations = st.thingRegistry.Locations[:0]
	//return st.SaveThingRegistry()
	return nil
}
