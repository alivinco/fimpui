package node

import "time"
import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow/model"
)
func WaitNode(node *model.MetaNode) error {
	delayMilisec, ok := node.Config.(int)
	if ok {
		log.Info("<Node> Waiting  for = ", delayMilisec)
		time.Sleep(time.Millisecond * time.Duration(delayMilisec))
	} else {
		log.Error("<Node> Wrong time format")
	}

	return nil
}

