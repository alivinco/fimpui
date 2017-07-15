package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/model"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"path/filepath"
	"os"
)

type Manager struct {
	flowRegistry  map[string]*Flow
	msgStreams    map[string]MsgPipeline
	msgTransport  *fimpgo.MqttTransport
	globalContext Context
	config        *model.FimpUiConfigs
}

type FlowListItem struct {
	Id string
	Name string
	Description string
	TriggerCounter int64
	ErrorCounter int64
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
	log.Info("<FlMan> Mqtt transport connected")
	if err != nil {
		log.Error("<FlMan> Error connecting to broker : ", err)
	}
	mg.msgTransport.SetMessageHandler(mg.onMqttMessage)

}

func (mg *Manager) onMqttMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	msg := Message{AddressStr: topic, Address: *addr, Payload: *iotMsg}
	// Message broadcast to all flows
	for id, stream := range mg.msgStreams {
		select {
		case stream <- msg:
			log.Debug("<FlMan> Message was sent to flow with id = ", id)
		default:
			log.Debug("<FlMan> Message is dropped (no listeners) for flow with id = ", id)
		}
	}
}

func (mg *Manager) GenerateNewFlow() Flow {
	fl := Flow {}
	fl.AddNode(MetaNode{Id:"1",Type:"trigger",Label:"no label"})
	fl.Id = GenerateId(10)

	return fl
}

func (mg *Manager) InitNewFlow(flow *Flow) string {
	msgStream := make(MsgPipeline,10)
	flow.SetMessageStream(msgStream)
	mg.msgStreams[flow.Id] = msgStream
	mg.flowRegistry[flow.Id] = flow
	return flow.Id
}

func (mg *Manager) GetFlowFileNameById(id string ) string {
	return filepath.Join(mg.config.FlowStorageDir,id+".json")
}

func (mg *Manager) LoadAllFlowsFromStorage () error {
	files, err := ioutil.ReadDir(mg.config.FlowStorageDir)
	if err != nil {
		log.Error(err)
		return err
	}

	for _, file := range files {
		mg.LoadFlowFromFile(filepath.Join(mg.config.FlowStorageDir,file.Name()))
	}
	return nil
}

func (mg *Manager) LoadFlowFromFile(fileName string) error{
	log.Info("<FlMan> Loading flow from file : ",fileName)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Error("<FlMan> Can't open Flow file.")
		return err
	}
	mg.LoadFlowFromJson(file)
	return nil
}

func (mg *Manager) LoadFlowFromJson(flowJsonDef []byte) error{
	flow := NewFlow("", &mg.globalContext, mg.msgTransport)
	err := json.Unmarshal(flowJsonDef, flow)
	if err != nil {
		fmt.Println("<FlMan> Can't unmarshel DB file.")
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
				log.Error("<FlMan> Can't cast Wait node delay value")
			}

		}
	}
	mg.InitNewFlow(flow)

	flow.Start()
	return nil
}
func (mg *Manager) UpdateFlowFromJson(id string, flowJsonDef []byte) error {
	mg.UnloadFlow(id)
	err := mg.LoadFlowFromJson(flowJsonDef)
	return err
}

func (mg *Manager) ReloadFlowFromStorage(id string ) error {
	mg.UnloadFlow(id)
	return mg.LoadFlowFromFile(mg.GetFlowFileNameById(id))
}

//
func (mg *Manager) UpdateFlowFromJsonAndSaveToStorage(id string, flowJsonDef []byte) (error) {
	fileName := mg.GetFlowFileNameById(id)
	log.Debugf("<FlMan> Saving flow to file %s , data size %d :",fileName,len(flowJsonDef))

	err := ioutil.WriteFile(fileName, flowJsonDef, 0644)
	if err != nil {
		log.Error("Can't save flow to file . Error : ",err)
		return err
	}
	err = mg.UpdateFlowFromJson(id,flowJsonDef)
	return err

}

func (mg *Manager) GetFlowById(id string) *Flow{
	for i := range mg.flowRegistry {
		if mg.flowRegistry[i].Id == id {
			return mg.flowRegistry[i]
		}
	}
	return nil
}

func (mg *Manager) GetFlowList() []FlowListItem{
	response := make ([]FlowListItem,len(mg.flowRegistry))
	var c int
	for _,flow := range mg.flowRegistry {
		response[c] = FlowListItem{Id:flow.Id,Name:flow.Name,Description:flow.Description,TriggerCounter:flow.TriggerCounter,ErrorCounter:flow.ErrorCounter}
		c++
	}
	return response
}

func (mg *Manager) UnloadFlow(id string) {
	if mg.GetFlowById(id) == nil {
		return
	}
	mg.flowRegistry[id].Stop()
	delete(mg.flowRegistry, id)
	close(mg.msgStreams[id])
	delete(mg.msgStreams,id)
	log.Infof("Flow with Id = %s is unloaded",id)
}

func (mg *Manager) DeleteFlow(id string) {
	if mg.GetFlowById(id) == nil {
		return
	}
	mg.UnloadFlow(id)
	os.Remove(mg.GetFlowFileNameById(id))
}

