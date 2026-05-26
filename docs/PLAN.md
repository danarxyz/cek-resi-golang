# cek-resi-golang — Native Go Rewrite Plan

Status as of 2026-05-26. Branch: `feature/rewrite-native-go`.

Goal: rewrite the project off Goravel into a native, layered, production-grade Go service while preserving the existing domain (6 expedition trackers, REST + gRPC, email notify on status change).

This document is meant to be opened in any session as the source of truth for plan + open decisions. Each phase below is self-contained enough to start work from cold.

---

## 0. Locked decisions (do not re-litigate without explicit cause)

### Architecture
- **Layered**: `handler → service → repository`.
- Layout: `cmd/server/main.go`, `internal/{config,errs,log,metrics,db,domain,repository,expedition,service,handler/http,handler/grpc,queue,mail}/`, `migrations/`, `proto/`, `docs/`.
- No DI framework. Wiring in `main.go`.

### Stack
| Concern | Pick |
|---|---|
| Router | `go-chi/chi/v5` |
| DB driver | `jackc/pgx/v5` + `pgxpool` |
| Queries | `sqlc` (generated, used directly — no thin repo wrapper) |
| Migrations | `golang-migrate/migrate` |
| Validation | `go-playground/validator/v10` |
| Logging | stdlib `log/slog`, JSON handler in prod |
| Config | `kelseyhightower/envconfig` |
| Queue | `hibiken/asynq` (Redis) |
| Scheduler | `asynq.PeriodicTaskManager` |
| Mail | `wneessen/go-mail` |
| Circuit breaker | `sony/gobreaker` |
| Rate limit | `go-chi/httprate` (per-replica Phase 1) |
| Metrics | `prometheus/client_golang` |
| gRPC | `google.golang.org/grpc` |
| Test | stdlib + `stretchr/testify` |

### Error contract (locked)
- Single application error type: `errs.Error{Code, PublicDetail, InternalCtx, Fields, Cause}`.
- `Code` is a string enum, mapped via single `errs.Lookup(Code) Spec` → `(HTTPStatus, GRPCCode, Title)`.
- `Retryable` is **not** in the registry — decided per call site at the queue handler.
- `PublicDetail` exposure rule:
  - 5xx → MUST be empty.
  - 4xx → opt-in, only for whitelist: `{VALIDATION_FAILED, INVALID_INPUT, UNSUPPORTED_CARRIER, MALFORMED_REQUEST}`.
- `Fields` is `[]FieldError{Field, Code, Reason}` (typed, not map).
- Code boundaries:
  - `MALFORMED_REQUEST` — body parse failure, content-type mismatch.
  - `VALIDATION_FAILED` — schema valid, field-level constraint violation. `Fields` populated.
  - `INVALID_INPUT` — schema valid, business-rule violation.

### Boundary mapping
- HTTP error body: RFC 9457 Problem Details (`application/problem+json`). Mandatory members: `type, title, status, code, instance, request_id`. Optional: `detail, fields`.
- gRPC error: `status.Status` with `ErrorInfo` extension carrying `reason=Code`, `request_id`, `retry_after` (when applicable). No raw detail leakage.
- Single source of truth = registry. Boundary-specific behavior (headers, metadata, extensions) stays at each boundary.

### Logging
- Structured JSON via `slog`. Required fields every entry: `time, level, msg, service=cek-resi, version=$GIT_SHA, stream`.
- Stream taxonomy (locked):
  - `request` — HTTP/gRPC lifecycle, sampled (10% for 2xx, 100% non-2xx).
  - `app_error` — 5xx and worker failure paths, 100%.
  - `audit` — auth and sensitive mutations, 100%.
  - `task` — asynq lifecycle. Completed → metric only (no log); failed/archived → 100% log.
- Request correlation: `request_id` (HTTP) or `task_id` (worker) injected into ctx, echoed in response header `X-Request-Id`.
- Log sink routing and retention are **infra concerns**, not locked here.

