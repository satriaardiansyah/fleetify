package mysql

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/yourorg/maintenance/internal/domain"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// CreateWithItems insert report header + items dalam satu transaksi atomik.
// Price snapshot diambil dari master_items saat insert — bukan dari input user.
func (r *ReportRepository) CreateWithItems(
	ctx context.Context,
	report *domain.MaintenanceReport,
	items []domain.ReportItem,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Insert report header
		if err := tx.Create(report).Error; err != nil {
			return fmt.Errorf("insert maintenance_report: %w", err)
		}

		// 2. Loop item: ambil harga dari master, lalu insert detail
		for i := range items {
			items[i].ReportID = report.ID

			var master domain.MasterItem
			if err := tx.First(&master, items[i].ItemID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("master_item id=%d tidak ditemukan", items[i].ItemID)
				}
				return fmt.Errorf("query master_item id=%d: %w", items[i].ItemID, err)
			}

			// Price Snapshot: salin harga saat ini, bukan dari input user
			items[i].PriceSnapshot = master.Price
			items[i].MasterItem = &master

			if err := tx.Create(&items[i]).Error; err != nil {
				return fmt.Errorf("insert report_item: %w", err)
			}
		}

		report.ReportItems = items
		return nil
	})
}

// UpdateStatus mengubah status laporan.
// Validasi transisi status dilakukan di service layer.
func (r *ReportRepository) UpdateStatus(
	ctx context.Context,
	id uint64,
	status domain.ReportStatus,
	proofPhoto string, // hanya diisi untuk status DONE, kosong untuk APPROVED
) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if proofPhoto != "" {
		updates["proof_photo"] = proofPhoto
	}

	result := r.db.WithContext(ctx).
		Model(&domain.MaintenanceReport{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID mengembalikan satu laporan beserta semua asosiasi.
func (r *ReportRepository) FindByID(ctx context.Context, id uint64) (*domain.MaintenanceReport, error) {
	var report domain.MaintenanceReport
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Creator").
		Preload("ReportItems.MasterItem").
		First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// FindAll mengembalikan seluruh laporan untuk F-04.
func (r *ReportRepository) FindAll(ctx context.Context) ([]domain.MaintenanceReport, error) {
	var reports []domain.MaintenanceReport
	err := r.db.WithContext(ctx).
		Preload("Vehicle").
		Preload("Creator").
		Preload("ReportItems.MasterItem").
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}
