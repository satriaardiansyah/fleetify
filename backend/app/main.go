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
	appMiddleware "github.com/yourorg/maintenance/internal/middleware"
	infraMySQL "github.com/yourorg/maintenance/internal/infra/mysql"
	repoMysql "github.com/yourorg/maintenance/internal/repository/mysql"
	"github.com/yourorg/maintenance/internal/service"
)

func main() {
	// ── Config ────────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("gagal load config: %v", err)
	}

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := infraMySQL.New(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}

	// ── Dependency Injection ──────────────────────────────────────────────────
	// WEBHOOK_URL dari env, kosong = fitur webhook dimatikan
	webhookURL := os.Getenv("WEBHOOK_URL")

	reportRepo    := repoMysql.NewReportRepository(db)
	reportSvc     := service.NewReportService(reportRepo, webhookURL)
	reportHandler := handler.NewReportHandler(reportSvc)

	// ── Router ────────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)

	// Health check — tanpa auth
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		// Semua route /api/* butuh X-User-ID yang valid
		r.Use(appMiddleware.Authenticate(db))

		// F-01: Buat laporan — hanya SA
		r.With(appMiddleware.RequireRole(domain.RoleSA)).
			Post("/reports", reportHandler.Create)

		// F-04: Daftar semua laporan — semua role
		r.Get("/reports", reportHandler.List)

		r.Route("/reports/{id}", func(r chi.Router) {
			// F-02: Approve laporan — hanya APPROVAL
			r.With(appMiddleware.RequireRole(domain.RoleApproval)).
				Put("/approve", reportHandler.Approve)

			// F-03: Complete laporan — hanya SA
			r.With(appMiddleware.RequireRole(domain.RoleSA)).
				Put("/complete", reportHandler.Complete)
		})
	})

	// ── Server ────────────────────────────────────────────────────────────────
	addr := ":" + cfg.AppPort
	log.Printf("[server] berjalan di %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
