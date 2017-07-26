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
	clientId := mg.config.MqttClientIdPrefix+"things_registry"
	mg.msgTransport = fimpgo.NewMqttTransport(mg.config.MqttServerURI,clientId , "", "", true, 1, 1)
	err := mg.msgTransport.Start()
	log.Info("<MqRegInt> Mqtt transport connected")
	if err != nil {
		log.Error("<MqRegInt> Error connecting to broker : ", err)
	}
	mg.msgTransport.SetMessageHandler(mg.onMqttMessage)
	mg.msgTransport.Subscribe("pt:j1/mt:evt/rt:ad/+/+")

}

func (mg *MqttIntegration) onMqttMessage(topic string, addr *fimpgo.Address, iotMsg *fimpgo.FimpMessage, rawMessage []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<MqRegInt> MqttIntegration process CRASHED with error : ",r)
			log.Errorf("<MqRegInt> Crashed while processing message from topic = %s msgType = ",r,addr.MsgType)
		}
	}()
	if iotMsg.Type == "evt.thing.inclusion_report" {
		mg.processInclusionReport(iotMsg)
	}else if iotMsg.Type == "evt.thing.exclusion_report" {
		mg.processExclusionReport(iotMsg)
	}
}
func (mg *MqttIntegration) processInclusionReport(msg *fimpgo.FimpMessage) error {
	log.Info("<MqRegInt> New inclusion report")
	newThing := Thing{}
	err := msg.GetObjectValue(&newThing)
	log.Debugf("%+v\n",err)
	log.Debugf("%+v\n",newThing)
	if newThing.CommTechnology != "" && newThing.Address != "" 	{
		thing , err := mg.registry.GetThingByAddress(newThing.CommTechnology,newThing.Address)
		if err == nil {
			if thing == nil {
				err = mg.registry.UpsertThing(newThing)
				if err != nil {
					log.Error("")
				}
			} else {
				// updating existing node
				log.Info("<MqRegInt> Thing already in registry . Skipped.")
				//if thing.ProductHash == "" {
				//	newThing.Id = thing.Id
				//	err = mg.registry.UpsertThing(newThing)
				//}
			}
		}else {
			log.Error("<MqRegInt> Can't find the thing by its address . Error:",err)
			return err
		}
	}else {
		log.Error("<MqRegInt> Either address or commTech is empty ")
	}
	return nil

}

func (mg *MqttIntegration) processExclusionReport(inclusionReport *fimpgo.FimpMessage) error {
	return nil
}