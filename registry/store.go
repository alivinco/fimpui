package registry

import (
	"encoding/gob"

	log "github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	gobcodec "github.com/asdine/storm/codec/gob"
	"github.com/asdine/storm/q"
)

type ThingRegistryStore struct {
	thingRegistryStoreFile string
	db                     *storm.DB
}

func NewThingRegistryStore(storeFile string) *ThingRegistryStore {
	store := ThingRegistryStore{thingRegistryStoreFile: storeFile}
	store.Connect()
	return &store
}

func (st *ThingRegistryStore) Connect() error {
	var err error
	gob.Register([]interface{}{})
	st.db, err = storm.Open(st.thingRegistryStoreFile, storm.Codec(gobcodec.Codec))
	if err != nil {
		log.Error("<Reg> Can't open DB file . Error : ", err)
		return err
	}

	err = st.db.Init(&Thing{})
	if err != nil {
		log.Error("<Reg> Can't Init Things . Error : ", err)
		return err
	}

	err = st.db.Init(&Service{})
	if err != nil {
		log.Error("<Reg> Can't Init Services . Error : ", err)
		return err
	}

	err = st.db.Init(&Location{})
	if err != nil {
		log.Error("<Reg> Can't Init Things . Error : ", err)
		return err
	}

	return nil

}

func (st *ThingRegistryStore) Disconnect() {
	st.db.Close()
}

func (st *ThingRegistryStore) GetThingById(Id ID) (*Thing, error) {
	var thing Thing
	err := st.db.One("ID", Id, &thing)
	return &thing, err
}

func (st *ThingRegistryStore) GetServiceById(Id ID) (*Service, error) {
	var service Service
	err := st.db.One("ID", Id, &service)
	return &service, err
}

func (st *ThingRegistryStore) GetLocationById(Id ID) (*Location, error) {
	var location Location
	err := st.db.One("ID", Id, &location)
	return &location, err
}

func (st *ThingRegistryStore) GetAllThings() ([]Thing, error) {
	var things []Thing
	err := st.db.All(&things)
	return things, err
}

func (st *ThingRegistryStore) ExtendThingsWithLocation(things []Thing) ([]ThingWithLocationView) {
	response := make([]ThingWithLocationView,len(things))
	for i := range things {
		response[i].Thing = things[i]
		loc , _ := st.GetLocationById(things[i].LocationId)
		response[i].LocationAlias = loc.Alias
	}
	return response
}

func (st *ThingRegistryStore) GetAllServices() ([]Service, error) {
	var services []Service
	err := st.db.All(&services)
	return services, err
}
// GetThingExtendedViewById return thing enhanced with linked services and location Alias
func (st *ThingRegistryStore) GetThingExtendedViewById(Id ID) (*ThingExtendedView, error) {
	var thing Thing
	var thingExView  ThingExtendedView
	err := st.db.One("ID", Id, &thing)
	thingExView.Thing = thing
	services , err := st.GetExtendedServices("",false,Id,IDnil)
	thingExView.Services = make([]ServiceExtendedView,len(services))
	for i := range services {
		thingExView.Services[i] = services[i]
	}
	location , _ := st.GetLocationById(thing.LocationId)
	if location != nil {
		thingExView.LocationAlias = location.Alias
	}
	return &thingExView, err
}

func (st *ThingRegistryStore) GetServiceByAddress(serviceName string ,serviceAddress string) (*Service, error) {
	service := Service{}
	err := st.db.Select(q.And(q.Eq("Name", serviceName), q.Eq("Address", serviceAddress))).First(&service)
	return &service, err
}


// GetExtendedServices return services enhanced with location Alias
func (st *ThingRegistryStore) GetExtendedServices(serviceNameFilter string,filterWithoutAlias bool,thingIdFilter ID,locationIdFilter ID) ([]ServiceExtendedView, error) {
	var services []Service
	var matcher []q.Matcher
	if serviceNameFilter != "" {
		match := q.Eq("Name", serviceNameFilter)
		matcher = append(matcher, match)
	}
	if locationIdFilter != IDnil {
		match := q.Eq("LocationId", locationIdFilter)
		matcher = append(matcher, match)
	}
	if filterWithoutAlias {
		match := q.Not(q.Eq("Alias", ""))
		matcher = append(matcher, match)
	}
	if thingIdFilter != IDnil {
		matcher = append(matcher, q.Eq("ParentContainerId", thingIdFilter))
		matcher = append(matcher, q.Eq("ParentContainerType", ThingContainer))
	}

	err := st.db.Select(matcher...).Find(&services)
	if err != nil {
		log.Error("<Reg> Can't fetch services . Error : ", err)
		return nil, err
	}
	var result []ServiceExtendedView
	for si := range services {
			serviceResponse := ServiceExtendedView{Service:services[si]}
			location , _ := st.GetLocationById(serviceResponse.LocationId)
			if location != nil {
						serviceResponse.LocationAlias = location.Alias
			}
			result = append(result,serviceResponse)
	}
	return result, nil
}

