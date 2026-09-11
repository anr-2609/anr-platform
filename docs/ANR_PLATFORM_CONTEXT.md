# ANR Platform --- Master Context

> **Single Source of Truth (SSOT)** cho ANR Platform.
>
> Developer, AI assistant và coding agent phải đọc file này trước khi
> thay đổi kiến trúc, hạ tầng, deployment hoặc cấu trúc repository.
>
> **Không tự ý redesign, rename, move hoặc thêm top-level architecture
> nếu chưa được yêu cầu.**
>
> Last updated: 2026-09-07

------------------------------------------------------------------------

## 1. Tổng quan

**ANR Platform** là nền tảng kỹ thuật dùng chung cho hệ sinh thái ứng
dụng của ANR Studio.

Mục tiêu:

-   Một backend dùng chung cho nhiều app.
-   Auth first-party tự xây bằng Go.
-   Một Admin Dashboard quản lý toàn bộ app.
-   Common API + API/module riêng theo từng app.
-   PostgreSQL và Redis dùng chung ở tầng platform.
-   Production chạy bằng Docker.
-   Source code và infrastructure quản lý bằng Git.
-   Sẵn sàng CI/CD.
-   Có thể scale từ 1 VPS lên nhiều node/service mà không phải đập lại
    repository.

------------------------------------------------------------------------

# 2. Quyết định kiến trúc đã CHỐT

1.  ANR Platform là **monorepo**.
2.  Backend production dùng **Go**.
3.  Auth tự xây bằng Go, không dùng Firebase Auth/Auth0/Clerk/Supabase
    Auth làm auth lõi.
4.  Có common platform services dùng chung cho các app.
5.  Mỗi app có thể có module/API/config riêng.
6.  App ID theo format `anr-NNN-slug`, ví dụ `anr-001-wallpaper`.
7.  Admin Dashboard quản lý **toàn bộ portfolio app**.
8.  Landing, Admin, Backend và Infrastructure nằm trong cùng repository.
9.  Production containerized bằng Docker.
10. Docker Compose là orchestration hiện tại.
11. Portainer là Docker management UI, **không phải** Admin Dashboard.
12. PostgreSQL là relational database chính.
13. Redis là cache/ephemeral service.
14. Nginx là gateway/reverse proxy dự kiến.
15. GitHub là source repository và nguồn CI/CD.
16. Infrastructure configuration phải được version-control khi hợp lý.
17. Không đổi cấu trúc chỉ vì có một layout khác cũng hợp lệ.

------------------------------------------------------------------------

# 3. Repository Structure --- LOCKED

``` text
anr-platform/
├── .github/
│   └── workflows/
├── admin/
├── backend/
├── docs/
├── infrastructure/
│   ├── docker/
│   ├── gateway/
│   ├── database/
│   ├── cache/
│   ├── storage/
│   ├── monitoring/
│   ├── management/
│   ├── security/
│   └── scripts/
├── landing/
├── tools/
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Makefile
└── README.md
```

Không tự ý:

-   thêm wrapper `apps/`;
-   di chuyển `admin`, `backend`, `landing`;
-   tách mỗi component thành repo riêng;
-   đổi tên `infrastructure`;
-   tạo một top-level layout mới;
-   di chuyển root `docker-compose.yml` nếu chưa có quyết định mới.

Có thể tạo folder con bên trong các component khi implementation yêu
cầu.

------------------------------------------------------------------------

# 4. Trách nhiệm từng component

## `backend/`

Go backend cho toàn bộ ANR ecosystem.

Phạm vi:

-   Authentication
-   Authorization
-   Users
-   Devices
-   Apps
-   App metadata & versioning (Feature flags & remote config dùng trực tiếp Firebase Remote Config trên mobile app)
-   Subscription state
-   API access
-   Admin operations
-   Analytics ingestion/aggregation khi phù hợp
-   Shared platform services
-   App-specific modules
-   Internal operational APIs

Backend phải giữ được ranh giới rõ giữa common platform và business
logic riêng từng app.

## `admin/`

