package database

import (
	"backend-agent-demo/adapter/repository"
	"errors"
)

var (
	errInvalidNoSQLDatabaseInstance = errors.New("invalid nosql database instace")
)

const (
	InstanceMongoDB int = iota
)

func NewDatabaseNoSQLFactory(instance int) (repository.NoSQL, error) {
	switch instance {
	case InstanceMongoDB:
		return NewMongoHandler(newConfigMongoDB())
	default:
		return nil, errInvalidNoSQLDatabaseInstance
	}
}