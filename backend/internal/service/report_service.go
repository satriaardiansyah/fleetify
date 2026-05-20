package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/yourorg/maintenance/internal/domain"
	"github.com/yourorg/maintenance/internal/dto"
	repoMysql "github.com/yourorg/maintenance/internal/repository/mysql"
	"github.com/yourorg/maintenance/pkg/webhook"
)

type ReportService struct {
	repo       *repoMysql.ReportRepository
	webhookURL string // kosong = fitur webhook dimatikan
}

func NewReportService(repo *repoMysql.ReportRepository, webhookURL string) *ReportService {
	return &ReportService{repo: repo, webhookURL: webhookURL}
}

// ── F-01: Buat laporan baru ────────────────────────────────────────────────────

func (s *ReportService) Create(
	ctx context.Context,
	creator *domain.User,
	req dto.CreateReportRequest,
) (*dto.CreateReportResponse, error) {

	if req.VehicleID == 0 {
		return nil, errors.New("vehicle_id wajib diisi")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("minimal satu item harus dimasukkan")
	}

	report := &domain.MaintenanceReport{
		VehicleID:    req.VehicleID,
		CreatedBy:    creator.ID,
		Odometer:     req.Odometer,
		Complaint:    req.Complaint,
		Status:       domain.StatusPending, // selalu PENDING saat dibuat
		InitialPhoto: req.InitialPhoto,
	}

	items := make([]domain.ReportItem, len(req.Items))
	for i, inp := range req.Items {
		items[i] = domain.ReportItem{
			ItemID:   inp.ItemID,
			Quantity: inp.Quantity,
		}
	}

	if err := s.repo.CreateWithItems(ctx, report, items); err != nil {
		return nil, err
	}

	return buildCreateResponse(report, creator), nil
}

// ── F-02: Approve laporan (hanya role APPROVAL) ────────────────────────────────

func (s *ReportService) Approve(
	ctx context.Context,
	actor *domain.User,
	reportID uint64,
) (*dto.UpdateReportResponse, error) {

	// Ambil laporan saat ini untuk validasi transisi status
	report, err := s.repo.FindByID(ctx, reportID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("laporan tidak ditemukan")
		}
		return nil, err
	}

	// Guard: hanya laporan berstatus PENDING yang bisa di-approve
	if report.Status != domain.StatusPending {
		return nil, fmt.Errorf(
			"laporan tidak bisa di-approve, status saat ini: %s (harus PENDING)",
			report.Status,
		)
	}

	if err := s.repo.UpdateStatus(ctx, reportID, domain.StatusApproved, ""); err != nil {
		return nil, err
	}

	// B-02: Kirim webhook secara asinkronus (goroutine) — tidak memblokir response
	webhook.Fire(s.webhookURL, "report.approved", reportID, string(domain.StatusApproved))

	return &dto.UpdateReportResponse{
		ID:      reportID,
		Status:  string(domain.StatusApproved),
		Message: "laporan berhasil disetujui",
	}, nil
}

// ── F-03: Complete laporan (hanya role SA) ─────────────────────────────────────

func (s *ReportService) Complete(
	ctx context.Context,
	actor *domain.User,
	reportID uint64,
	req dto.CompleteReportRequest,
) (*dto.UpdateReportResponse, error) {

	if req.ProofPhoto == "" {
		return nil, errors.New("proof_photo wajib diisi")
	}

	// Ambil laporan untuk validasi transisi status
	report, err := s.repo.FindByID(ctx, reportID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("laporan tidak ditemukan")
		}
		return nil, err
	}

	// Guard: hanya laporan berstatus APPROVED yang bisa di-complete
	if report.Status != domain.StatusApproved {
		return nil, fmt.Errorf(
			"laporan tidak bisa diselesaikan, status saat ini: %s (harus APPROVED)",
			report.Status,
		)
	}

	if err := s.repo.UpdateStatus(ctx, reportID, domain.StatusDone, req.ProofPhoto); err != nil {
		return nil, err
	}

	// B-02: Kirim webhook secara asinkronus (goroutine) — tidak memblokir response
	webhook.Fire(s.webhookURL, "report.completed", reportID, string(domain.StatusDone))

	return &dto.UpdateReportResponse{
		ID:         reportID,
		Status:     string(domain.StatusDone),
		ProofPhoto: req.ProofPhoto,
		Message:    "laporan berhasil diselesaikan",
	}, nil
}

// ── F-04: Daftar semua laporan ─────────────────────────────────────────────────

func (s *ReportService) ListAll(ctx context.Context) ([]domain.MaintenanceReport, error) {
	return s.repo.FindAll(ctx)
}

// ── Helper ─────────────────────────────────────────────────────────────────────

func buildCreateResponse(
	report *domain.MaintenanceReport,
	creator *domain.User,
) *dto.CreateReportResponse {

	resp := &dto.CreateReportResponse{
		ID:           report.ID,
		VehicleID:    report.VehicleID,
		CreatedBy:    creator.ID,
		CreatorName:  creator.Username,
		Odometer:     report.Odometer,
		Complaint:    report.Complaint,
		Status:       string(report.Status),
		InitialPhoto: report.InitialPhoto,
	}

	if report.Vehicle != nil {
		resp.LicensePlate = report.Vehicle.LicensePlate
	}

	var total float64
	for _, item := range report.ReportItems {
		subtotal := item.PriceSnapshot * float64(item.Quantity)
		total += subtotal

		ri := dto.ReportItemResponse{
			ID:            item.ID,
			ItemID:        item.ItemID,
			Quantity:      item.Quantity,
			PriceSnapshot: item.PriceSnapshot,
			Subtotal:      subtotal,
		}
		if item.MasterItem != nil {
			ri.ItemName = item.MasterItem.ItemName
			ri.ItemType  = string(item.MasterItem.Type)
		}
		resp.Items = append(resp.Items, ri)
	}
	resp.TotalPrice = total

	return resp
}
