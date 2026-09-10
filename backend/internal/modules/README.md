# App Modules (anr-NNN-slug)

Thư mục này chứa logic nghiệp vụ đặc thù cho từng ứng dụng cụ thể trong hệ sinh thái ANR Studio.

## Quy ước đặt tên module:
- Bắt buộc theo format: `anr-NNN-slug` (ví dụ: `anr-001-wallpaper`, `anr-002-vpn`).
- Mỗi module là một domain package độc lập.
- Các module có thể sử dụng các dịch vụ từ `internal/platform/` (như Auth, User, Device) nhưng không nên phụ thuộc chéo lẫn nhau.
