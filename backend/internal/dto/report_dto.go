package dto

// ── F-01 Request ──────────────────────────────────────────────────────────────

type CreateReportRequest struct {
	VehicleID    uint64            `json:"vehicle_id"`
	Odometer     int64             `json:"odometer"`
	Complaint    string            `json:"complaint"`
	InitialPhoto string            `json:"initial_photo"`
	Items        []ReportItemInput `json:"items"`
}

type ReportItemInput struct {
	ItemID   uint64 `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// ── F-02 Request ──────────────────────────────────────────────────────────────

// ApproveReportRequest - body opsional, catatan dari approval
type ApproveReportRequest struct {
	Note string `json:"note"`
}

// ── F-03 Request ──────────────────────────────────────────────────────────────

type CompleteReportRequest struct {
	ProofPhoto string `json:"proof_photo"` // URL / path foto bukti pengerjaan
}

// ── Response shapes ───────────────────────────────────────────────────────────

type ReportItemResponse struct {
	ID            uint64  `json:"id"`
	ItemID        uint64  `json:"item_id"`
	ItemName      string  `json:"item_name"`
	ItemType      string  `json:"item_type"`
	Quantity      int     `json:"quantity"`
	PriceSnapshot float64 `json:"price_snapshot"`
	Subtotal      float64 `json:"subtotal"`
}

type CreateReportResponse struct {
	ID           uint64               `json:"id"`
	VehicleID    uint64               `json:"vehicle_id"`
	LicensePlate string               `json:"license_plate"`
	CreatedBy    uint64               `json:"created_by"`
	CreatorName  string               `json:"creator_name"`
	Odometer     int64                `json:"odometer"`
	Complaint    string               `json:"complaint"`
	Status       string               `json:"status"`
	InitialPhoto string               `json:"initial_photo"`
	Items        []ReportItemResponse `json:"items"`
	TotalPrice   float64              `json:"total_price"`
}

type UpdateReportResponse struct {
	ID         uint64 `json:"id"`
	Status     string `json:"status"`
	ProofPhoto string `json:"proof_photo,omitempty"`
	Message    string `json:"message"`
}