### Selective 4xx telemetry — classification locked
| Signal | Transport |
|---|---|
| 401/403 repeated same IP | metric counter + audit log per event |
| 429 | metric counter only |
| 422 spike | metric counter + sampled log |
| Noise 4xx (404/400 scanner) | metric counter only |

Thresholds are draft until we have prod baseline.

### Concurrency capacity — `STARTING_DEFAULTS_v0`
Re-derive when: SLO changes, upstream p95 changes >50%, archived rate doubles week-over-week, DB pool saturation >80% sustained.

| Component | Calc | Setting |
|---|---|---|
| HTTP in-flight (target p99 500ms, λ=100/s, W=300ms budget) | 100 × 0.3 = 30 | no cap, DB pool gates it |
| Tracking workers (1000/5min, W=5s) | 3.3 × 5 = 17 | `asynq.Concurrency` queue `tracking` = 20 |
| Mail workers (100/h, W=10s) | 0.3 | `asynq.Concurrency` queue `mail` = 5 |
| pgx pool (HTTP + workers + scheduler + headroom) | 40 + 25 + 1 + headroom | `POOL_MAX_CONNS=60` |

Multi-replica (≥2) → `pgbouncer` transaction-mode required. Documented in README.

### HTTP client (tracker) — mandatory tuning
```
Timeout: 5s
DialContext.Timeout: 2s
TLSHandshakeTimeout: 3s
IdleConnTimeout: 90s
MaxIdleConnsPerHost: 50
```

### External provider (locked)
- Tracking backend: **BinderByte** aggregator API (single upstream covering 24 carriers including all currently supported).
- Base URL: `https://api.binderbyte.com`. Endpoint: `GET /v1/track?api_key&courier&awb`.
- Auth via `BINDERBYTE_API_KEY` env var (required at boot).
- Quota poller hits `GET /v1/checkQuota` periodically (5 min default) → `binderbyte_quota_remaining` metric.
- Single circuit breaker per upstream (not per-carrier — single vendor).
- Response cache mandatory: Redis-backed, TTL by terminal status (DELIVERED → 24h, in-transit → 5–10 min). Cache hits do not count against BB quota.
- `SPX_TOKEN` removed — SPX tracking now routes through BinderByte's `spx` courier code.

### Out of scope Phase 1 cycle (deferred)
- OpenTelemetry tracing (defer to Phase 2 unless multi-hop ambiguity emerges).
- Global cluster rate limit (Phase 2; in-memory per-replica is documented as non-global).
- `dead_tasks` persistence table (decision boundary locked; schema not designed until archived inspection proves insufficient).
- Alert thresholds for archived/error rate (need 1-2 weeks of prod baseline).
- Cek ongkir endpoint (`POST /v1/cost`) and wilayah reference data (`GET /wilayah/*`) — both BinderByte features. Tracked as Phase 8+.

---

## Phase 1 — Foundation (errs, problem details, slog, request_id, metrics skeleton)

**Goal**: ship the cross-cutting infrastructure that every later phase depends on. No business logic.

**Branch**: continue on `feature/rewrite-native-go`. PR title: `feat(phase-1): foundation — errs, problem details, slog, metrics skeleton`.

### Scope (in)
- `cmd/server/main.go` minimal HTTP boot with chi
- `internal/config/` — envconfig with defaults
- `internal/errs/` — `Code` enum, registry, `Error` type, `New/Wrap/CodeOf/Lookup`
- `internal/log/` — slog setup, ctx logger, stream attribute helper
- `internal/handler/http/middleware/` — `RequestID`, `RequestLogger`, `Recoverer`
- `internal/handler/http/problem/` — RFC 9457 writer
- `internal/metrics/` — Prometheus registry, basic counters: `http_requests_total{method,path,status}`, `errors_total{code,stream}`
- `/healthz`, `/readyz`, `/metrics` endpoints
- Sampling for `request` stream

