package model


type NodeID string

type MetaNode struct {
	Id                NodeID
	Type              string
	Label             string
	SuccessTransition NodeID
	TimeoutTransition NodeID
	ErrorTransition   NodeID
	Address           string
	Service           string
	ServiceInterface  string
	Config            interface{}
}







