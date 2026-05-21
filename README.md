# Fleetify

Sistem manajemen laporan perawatan kendaraan berbasis REST API menggunakan Go + MySQL.

---

## Tech Stack

* Go
* Chi Router
* GORM
* MySQL 8
* Docker & Docker Compose
* Vanilla JS + Bootstrap 5

---

## Menjalankan Project

### 1. Clone Repository

```bash
git clone <repository-url>
cd backend
```

### 2. Buat File Environment

```bash
cp env.example .env
```

Default `.env` sudah siap digunakan.

---

### 3. Jalankan Docker

```bash
docker-compose up --build
```

Service yang akan berjalan:

* Backend Go
* MySQL
* Seeder data otomatis

Backend tersedia di:

```txt
http://localhost:8081
```

---

### 4. Jalankan Frontend

Buka:

```txt
frontend/index.html
```

Disarankan menggunakan extension **Live Server** di VS Code.

---

## Testing

### Health Check

```bash
curl http://localhost:8081/health
```

Response:

```json
{
  "status": "ok"
}
```

---

## Akun Dummy

| Role     | User ID |
| -------- | ------- |
| SA       | 1       |
| APPROVAL | 2       |

Semua endpoint `/api/*` wajib menggunakan header:

```txt
X-User-ID: <id>
```

Contoh:

```bash
-H "X-User-ID: 1"
```

---

## Flow Laporan

```txt
PENDING -> APPROVED -> DONE
```

Keterangan:

* `SA` membuat & menyelesaikan laporan
* `APPROVAL` menyetujui laporan

---

# API Endpoint

## Master Data

### Get Users

```http
GET /api/users
```

---

### Get Vehicles

```http
GET /api/vehicles
```

---

### Get Master Items

```http
GET /api/master-items
```

---

## Reports

### Buat Laporan

```http
POST /api/reports
```

Request:

```json
{
  "vehicle_id": 1,
  "odometer": 52000,
  "complaint": "Oli bocor",
  "initial_photo": "base64-image",
  "items": [
    {
      "item_id": 1,
      "quantity": 1
    }
  ]
}
```

Catatan:

* Upload foto sudah menggunakan format base64
* Tidak menggunakan URL image lagi
* Minimal harus memiliki 1 item

---

### Get Reports

```http
GET /api/reports
```

---

### Approve Report

```http
PUT /api/reports/:id/approve
```

Role:

```txt
APPROVAL
```

---

### Complete Report

```http
PUT /api/reports/:id/complete
```

Request:

```json
{
  "proof_photo": "base64-image"
}
```

Catatan:

* Foto bukti juga menggunakan base64

---

## Contoh Request

### Create Report

```bash
curl -X POST http://localhost:8081/api/reports \
  -H "X-User-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_id": 1,
    "odometer": 52000,
    "complaint": "Oli bocor",
    "initial_photo": "base64-image",
    "items": [
      {
        "item_id": 1,
        "quantity": 1
      }
    ]
  }'
```

## Notes

* Seeder otomatis berjalan saat pertama kali container dibuat
* Jika mengubah schema/seeder:

```bash
docker-compose down -v
docker-compose up --build
```

* Semua upload foto disimpan dalam bentuk base64