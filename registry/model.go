package registry

type ID int

type ThingRegistry []Thing

type Thing struct {
	Id             ID     `json:"id"`
	Address        string `json:"address"`
	ProductHash    string `json:"productHash"`
	Alias          string `json:"alias"`
	CommTechnology string `json:"commTech"`
	Tags           []string
	Type           string
	Location       string
}
