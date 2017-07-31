package flow

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"testing"
	"time"
	"io/ioutil"
	//"github.com/alivinco/fimpui/flow/node"
	"github.com/alivinco/fimpui/flow/model"
	flownode "github.com/alivinco/fimpui/flow/node"
)

var msgChan = make(model.MsgPipeline)

func onMsg(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	log.Info("New message from topic = ", topic)

	fMsg := model.Message{AddressStr: topic, Address: *addr, Payload: *iotMsg}
	select {
	case msgChan <- fMsg:
		log.Info("Message was sent")
	default:
		log.Info("Message dropped , no receiver ")
	}
}

func sendMsg(mqtt *fimpgo.MqttTransport) {
	msg := fimpgo.NewBoolMessage("evt.binary.report", "out_bin_switch", true, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "out_bin_switch", ServiceAddress: "199_0"}
	mqtt.Publish(&adr, msg)
}

func TestNewFlow(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://localhost:1883", "flow_test", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(time.Second * 1)

	ctx := model.Context{}
	flow := NewFlow("1", &ctx, mqtt)
	flow.SetMessageStream(msgChan)

	node := model.MetaNode{Id: "1", Label: "Lux trigger", Type: "trigger", Address: "pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:199_0", Service: "out_bin_switch", ServiceInterface: "evt.binary.report", SuccessTransition: "2"}
	flow.AddNode(node)
	node = model.MetaNode{Id: "1.1", Label: "Button trigger 2", Type: "trigger", Address: "pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:299_0", Service: "out_bin_switch", ServiceInterface: "evt.binary.report", SuccessTransition: "2"}
	flow.AddNode(node)
	node = model.MetaNode{Id: "2", Label: "Bulb 1", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: "2.1"}
	flow.AddNode(node)
	node = model.MetaNode{Id: "2.1", Label: "Waiting for 500mil", Type: "wait", SuccessTransition: "3", Config: 200}
	flow.AddNode(node)
	node = model.MetaNode{Id: "3", Label: "Bulb 2", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: ""}
	flow.AddNode(node)
	flow.Start()
	time.Sleep(time.Second * 1)
	sendMsg(mqtt)
	time.Sleep(time.Second * 5)

}

func TestNewFlow2(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://localhost:1883", "flow_test", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(time.Second * 1)

	ctx := model.Context{IsFlowRunning:true}
	flow := NewFlow("1", &ctx, mqtt)
	flow.SetMessageStream(msgChan)

	node := model.MetaNode{Id: "1", Label: "Button trigger 1", Type: "trigger", Address: "pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:sensor_lumin/ad:199_0", Service: "sensor_lumin", ServiceInterface: "evt.sensor.report", SuccessTransition: "1.1"}
	flow.AddNode(node)
	node = model.MetaNode{Id: "1.1", Label: "IF node", Type: "if", Config: flownode.IFExpressions{TrueTransition: "2", FalseTransition: "3", Expression: []flownode.IFExpression{
		{Value: int64(100), ValueType: "int", Operand: "gt" ,BooleanOperator:"and"},
		{Value: int64(200), ValueType: "int", Operand: "lt" }}}}
	flow.AddNode(node)
	b, err := json.Marshal(node)
	log.Info(string(b))
	node = model.MetaNode{Id: "2", Label: "Bulb 1.Room light intensity is > 100 lux", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: "",
		Config: flownode.DefaultValue{ValueType: "bool", Value: true}}
	flow.AddNode(node)
	node = model.MetaNode{Id: "3", Label: "Bulb 2.Room light intensity is < 100 lux", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: "",
		Config: flownode.DefaultValue{ValueType: "bool", Value: true}}
	flow.AddNode(node)

	data, err := json.Marshal(flow)
	if err == nil {
		ioutil.WriteFile("testflow.json", data, 0644)
	}


	flow.Start()
	time.Sleep(time.Second * 1)
	// send msg

	msg := fimpgo.NewIntMessage("evt.sensor.report", "sensor_lumin", 150, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "sensor_lumin", ServiceAddress: "199_0"}
	mqtt.Publish(&adr, msg)

	// end
	time.Sleep(time.Second * 5)

}

func TestNewFlow3(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	mqtt := fimpgo.NewMqttTransport("tcp://localhost:1883", "flow_test", "", "", true, 1, 1)
	err := mqtt.Start()
	t.Log("Connected")
	if err != nil {
		t.Error("Error connecting to broker ", err)
	}

	mqtt.SetMessageHandler(onMsg)
	time.Sleep(time.Second * 1)

	ctx := model.Context{}
	flow := NewFlow("2", &ctx, mqtt)
	flow.SetMessageStream(msgChan)

	node := model.MetaNode{Id: "1", Label: "Button trigger", Type: "trigger", Address: "pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:199_0", Service: "out_bin_switch", ServiceInterface: "evt.binary.report", SuccessTransition: "1.1"}
	flow.AddNode(node)
	node = model.MetaNode{Id: "1.1", Label: "IF node", Type: "if", Config: flownode.IFExpressions{TrueTransition: "2", FalseTransition: "3", Expression: []flownode.IFExpression{{Value: false, ValueType: "bool", Operand: "eq"}}}}
	flow.AddNode(node)
	node = model.MetaNode{Id: "2", Label: "Lights ON", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: "",
		Config: flownode.DefaultValue{ValueType: "bool", Value: true}}
	flow.AddNode(node)
	node = model.MetaNode{Id: "3", Label: "Lights OFF", Type: "action", Address: "pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0", Service: "out_bin_switch", ServiceInterface: "cmd.binary.set", SuccessTransition: "",
		Config: flownode.DefaultValue{ValueType: "bool", Value: false}}
	flow.AddNode(node)
	flow.Start()
	time.Sleep(time.Second * 1)
	// send msg

	msg := fimpgo.NewBoolMessage("evt.binary.report", "out_bin_switch", true, nil, nil, nil)
	adr := fimpgo.Address{MsgType: fimpgo.MsgTypeEvt, ResourceType: fimpgo.ResourceTypeDevice, ResourceName: "test", ResourceAddress: "1", ServiceName: "out_bin_switch", ServiceAddress: "199_0"}
	mqtt.Publish(&adr, msg)
	time.Sleep(time.Second * 1)
	flow.Stop()
	// end
	time.Sleep(time.Second * 2)

}
