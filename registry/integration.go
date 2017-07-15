package registry

import (
	"github.com/alivinco/fimpui/model"
	"github.com/alivinco/fimpgo"
	log "github.com/Sirupsen/logrus"
)



type MqttIntegration struct  {
	msgTransport  *fimpgo.MqttTransport
	config        *model.FimpUiConfigs
	registry *ThingRegistryStore
}

func NewMqttIntegration(config *model.FimpUiConfigs,registry* ThingRegistryStore) *MqttIntegration {
	int := MqttIntegration{config:config,registry:registry}
	return &int
}
func (mg *MqttIntegration) InitMessagingTransport() {
	mg.msgTransport = fimpgo.NewMqttTransport(mg.config.MqttServerURI, "flow_manager", "", "", true, 1, 1)
	err := mg.msgTransport.Start()
	log.Info("<FlMan> Mqtt transport connected")
	if err != nil {
		log.Error("<FlMan> Error connecting to broker : ", err)
	}
	mg.msgTransport.SetMessageHandler(mg.onMqttMessage)
	mg.msgTransport.Subscribe("pt:j1/mt:evt/rt:ad/+/+")

}

func (mg *MqttIntegration) onMqttMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {

}
func processInclusionReport(inclusionReport *fimpgo.FimpMessage) error {

}