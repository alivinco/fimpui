package model

import "github.com/alivinco/fimpgo"

type MsgPipeline chan Message

type Message struct {
	AddressStr string
	Address    fimpgo.Address
	Payload    fimpgo.FimpMessage
	Header     map[string]string
}