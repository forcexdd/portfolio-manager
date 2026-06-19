package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forcexdd/portfoliomanager/internal/adapter/moex"
	"github.com/forcexdd/portfoliomanager/internal/config"
	"github.com/forcexdd/portfoliomanager/internal/logger"
	"github.com/forcexdd/portfoliomanager/internal/repository"
	"github.com/forcexdd/portfoliomanager/internal/service"
	httptransport "github.com/forcexdd/portfoliomanager/internal/transport/http"
	"github.com/forcexdd/portfoliomanager/internal/transport/http/api"
	"github.com/forcexdd/portfoliomanager/internal/worker"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	log := logger.New()

	db, err := sql.Open("postgres", cfg.DBConnStr)
	if err != nil {
		log.Fatal("DB connection failed", "error", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed", "error", err)
	}

	pRepo := repository.NewPostgresPortfolioRepository(db)
	aRepo := repository.NewPostgresAssetRepository(db)
	iRepo := repository.NewPostgresIndexRepository(db)
	sysRepo := repository.NewPostgresSystemRepository(db)

	pSvc := service.NewPortfolioService(pRepo, aRepo, log)
	aSvc := service.NewAssetService(aRepo, log)
	iSvc := service.NewIndexService(iRepo, log)
	anSvc := service.NewAnalysisService(pRepo, iRepo, log)
	sysSvc := service.NewSystemService(sysRepo, log)

	moexClient := moex.NewMoexClient(cfg.MoexAPIURL, &http.Client{Timeout: 10 * time.Second})
	parser := worker.NewParserWorker(aSvc, iSvc, sysSvc, moexClient, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go parser.Run(ctx)

	handler := httptransport.NewHandler(pSvc, aSvc, iSvc, anSvc, log)
	apiServer, err := api.NewServer(handler)
	if err != nil {
		log.Fatal("Failed to create ogen server", "error", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", apiServer)
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	middlewareChain := httptransport.ContextLoggerMiddleware(log)(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: middlewareChain,
	}

	go func() {
		log.Info("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}
}