### Scope (out)
- Any DB, queue, expedition, or business handler.

### Deliverables
```
cmd/server/main.go
internal/
  config/config.go
  errs/
    code.go              # Code constants
    registry.go          # registry map + Lookup
    error.go             # Error struct, New, Wrap, CodeOf
    error_test.go
  log/
    log.go               # NewLogger, WithLogger, From
    sampler.go           # request stream sampling
  metrics/
    metrics.go           # registry, counters
  handler/http/
    problem/
      problem.go         # RFC 9457 Problem + writer
      problem_test.go
    middleware/
      requestid.go
      logger.go
      recoverer.go
    health.go
go.mod / go.sum
docs/PLAN.md (this file)
```

### Acceptance / audit checklist
- [ ] `go vet ./...` clean, `go test ./...` passes.
- [ ] `govulncheck ./...` clean.
- [ ] Hitting any unmapped path returns a Problem Details body with `code=NOT_FOUND`, `status=404`, `request_id` present and non-empty, `Content-Type: application/problem+json`.
- [ ] Forcing a panic in a test handler returns 500 Problem with empty `detail` (no stack/internal leak); log shows `stream=app_error level=error stack=...`.
- [ ] `X-Request-Id` from incoming header is honored; absent → ULID generated; response header always present.
- [ ] `/metrics` exposes `http_requests_total` with `status` label populated.
- [ ] Sampling: 2xx logs hit ~10% of count under load test; non-2xx always logged.
- [ ] `errs.CodeOf(error_chain)` walks correctly through `fmt.Errorf("...: %w", ...)`.
- [ ] No `*pgxpool.Pool`, no `asynq`, no business types touched.

### What stays draft this phase
- Code list will grow with later phases. Lock additions in PR review, not in this phase.

---

## Phase 2 — Persistence (pgx, sqlc, migrations, domain, repo)

**Goal**: type-safe DB layer, schema versioned and runnable from app and CLI.

**Depends on**: Phase 1.

### Scope (in)
- `pgxpool.Pool` setup with config-driven pool size
- `migrations/0001_init_resi.{up,down}.sql` — `resi` table (id, tracking_num UNIQUE, expedition, status, details, email, created_at, updated_at)
- `golang-migrate` wired via go-embed for app-startup migration, plus CLI subcommand `cek-resi migrate {up,down,version}`
- `sqlc` config + generated package in `internal/repository/sqlcgen`
- `internal/domain/resi.go` — domain type, mappers to/from sqlc row
- `internal/service/resi/service.go` — skeleton: `Get/List/Add/UpdateStatus/Delete`. Uses `sqlc.Queries` directly (no extra interface wrapper).

### Scope (out)
- HTTP/gRPC handlers (Phase 4/5).
- Asynq, mail (Phase 6).

### Deliverables
```
sqlc.yaml
migrations/
  0001_init_resi.up.sql
  0001_init_resi.down.sql
internal/
  db/
    pool.go                  # pgxpool setup with config
    migrate.go               # embed migrations + run on boot
  domain/
    resi.go
  repository/
    queries/                 # .sql files for sqlc
      resi.sql
    sqlcgen/                 # generated (DO NOT EDIT)
  service/resi/
    service.go
    service_test.go          # uses real postgres via testcontainers OR sqlmock — decide in PR
cmd/server/main.go           # add: db open, migrate, ping, service ctor
```

### Acceptance / audit checklist
- [ ] `sqlc generate` is idempotent in CI.
- [ ] App boots → migrations applied → version recorded in `schema_migrations`.
- [ ] `cek-resi migrate up/down/version` work standalone.
- [ ] Service returns `errs.New(CodeNotFound, ...)` mapping `pgx.ErrNoRows`. Other DB errors wrapped via `errs.Wrap(CodeInternal, err, "ctx")`.
- [ ] Pool size honors `POOL_MAX_CONNS`; `/readyz` checks `pool.Ping(ctx)` with 1s deadline.
- [ ] No business logic in repo. No SQL strings in service. No `sqlcgen` types leak past service boundary.
- [ ] `service_test.go` covers each method, including not-found path → `errs.CodeOf(err) == CodeNotFound`.

