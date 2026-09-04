package gbclog

import (
	internalproperties "goark.dev/gbc-log/internal/properties"
	coreenv "goark.dev/goark/core/env"
)

type loggingProperties = internalproperties.Properties

func readLoggingProperties(environment coreenv.Environment) (loggingProperties, error) {
	return internalproperties.Read(environment)
}
