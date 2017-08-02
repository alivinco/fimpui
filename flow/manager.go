package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	//"github.com/alivinco/fimpui/flow/node"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpui/flow/utils"
	fimpuimodel "github.com/alivinco/fimpui/model"
	"io/ioutil"
	"encoding/json"
	//"github.com/mitchellh/mapstructure"
	"path/filepath"
	"os"
)

type Manager struct {
	flowRegistry  map[string]*Flow
	msgStreams    map[string]model.MsgPipeline
	msgTransport  *fimpgo.MqttTransport
	globalContext model.Context
	config        *fimpuimodel.FimpUiConfigs
}

type FlowListItem struct {
	Id string
	Name string
	Description string
	State string
	TriggerCounter int64
	ErrorCounter int64
}

func NewManager(config *fimpuimodel.FimpUiConfigs) *Manager {
	man := Manager{config:config}
	man.msgStreams = make(map[string]model.MsgPipeline)
	man.flowRegistry = make(map[string]*Flow)
	man.globalContext = *model.NewContext(nil)
	return &man
}

func (mg *Manager) InitMessagingTransport() {
	clientId := mg.config.MqttClientIdPrefix+"flow_manager"
	mg.msgTransport = fimpgo.NewMqttTransport(mg.config.MqttServerURI, clientId, "", "", true, 1, 1)
	err := mg.msgTransport.Start()
	log.Info("<FlMan> Mqtt transport connected")
	if err != nil {
		log.Error("<FlMan> Error connecting to broker : ", err)
	}
	mg.msgTransport.SetMessageHandler(mg.onMqttMessage)

}

func (mg *Manager) onMqttMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	msg := model.Message{AddressStr: topic, Address: *addr, Payload: *iotMsg}
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

func (mg *Manager) GenerateNewFlow() model.FlowMeta {
	fl := model.FlowMeta{}
	fl.Nodes = []model.MetaNode{{Id:"1",Type:"trigger",Label:"no label"}}
	fl.Id = utils.GenerateId(10)
	return fl
}

func (mg *Manager) GetNewStream(Id string) model.MsgPipeline {
	msgStream := make(model.MsgPipeline,10)
	mg.msgStreams[Id] = msgStream
	return msgStream
}

func (mg *Manager) GetGlobalContext() *model.Context {
	return &mg.globalContext
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
	flowMeta := model.FlowMeta{}
	err := json.Unmarshal(flowJsonDef, &flowMeta)
	if err != nil {
		log.Error("<FlMan> Can't unmarshel DB file.")
		return err
	}

	flow := NewFlow(flowMeta, &mg.globalContext, mg.msgTransport)
	flow.SetMessageStream(mg.GetNewStream(flow.Id))
	flow.InitAllNodes()
	mg.flowRegistry[flow.Id] = flow
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
		log.Info("Adding flow with id = ",flow.Id)
		response[c] = FlowListItem{Id:flow.Id,Name:flow.Name,Description:flow.Description,TriggerCounter:flow.TriggerCounter,ErrorCounter:flow.ErrorCounter,State:flow.State}
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

