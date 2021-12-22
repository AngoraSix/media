package media

import (
	"log"

	static "angorasix.com/media/gen/static"
)

// static service example implementation.
// The example methods log the requests and return zero values.
type staticsrvc struct {
	logger *log.Logger
}

// NewStatic returns the static service implementation.
func NewStatic(logger *log.Logger) static.Service {
	return &staticsrvc{logger}
}