func (st *ThingRegistryStore) GetAllLocations() ([]Location, error) {
	var locations []Location
	err := st.db.All(&locations)
	return locations, err
}

func (st *ThingRegistryStore) GetThingByAddress(technology string, address string) (*Thing, error) {
	var thing Thing
	err := st.db.Select(q.And(q.Eq("Address", address), q.Eq("CommTechnology", technology))).First(&thing)
	return &thing,err
}

func (st *ThingRegistryStore) GetThingExtendedViewByAddress(technology string, address string) (*ThingExtendedView, error) {
	thing , err := st.GetThingByAddress(technology,address)
	if err != nil {
		return nil,err
	}
	return st.GetThingExtendedViewById(thing.ID)
}

func (st *ThingRegistryStore) GetThingsByLocationId(locationId ID) ([]Thing, error) {
	var things []Thing
	err := st.db.Select(q.Eq("LocationId", locationId)).Find(&things)
	return things, err
}

func (st *ThingRegistryStore) GetThingByIntegrationId(id string) (*Thing, error) {
	var thing Thing
	err := st.db.Select(q.Eq("IntegrationId", id)).First(&thing)
	return &thing, err
}

func (st *ThingRegistryStore) GetLocationByIntegrationId(id string) (*Location, error) {
	var location Location
	err := st.db.Select(q.Eq("IntegrationId", id)).First(&location)
	return &location, err
}

//func (st *ThingRegistryStore) GetFlatInterfaces(thingAddr string, thingTech string, serviceName string, intfMsgType string, locationId ID, thingId ID) ([]InterfaceFlatView, error) {
//	var result []InterfaceFlatView
//	//things, err  := st.GetAllThings()
//	var things []Thing
//	var matcher []q.Matcher
//	if thingAddr != "" {
//		match := q.Eq("Address", thingAddr)
//		matcher = append(matcher, match)
//	}
//	if thingTech != "" {
//		match := q.Eq("CommTechnology", thingTech)
//		matcher = append(matcher, match)
//	}
//	if thingId != 0 {
//		match := q.Eq("ID", thingId)
//		matcher = append(matcher, match)
//	}
//	err := st.db.Select(matcher...).Find(&things)
//	if err != nil {
//		return nil, err
//	}
//	for thi := range things {
//		for si := range things[thi].Services {
//			for inti := range things[thi].Services[si].Interfaces {
//				if (serviceName == "" || things[thi].Services[si].Name == serviceName) &&
//					(intfMsgType == "" || things[thi].Services[si].Interfaces[inti].MsgType == intfMsgType) &&
//					(locationId == 0 || things[thi].Services[si].LocationId == locationId) {
//
//					flatIntf := InterfaceFlatView{}
//					flatIntf.ThingId = things[thi].ID
//					flatIntf.ThingAddress = things[thi].Address
//					flatIntf.ThingTech = things[thi].CommTechnology
//					flatIntf.ThingAlias = things[thi].Alias
//					flatIntf.ServiceId = things[thi].Services[si].ID
//					flatIntf.ServiceName = things[thi].Services[si].Name
//					flatIntf.ServiceAlias = things[thi].Services[si].Alias
//					flatIntf.ServiceAddress = things[thi].Services[si].Address
//					flatIntf.InterfaceType = things[thi].Services[si].Interfaces[inti].Type
//					flatIntf.InterfaceMsgType = things[thi].Services[si].Interfaces[inti].MsgType
//					flatIntf.InterfaceValueType = things[thi].Services[si].Interfaces[inti].ValueContainer
//
//					//pt:j1/mt:evt/rt:dev/rn:zw/ad:1/sv:meter_elec/ad:21_0
//					prefix := "pt:j1/mt:evt"
//					if strings.Contains(prefix+things[thi].Services[si].Interfaces[inti].MsgType, "cmd") {
//						prefix = "pt:j1/mt:cmd"
//					}
//					flatIntf.InterfaceAddress = prefix + things[thi].Services[si].Address
//					location, _ := st.GetLocationById(things[thi].Services[si].LocationId)
//					if location != nil {
//						flatIntf.LocationId = location.ID
//						flatIntf.LocationAlias = location.Alias
//						flatIntf.LocationType = location.Type
//
//					}
//					location, _ = st.GetLocationById(things[thi].LocationId)
//					if location != nil {
//						if location.Alias != flatIntf.LocationAlias {
//							flatIntf.LocationAlias = location.Alias + " " + flatIntf.LocationAlias
//						}
//						if flatIntf.LocationType == "" {
//							flatIntf.LocationType = location.Type
//						}
//					}
//
//					result = append(result, flatIntf)
//				}
//
//			}
//		}
//	}
//	return result, nil
//}

