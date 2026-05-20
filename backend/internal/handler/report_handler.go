package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/yourorg/maintenance/internal/dto"
	"github.com/yourorg/maintenance/internal/middleware"
	"github.com/yourorg/maintenance/internal/response"
	"github.com/yourorg/maintenance/internal/service"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// ── F-01: POST /api/reports ───────────────────────────────────────────────────

func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	var req dto.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid: "+err.Error())
		return
	}

	result, err := h.svc.Create(r.Context(), user, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, "laporan berhasil dibuat", result)
}

// ── F-02: PUT /api/reports/:id/approve ───────────────────────────────────────

func (h *ReportHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "id laporan tidak valid")
		return
	}

	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	result, err := h.svc.Approve(r.Context(), user, id)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, result.Message, result)
}

// ── F-03: PUT /api/reports/:id/complete ──────────────────────────────────────

func (h *ReportHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		response.BadRequest(w, "id laporan tidak valid")
		return
	}

	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	var req dto.CompleteReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "body tidak valid: "+err.Error())
		return
	}

	result, err := h.svc.Complete(r.Context(), user, id, req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.OK(w, result.Message, result)
}

// ── F-04: GET /api/reports ────────────────────────────────────────────────────

func (h *ReportHandler) List(w http.ResponseWriter, r *http.Request) {
	reports, err := h.svc.ListAll(r.Context())
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}
	response.OK(w, "daftar laporan berhasil diambil", reports)
}

// ── Helper ────────────────────────────────────────────────────────────────────

func parseIDParam(r *http.Request, param string) (uint64, error) {
	return strconv.ParseUint(chi.URLParam(r, param), 10, 64)
}
