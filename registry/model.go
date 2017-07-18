package registry

type ID int

type ThingRegistry struct {
	Things []Thing
	Locations []Location
}

type Thing struct {
	Id             ID        `json:"id"`
	IntegrationId  string    `json:"integr_id"`
	Address        string    `json:"address"`
	Type           string    `json:"type"`
	ProductHash    string    `json:"product_hash"`
	Alias          string    `json:"alias"`
	CommTechnology string    `json:"comm_tech"`
	ProductId      string    `json:"product_id"`
	DeviceId       string    `json:"device_id"`
	HwVersion      string    `json:"hw_ver"`
	SwVersion      string    `json:"sw_ver"`
	PowerSource    string    `json:"power_source"`
	Tags           []string  `json:"tags"`
	LocationId     string    `json:"locationId"`
	Services       []Service `json:"services"`
	Props          []string  `json:"props"`
}

type Service struct {
	Id            ID                     `json:"id"`
	IntegrationId string                 `json:"integr_id"`
	Name          string                 `json:"name"`
	Alias         string                 `json:"alias"`
	Address       string                 `json:"address"`
	Groups        []string               `json:"groups"`
	LocationId    string                 `json:"location_id"`
	Props         map[string]interface{} `json:"props"`
	Tags          []string               `json:"tags"`
	Interfaces    []Interface            `json:"interfaces"`
}

type Interface struct {
	Type      string      `json:"intf_t"`
	MsgType   string      `json:"msg_t"`
	valueType string      `json:"val_t"`
	lastValue interface{} `json:"last_val"`
	version   string      `json:"ver"`
}

type Location struct {
	Id             ID         `json:"id"`
	IntegrationId  string     `json:"integr_id"`
	Type           string     `json:"type"`
	Alias          string     `json:"alias"`
	Address        string     `json:"address"`
	Image          string     `json:"image"`
	ChildLocations []Location `json:"child_locations"`
	State          string     `json:"state"`
}
