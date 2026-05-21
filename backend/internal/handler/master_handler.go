package handler

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/yourorg/maintenance/internal/domain"
	"github.com/yourorg/maintenance/internal/response"
)

// MasterHandler menyajikan data master (vehicles, items, users) untuk keperluan frontend.
type MasterHandler struct {
	db *gorm.DB
}

func NewMasterHandler(db *gorm.DB) *MasterHandler {
	return &MasterHandler{db: db}
}

// GET /api/vehicles
func (h *MasterHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	var vehicles []domain.Vehicle
	if err := h.db.Order("license_plate ASC").Find(&vehicles).Error; err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.OK(w, "vehicles retrieved", vehicles)
}

// GET /api/master-items
func (h *MasterHandler) ListMasterItems(w http.ResponseWriter, r *http.Request) {
	var items []domain.MasterItem
	if err := h.db.Order("type ASC, item_name ASC").Find(&items).Error; err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.OK(w, "master items retrieved", items)
}

// GET /api/users
func (h *MasterHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []domain.User
	if err := h.db.Order("id ASC").Find(&users).Error; err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.OK(w, "users retrieved", users)
}
