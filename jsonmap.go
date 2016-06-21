package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

//MetricMap maps name to float64 metric
type MetricMap map[string]float64

//MakeQueueMap creates a map of queues from json input. Only keys with float values are mapped.
func MakeQueueMap(d *json.Decoder) map[string]MetricMap {
	queueMap := make(map[string]MetricMap)
	var jsonArr []map[string]interface{}

	if d == nil {
		log.Error("JSON decoder not iniatilized")
		return queueMap
	}

	if err := d.Decode(&jsonArr); err != nil {
		log.WithField("error", err).Error("Error while decoding json")
		return queueMap
	}
	for _, el := range jsonArr {
		log.WithField("element", el).WithField("name", el["name"]).Debug("Iterate over array")
		if name, ok := el["name"]; ok {
			flMap := make(MetricMap)
			addFields(&flMap, "", el)
			queueMap[name.(string)] = flMap
		}
	}

	return queueMap
}

//MakeMap creates a map from json input. Only keys with float values are mapped.
func MakeMap(d *json.Decoder) MetricMap {
	flMap := make(MetricMap)
	var output map[string]interface{}

	if d == nil {
		log.Error("JSON decoder not iniatilized")
		return flMap
	}

	if err := d.Decode(&output); err != nil {
		log.WithField("error", err).Error("Error while decoding json")
		return flMap
	}

	addFields(&flMap, "", output)

	return flMap
}

func addFields(toMap *MetricMap, basename string, source map[string]interface{}) {
	prefix := ""
	if basename != "" {
		prefix = basename + "."
	}
	for k, v := range source {
		switch value := v.(type) {
		case float64:
			(*toMap)[prefix+k] = value
		case map[string]interface{}:
			addFields(toMap, prefix+k, value)
		case []interface{}:
			(*toMap)[prefix+k] = float64(len(value))
		}
	}
}
