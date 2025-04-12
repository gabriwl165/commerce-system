package services

import "log"

func ProcessConsumption(consumptionChan <-chan map[string]interface{}) {
	for consumption := range consumptionChan {
		log.Print(consumption)
	}
}
