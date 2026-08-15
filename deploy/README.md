# Batam MedHub — Production Deployment (Docker Compose + nginx + HTTPS)

This guide deploys the **full Batam MedHub stack** (core backend API + all four
provider services + PostgreSQL) on a Linux VPS with Docker and Compose, then
exposes the API publicly at **`https://api.bayumaulana.my.id`** through nginx
with a Let's Encrypt HTTPS certificate.

```
                        HTTPS
   Mobile App  ────────────▶  api.bayumaulana.my.id
                                  │ nginx (host, :443)
                                  ▼
                         127.0.0.1:8080
                                  │
                        ┌─────────▼──────────┐
                        │   backend (Go)     │
                        └────┬─────┬────┬────┘
                        hospital ferry hotel transport   (internal network)
                             └───┬──────┘
                             postgres:17 (5 databases)
```

---

## 1. Prerequisites

On the VPS (Debian/Ubuntu assumed):

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-v2 nginx git curl
sudo systemctl enable --now docker nginx
```

Check that both are present:

```bash
docker --version
docker compose version
nginx -v
```

> The backend talks to all four providers **over the Docker internal network by
> service name** — only nginx (running on the host) reaches the backend, and it
> does so via `127.0.0.1:8080`.

---

## 2. DNS — point the domain at the VPS

Create an **A record** in your DNS provider (e.g. Cloudflare DNS / registrar):

| Type | Name                 | Value      | TTL  |
|------|----------------------|------------|------|
| A    | `api`                | `<VPS_IP>` | Auto |

So `api.bayumaulana.my.id` resolves to the VPS public IP. Verify:

```bash
dig +short api.bayumaulana.my.id
# should print your VPS IP
```

> If you use Cloudflare DNS with the orange (proxied) cloud, ensure SSL mode is
> **Full (strict)** so the origin nginx can present its Let's Encrypt cert.

---

## 3. Get the code on the VPS

```bash
git clone <your-batam-medhub-repo-url> /opt/batam-medhub
cd /opt/batam-medhub/deploy
```

---

## 4. Configure environment secrets

```bash
cp .env.example .env
nano .env
```

Generate strong random values with `openssl rand -hex 24` and paste them into
the `CHANGE_ME_*` fields:

- `PROVIDER_POSTGRES_PASSWORD`, `HOSPITAL_DB_PASSWORD`, `FERRY_DB_PASSWORD`,
  `HOTEL_DB_PASSWORD`, `TRANSPORT_DB_PASSWORD` — PostgreSQL credentials.
- `HOSPITAL_INTEGRATION_KEY` … `TRANSPORT_INTEGRATION_KEY` — shared secrets the
  backend uses to authenticate each provider.
- `JWT_SIGNING_SECRET` — **at least 32 characters** (validated at startup).
- `DEMO_SECRET` — authenticates `POST /v1/demo/reset`.
- `CLOUDFLARE_ACCOUNT_ID`, `CLOUDFLARE_API_TOKEN` — your Cloudflare Workers AI
  credentials (used for LLM intent extraction).
- `CLOUDFLARE_AI_BASE_URL` — the OpenAI-compatible gateway URL; the backend
  auto-detects the `/ai/v1` suffix and calls `/chat/completions`.

> ⚠️ `.env` is git-ignored. Never commit it.

---

## 5. Build and start the stack

```bash
docker compose up -d --build
```

Compose starts, in order:

1. `postgres` — creates the four provider databases/roles and seeds them.
2. `core-db-init` — creates `core_db` for the backend.
3. `hospital-migrate` / `ferry-migrate` / `hotel-migrate` / `transport-migrate` —
   apply provider schemas + golden seed data (one-shot).
4. `backend-migrate` — applies the core schema migrations (one-shot).
5. `hospital` / `ferry` / `hotel` / `transport` — the provider APIs.
6. `backend` — the core API, published on `127.0.0.1:8080`.

Watch startup:

```bash
docker compose ps
docker compose logs -f backend
```

Healthy outcome: `postgres`, the four providers, and `backend` are `running
(healthy)`; the migrate/init services exited with code `0`.

---

## 6. Verify the stack locally (before HTTPS)

```bash
curl -s http://127.0.0.1:8080/healthz      # liveness  -> 200
curl -s http://127.0.0.1:8080/readyz       # readiness -> 200 (DB reachable)
```

Run the full smoke flow (register → AI trip request → plan → book):

```bash
# 1. Optional: reset to clean demo state
curl -s -X POST http://127.0.0.1:8080/v1/demo/reset \
  -H "X-Demo-Key: <DEMO_SECRET>" -H "Content-Type: application/json" \
  -d '{"scenario":"DEFAULT","confirm":true}'

# 2. Register a synthetic patient
curl -s -X POST http://127.0.0.1:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Eleanor Vance","email":"eleanor.vance@example.com",\
       "password":"Password123!","preferred_currency":"SGD","nationality":"SG"}'
# -> copy the access token
TOKEN="<access_token>"

# 3. Create a trip request (uses Cloudflare Workers AI)
curl -s -X POST http://127.0.0.1:8080/v1/trip-requests \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"prompt":"I need a basic medical check-up in Batam on August 22, 2026 for 1 person.","preferred_currency":"SGD"}'
# -> copy trip_request_id
TRIP_ID="<trip_request_id>"

