package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Payload adalah body yang dikirim ke URL eksternal.
type Payload struct {
	Event    string `json:"event"`     // contoh: "report.approved" / "report.completed"
	ReportID uint64 `json:"report_id"`
	Status   string `json:"status"`
	OccurredAt string `json:"occurred_at"`
}

// Fire mengirim HTTP POST ke webhookURL secara asinkronus di goroutine terpisah.
// Fungsi ini langsung return — caller tidak perlu menunggu.
//
// Contoh pemakaian:
//
//	webhook.Fire("https://hooks.example.com/notify", "report.approved", report.ID, "APPROVED")
func Fire(webhookURL, event string, reportID uint64, status string) {
	if webhookURL == "" {
		return // webhook tidak dikonfigurasi, skip
	}

	payload := Payload{
		Event:      event,
		ReportID:   reportID,
		Status:     status,
		OccurredAt: time.Now().Format(time.RFC3339),
	}

	// Jalankan di goroutine agar tidak memblokir response ke client
	go func() {
		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[webhook] marshal error: %v", err)
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}

		resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("[webhook] POST %s error: %v", webhookURL, err)
			return
		}
		defer resp.Body.Close()

		log.Printf("[webhook] event=%s report_id=%d -> %s %d",
			event, reportID, webhookURL, resp.StatusCode)
	}()
}
