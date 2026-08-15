# Panduan Deployment Batam MedHub di VPS

Panduan lengkap untuk men-deploy backend **Batam MedHub** (Core Backend Orchestrator, 4 Mock Providers, PostgreSQL, dan Caddy SSL Reverse Proxy) ke Virtual Private Server (VPS) berbasis Linux (Ubuntu/Debian) menggunakan Docker Compose.

---

## 1. Arsitektur Deployment

Seluruh stack berjalan dalam container Docker terisolasi:

```
[ Internet / Mobile App ]
           │
           │ HTTPS (Port 443) / HTTP (Port 80)
           ▼
   +---------------+
   | Caddy Ingress |  (Otomatis Let's Encrypt SSL, Reverse Proxy, Security Headers)
   +-------┬-------+
           │ HTTP (Internal Docker Network)
           ▼
   +---------------+
   | Core Backend  |  (Go API Orchestrator :8080)
   +───┬───┬───┬───+
       │   │   │   │
       │   │   │   └───────────────┐
       │   │   └──────────────┐    │
       ▼   ▼                  ▼    ▼
   +---------------+      +-------------------------+
   | Core DB       |      | 4 Mock Provider APIs    |
   | (PostgreSQL)  |      | - Hospital (:8080)      |
   +---------------+      | - Ferry    (:8080)      |
                          | - Hotel    (:8080)      |
                          | - Transport(:8080)      |
                          +────────────┬────────────+
                                       │
                                       ▼
                          +-------------------------+
                          | Provider Databases      |
                          | (hospital_db, ferry_db, |
                          |  hotel_db, transport_db)|
                          +-------------------------+
```

---

## 2. Prasyarat VPS

- **OS**: Ubuntu 22.04 LTS / 24.04 LTS atau Debian 12
- **Spesifikasi Minimal**: 1 vCPU, 2 GB RAM, 20 GB Storage
- **Akses**: Root / User dengan hak `sudo` via SSH
- **Domain / Subdomain**: A Record DNS mengarah ke IP Public VPS (contoh: `api.domainkamu.com` -> `103.xxx.xxx.xxx`)

---

## 3. Persiapan Server VPS (Langkah Awal)

### A. Update Paket Server
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y git curl wget ufw
```

### B. Konfigurasi Firewall (UFW)
Pastikan port SSH (22), HTTP (80), dan HTTPS (443) dibuka:
```bash
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### C. Instalasi Docker & Docker Compose
Gunakan skrip resmi Docker:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Tambahkan user aktif ke grup docker agar tidak perlu sudo setiap saat:
sudo usermod -aG docker $USER
newgrp docker