# 4. Generate plan options
curl -s -X POST http://127.0.0.1:8080/v1/trip-requests/$TRIP_ID/plans \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{}'
# -> copy plan_option_id
PLAN_OPTION_ID="<plan_option_id>"

# 5. Book it (multi-provider saga)
curl -s -X POST http://127.0.0.1:8080/v1/plan-options/$PLAN_OPTION_ID/confirm \
  -H "Authorization: Bearer $TOKEN" -H "Idempotency-Key: idem-smoke-001" \
  -H "Content-Type: application/json" -d '{"approved":true}'
# -> active journey JSON with confirmed reservations
```

If trip planning returns `NO_MATCH` for the date, pin the year in the prompt
(e.g. "August 22, **2026**") — the demo only seeds 2026-08-22 availability.

---

## 7. nginx reverse proxy + HTTPS

### 7.1 Install the site config

```bash
sudo cp deploy/nginx/api.bayumaulana.my.id.conf /etc/nginx/sites-available/
sudo ln -s /etc/nginx/sites-available/api.bayumaulana.my.id.conf /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

(If using a distro that loads `/etc/nginx/conf.d/*.conf` instead of
`sites-available`, copy the file there instead.)

Now `http://api.bayumaulana.my.id` redirects to HTTPS — but the cert doesn't
exist yet, so the 443 block is inactive until step 7.2.

### 7.2 Obtain a Let's Encrypt certificate

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d api.bayumaulana.my.id
```

Follow the prompts (agree to terms, provide an email). Certbot will install the
certificate into the config and reload nginx automatically.

Verify the cert auto-renews (a timer is installed by the package):

```bash
sudo certbot renew --dry-run
```

### 7.3 Final check

```bash
curl -s https://api.bayumaulana.my.id/healthz   # 200 over HTTPS
curl -sI https://api.bayumaulana.my.id/readyz   # HTTP/2 200
```

Optionally test with an actual registration:

```bash
curl -s -X POST https://api.bayumaulana.my.id/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Bayu Test","email":"bayu.test@example.com",\
       "password":"Password123!","preferred_currency":"SGD","nationality":"ID"}'
```

---

## 8. Point the mobile app at the deployed API

The Flutter app lets you override the API base URL at build time. The
`API_BASE_URL` dart-define wins over any platform default:

```bash
cd mobile
flutter run --dart-define=API_BASE_URL=https://api.bayumaulana.my.id
```

For a release build:

```bash
flutter build apk --release --dart-define=API_BASE_URL=https://api.bayumaulana.my.id
```

---

## 9. Day-to-day operations

| Task | Command |
|------|---------|
| Show stack status | `docker compose ps` |
| Tail backend logs | `docker compose logs -f backend` |
| Tail a provider | `docker compose logs -f ferry` |
| Restart a service | `docker compose restart backend` |
| Rebuild + restart after a code change | `docker compose up -d --build` |
| Stop everything | `docker compose down` |
| Stop + wipe databases | `docker compose down -v` |
| Reset demo data (keeps reference data) | `curl -X POST http://127.0.0.1:8080/v1/demo/reset -H "X-Demo-Key: <DEMO_SECRET>" -d '{"scenario":"DEFAULT","confirm":true}'` |
| Full provider re-seed (drop + migrate) | `docker compose up -d --force-recreate` then run the `*-migrate` services: `docker compose run --rm hospital-migrate -path=/migrations -database="postgres://hospital_user:<pass>@postgres:5432/hospital_db?sslmode=disable" down -all` … `up` |

### Backups (recommended)

Back up the PostgreSQL volume before any major change:

```bash
docker compose exec postgres pg_dumpall -U provider_admin > backup_$(date +%F).sql
```

Restore:

```bash
cat backup_YYYY-MM-DD.sql | docker compose exec -T postgres psql -U provider_admin -d postgres
```

---

## 10. Troubleshooting

| Symptom | Likely cause / fix |
|---------|--------------------|
| `docker compose up` fails with "required variable ... not set" | `.env` is missing or a `CHANGE_ME_*` value was left in place — fix `.env` and rerun. |
| `backend` not healthy | DB migration didn't finish. Check `docker compose logs backend-migrate core-db-init postgres`. |
| Trip request returns `NO_MATCH` | LLM parsed the date in the past. Pin the year: "August 22, **2026**". |
| Booking fails `CAPACITY_CONFLICT` | The tiny synthetic demo slots (capacity 1) are already consumed. Run `/v1/demo/reset` or re-seed, then re-plan. |
| `503` / TLS handshake errors over HTTPS | Cert not installed yet — run `sudo certbot --nginx -d api.bayumaulana.my.id`. |
| Cloudflare proxied domain won't load | Set SSL/TLS mode to **Full (strict)** in the Cloudflare dashboard. |
| Need to expose a provider for debugging | Temporarily uncomment the `ports:` block for that service in `docker-compose.yml`. |

---

## 11. Layout reference

```
deploy/
├── docker-compose.yml          # full stack (backend + providers + postgres)
├── .env.example                # copy to .env and fill in secrets
├── .gitignore                  # ignores .env
├── postgres/
│   └── core-db-init.sh         # creates core_db (one-shot)
└── nginx/
    └── api.bayumaulana.my.id.conf   # nginx reverse proxy + HTTPS
```

Backend image: `backend/Dockerfile` (builds both `api` and `migrate`).
Provider images: built from `providers/Dockerfile` (unchanged).
