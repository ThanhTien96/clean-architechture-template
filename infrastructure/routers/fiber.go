package routers

import (
	"backend-agent-demo/adapter/logger"
	"backend-agent-demo/adapter/repository"
	"backend-agent-demo/adapter/validator"
	"fmt"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_logger "github.com/gofiber/fiber/v2/middleware/logger"
)

type fiberEngine struct {
	app        *fiber.App
	log        logger.Logger
	db         repository.NoSQL
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newFiberServer(
	log logger.Logger,
	db repository.NoSQL,
	validator validator.Validator,
	port Port,
	ctxTimeout time.Duration,
) *fiberEngine {
	return &fiberEngine{
		app:        fiber.New(),
		log:        log,
		db:         db,
		validator:  validator,
		port:       port,
		ctxTimeout: ctxTimeout,
	}
}

func (f *fiberEngine) Listen() {
	app := f.app

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Authorization, Content-Length",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH",
		AllowCredentials: false,
	}))

	app.Use(fiber_logger.New(fiber_logger.ConfigDefault))

	f.setAppHandlers(app)

	// routes
	f.GracefulShutdown(fmt.Sprintf("%v", f.port))
}

func (f *fiberEngine) GracefulShutdown(port string) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := f.app.Listen(":" + port); err != nil {
			f.log.Fatalln("error when listening to :%s, %s", port, err)
		}
	}()

	f.log.Infof("server is running on :%s", port)

	<-stop

	f.log.Infof("server gracefully shutdown")

	if err := f.app.Shutdown(); err != nil {
		f.log.Fatalln("error when shutting down the server, %s", err)
	}

	f.log.Infof("process clean up...")
}

/* TODO ADD MIDDLEWARE */
func (f *fiberEngine) setAppHandlers(app *fiber.App) {

	v1 := app.Group("/v1")
	{
		v1.Get("/health", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})
	}

}