### Open questions for this phase
- Test infra: testcontainers vs sqlmock. Recommend testcontainers if CI has Docker, sqlmock otherwise. Decide in PR.

---

## Phase 3 — Expedition: BinderByte client + breaker + cache

**Goal**: single tracking backend (BinderByte) reachable through a `Tracker` interface, guarded by one circuit breaker on the BB upstream, with Redis-backed response cache to conserve quota. 24 carriers supported via courier-code dispatch.

**Depends on**: Phase 1 (errs, slog, metrics).

### Scope (in)
- `internal/expedition/types.go` — `Carrier` enum (24 entries) and `Tracking` domain type (uniform shape: summary + history). `Carrier` is a typed string mapped to BB courier code via a constant table.
- `internal/expedition/tracker.go` — `Tracker` interface with one method: `Track(ctx, Carrier, awb) (*Tracking, error)`.
- `internal/expedition/binderbyte/client.go` — single HTTP client. Uses shared `http.Client` from §0 mandatory tuning. Builds `GET /v1/track`, parses uniform response.
- `internal/expedition/binderbyte/mapper.go` — BB JSON → domain `Tracking`. Status-field parser:
  - `body.status == 200` → success
  - `body.status == 400 && message ~ "Data not found"` → `errs.New(CodeNotFound, "")`
  - `body.status == 401/403` → `errs.New(CodeUnauthenticated, "")` + `InternalCtx="bb api_key invalid or quota exhausted"`
  - HTTP timeout or 5xx → `errs.New(CodeCarrierUnavailable, "")` + `InternalCtx="bb http=<code>"` (retry-able)
  - Otherwise → `errs.Wrap(CodeInternal, raw, "bb unknown response")`
- `internal/expedition/breaker.go` — single `gobreaker` instance for the BB upstream. Settings labeled `BOOT_DEFAULTS_v0`: `MaxRequests=1, Interval=60s, Timeout=30s, ReadyToTrip=ConsecutiveFailures>=5`. Breaker-open → `CodeCarrierUnavailable` + `InternalCtx="bb breaker open"`.
- `internal/expedition/cache.go` — Redis-backed response cache (uses same Redis as asynq):
  - Key: `track:<carrier>:<awb>`
  - TTL by terminal status: `DELIVERED|RETURNED` → 24h; in-transit → 10m; not-found → 1m (negative cache, short).
  - Cache hits never call BB. Cache miss → call BB → store on success only.
  - On cache miss + breaker-open → return last cached entry if any, marked stale (else error).
- `internal/expedition/dispatcher.go` — orders Cache → Breaker → BinderByte client. Composition done in `main.go`, no runtime DI.
- `internal/expedition/quota.go` — periodic poller hitting `GET /v1/checkQuota` every 5 min. Result → metric `binderbyte_quota_remaining`. Owned by Phase 6 scheduler in practice; stub interface here.
- Metrics: `expedition_calls_total{carrier,result}`, `expedition_breaker_state` (gauge), `binderbyte_cache_total{result=hit|miss|stale}`, `binderbyte_quota_remaining`.

### Scope (out)
- Wiring into HTTP/gRPC handlers (Phase 4/5).
- Async polling job that consumes the tracker (Phase 6).
- BB ongkir (`/v1/cost`) and wilayah (`/wilayah/*`) endpoints — Phase 8+.