Dashboard quản trị toàn bộ app ANR.

Dự kiến quản lý:

-   Apps
-   Users
-   Devices
-   App versions
-   Subscriptions
-   API usage
-   Analytics/revenue integrations
-   Operational status
-   Logs/audit
-   App-specific management screens

Admin gọi trực tiếp Go backend.

**Không cần `admin-api` riêng ở giai đoạn kiến trúc hiện tại.** Admin
endpoints nằm trong Go backend và được bảo vệ bằng auth/authorization
phù hợp.

## `landing/`

Public website:

-   Main landing
-   App portfolio
-   Support
-   Privacy Policy
-   Terms
-   Contact
-   Blog/content nếu cần

## `tools/`

CLI, import/export, maintenance utilities và tooling nội bộ.

## `docs/`

Tài liệu kỹ thuật và vận hành.

File này:

``` text
docs/ANR_PLATFORM_CONTEXT.md
```

là context/architecture source of truth chính.

------------------------------------------------------------------------

# 5. Infrastructure

``` text
infrastructure/
├── docker/
├── gateway/
├── database/
├── cache/
├── storage/
├── monitoring/
├── management/
├── security/
└── scripts/
```

-   `docker/`: Docker assets/config dùng chung.
-   `gateway/`: Nginx gateway/reverse proxy.
-   `database/`: PostgreSQL config, init, backup-related infrastructure.
-   `cache/`: Redis.
-   `storage/`: object/file storage nếu sau này có nhu cầu thực tế.
-   `monitoring/`: Prometheus, Grafana, Uptime Kuma.
-   `management/`: Portainer.
-   `security/`: security configuration/documentation, không chứa
    secret/private key.
-   `scripts/`: deploy, backup, restore, maintenance.

------------------------------------------------------------------------

# 6. App Module Convention

Stable identifier:

``` text
anr-NNN-slug
```

Ví dụ:

``` text
anr-001-wallpaper
```

Quy tắc:

-   `anr`: namespace.
-   `NNN`: số thứ tự zero-padded.
-   `slug`: lowercase, ngắn gọn.
-   ID giữ ổn định kể cả public app name thay đổi.
-   Không tái sử dụng ID cho app khác.

Mental model:

``` text
ANR Platform
├── Common Platform
│   ├── Auth
│   ├── Users
│   ├── Devices
│   ├── Config
│   ├── Subscriptions
│   └── Analytics
└── App Modules
    ├── anr-001-wallpaper
    ├── anr-002-...
    └── anr-003-...
```

------------------------------------------------------------------------

# 7. Authentication

Auth là **first-party Go authentication**.

Không thay core auth bằng Firebase Authentication hoặc hosted identity
provider khác nếu chưa có yêu cầu đổi kiến trúc.

Concept:

``` text
Mobile App / Admin
        │
        ▼
     Go API
        │
        ▼
 Authentication
        ├── Identity
        ├── Token / Session
        ├── Device
        ├── Authorization
        └── App Access
```

Token/session model, refresh strategy, admin auth, device binding và
credential model phải được thiết kế trước khi production hóa.

------------------------------------------------------------------------

# 8. API Architecture

Target:

``` text
api.anr-studio.com
```

Concept:

``` text
api.anr-studio.com
        │
        ▼
      Nginx
        │
        ▼
    Go Backend
        ├── Common API
        ├── Admin API
        └── App APIs
            ├── anr-001-wallpaper
            ├── anr-002-...
            └── ...
```

Common auth/authorization phải được enforce nhất quán.

------------------------------------------------------------------------

# 9. Data Layer

## PostgreSQL

Primary durable relational database.

Dữ liệu có thể gồm:

-   users
-   devices
-   applications
-   app versions
-   app config
-   auth/session data phù hợp
-   subscriptions
-   feature config
-   audit records
-   operational metadata

**Schema chưa chốt.** Không coi tên table minh họa là schema locked.

## Redis

Dùng cho:

-   cache
-   rate limiting
-   short-lived state
-   distributed locks
-   temporary tokens
-   queue khi phù hợp

