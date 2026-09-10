# Shared Core Platform

Thư mục này chứa các domain và usecase cốt lõi dùng chung cho toàn bộ ứng dụng trong hệ sinh thái ANR Studio:

- **`auth/`**: First-party Go Authentication (Session, JWT tokens, credential validation, refresh tokens).
- **`user/`**: Quản lý tài khoản người dùng dùng chung.
- **`device/`**: Quản lý thiết bị (device registration, fingerprint, push token binding).
- **`app/`**: Quản lý danh mục app (`anr-NNN-slug`), app versioning và Remote Config.
- **`subscription/`**: Trạng thái thanh toán và subscription chéo app.

## Clean Architecture Boundary
- **Entities & Ports**: Định nghĩa domain model và interface (repository port) tại đây.
- **Dependency Rule**: Tầng này **không** phụ thuộc vào tầng HTTP delivery hay trực tiếp vào các database driver cụ thể.
