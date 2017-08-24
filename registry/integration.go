package registry

import (
	"github.com/alivinco/fimpui/model"
	"github.com/alivinco/fimpgo"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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
	mg.msgTransport.Subscribe("pt:j1/mt:cmd/rt:app/rn:registry/ad:1")

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
		tech :=  addr.ResourceName
		mg.processExclusionReport(iotMsg,tech)
	}
}
func (mg *MqttIntegration) processInclusionReport(msg *fimpgo.FimpMessage) error {
	log.Info("<MqRegInt> New inclusion report")
	newThing := Thing{}
	err := msg.GetObjectValue(&newThing)
	log.Debugf("%+v\n",err)
	log.Debugf("%+v\n",newThing)
	if newThing.CommTechnology != "" && newThing.Address != "" 	{
		_ , err := mg.registry.GetThingByAddress(newThing.CommTechnology,newThing.Address)
		if err != nil {
				_,err = mg.registry.UpsertThing(&newThing)
				if err != nil {
					log.Error("<MqRegInt> Can't insert new Thing . Error: ",err)
				}
		} else {
				// updating existing node
				log.Info("<MqRegInt> Thing already in registry . Skipped.")
				//if thing.ProductHash == "" {
				//	newThing.ID = thing.ID
				//	err = mg.registry.UpsertThing(newThing)
				//}
		}

	}else {
		log.Error("<MqRegInt> Either address or commTech is empty ")
	}
	return nil

}

func (mg *MqttIntegration) processExclusionReport(msg *fimpgo.FimpMessage,technology string) error {
	valMap,err := msg.GetStrMapValue()
	if err != nil {
		return  err
	}
	address,ok := valMap["address"]
	if ok {
		thing ,err := mg.registry.GetThingByAddress(technology,address)
		if err != nil {
			return errors.New("Can't find the thing to be deleted")
		}
		mg.registry.DeleteThing(thing.ID)
		log.Infof("Thing with address = %s , tech = %s was deleted.",address,technology)
	}else {
		log.Error("<MqRegInt>Exclusion Message doesn't contain address. ")
	}
	return nil
}