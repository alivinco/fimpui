package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpgo/fimptype"
	//"github.com/alivinco/fimpui/flow/node"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpui/flow/utils"
	fimpuimodel "github.com/alivinco/fimpui/model"
	"io/ioutil"
	"encoding/json"
	//"github.com/mitchellh/mapstructure"
	"path/filepath"
	"os"
	"strings"
	"github.com/alivinco/fimpui/flow/node"
)

type Manager struct {
	flowRegistry  map[string]*Flow
	msgStreams    map[string]model.MsgPipeline
	msgTransport  *fimpgo.MqttTransport
	globalContext *model.Context
	config        *fimpuimodel.FimpUiConfigs
}

type FlowListItem struct {
	Id string
	Name string
	Description string
	State string
	TriggerCounter int64
	ErrorCounter int64
	Stats *model.FlowStatsReport
}

func NewManager(config *fimpuimodel.FimpUiConfigs) (*Manager,error) {
	var err error
	man := Manager{config:config}
	man.msgStreams = make(map[string]model.MsgPipeline)
	man.flowRegistry = make(map[string]*Flow)
	man.globalContext,err = model.NewContextDB(config.ContextStorageDir)
	man.globalContext.RegisterFlow("global")
	return &man,err
}

func (mg *Manager) InitMessagingTransport() {
	clientId := mg.config.MqttClientIdPrefix+"flow_manager"
	mg.msgTransport = fimpgo.NewMqttTransport(mg.config.MqttServerURI, clientId, mg.config.MqttUsername, mg.config.MqttPassword, true, 1, 1)
	mg.msgTransport.SetGlobalTopicPrefix(mg.config.MqttTopicGlobalPrefix)
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
		if mg.GetFlowById(id).IsFlowInterestedInMessage(topic) {
			select {
			case stream <- msg:
				log.Debug("<FlMan> Message was sent to flow with id = ", id)
			default:
				log.Debug("<FlMan> Message is dropped (no listeners) for flow with id = ", id)
			}
		}
	}
}

func (mg *Manager) GenerateNewFlow() model.FlowMeta {
	fl := model.FlowMeta{}
	fl.Nodes = []model.MetaNode{{Id:"1",Type:"trigger",Label:"no label",Config:node.TriggerConfig{Timeout:0,ValueFilter:model.Variable{},IsValueFilterEnabled:false}}}
	fl.Id = utils.GenerateId(10)
	return fl
}

func (mg *Manager) GetNewStream(Id string) model.MsgPipeline {
	msgStream := make(model.MsgPipeline,10)
	mg.msgStreams[Id] = msgStream
	return msgStream
}

func (mg *Manager) GetGlobalContext() *model.Context {
	return mg.globalContext
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
		if strings.Contains(file.Name(),".json"){
			mg.LoadFlowFromFile(filepath.Join(mg.config.FlowStorageDir,file.Name()))
		}
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

	flow := NewFlow(flowMeta, mg.globalContext, mg.msgTransport)
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
	for i := range mg.flowRegistry {
		response[c] = FlowListItem{
			Id:mg.flowRegistry[i].Id,
			Name:mg.flowRegistry[i].Name,
			Description:mg.flowRegistry[i].Description,
			TriggerCounter:mg.flowRegistry[i].TriggerCounter,
			ErrorCounter:mg.flowRegistry[i].ErrorCounter,
			State:mg.flowRegistry[i].opContext.State,
			Stats:mg.flowRegistry[i].GetFlowStats()}
		c++
	}
	return response
}

func (mg *Manager) ControlFlow(cmd string , flowId string) error {
	switch cmd {
	case "START":
		return mg.GetFlowById(flowId).Start()
	case "STOP":
		return mg.GetFlowById(flowId).Stop()

	}
	return nil
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
	flow := mg.GetFlowById(id)
	if flow == nil {
		return
	}
	flow.CleanupBeforeDelete()
	mg.UnloadFlow(id)

	os.Remove(mg.GetFlowFileNameById(id))
}

func (mg *Manager) SendInclusionReport(id string) {
	flow := mg.GetFlowById(id)
	report := fimptype.ThingInclusionReport{}
	report.Type = "flow"
	report.Address = id
	report.Alias = flow.FlowMeta.Name
	report.CommTechnology = "flow"
	report.PowerSource = "ac"
	report.ProductName = flow.FlowMeta.Name
	report.ProductHash = "flow_"+id
	report.SwVersion = "1.0"
	report.Groups = []string{}
	report.ProductId = "flow_1"
	report.ManufacturerId = "fh"
	report.Security = "tls"

   var services []fimptype.Service

	for i := range flow.Nodes {
		if flow.Nodes[i].IsStartNode() {
			service := fimptype.Service{}
			service.Name = flow.Nodes[i].GetMetaNode().Service
			service.Alias = flow.Nodes[i].GetMetaNode().Label
			service.Enabled = true
			address := strings.Replace( flow.Nodes[i].GetMetaNode().Address,"pt:j1/mt:cmd","",-1)
			address = strings.Replace( address,"pt:j1/mt:evt","",-1)
			service.Address = address
			service.Groups = []string{string(flow.Nodes[i].GetMetaNode().Id)}
			report.Groups = append(report.Groups,string(flow.Nodes[i].GetMetaNode().Id))
			intf := fimptype.Interface{}
			intf.Type = "in"
			intf.MsgType = flow.Nodes[i].GetMetaNode().ServiceInterface
			intf.ValueType = "bool"
			intf.Version = "1"
			service.Interfaces = []fimptype.Interface{intf}
			service.Props = map[string]interface{}{}
			service.Tags = []string{}
			services = append(services,service)
		}

	}
	report.Services = services
	msg := fimpgo.NewMessage("evt.thing.inclusion_report", "flow","object", report, nil,nil,nil)
	addrString := "pt:j1/mt:evt/rt:ad/rn:flow/ad:1"
	addr, _ := fimpgo.NewAddressFromString(addrString)
	mg.msgTransport.Publish(addr,msg)
}

func (mg *Manager) SendExclusionReport(id string) {
	report := fimptype.ThingExclusionReport{Address:id}
	msg := fimpgo.NewMessage("evt.thing.exclusion_report", "flow","object", report, nil,nil,nil)
	addrString := "pt:j1/mt:evt/rt:ad/rn:flow/ad:1"
	addr, _ := fimpgo.NewAddressFromString(addrString)
	mg.msgTransport.Publish(addr,msg)
}


