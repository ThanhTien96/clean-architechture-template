package infrastructure

import (
	"backend-agent-demo/adapter/logger"
	"backend-agent-demo/adapter/repository"
	"backend-agent-demo/adapter/validator"
	database "backend-agent-demo/infrastructure/databases"
	log "backend-agent-demo/infrastructure/logs"
	"backend-agent-demo/infrastructure/routers"
	"backend-agent-demo/infrastructure/validation"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	appName        string
	logger         logger.Logger
	validator      validator.Validator
	dbNoSQL        repository.NoSQL
	dbSQL          repository.SQL
	httpServerPort routers.Port
	httpServer     routers.Server
	ctxTimeout     time.Duration
}

func NewConfig() *config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &config{}
}

func (cfg *config) ContextWithTimeout(t time.Duration) *config {
	cfg.ctxTimeout = t
	return cfg
}

func (cfg *config) Name(name string) *config {
	cfg.appName = name
	return cfg
}

func (cfg *config) Logger(instance int) *config {
	log, err := log.NewLoggerFactory(instance)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.logger = log
	cfg.logger.Infof("Successfully configured log")
	return cfg
}

// db SQL
func (cfg *config) DBSQL(instance int) *config {
	db, err := database.NewDatabaseSQLFactory(instance)
	if err != nil {
		cfg.logger.Fatalln(err, "Could not make a connection to the SQL database")
	}

	cfg.logger.Infof("Successfully connected to the SQL database")

	cfg.dbSQL = db
	return cfg
}


// db NoSQL
func (cfg *config) DBNoSQL(instance int) *config {
	db, err := database.NewDatabaseNoSQLFactory(instance)
	if err != nil {
		cfg.logger.Fatalln(err, "Could not make a connection to the database")
	}

	cfg.logger.Infof("Successfully connected to the NoSQL database")

	cfg.dbNoSQL = db

	return cfg
}

func (cfg *config) Validator(instance int) *config {
	v, err := validation.NewValidatorFactory(instance)
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	cfg.logger.Infof("Successfully configured validator")

	cfg.validator = v
	return cfg
}

func (cfg *config) HttpServer(instance int) *config {
	s, err := routers.NewWebServerFactory(
		instance,
		cfg.logger,
		cfg.dbNoSQL,
		cfg.validator,
		cfg.httpServerPort,
		cfg.ctxTimeout,
	)

	if err != nil {
		cfg.logger.Fatalln(err)
	}

	cfg.logger.Infof("Successfully configured router server")

	cfg.httpServer = s
	return cfg
}

func (cfg *config) HttpServerPorts(port string) *config {
	p, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		cfg.logger.Fatalln(err)
	}

	cfg.httpServerPort = routers.Port(p)

	return cfg
}

func (cfg *config) Start() {
	cfg.httpServer.Listen()
}
