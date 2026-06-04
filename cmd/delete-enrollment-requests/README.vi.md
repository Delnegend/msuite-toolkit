# Công cụ xóa enrollment request

## Mô tả
Công cụ này lấy toàn bộ enrollment request cho một organizational unit (OU) từ Admin Portal rồi xóa chúng bằng bulk delete endpoint. Không giống `delete-pending-enrollment-requests`, công cụ này **không** lọc theo trạng thái — nó xóa tất cả request (Pending, Approved, Rejected, v.v.) tìm thấy cho OU đó.

> `organizational_unit_id` là **bắt buộc** trong `config.toml`.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal có trong danh sách ứng dụng.
- Điền `config.toml` ở thư mục gốc (xem hướng dẫn bên dưới).

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./delete-enrollment-requests.exe
```

Dùng `-config` để chỉ file config khác.

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: lấy từ local storage của Admin Portal, khóa `admin_portal_access_token`.
- `admin_user_id`: User ID của admin trong Admin Portal.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `organizational_unit_id`: **bắt buộc** — OU mà enrollment request sẽ bị xóa.
- `dry_run`: đặt `true` để xem trước (không xóa thật, ghi `to-be-deleted-enrollment-requests.csv`); đặt `false` để xóa thật (ghi `deleted-enrollment-requests.csv`).
- `exclude_emails`: danh sách email tùy chọn — các user này sẽ được giữ lại (không bị xóa). Ví dụ: `exclude_emails = ["admin@example.com"]`.
- `worker_count`: cấu hình số luồng cho các request phân trang (mặc định: `100`).

Báo cáo (CSV)
- `deleted-enrollment-requests.csv` / `to-be-deleted-enrollment-requests.csv` — các request đã (hoặc sẽ) bị xóa.
- `excluded-enrollment-requests.csv` — các request được giữ lại vì email nằm trong `exclude_emails` hoặc loại thiết bị không phải desktop (Windows, macOS, Linux).

## Ví dụ
Chạy từ thư mục gốc của dự án:

```
./delete-enrollment-requests.exe -config ./config.toml
```