package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wayla99/go_clean/src/interface/fiber_server"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wayla99/go_clean/src/repository/staff_repository"

	"github.com/go-pg/pg/v10"

	"github.com/wayla99/go_clean/src/use_case"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.uber.org/zap/zapcore"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logger *zap.Logger

type config struct {
	AppName                   string `env:"APP_NAME" envDefault:"staff-test"`
	AppVersion                string `env:"APP_VERSION" envDefault:"v0.0.0"`
	Environment               string `env:"ENVIRONMENT" envDefault:"development"`
	Port                      uint   `env:"PORT" envDefault:"8080"`
	Debuglog                  bool   `env:"DEBUG_LOG" envDefault:"false"`
	JaegerEndpoint            string `env:"JAEGER_ENDPOINT" envDefault:"http://localhost:14268/api/traces"`
	PostgresDbStaffServiceUri string `env:"POSTGRES_DB_STAFF_SERVICE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/staff?sslmode=disable"`
}

func main() {
	cfg := initEnvironment()

	initLogger(cfg)
	//initTracer(cfg)

	staffRepo := initRepositories(cfg)
	useCase := use_case.New(staffRepo)
	initInterfaces(cfg, useCase)

}

func initEnvironment() config {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}

	var cfg config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Error parse env: %s\n", err)
	}

	return cfg
}

func initLogger(cfg config) {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	LogLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	if cfg.Debuglog {
		LogLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	conf.Level = LogLevel

	lg, err := conf.Build()
	if err != nil {
		log.Fatalf("Error build logger: %s\n", err)
	}
	defer lg.Sync()

	zap.ReplaceGlobals(lg)
	logger = zap.L().Named("bootstrap")
	logger.Info("Logger initialized")
}

func initTracer(cfg config) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		logger.Fatal("Error init Jeager exporter", zap.Error(err))
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.AppName),
			semconv.ServiceVersionKey.String(cfg.AppVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		logger.Fatal("Error init Jaeger resource", zap.Error(err))
	}

	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(r),
	)

	otel.SetTracerProvider(tp)
	logger.Info("Tracer initialized")
}

func initRepositories(cfg config) use_case.StaffRepository {
	opt, err := pg.ParseURL(cfg.PostgresDbStaffServiceUri)
	if err != nil {
		logger.Fatal("Error parse postgres uri", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := pg.Connect(opt)
	if err := db.Ping(ctx); err != nil {
		logger.Fatal("Error connect to postgres", zap.Error(err))
	}

	staffRepo, err := staff_repository.NewGoPg(ctx, db)
	if err != nil {
		logger.Fatal("Error init staff repository", zap.Error(err))
	}
	logger.Info("Staff repository initialized")
	logger.Info("Repositories initialized")

	return staffRepo
}

func initInterfaces(cfg config, useCase *use_case.UseCase) {
	wg := new(sync.WaitGroup)
	prom := prometheus.NewRegistry()

	serv := fiber_server.New(useCase, prom, &fiber_server.ServerConfig{
		AppVersion:    cfg.AppVersion,
		ListenAddress: fmt.Sprintf(":%d", cfg.Port),
		RequestLog:    true,
	})
	logger.Info("Fiber server initialized")

	serv.Start(wg)
	logger.Info("Fiber server started")

	wg.Wait()
	logger.Info("Application stopped")
}