### Deliverables
```
internal/expedition/
  types.go                   # Carrier enum + BB code table + Tracking domain
  tracker.go                 # Tracker interface
  dispatcher.go              # Cache → Breaker → BB chain ctor
  breaker.go                 # single gobreaker for BB
  cache.go                   # Redis response cache w/ status-aware TTL
  quota.go                   # quota poller stub + metric
  binderbyte/
    client.go                # HTTP client; calls GET /v1/track
    mapper.go                # BB JSON → domain Tracking + error mapping
    client_test.go           # golden tests per carrier (24 fixtures from postman)
    mapper_test.go           # body.status parsing matrix
internal/config/config.go    # add BINDERBYTE_API_KEY, BINDERBYTE_BASE_URL
```

### Acceptance / audit checklist
- [ ] All 24 BB courier codes mapped in `Carrier` enum; round-trip `Carrier ↔ BB code` covered by table-driven test.
- [ ] Golden fixtures (taken from `BinderByte.postman_collection.json` Success samples — file itself gitignored) parse into domain `Tracking` without data loss for at least the 6 current carriers, plus 3 representative new ones (POS, TIKI, Anteraja).
- [ ] `body.status` parsing matrix: 200/400-not-found/401/403/malformed/5xx all produce expected `errs.Code`.
- [ ] Upstream timeout → call returns within `Timeout`, no goroutine leak (verified by leak-check in test).
- [ ] 5 consecutive upstream failures → breaker opens; next call fails fast with `CodeCarrierUnavailable`, `InternalCtx="bb breaker open"`. No 6th BB hit observed in the test transport.
- [ ] Cache hit path issues 0 BB calls; cache miss + success populates cache; not-found populates negative cache with 1m TTL.
- [ ] Breaker-open + stale cache present → returns cached entry with `Tracking.Stale=true` (or sentinel) — degraded mode behavior is explicit, not silent.
- [ ] No `PublicDetail` set on any `CodeCarrierUnavailable` / `CodeInternal` return path.
- [ ] No leak of `BINDERBYTE_API_KEY` into logs, error messages, or metric labels.
- [ ] `SPX_TOKEN` removed from `config.Config`; references in code = 0 (`grep` check).

### What stays draft
- Breaker thresholds `BOOT_DEFAULTS_v0` — tune from `expedition_calls_total` after baseline.
- Cache TTL values are `STARTING_DEFAULTS_v0`; re-derive once we see ratio cache-hit vs quota-burn in prod.
- Quota alert threshold: `DRAFT`. Initial value documented as placeholder.
- Stale-cache return behavior (degrade vs error) — locked direction (return stale, mark stale), but the marker mechanism (`Tracking.Stale bool` vs separate response code) finalized in PR.

---

## Phase 4 — HTTP API (handlers, validator, rate limit)

**Goal**: expose REST endpoints with Problem Details errors, request validation, per-replica rate limit.

**Depends on**: Phase 1, 2, 3.

### Scope (in)
- `internal/handler/http/router.go` — chi router wiring (middleware order: Recoverer → RequestID → RequestLogger → RateLimit → Routes).
- Endpoints (match current Goravel routes; finalize in PR):
  - `GET /v1/track?awb=...&carrier=...` — sync tracker call
  - `GET /v1/resi` — list watched resi
  - `POST /v1/resi` — add to watch (validation: awb required, carrier in enum, email required+email)
  - `PATCH /v1/resi/{awb}/status` — manual status update
  - `DELETE /v1/resi/{awb}`
- `internal/handler/http/track.go`, `resi.go` — handlers (thin: parse → validate → service → respond).
- Validator setup with sentinel translation → `errs.New(CodeValidationFailed, "", fieldErrors)`.
- `httprate` per-IP limit, configurable RPS, label `STARTING_DEFAULTS`. Documented as **per-replica non-global**.

### Scope (out)
- gRPC (Phase 5).
- Async path (Phase 6).

### Deliverables
```
internal/handler/http/
  router.go
  track.go
  track_test.go
  resi.go
  resi_test.go
  validate.go                # validator instance + field-error translator
  ratelimit.go               # httprate config
```

