package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/pandakn/cafe-beans/config"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	// middleware
	middleware := InitMiddleware(s)
	s.app.Use(middleware.Logger())
	s.app.Use(middleware.Cors())

	// Modules
	v1 := s.app.Group("/v1")
	modules := InitModule(v1, s, middleware)
	modules.MonitorModule()

	// RouterCheck
	s.app.Use(middleware.RouterCheck())

	// "Gracefully Shutdown" the server by releasing all resources
	// before terminating the application to prevent any lingering items in memory.
	// Create a channel to capture OS interrupt signals (e.g., Ctrl+C)
	c := make(chan os.Signal, 1)

	// Notify the channel when an interrupt signal is received
	signal.Notify(c, os.Interrupt) // for check interrupt
	go func() {
		// Wait for an interrupt signal
		<-c //read from channel
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen server host:port
	log.Printf("server listening on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
