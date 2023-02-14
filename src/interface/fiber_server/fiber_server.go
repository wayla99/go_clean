package fiber_server

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wayla99/go_clean/src/use_case"
)

type FiberServer struct {
	useCase *use_case.UseCase
	server  *fiber.App
	prom    *prometheus.Registry
	config  *ServerConfig
}

type ServerConfig struct {
	AppVersion    string
	RequestLog    bool
	ListenAddress string
}

var (
	ErrInvalidPayload   = errors.New("invalid payload")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrUnauthenticated  = errors.New("unauthenticated")
)

func New(uc *use_case.UseCase, prom *prometheus.Registry, sc *ServerConfig) *FiberServer {
	server := fiber.New(fiber.Config{
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
	})

	f := &FiberServer{
		useCase: uc,
		server:  server,
		prom:    prom,
		config:  sc,
	}

	server.Use(f.recover)
	server.Use(cors.New())
	f.addMetrics(sc.AppVersion)
	f.addRouteSystem(server)

	if sc.RequestLog {
		server.Use(loggerMiddleware)
	}

	server.Use(tracerMiddleware)

	f.addRouteStaff(server)
	f.addRouteSwagger(server)

	return f
}

func (f *FiberServer) Start(wg *sync.WaitGroup) {
	wg.Add(2)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		defer wg.Done()
		<-exit
		zap.L().Info("Shutting down server...")

		err := f.server.Shutdown()
		if err != nil {
			zap.L().Info("Server shutdown with error", zap.Error(err))
		} else {
			zap.L().Info("Server gracefully shutdown")
		}
	}()

	go func() {
		defer wg.Done()
		zap.L().Info("Server is starting...")
		err := f.server.Listen(f.config.ListenAddress)
		if err != nil {
			zap.L().Info("Server error", zap.Error(err))
		}
		zap.L().Info("Server has been shutdown")
	}()
}
