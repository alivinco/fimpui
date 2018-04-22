package flow

import (
	"github.com/alivinco/fimpgo/fimptype"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow/node"
	"github.com/mitchellh/mapstructure"
	"github.com/alivinco/fimpgo"
)

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
	report.Groups = []string{}

	var services []fimptype.Service

	addGroup := func(group string) {
		for i := range report.Groups {
			if report.Groups[i] == group {
				return
			}
		}
		report.Groups = append(report.Groups, group)
	}

	getService := func(name string,group string) (*fimptype.Service,bool) {
		for i := range services {
			if services[i].Name == name {
				if services[i].Groups[0] == group {
					return &services[i],false
				}

			}
		}
		service := fimptype.Service{}
		service.Name = name
		service.Groups = []string{group}
		service.Enabled = true
		service.Tags = []string{}
		service.Props = map[string]interface{}{}
		addGroup(group)
		return &service,true
	}

	for i := range flow.Nodes {
		if flow.Nodes[i].IsStartNode() {
			var config node.TriggerConfig
			err := mapstructure.Decode(flow.Nodes[i].GetMetaNode().Config,&config)
			if err==nil {
				if config.RegisterAsVirtualService{
					log.Debug("New trigger to add ")
					group := config.VirtualServiceGroup
					if group == "" {
						group = string(flow.Nodes[i].GetMetaNode().Id)
					}
					service,new := getService(flow.Nodes[i].GetMetaNode().Service,group)
					intf := fimptype.Interface{}
					intf.Type = "in"
					intf.MsgType = flow.Nodes[i].GetMetaNode().ServiceInterface
					intf.ValueType = config.InputVariableType
					intf.Version = "1"
					if new {
						log.Debug("Adding new trigger ")
						service.Alias = flow.Nodes[i].GetMetaNode().Label
						address := strings.Replace(flow.Nodes[i].GetMetaNode().Address, "pt:j1/mt:cmd", "", -1)
						address = strings.Replace(address, "pt:j1/mt:evt", "", -1)
						service.Address = address
						service.Interfaces = []fimptype.Interface{intf}
						services = append(services,*service)
					}else {
						service.Interfaces = append(service.Interfaces,intf)
					}


				}
			}else {
				log.Error("<FlMan> Fail to register trigger.Error ",err)
			}
		}
		if flow.Nodes[i].GetMetaNode().Type == "action" {
			//config,ok := flow.Nodes[i].GetMetaNode().Config.(node.ActionNodeConfig)
			config := node.ActionNodeConfig{}
			err := mapstructure.Decode(flow.Nodes[i].GetMetaNode().Config,&config)
			if err==nil {
				if config.RegisterAsVirtualService {
					group := config.VirtualServiceGroup
					if group == "" {
						group = string(flow.Nodes[i].GetMetaNode().Id)
					}
					service,new := getService(flow.Nodes[i].GetMetaNode().Service,group)

					intf := fimptype.Interface{}
					intf.Type = "out"
					intf.MsgType = flow.Nodes[i].GetMetaNode().ServiceInterface
					intf.ValueType = config.VariableType
					intf.Version = "1"

					if new {
						service.Alias = flow.Nodes[i].GetMetaNode().Label
						address := strings.Replace( flow.Nodes[i].GetMetaNode().Address,"pt:j1/mt:cmd","",-1)
						address = strings.Replace( address,"pt:j1/mt:evt","",-1)
						service.Address = address
						service.Interfaces = []fimptype.Interface{intf}
						services = append(services,*service)
					}
					service.Interfaces = append(service.Interfaces,intf)
				}

			}else {
				log.Error("<FlMan> Fail to register action .Error  ",err)
			}
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
