package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/yourorg/maintenance/config"
	"github.com/yourorg/maintenance/internal/domain"
	"github.com/yourorg/maintenance/internal/handler"
	infraMySQL "github.com/yourorg/maintenance/internal/infra/mysql"
	appMiddleware "github.com/yourorg/maintenance/internal/middleware"
	repoMysql "github.com/yourorg/maintenance/internal/repository/mysql"
	"github.com/yourorg/maintenance/internal/service"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("gagal load config: %v", err)
	}

	db, err := infraMySQL.New(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}

	webhookURL := os.Getenv("WEBHOOK_URL")

	reportRepo := repoMysql.NewReportRepository(db)
	reportSvc := service.NewReportService(reportRepo, webhookURL)
	reportHandler := handler.NewReportHandler(reportSvc)
	masterHandler := handler.NewMasterHandler(db)

	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(appMiddleware.Authenticate(db))

		// ── Master data (semua role) ───────────────────────────────────────
		r.Get("/users",        masterHandler.ListUsers)
		r.Get("/vehicles",     masterHandler.ListVehicles)
		r.Get("/master-items", masterHandler.ListMasterItems)

		// ── Reports ───────────────────────────────────────────────────────
		r.With(appMiddleware.RequireRole(domain.RoleSA)).
			Post("/reports", reportHandler.Create)

		r.Get("/reports", reportHandler.List)

		r.Route("/reports/{id}", func(r chi.Router) {
			r.With(appMiddleware.RequireRole(domain.RoleApproval)).
				Put("/approve", reportHandler.Approve)

			r.With(appMiddleware.RequireRole(domain.RoleSA)).
				Put("/complete", reportHandler.Complete)
		})
	})

	addr := ":" + cfg.AppPort
	log.Printf("[server] berjalan di %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}