# Verifikasi instalasi:
docker --version
docker compose version
```

---

## 4. Langkah-Langkah Deployment

### Langkah 1: Clone Repository ke VPS
```bash
cd ~
git clone https://github.com/bayuzzs/batam-medhub.git
cd batam-medhub
```

### Langkah 2: Buat & Konfigurasi File `.env.production`
Salin template konfigurasi:
```bash
cp .env.production.example .env.production
nano .env.production
```

Isi dan sesuaikan variabel penting berikut:

1. **`DOMAIN`**: Masukkan domain/subdomain kamu (contoh: `api.medhub.example.com`).
2. **`JWT_SIGNING_SECRET`**: Buat secret acak 32+ karakter:
   ```bash
   openssl rand -hex 32
   ```
3. **Database Passwords & Provider Secrets**: Buat password acak:
   ```bash
   openssl rand -base64 24
   ```
   Ganti nilai `POSTGRES_PASSWORD`, `CORE_DB_PASSWORD`, `HOSPITAL_DB_PASSWORD`, `FERRY_DB_PASSWORD`, `HOTEL_DB_PASSWORD`, `TRANSPORT_DB_PASSWORD`, `DEMO_SECRET`, dan seluruh `*_INTEGRATION_KEY`.
4. **Cloudflare AI (Opsional)**: Masukkan `CLOUDFLARE_ACCOUNT_ID` dan `CLOUDFLARE_API_TOKEN` jika ingin menggunakan model LLM Workers AI (jika kosong, backend otomatis memakai rule-based engine).

Simpan file dengan menekan `Ctrl + O` lalu `Enter`, kemudian keluar dengan `Ctrl + X`.

### Langkah 3: Eksekusi Deployment
Jalankan skrip deployment otomatis:
```bash
./scripts/deploy.sh
```

Skrip ini akan secara otomatis:
1. Memvalidasi environment dan dependency Docker.
2. Melakukan build container untuk Core Backend dan 4 Mock Providers.
3. Menginisialisasi 5 logical database PostgreSQL (`core_db`, `hospital_db`, `ferry_db`, `hotel_db`, `transport_db`).
4. Menjalankan seluruh database migration secara otomatis (`core-migrate`, `hospital-migrate`, dll).
5. Menyalakan Caddy reverse proxy yang otomatis mendaftarkan sertifikat SSL Let's Encrypt.
6. Melakukan automated health check ke Core Backend.

---

## 5. Verifikasi & Pengujian

### A. Cek Status Container
```bash
docker compose --env-file .env.production -f docker-compose.prod.yml ps
```
Semua container harus berstatus `Up (healthy)` atau `Up`.

### B. Health Check via Domain (HTTPS)
Uji endpoint health check dari komputer lokal atau terminal:
```bash
curl -i https://<DOMAIN_KAMU>/healthz
```
*Respons yang diharapkan:* `HTTP/2 200` dengan JSON `{"status":"ok"}`.

### C. End-to-End Smoke Test
Uji alur lengkap registrasi patient & trip request:

```bash
# 1. Register Patient Baru
curl -s -X POST https://<DOMAIN_KAMU>/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Eleanor Vance",
    "email": "eleanor.vance@example.com",
    "password": "Password123!",
    "preferred_currency": "SGD",
    "nationality": "SG"
  }'

# 2. Ambil token dari respons dan buat trip request
TOKEN="<MASUKKAN_ACCESS_TOKEN_DARI_STEP_1>"

curl -s -X POST https://<DOMAIN_KAMU>/v1/trip-requests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "I need a comprehensive health screening in Batam on 22 August 2026 for 1 person.",
    "preferred_currency": "SGD"
  }'
```

---

## 6. Operasional & Pemeliharaan (Day-2 Operations)

### Melihat Log Realtime
```bash
# Semua service:
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f

# Khusus Core Backend:
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f backend

# Khusus Caddy (Traffic & SSL):
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f caddy
```

### Melakukan Update Kode ke Versi Terbaru
Ketika ada perubahan kode di repository GitHub:
```bash
cd ~/batam-medhub
git pull origin main
./scripts/deploy.sh
```

### Menghentikan & Menyalakan Ulang Stack
```bash
# Restart stack:
docker compose --env-file .env.production -f docker-compose.prod.yml restart

# Stop stack (data tetap aman di Docker Volume):
docker compose --env-file .env.production -f docker-compose.prod.yml down
```

### Backup Database PostgreSQL
Untuk membuat backup seluruh database:
```bash
docker compose --env-file .env.production -f docker-compose.prod.yml exec -T postgres pg_dumpall -U provider_admin > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Reset Demo State
Jika ingin mereset state transaksi dan demo data ke kondisi awal:
```bash
curl -s -X POST https://<DOMAIN_KAMU>/v1/demo/reset \
  -H "X-Demo-Key: <DEMO_SECRET>" \
  -H "Content-Type: application/json" \
  -d '{
    "scenario": "DEFAULT",
    "confirm": true
  }'
```

---

## 7. Troubleshooting

| Masalah | Penyebab Umum | Solusi |
|---|---|---|
| **SSL Certificate Error (Caddy)** | DNS Domain belum mengarah ke IP VPS atau Port 80/443 diblokir firewall. | Pastikan DNS A record sudah terpropagasi (`dig +short A domain.com`) dan port 80/443 terbuka di UFW (`sudo ufw status`). |
| **Backend 500 / Database Error** | Database belum selesai inisialisasi atau password salah. | Cek log backend: `docker compose logs backend` dan log postgres: `docker compose logs postgres`. |
| **Port Conflict (80/443 already in use)** | Ada Nginx/Apache bawaan OS yang masih berjalan di VPS. | Hentikan webserver host: `sudo systemctl stop nginx apache2 && sudo systemctl disable nginx apache2`. |
