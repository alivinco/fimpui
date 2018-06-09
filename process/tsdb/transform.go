package tsdb

import (
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/alivinco/fimpgo"
)

// DefaultTransform - transforms IotMsg into InfluxDb datapoint
func DefaultTransform(context *MsgContext, topic string, iotMsg *fimpgo.FimpMessage, domain string) (*influx.Point, error) {
	tags := map[string]string{
		"topic":  topic,
		"domain": domain,
		"mtype":    iotMsg.Type,
		"serv": iotMsg.Service,
	}
	var fields map[string]interface{}
	var vInt int64
	var err error
	switch iotMsg.ValueType {
	case "float":
		val ,err := iotMsg.GetFloatValue()
		if err == nil {
			fields = map[string]interface{}{
				"value": val,
				"unit":  iotMsg.Properties["unit"],
			}
		}

	case "bool":
		val ,err := iotMsg.GetBoolValue()
		if err == nil {
			fields = map[string]interface{}{
				"value": val,
			}
		}

	case "int":
		vInt, err = iotMsg.GetIntValue()
		if err == nil {
			return nil, err
		}
		fields = map[string]interface{}{
			"value": vInt,
		}
	case "string":
		vStr, err := iotMsg.GetStringValue()
		if err == nil {
			return nil, err
		}
		fields = map[string]interface{}{
			"value": vStr,
		}
	case "null":
		fields = map[string]interface{}{
			"value": 0,
		}
	default:
		fields = map[string]interface{}{
			"value": iotMsg.Value,
		}

	}

	if fields != nil {
		point, err := influx.NewPoint(context.measurementName, tags, fields, context.time)
		return point, err
	}

	return nil, err

}