### Acceptance / audit checklist
- [ ] Each endpoint integration-tested with `httptest`, asserting status, `code`, `request_id` presence in error responses.
- [ ] Validation error returns 422, `code=VALIDATION_FAILED`, `fields` array typed.
- [ ] Malformed body returns 400, `code=MALFORMED_REQUEST`.
- [ ] Unknown carrier returns 400, `code=UNSUPPORTED_CARRIER`.
- [ ] Rate limit triggers 429 with `code=RATE_LIMITED`; metric counter increments.
- [ ] No internal carrier detail / SQL error / stack in any response body.

---

## Phase 5 — gRPC service

**Goal**: gRPC server exposing equivalent of HTTP `/v1/track` and CRUD, sharing service layer.

**Depends on**: Phase 1, 2, 3.

### Scope (in)
- `proto/resi.proto` — review existing, regenerate if changed.
- `internal/handler/grpc/server.go` — implements generated `ResiServiceServer`.
- `internal/handler/grpc/interceptor.go` — UnaryServerInterceptor: request_id (from `x-request-id` metadata), logger ctx, recover, `errs → status` translation via `errs.Lookup`.
- `internal/handler/grpc/errmap.go` — `toStatus(ctx, err)` builds `Status` with `ErrorInfo{Reason=Code, Metadata=request_id, retry_after}`.

### Scope (out)
- gRPC-gateway (not currently used).
- TLS termination (handled by reverse proxy or Phase 7).

### Deliverables
```
proto/
  resi.proto
  resi.pb.go               # regenerated
  resi_grpc.pb.go
internal/handler/grpc/
  server.go
  server_test.go
  interceptor.go
  errmap.go
```

### Acceptance / audit checklist
- [ ] `bufbuild/buf` or `protoc` generation reproducible.
- [ ] Returned `Status` carries `ErrorInfo.Reason` matching `errs.Code`.
- [ ] Same service errors → same `Code`, mapped to matching `codes.X` from registry.
- [ ] Request ID echoed in trailing metadata.
- [ ] No `PublicDetail` leak; `InternalCtx` only in logs.

---

## Phase 6 — Async pipeline (asynq + scheduler + mail)

**Goal**: scheduled status polling + email notify on change, with retry and archived inspection.

**Depends on**: Phase 1, 2, 3.

### Scope (in)
- `internal/queue/client.go` — asynq client.
- `internal/queue/server.go` — asynq server with queues `tracking` (concurrency 20) and `mail` (concurrency 5). Middleware: TaskLogger, Recoverer.
- `internal/queue/tasks.go` — task type constants, payload structs, enqueue helpers with per-task `[]asynq.Option` (Queue, MaxRetry, Timeout, Retention, Unique).
- `internal/queue/handlers/check_tracking.go` — idempotent handler:
  1. fetch current row → if not found, SkipRetry
  2. tracker call → on `CodeCarrierUnavailable`, return err (retry); on parse/unsupported, SkipRetry
  3. compare to last known status → no-op if unchanged
  4. UPDATE status
  5. enqueue mail with `Unique(10m)` for dedupe
- `internal/queue/handlers/send_mail.go` — wneessen/go-mail SMTP, idempotent (rely on `Unique` and recipient mailbox dedupe).
- `internal/queue/scheduler.go` — `PeriodicTaskManager` enqueues a "check all tracking" fan-out task every N minutes, which itself enqueues per-resi `check_tracking` tasks. Avoids worker starvation from one mega-task.
- Server `ErrorHandler` → archived event logged as `app_error` with payload preview (no PII in log).

### Scope (out)
- `dead_tasks` table (decision boundary only; not built yet).
- Replay CLI (Phase 7+ if needed).

### Deliverables
```
internal/queue/
  client.go
  server.go
  tasks.go
  scheduler.go
  middleware.go              # TaskLogger
  handlers/
    check_tracking.go
    check_tracking_test.go
    send_mail.go
    send_mail_test.go
internal/mail/
  smtp.go                    # go-mail client wrapper
  smtp_test.go
```

