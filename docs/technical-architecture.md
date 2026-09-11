# ANR Platform - Technical Architecture

## 1. Mục tiêu

ANR Platform là monorepo quản lý toàn bộ hạ tầng, backend, website và
công cụ nội bộ của ANR Studio.

Mục tiêu:

-   Một backend dùng chung cho nhiều ứng dụng.
-   Một hệ thống Auth riêng.
-   Quản lý toàn bộ app từ Admin Dashboard.
-   Deploy bằng Docker.
-   CI/CD qua GitHub Actions.
-   Dễ scale.

------------------------------------------------------------------------

# 2. Repository

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

------------------------------------------------------------------------

# 3. Thành phần

## backend

Go REST API.

Chịu trách nhiệm:

-   Authentication
-   User
-   Device
-   Subscription
-   Analytics
-   App Management
-   API cho tất cả app

------------------------------------------------------------------------

## admin

Dashboard quản trị.

Chức năng:

-   Quản lý người dùng
-   Quản lý ứng dụng
-   Subscription
-   Analytics
-   Nhật ký hệ thống
-   Quản trị hệ thống

------------------------------------------------------------------------

## landing

Website chính:

-   Landing Page
-   Blog
-   Chính sách
-   Hỗ trợ

------------------------------------------------------------------------

## tools

Các CLI và script nội bộ.

------------------------------------------------------------------------

# 4. Infrastructure

## docker

Dockerfile và tài nguyên Docker.

## gateway

Reverse Proxy (Nginx).

## database

PostgreSQL.

## cache

Redis.

## storage

Không gian cho object storage nếu cần (MinIO...).

## monitoring

-   Grafana
-   Prometheus
-   Uptime Kuma

## management

Portainer.

## security

Chứng chỉ, chính sách bảo mật, SSH.

## scripts

Script deploy, backup, restore.

------------------------------------------------------------------------

# 5. Docker Services

-   gateway
-   backend
-   admin
-   landing
-   postgres
-   redis
-   portainer
-   grafana
-   prometheus
-   uptime-kuma

------------------------------------------------------------------------

# 6. Domain

-   anr-studio.com
-   api.anr-studio.com
-   admin.anr-studio.com
-   portainer.anr-studio.com
-   monitor.anr-studio.com

------------------------------------------------------------------------

# 7. Deployment

Mac

↓

Git Push

↓

GitHub

↓

GitHub Actions

↓

SSH

↓

Docker Compose

↓

VPS

------------------------------------------------------------------------

# 8. Bảo mật

-   SSH Key
-   Root login disabled
-   Password login disabled (production)
-   UFW
-   Fail2Ban

------------------------------------------------------------------------

# 9. Công nghệ

Backend: - Go

Frontend: - Next.js

Database: - PostgreSQL

Cache: - Redis

Gateway: - Nginx

Container: - Docker - Docker Compose - Portainer

Monitoring: - Grafana - Prometheus - Uptime Kuma

CI/CD: - GitHub Actions

------------------------------------------------------------------------

# 10. Roadmap

1.  Backend
2.  Docker Compose
3.  PostgreSQL
4.  Redis
5.  Gateway
6.  Landing
7.  Admin
8.  Monitoring
9.  CI/CD
10. Production
