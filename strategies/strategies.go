package strategies

import (
	"fmt"

	"angorasix.com/media/config"
)

// StrategiesMap ...
var StrategiesMap = map[string]func(config *config.ServiceConfig) (StorageStrategy, error){
	"local":  CreateLocalStrategy,
	"google": CreateGoogleCloudStrategy,
}

// CreateStrategyFromConfig ...
func CreateStrategyFromConfig(config *config.ServiceConfig) (StorageStrategy, error) {
	creationStrategyFunction, ok := StrategiesMap[config.Strategy]
	if !ok {
		panic(fmt.Sprintf("%s is not a valid strategy.", config.Strategy))
	}
	return creationStrategyFunction(config)
}