Không dùng Redis làm durable source duy nhất cho business data quan
trọng.

------------------------------------------------------------------------

# 10. Docker Architecture

Target services:

``` text
Docker Host (VPS)
├── gateway
├── backend
├── admin
├── landing
├── postgres
├── redis
├── portainer
├── prometheus
├── grafana
└── uptime-kuma
```

Root `docker-compose.yml` là Compose entry point chính trừ khi có quyết
định kiến trúc mới.

------------------------------------------------------------------------

# 11. Portainer

Portainer CE đã cài trên VPS.

Dùng để:

-   xem containers;
-   logs;
-   start/stop/restart;
-   images;
-   networks;
-   volumes;
-   container resource usage;
-   operational stack management.

Portainer **không phải**:

-   ANR Admin;
-   backend;
-   database;
-   source of truth của infrastructure;
-   replacement cho monitoring stack.

Infrastructure definitions phải nằm trong Git.

Current direct HTTPS port:

``` text
9443
```

Planned domain:

``` text
portainer.anr-studio.com
```

Public exposure của management UI cần được harden trước production.

------------------------------------------------------------------------

# 12. Domain Plan

``` text
anr-studio.com
        └── Landing

admin.anr-studio.com
        └── Admin Dashboard

api.anr-studio.com
        └── Go Backend

portainer.anr-studio.com
        └── Portainer

monitor.anr-studio.com
        └── Grafana / Monitoring
```

Nginx sẽ reverse proxy. Production-facing traffic phải dùng HTTPS.

------------------------------------------------------------------------

# 13. Production VPS

Current deployment: **1 Contabo VPS tại Singapore**.

``` text
OS: Ubuntu 24.04 LTS
IPv4: 185.227.135.234
Primary Linux user: longtq
Deployment root: /opt/anr-platform
```

Hiện tại cố ý dùng 1 VPS. Không tự đưa Kubernetes/distributed
orchestration vào khi chưa có nhu cầu scale thực tế.

------------------------------------------------------------------------

# 14. SSH State

Daily administration:

``` text
User: longtq
```

Root SSH đã disabled.

Model:

``` text
Work Mac ── SSH key ──┐
                      ├──> longtq@VPS
Home Mac ── SSH key ──┘
```

Mỗi máy có private key riêng. VPS lưu public keys tại:

``` text
/home/longtq/.ssh/authorized_keys
```

Alias:

``` bash
ssh contabo
```

Đã thực hiện:

-   xóa root `authorized_keys`;
-   `PermitRootLogin no`;
-   public-key auth hoạt động;
-   Work Mac SSH key hoạt động;
-   Home Mac SSH key hoạt động.

Password authentication đã được **bật lại tạm thời** để bootstrap máy
nhà.

### Security TODO

Đưa production SSH về:

``` text
PermitRootLogin no
PubkeyAuthentication yes
PasswordAuthentication no
```

------------------------------------------------------------------------

# 15. VPS Software State

Đã hoàn thành:

-   Ubuntu 24.04 LTS
-   user `longtq`
-   sudo
-   SSH keys
-   Docker Engine
-   Docker Compose
-   Portainer CE
-   `/opt/anr-platform`

Observed setup versions:

``` text
Docker 29.7.2
Docker Compose v5.5.0
Portainer CE 2.39.6
```

Đây chỉ là snapshot version, không phải version lock.

------------------------------------------------------------------------

# 16. Git / GitHub

Repository:

``` text
anr-platform
```

Primary branch:

``` text
main
```

Repo đã:

-   init Git;
-   tạo base monorepo;
-   commit;
-   push GitHub.

Git là source of truth cho code, infra config, Docker config, workflows
và docs.

------------------------------------------------------------------------

# 17. Environment & Secrets

Repo có:

``` text
.env.example
```

Không commit:

``` text
.env
SSH private keys
DB passwords
JWT signing secrets
API secrets
service account credentials
private TLS keys
Portainer credentials
```

`.env.example` chỉ chứa tên biến và placeholder an toàn.

------------------------------------------------------------------------

