package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"testing"
	"time"
)

var msgChan = make(MsgPipeline)

func onMsg(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	log.Info("New message from topic = ",topic)

	fMsg := Message{AddressStr: topic, Address: *addr, Payload: *iotMsg}
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
	mqtt.Publish(&adr,msg)
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

	ctx := Context{}
	flow := NewFlow(&ctx,mqtt);
	flow.SetMessageStream(msgChan);

	node := Node{Id:"1",Label:"Button trigger 1",Type:"trigger",Address:"pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:199_0",Service:"out_bin_switch",ServiceInterface:"evt.binary.report",SuccessTransition:"2"}
	flow.AddNode(node)
	node = Node{Id:"1.1",Label:"Button trigger 2",Type:"trigger",Address:"pt:j1/mt:evt/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:299_0",Service:"out_bin_switch",ServiceInterface:"evt.binary.report",SuccessTransition:"2"}
	flow.AddNode(node)
	node = Node{Id:"2",Label:"Bulb 1",Type:"action",Address:"pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0",Service:"out_bin_switch",ServiceInterface:"cmd.binary.set",SuccessTransition:"2.1"}
	flow.AddNode(node)

	node = Node{Id:"2.1",Label:"Waiting for 500mil",Type:"wait",SuccessTransition:"3",Config:2000}
	flow.AddNode(node)
	node = Node{Id:"3",Label:"Bulb 2",Type:"action",Address:"pt:j1/mt:cmd/rt:dev/rn:test/ad:1/sv:out_bin_switch/ad:200_0",Service:"out_bin_switch",ServiceInterface:"cmd.binary.set",SuccessTransition:""}
	flow.AddNode(node)
	flow.Start()
	time.Sleep(time.Second*1)
	sendMsg(mqtt)
	time.Sleep(time.Second*5)

}
