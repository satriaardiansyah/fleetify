package domain

import "time"

// ── Users ─────────────────────────────────────────────────────────────────────

type UserRole string

const (
	RoleSA       UserRole = "SA"
	RoleApproval UserRole = "APPROVAL"
)

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:100;not null;unique"  json:"username"`
	Role      UserRole  `gorm:"type:enum('SA','APPROVAL');not null" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Vehicles ──────────────────────────────────────────────────────────────────

type Vehicle struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	LicensePlate string    `gorm:"size:50;not null;unique"  json:"license_plate"`
	Model        string    `gorm:"size:100;not null"        json:"model"`
	CreatedAt    time.Time `json:"created_at"`
}

// ── Master Items ──────────────────────────────────────────────────────────────

type ItemType string

const (
	ItemTypePart    ItemType = "PART"
	ItemTypeService ItemType = "SERVICE"
)

type MasterItem struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemName  string    `gorm:"size:255;not null"        json:"item_name"`
	Type      ItemType  `gorm:"type:enum('PART','SERVICE');not null" json:"type"`
	Price     float64   `gorm:"type:decimal(15,2);not null;default:0" json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// ── Maintenance Reports ───────────────────────────────────────────────────────

type ReportStatus string

const (
	StatusPending  ReportStatus = "PENDING"
	StatusApproved ReportStatus = "APPROVED"
	StatusRejected ReportStatus = "REJECTED"
	StatusDone     ReportStatus = "DONE"
)

type MaintenanceReport struct {
	ID           uint64       `gorm:"primaryKey;autoIncrement" json:"id"`
	VehicleID    uint64       `gorm:"not null"                 json:"vehicle_id"`
	CreatedBy    uint64       `gorm:"not null"                 json:"created_by"`
	Odometer     int64        `gorm:"not null"                 json:"odometer"`
	Complaint    string       `gorm:"type:text"                json:"complaint"`
	Status       ReportStatus `gorm:"type:enum('PENDING','APPROVED','REJECTED','DONE');default:'PENDING'" json:"status"`
	InitialPhoto string       `gorm:"size:255"                 json:"initial_photo"`
	ProofPhoto   string       `gorm:"size:255"                 json:"proof_photo"`
	CreatedAt    time.Time    `json:"created_at"`

	// Preloadable associations
	Vehicle     *Vehicle      `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Creator     *User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	ReportItems []ReportItem  `gorm:"foreignKey:ReportID"  json:"report_items,omitempty"`
}

// ── Report Items ──────────────────────────────────────────────────────────────

type ReportItem struct {
	ID            uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ReportID      uint64      `gorm:"not null"                 json:"report_id"`
	ItemID        uint64      `gorm:"not null"                 json:"item_id"`
	Quantity      int         `gorm:"not null;default:1"       json:"quantity"`
	PriceSnapshot float64     `gorm:"type:decimal(15,2);not null;default:0" json:"price_snapshot"`
	CreatedAt     time.Time   `json:"created_at"`

	MasterItem *MasterItem `gorm:"foreignKey:ItemID" json:"master_item,omitempty"`
}