func (st *ThingRegistryStore) UpsertThing(thing *Thing) (ID, error) {
	var err error
	if thing.ID == IDnil {
		err = st.db.Save(thing)
	} else {
		err = st.db.Update(thing)
	}

	if err != nil {
		log.Error("<Reg> Can't save thing . Error :", err)
		return 0, err
	} else {
		log.Debug("<Reg> Thing saved ")
	}

	// Updating linked services
	services ,_ := st.GetExtendedServices("",false,thing.ID,IDnil)
	var isChanged bool;
	for i := range services{
		isChanged = false
		if services[i].Alias == "" {
			services[i].Alias = thing.Alias
			isChanged = true
		}
		if services[i].LocationId == IDnil {
			services[i].LocationId = thing.LocationId
			isChanged = true
		}
		if isChanged {
			st.UpsertService(&services[i].Service)
		}
	}
	return thing.ID, nil
}

func (st *ThingRegistryStore) UpsertService(service *Service) (ID, error) {
	var err error
	// Check if service is already registered in system , if record already exits , updating the record
	if service.ID == IDnil {
		serviceCheck := Service{}
		err = st.db.Select(q.And(q.Eq("Name", service.Name), q.Eq("Address", service.Address))).First(&serviceCheck)
		if err == nil {
			service.ID = serviceCheck.ID
		}
	}
	if service.ID == IDnil {
		// Create new service
		err = st.db.Save(service)
	} else {
		err = st.db.Update(service)
	}

	if err != nil {
		log.Error("<Reg> Can't save service . Error :", err)
		return 0, err
	} else {
		log.Debug("<Reg> Service saved ")
	}

	return service.ID, nil
}


func (st *ThingRegistryStore) UpsertLocation(location *Location) (ID, error) {
	var err error
	if location.ID == 0 {
		err = st.db.Save(location)
	} else {
		err = st.db.Update(location)
	}

	if err != nil {
		log.Error("Can't save location . Error :", err)
		return 0, err
	} else {
		log.Debug("Location saved ")
	}

	return location.ID, nil
}

func (st *ThingRegistryStore) DeleteThing(id ID) error {
	thing, err := st.GetThingById(id)
	log.Debug("<Reg> Deleting thing ", thing.ID)
	if err != nil {
		return err
	}
	st.db.DeleteStruct(thing)
	// Deleting all linked services
	services ,_ := st.GetExtendedServices("",false,id,IDnil)
	var servIDs []ID
	for i := range services{
		servIDs = append(servIDs,services[i].ID)
	}
	for _,id := range servIDs{
		st.DeleteService(id)
	}

	return nil
}

func (st *ThingRegistryStore) DeleteService(id ID) error {
	service, err := st.GetServiceById(id)
	log.Debug("<Reg> Deleting service = ", service.ID)
	if err != nil {
		return err
	}
	return st.db.DeleteStruct(service)
}

func (st *ThingRegistryStore) DeleteLocation(id ID) error {
	location, err := st.GetLocationById(id)
	log.Debug("<Reg> Deleting location = ", location.ID)
	if err != nil {
		return err
	}
	return st.db.DeleteStruct(location)
}

func (st *ThingRegistryStore) ReindexAll() error {
	log.Info("Starting reindex")
	err:=st.db.ReIndex(&Thing{})
	err =st.db.ReIndex(&Location{})
	err =st.db.ReIndex(&Service{})
	log.Info("Reindex is complete")
	return err
}

func (st *ThingRegistryStore) ClearAll() error {
	thing := Thing{}
	location := Location{}
	service := Service{}
	st.db.Drop(thing)
	st.db.Drop(location)
	st.db.Drop(service)

	err := st.db.Init(&thing)
	if err != nil {
		log.Error("<Reg> Can't Init Things . Error : ", err)
		return err
	}

	err = st.db.Init(&location)
	if err != nil {
		log.Error("<Reg> Can't Init Locations . Error : ", err)
		return err
	}

	err = st.db.Init(&service)
	if err != nil {
		log.Error("<Reg> Can't Init Service . Error : ", err)
		return err
	}
	return nil
}
