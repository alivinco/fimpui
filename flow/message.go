package flow

import "github.com/alivinco/fimpgo"

type Message struct {
	AddressStr string
	Address fimpgo.Address
	Payload fimpgo.FimpMessage
	Header map[string]string
}
