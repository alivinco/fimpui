package registry

type ID int

type ThingRegistry struct {
	Things    []Thing
	Locations []Location
}

type ValueType struct {
	Value     interface{}
	ValueType string
}

type Thing struct {
	ID             ID        `json:"id" storm:"id,increment"`
	IntegrationId  string    `json:"integr_id" storm:"index"`
	Address        string    `json:"address" storm:"index"`
	Type           string    `json:"type"`
	ProductHash    string    `json:"product_hash"`
	Alias          string    `json:"alias"`
	CommTechnology string    `json:"comm_tech" storm:"index"`
	ProductId      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	ManufacturerId string    `json:"manufacturer_id"`
	DeviceId       string    `json:"device_id"`
	HwVersion      string    `json:"hw_ver"`
	SwVersion      string    `json:"sw_ver"`
	PowerSource    string    `json:"power_source"`
	Tags           []string  `json:"tags"`
	LocationId     ID        `json:"location_id"`
	Services       []Service `json:"services"`
	Props          []string  `json:"props"`
}

type Service struct {
	ID            ID                     `json:"id" , storm:"id,increment"`
	IntegrationId string                 `json:"integr_id" storm:"index"`
	Name          string                 `json:"name"`
	Alias         string                 `json:"alias"`
	Address       string                 `json:"address"`
	Groups        []string               `json:"groups"`
	LocationId    ID                     `json:"location_id"`
	Props         map[string]interface{} `json:"props"`
	Tags          []string               `json:"tags"`
	Interfaces    []Interface            `json:"interfaces"`
	Attributes    map[string]ValueType   `json:"attributes"`
}

type Interface struct {
	Type      string `json:"intf_t"`
	MsgType   string `json:"msg_t"`
	ValueType string `json:"val_t"`
	Version   string `json:"ver"`
}

type Location struct {
	ID             ID         `json:"id" storm:"id,increment"`
	IntegrationId  string     `json:"integr_id"`
	Type           string     `json:"type"`
	Alias          string     `json:"alias"`
	Address        string     `json:"address"`
	Longitude      float64    `json:"long"`
	Latitude       float64    `json:"lat"`
	Image          string     `json:"image"`
	ChildLocations []Location `json:"child_locations"`
	State          string     `json:"state"`
}

type ServiceResponse struct {
	Service
	LocationAlias string                 `json:"location_alias"`
}

type InterfaceFlatView struct {
	ThingId            ID       `json:"thing_id"`
	ThingAddress       string   `json:"thing_address"`
	ThingTech          string   `json:"thing_tech"`
	ThingAlias         string   `json:"thing_alias"`
	ServiceId          ID       `json:"service_id"`
	ServiceName        string   `json:"service_name"`
	ServiceAlias       string   `json:"service_alias"`
	ServiceAddress     string   `json:"service_address"`
	InterfaceType      string   `json:"intf_type"`
	InterfaceMsgType   string   `json:"intf_msg_type"`
	InterfaceAddress   string   `json:"intf_address"`
	InterfaceValueType string   `json:"intf_val_type"`
	LocationId         ID       `json:"location_id"`
	LocationAlias      string   `json:"location_alias"`
	LocationType       string   `json:"location_type"`
	Groups             []string `json:"groups"`
}
