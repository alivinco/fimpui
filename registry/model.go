package registry

type ID int

type ThingRegistry []Thing

type Thing struct {
	Id             ID     `json:"id"`
	Address        string `json:"address"`
	ProductHash    string `json:"productHash"`
	Alias          string `json:"alias"`
	CommTechnology string `json:"commTech"`
	ProductId      string `json:"productId"`
	DeviceId       string `json:"deviceId"`
	HwVersion      string `json:"hwVersion"`
	SwVersion      string `json:"swVersion"`
	PowerSource    string  `json:"powerSource"`
	Tags           []string
	Type           string
	Location       string
	Services       []Service `json:"services"`
	Props          []string
}

type Service struct {
	Id         ID
	Name       string `json:"name"`
	Alias      string
	Address    string  `json:"address"`
	Groups     []string `json:"groups"`
	Location   string
	Props      map[string]string `json:"props"`
	Tags       []string
	Interfaces []Interface `json:"interfaces"`
}

type Interface struct {
	Type      string `json:"type"`
	MsgType   string `json:"msgType"`
	valueType string	`json:"valueType"`
	lastValue interface{}
	version   string
}
