# ANR Platform — Disaster Recovery & Migration Guide

Tài liệu này ghi lại toàn bộ quy trình sao lưu và khôi phục hệ thống ANR Platform khi ngừng sử dụng VPS hiện tại và chuyển sang VPS mới.

---

## 1. Dữ liệu đã được sao lưu an toàn về máy cá nhân

Toàn bộ dữ liệu của hệ thống đã được tự động tải về máy tại:

1. **Database PostgreSQL** (`users`, `devices`, `applications`):
   - File: `infrastructure/storage/backups/anr_platform_prod_20260916_200401.sql.gz`
   - Đã kiểm tra tính toàn vẹn (integrity check passed).
2. **Cấu hình môi trường Production**:
   - File: `.env.production.backup`
   - Chứa toàn bộ JWT Secret, mật khẩu PostgreSQL, cấu hình cổng và môi trường.
3. **Toàn bộ Source Code & Cấu hình Hạ tầng**:
   - Nằm 100% trên GitHub repository: `github.com/anr-2609/anr-platform` (nhánh `main`).

*(Cả 2 file backup trên đều đã được cấu hình trong `.gitignore`, đảm bảo không bị lộ bí mật lên GitHub).*

---

## 2. Các bước xử lý khi hủy VPS hiện tại

Khi VPS hết hạn hoặc chủ động hủy:
- Bạn chỉ việc tắt hoặc xóa instance trên nhà cung cấp (Contabo / Hetzner...).
- **Lưu ý về Tên miền**: Bạn không cần làm gì với tên miền `anr-studio.com`, chỉ cần duy trì gia hạn tên miền trên registrar (Cloudflare / Namecheap / GoDaddy).

---

## 3. Quy trình khôi phục 100% trên VPS mới (Chỉ mất ~10 phút)

Khi thuê lại VPS mới (Khuyến nghị: Ubuntu 22.04 LTS hoặc 24.04 LTS, RAM từ 2GB trở lên):

### Bước 1: Trỏ lại DNS tên miền
Đăng nhập vào trang quản lý DNS (Cloudflare / Namecheap) và cập nhật địa chỉ IPv4 mới cho các bản ghi sau:
- `@` (`anr-studio.com`)
- `api.anr-studio.com`
- `admin.anr-studio.com`
- `status.anr-studio.com`
- `grafana.anr-studio.com`

### Bước 2: Cài đặt Docker & Docker Compose trên VPS mới
SSH vào VPS mới và chạy:
```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

### Bước 3: Kéo mã nguồn từ GitHub
```bash
git clone https://github.com/anr-2609/anr-platform.git /opt/anr-platform
cd /opt/anr-platform
```

### Bước 4: Đưa file cấu hình `.env` lên VPS
Từ máy cá nhân, tải file cấu hình lên VPS mới:
```bash
scp .env.production.backup <user>@<IP_VPS_MOI>:/opt/anr-platform/.env
```

### Bước 5: Khởi động các dịch vụ
Trên VPS mới, chạy lệnh:
```bash
cd /opt/anr-platform
docker compose up -d
```
Hệ thống sẽ tự động tải các image cần thiết (PostgreSQL, Redis, Nginx, Go Backend, Landing, Admin, Prometheus, Grafana, Uptime Kuma) và khởi chạy container.

### Bước 6: Khôi phục dữ liệu Database
Tải file backup từ máy cá nhân lên VPS mới:
```bash
scp infrastructure/storage/backups/anr_platform_prod_20260916_200401.sql.gz <user>@<IP_VPS_MOI>:/opt/anr-platform/backup.sql.gz
```

Trên VPS mới, chạy lệnh khôi phục:
```bash
cd /opt/anr-platform
./infrastructure/scripts/restore-db.sh /opt/anr-platform/backup.sql.gz
```

### Bước 7: Cấp chứng chỉ SSL Let's Encrypt
```bash
cd /opt/anr-platform
./infrastructure/scripts/init-ssl.sh
```

---

## 4. Kiểm tra sau khi khôi phục

Truy cập kiểm tra lại các địa chỉ:
- Landing Page: `https://anr-studio.com`
- Status Page: `https://anr-studio.com/status`
- Admin Console: `https://admin.anr-studio.com/login`
- Backend API: `https://api.anr-studio.com/livez`
