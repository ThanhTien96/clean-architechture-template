package main

import (
	"backend-agent-demo/infrastructure"
	database "backend-agent-demo/infrastructure/databases"
	log "backend-agent-demo/infrastructure/logs"
	"backend-agent-demo/infrastructure/routers"
	"backend-agent-demo/infrastructure/validation"
	"os"
	"time"
)

func main() {
	
	var app = infrastructure.NewConfig().
	Name(os.Getenv("APP_NAME")).
	ContextWithTimeout(10 * time.Second).
	Logger(log.InstanceLogrusLogger).
	Validator(validation.InstanceGoPlayground).
	DBSQL(database.InstancePostgres)
	// .DBNoSQL(database.InstanceMongoDB)


	app.HttpServerPorts(os.Getenv("APP_PORT")).
	HttpServer(routers.InstanceFiber).
	Start()
}