# 18. CI/CD Target

``` text
Mac
 │ git push
 ▼
GitHub
 │
 ▼
GitHub Actions
 ├── Test
 ├── Build
 ├── Docker image build
 └── Deploy
       │
       ▼
      VPS
       │
       ▼
 Docker Compose
```

CI/CD **chưa hoàn thành**.

Mục tiêu cuối là giảm tối đa manual SSH deployment.

------------------------------------------------------------------------

# 19. Monitoring Target

Planned:

-   Prometheus: metrics.
-   Grafana: visualization/dashboard.
-   Uptime Kuma: availability checks.
-   Portainer: container operations.

Monitoring stack chưa được đánh dấu completed.

------------------------------------------------------------------------

# 20. Scaling Strategy

Hiện tại:

``` text
Single VPS
├── Gateway
├── Backend
├── Web
├── Database
├── Cache
└── Operations
```

Sau này có thể tách:

``` text
Gateway / Load Balancer
        │
        ├── Backend Node 1
        ├── Backend Node 2
        └── Backend Node N

Dedicated/Managed PostgreSQL
Redis
Object Storage
Monitoring
```

Scale deployment không đồng nghĩa phải redesign monorepo.

------------------------------------------------------------------------

# 21. Development vs Production

## Mac

-   code;
-   test;
-   commit;
-   push;
-   local Docker khi cần.

## VPS

-   production containers;
-   production data;
-   public traffic;
-   monitoring;
-   backup;
-   operations.

Không dùng VPS làm coding workstation chính.

------------------------------------------------------------------------

# 22. Backup Requirements

Trước khi có production data quan trọng cần:

-   scheduled PostgreSQL backups;
-   off-server copy;
-   retention;
-   restore procedure;
-   restore test;
-   Docker persistent volume awareness.

Backup chưa restore-test thì chưa coi là hoàn chỉnh.

Status: **Pending**.

------------------------------------------------------------------------

# 23. Security Baseline

-   Non-root SSH administration.
-   Root SSH disabled.
-   SSH key riêng từng device.
-   Password SSH phải tắt sau bootstrap.
-   Firewall.
-   Minimal public ports.
-   TLS.
-   Rate limiting.
-   Secrets outside Git.
-   PostgreSQL không public.
-   Redis không public.
-   Harden Portainer exposure.
-   Security updates.
-   Backups.
-   Audit log cho sensitive Admin actions.

------------------------------------------------------------------------

# 24. Public Port Philosophy

Application traffic:

``` text
80/tcp
443/tcp
```

Administration:

``` text
22/tcp
```

Portainer `9443` hiện dùng trong bootstrap nhưng không mặc định coi việc
public trực tiếp port management là trạng thái production cuối.

PostgreSQL/Redis phải ưu tiên Docker/internal network.

------------------------------------------------------------------------

# 25. Backend Design Direction

Milestone lớn tiếp theo: Go backend.

Yêu cầu:

-   clean boundaries;
-   domain/module separation;
-   shared platform services;
-   app modules;
-   explicit dependencies;
-   central config;
-   structured logging;
-   DB migrations;
-   API versioning;
-   auth middleware;
-   admin authorization;
-   health/readiness endpoints;
-   graceful shutdown;
-   testability.

Clean Architecture không đồng nghĩa tạo tối đa folder. Mục tiêu là
boundary rõ và maintainable.

### Final Go Backend Package Layout (LOCKED)

``` text
backend/
├── cmd/
│   ├── server/                 # Composition root (main.go)
│   └── migrate/                # Migration runner CLI
├── internal/
│   ├── config/                 # 12-factor configuration (config.go)
│   ├── database/               # PostgreSQL pgxpool connection pool (postgres.go)
│   ├── cache/                  # Redis client wrapper (redis.go)
│   ├── platform/               # Shared Core Platform domains
│   │   ├── auth/               # First-party Auth domain & ports
│   │   ├── user/               # User domain & ports
│   │   └── device/             # Device registration & binding
│   ├── modules/                # App-specific modules (anr-NNN-slug)
│   └── transport/
│       └── http/               # HTTP Delivery layer
│           ├── router.go       # Chi router & route registration
│           ├── middleware/     # Slog, RequestID, CORS, Recovery
│           ├── handler/        # HTTP handlers (/livez, /readyz, /healthz)
│           └── response/       # Standard JSON envelopes & error mapping
├── migrations/                 # SQL migration files
├── Dockerfile                  # Multi-stage production build
├── Makefile                    # Build & run scripts
└── go.mod
```

