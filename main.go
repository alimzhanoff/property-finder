package main

import (
	"context"
	"fmt"
	"github.com/alimzhanoff/property-finder/internal/api/http/router"
	"github.com/alimzhanoff/property-finder/internal/config"
	"github.com/alimzhanoff/property-finder/internal/storage/property/postgres"
	"github.com/alimzhanoff/property-finder/pkg/logging"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/alimzhanoff/property-finder/docs"
)

func main() {
	// @title           Swagger Example API
	// @version         1.0
	// @description     This is a sample server celler server.
	// @termsOfService  http://swagger.io/terms/

	// @contact.name   API Support
	// @contact.url    http://www.swagger.io/support
	// @contact.email  support@swagg.io

	// @license.name  Apache 2.0
	// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

	// @host      localhost:3000
	// @BasePath  /api/v1

	// @securityDefinitions.basic  BasicAuth

	// @externalDocs.description  OpenAPI
	// @externalDocs.url          https://swagger.io/resources/open-api/
	cfg := config.MustLoad()

	log := logging.MustLoad(cfg)

	log.Info("starting property-finder")

	ctx := context.Background()

	pgIns := postgres.NewPG(ctx, cfg.Database)

	storage, err := postgres.New(ctx, pgIns)
	if err != nil {
		log.Error("failed to init storage: ", err)
		os.Exit(1)
	}
	defer storage.Close()
	r := router.InitRoutes(storage, &cfg.Server)

	log.Info("starting server", slog.String("host", cfg.Server.Host), slog.Int("port", cfg.Server.Port))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeouts.Read,
		WriteTimeout: cfg.Server.Timeouts.Write,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	defer postgres.Drop(ctx, pgIns)
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", err)

		return
	}

	log.Info("server stopped")
}