### Acceptance / audit checklist
- [ ] Re-running `check_tracking` for same AWB with unchanged status is a true no-op (no UPDATE issued, no mail enqueued). Verified by test.
- [ ] Carrier 5xx → handler returns retryable error; archived only after MaxRetry exhausted.
- [ ] Malformed payload → SkipRetry on first failure; logged once.
- [ ] Mail `Unique(10m)` dedupe verified — double enqueue inside window enqueues exactly one task.
- [ ] Worker shuts down cleanly on SIGTERM; in-flight tasks finish or return to queue.
- [ ] Metrics: `asynq_*` exposed via `metrics.PrometheusMetricsExporter`.
- [ ] Archived event log entry includes: `task_id, type, queue, retry, last_error, payload_preview` (preview = first 256 bytes, redacted email).

### What stays draft
- `dead_tasks` schema (only build if asynq retention/asynqmon proves insufficient).
- Alert thresholds for archived rate.

---

## Phase 7 — Operationalization (docker, asynqmon, sizing docs, cleanup)

**Goal**: deployable image, ops visibility, decommission Goravel.

**Depends on**: Phases 1–6.

### Scope (in)
- `Dockerfile` — multi-stage; distroless or alpine final; non-root user.
- `docker-compose.yml` — `app`, `postgres`, `redis`, `asynqmon`. Drop Goravel env block. Add minimal env example.
- `docker-compose.prod.yml` — production overrides (no port expose for redis/postgres, restart policies).
- `.env.example` — concrete vars used by `envconfig`.
- `README.md` — replace Goravel docs with: getting started, env vars, sizing notes (`STARTING_DEFAULTS_v0`), multi-replica → pgbouncer requirement, observability endpoints.
- Delete: `app/`, `bootstrap/`, `config/` (Goravel), `console/`, `providers/`, `routes/`, `storage/framework`, `.air.toml`, `artisan`, old `main.go`.
- Keep proto under `proto/`. Keep `migrations/`. Keep `LICENSE`, `README_zh.md` (translate later).

### Scope (out)
- OTel (Phase 2+).
- `dead_tasks` (deferred).

### Deliverables
```
Dockerfile
docker-compose.yml
docker-compose.prod.yml
.env.example
README.md
(deletions of Goravel artifacts)
```

### Acceptance / audit checklist
- [ ] `docker compose up` brings up clean stack; `/healthz` and `/readyz` pass; `/metrics` exposes counters.
- [ ] asynqmon reachable at documented port; queues `tracking` and `mail` visible.
- [ ] `curl /v1/track?awb=...&carrier=spx` returns valid Problem on bad input, valid JSON on happy path.
- [ ] gRPC reflection (dev only) responds; `grpcurl` happy-path works.
- [ ] Image size reasonable (<50MB for distroless build).
- [ ] No Goravel imports remain (`grep -r "goravel/framework"` returns nothing).

---

## Cross-cutting acceptance (final, before declaring done)

- [ ] `govulncheck ./...` clean.
- [ ] `go vet ./...` clean.
- [ ] Unit + integration tests pass in CI.
- [ ] No `PublicDetail` set on any 5xx code (grep + lint rule documented).
- [ ] No `Cause` or `InternalCtx` appears in any response body across HTTP/gRPC.
- [ ] Status table in this file updated: `STARTING_DEFAULTS_v0`/`BOOT_DEFAULTS_v0` markers untouched until baseline metrics revisit.
- [ ] README sizing section reflects current defaults and re-derivation triggers.

---

## How to use this plan

- One phase = one PR. Don't merge a phase until its audit checklist is satisfied.
- Re-open this file at the start of each session. The "Locked decisions" section is binding; "Out of scope Phase 1" and "What stays draft" mark non-binding intent.
- When a draft item gets locked in a later phase, edit this file in the same PR.
- If a locked decision needs to change, write the justification in the PR description and update §0 in the same commit.