------------------------------------------------------------------------

# 26. Admin Direction

Concept:

``` text
Admin
├── Overview
├── Apps
│   ├── anr-001-wallpaper
│   ├── anr-002-...
│   └── ...
├── Users
├── Devices
├── Subscriptions
├── Analytics
├── Operations
└── Settings
```

Technology direction:

``` text
Next.js
```

Chưa implemented.

------------------------------------------------------------------------

# 27. Landing Direction

Technology direction:

``` text
Next.js
```

Landing là public site, operationally separate với authenticated Admin
dù cùng monorepo.

------------------------------------------------------------------------

# 28. Current Status

## Completed

-   [x] Contabo VPS
-   [x] Ubuntu 24.04 LTS
-   [x] User `longtq`
-   [x] sudo
-   [x] Root SSH disabled
-   [x] Work Mac SSH key
-   [x] Home Mac SSH key
-   [x] Docker
-   [x] Docker Compose
-   [x] `/opt/anr-platform`
-   [x] Portainer CE
-   [x] Portainer admin initialized
-   [x] Monorepo structure
-   [x] Git repository
-   [x] GitHub repository/push
-   [x] Technical docs started
-   [x] Master Context
-   [x] Go backend bootstrap (Chi router, Slog, Graceful shutdown, Health endpoints)
-   [x] Final Go package architecture (Clean Architecture + Modular Monolith)
-   [x] Device Guest Authentication & JWT Session (First-party Auth)

## Cleanup

-   [ ] Disable SSH password authentication again

## Pending

-   [ ] First-party Auth
-   [ ] PostgreSQL
-   [ ] Database schema
-   [ ] Migrations
-   [ ] Redis
-   [ ] Nginx
-   [ ] Domain routing
-   [ ] TLS automation
-   [ ] Landing
-   [ ] Admin
-   [ ] Monitoring
-   [ ] CI/CD
-   [ ] Automated backups
-   [ ] Restore test
-   [ ] Firewall/security hardening review
-   [ ] Portainer exposure hardening

------------------------------------------------------------------------

# 29. Immediate Next Milestone

**Triển khai First-party Go Auth và Core Platform Data Layer.**

Tiếp theo:
1. Hoàn thiện entity & repository interface cho `platform/auth`, `platform/user`, `platform/app`.
2. Chạy migration schema khởi đầu lên PostgreSQL.
3. Token generation (Access Token + Refresh Token), password hashing (bcrypt/argon2id).
4. Auth HTTP middleware cho các route cần bảo vệ.
5. Setup Nginx reverse proxy hoặc Docker network kết nối services.

------------------------------------------------------------------------

# 30. Decision Log

  -----------------------------------------------------------------------
  Decision                Choice                  Status
  ----------------------- ----------------------- -----------------------
  Repository              Monorepo `anr-platform` LOCKED

  Backend                 Go                      LOCKED

  Backend Architecture    Clean Architecture +    LOCKED
                          Modular Monolith        

  Dependency Injection    Manual Constructor      LOCKED
                          Injection               

  HTTP Router             `go-chi/chi/v5`         LOCKED

  Database Driver         `pgx/v5` (pgxpool)      LOCKED

  Authentication          First-party Go Auth     LOCKED
                          (Guest Device Session)  

  Remote Config           Firebase Remote Config  LOCKED
                          trực tiếp trên mobile,  
                          không làm trên Go BE    

  App ID                  `anr-NNN-slug`          LOCKED

  Current deployment      Docker + Docker Compose LOCKED for current
                                                  stage

  Container management    Portainer CE            ACTIVE

  Database                PostgreSQL              LOCKED

  Cache                   Redis                   LOCKED

  Gateway                 Nginx                   LOCKED

  Admin API               Same Go backend,        LOCKED unless justified
                          protected admin         
                          endpoints               
  -----------------------------------------------------------------------

