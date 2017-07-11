package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/model"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type Manager struct {
	flowRegistry  map[string]*Flow
	msgStreams    map[string]MsgPipeline
	msgTransport  *fimpgo.MqttTransport
	globalContext Context
	config        *model.FimpUiConfigs
}

func NewManager(config *model.FimpUiConfigs) *Manager {
	man := Manager{config:config}
	man.msgStreams = make(map[string]MsgPipeline)
	man.flowRegistry = make(map[string]*Flow)
	return &man
}

func (mg *Manager) InitMessagingTransport() {
	mg.msgTransport = fimpgo.NewMqttTransport(mg.config.MqttServerURI, "flow_manager", "", "", true, 1, 1)
	err := mg.msgTransport.Start()
	log.Info("Mqtt transport connected")
	if err != nil {
		log.Error("Error connecting to broker : ", err)
	}
	mg.msgTransport.SetMessageHandler(mg.onMqttMessage)

}

func (mg *Manager) onMqttMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	msg := Message{AddressStr: topic, Address: *addr, Payload: *iotMsg}
	// Message broadcast to all flows
	for id, stream := range mg.msgStreams {
		select {
		case stream <- msg:
			log.Debug("Message was sent to flow with id = ", id)
		default:
			log.Debug("Message is dropped (no listeners) for flow with id = ", id)
		}
	}
}

func (mg *Manager) InitNewFlow(flow *Flow) {
	if flow.Id == "" {
		flow.Id = GenerateId(10)
	}
	msgStream := make(MsgPipeline,10)
	flow.SetMessageStream(msgStream)
	mg.msgStreams[flow.Id] = msgStream
	mg.flowRegistry[flow.Id] = flow
}

func (mg *Manager) LoadFlowFromFile(fileName string) error{
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Can't open Flow file.")
		return err
	}
	flow := NewFlow("", &mg.globalContext, mg.msgTransport)
	err = json.Unmarshal(file, flow)
	if err != nil {
		fmt.Println("Can't unmarshel DB file.")
		return err
	}
	for i := range flow.Nodes {
		switch flow.Nodes[i].Type  {
		case "if":
			exp := IFExpressions{}
			err = mapstructure.Decode(flow.Nodes[i].Config,&exp)
			if err != nil{
				log.Error(err)
			}else {
				flow.Nodes[i].Config = exp
			}
		case "action":
			defValue := DefaultValue{}
			err = mapstructure.Decode(flow.Nodes[i].Config,&defValue)
			if err != nil{
				log.Error(err)
			}else {
				flow.Nodes[i].Config = defValue
			}
		case "wait":
			delay ,ok := flow.Nodes[i].Config.(float64)
			if ok {
				flow.Nodes[i].Config = int(delay)
			}else {
				log.Error("Can't cast Wait node delay value")
			}

		}
	}
	mg.InitNewFlow(flow)

	flow.Start()
	return nil

}

//
//func (mg *Manager) UpsertFlow() (string,error) {
//
//
//}
//
func (mg *Manager) DeleteFlow(id string) {
	mg.flowRegistry[id].Stop()
	delete(mg.flowRegistry, id)
	close(mg.msgStreams[id])
	delete(mg.msgStreams,id)
}
//
//func (mg *Manager) PauseFlow(flowId string) {
//
//
//}
