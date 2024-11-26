package database

import (
	"backend-agent-demo/adapter/repository"
	"errors"
)

var (
	errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

const (
	InstancePostgres int = iota
)

func NewDatabaseSQLFactory(instance int) (repository.SQL, error) {
	switch instance {
	case InstancePostgres:
		return NewPostgresHandler(newConfigPostgres())
	default:
		return nil, errInvalidSQLDatabaseInstance
	}
}