------------------------------------------------------------------------

# 31. Rules for AI / Coding Agents

1.  Đọc file này trước.
2.  Locked decisions là constraints.
3.  Không đề xuất repository layout khác trong task không liên quan.
4.  Không redesign lại quyết định đã chốt.
5.  Nếu thật sự cần đổi, phải nêu lý do cụ thể và impact trước.
6.  Giữ naming convention.
7.  App IDs dùng `anr-NNN-slug`.
8.  Không thay first-party Go Auth bằng hosted Auth.
9.  Không public PostgreSQL/Redis mặc định.
10. Không commit secrets.
11. Infrastructure config ưu tiên nằm trong Git.
12. Tách local development và production concerns.
13. Khi đưa command vận hành, luôn ghi rõ **Mac** hay **VPS**.
14. Với thay đổi server có rủi ro, kiểm tra current state trước khi
    mutate.
15. Không xóa SSH access đang hoạt động trước khi replacement được test.
16. Sau milestone quan trọng phải cập nhật file này.
17. `Current Status` phải phản ánh đúng thực tế.
18. Không thêm infrastructure chỉ vì "có thể cần".
19. Khi API đã production phải cân nhắc backward compatibility.
20. Nếu không chắc một thứ đã chốt chưa, kiểm tra file này trước thay vì
    tự tạo convention mới.

------------------------------------------------------------------------

# 32. Working Convention

Khi mở chat/session mới:

``` text
Read docs/ANR_PLATFORM_CONTEXT.md first.
Treat it as the source of truth.
Continue from Current Status and do not redesign locked architecture unless I explicitly ask.
```

Sau milestone:

1.  Update `Current Status`.
2.  Update `Decision Log` nếu có decision mới.
3.  Update section kiến trúc liên quan.
4.  Commit docs cùng code/infra change.

Ví dụ:

``` bash
git add docs/ANR_PLATFORM_CONTEXT.md
git commit -m "docs: update ANR Platform context"
git push
```

------------------------------------------------------------------------

# 33. Platform Foundation --- Definition of Done

Foundation hoàn chỉnh khi có:

-   Stable Go backend skeleton
-   First-party Auth
-   PostgreSQL + migrations
-   Redis
-   Nginx gateway
-   HTTPS
-   Landing
-   Admin Dashboard
-   Ít nhất một `anr-NNN-*` app tích hợp
-   CI/CD
-   Monitoring
-   Automated backups
-   Restore test
-   Production security review
-   Documentation synchronized với hệ thống thật

------------------------------------------------------------------------

# 34. Core Mental Model

``` text
                         ANR PLATFORM

                             Internet
                                │
                                ▼
                         Nginx Gateway
                                │
              ┌─────────────────┼─────────────────┐
              │                 │                 │
              ▼                 ▼                 ▼
           Landing            Admin            Go API
                                                   │
                                   ┌───────────────┼───────────────┐
                                   │               │               │
                                   ▼               ▼               ▼
                              Common Core      App Modules      Integrations
                                   │               │
                                   │        ┌──────┴───────────┐
                                   │        │                  │
                                   │        ▼                  ▼
                                   │ anr-001-wallpaper    anr-002-...
                                   │
                          ┌────────┴────────┐
                          │                 │
                          ▼                 ▼
                     PostgreSQL           Redis


                      OPERATIONS / INFRA

                 Docker + Docker Compose
                          │
            ┌─────────────┼─────────────┐
            │             │             │
            ▼             ▼             ▼
        Portainer      Monitoring     Backups
                      Prometheus
                       Grafana
                     Uptime Kuma
```

**Mục tiêu:** một platform thống nhất phục vụ nhiều app mà không biến
mỗi app mới thành một dự án infrastructure mới.
