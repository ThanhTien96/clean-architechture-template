package routers

import (
	"backend-agent-demo/adapter/logger"
	"backend-agent-demo/adapter/validator"
	"backend-agent-demo/adapter/repository"
	"errors"
	"time"
)

type Port int64

type Server interface {
	Listen()
}

var (
	errInvalidWebServerInstance = errors.New("invalid router server instance")
)

const (
	InstanceFiber int = iota
	InstanceGin
)

func NewWebServerFactory(
	instance int,
	log logger.Logger,
	// dbSQL repository.SQL,
	dbNoSQL repository.NoSQL,
	validator validator.Validator,
	port Port,
	ctxTimeout time.Duration,
) (Server, error) {
	switch instance {
	case InstanceGin:
		return newGinServer(log, dbNoSQL, validator, port, ctxTimeout), nil
	case InstanceFiber:
		return newFiberServer(log, dbNoSQL, validator, port, ctxTimeout), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
