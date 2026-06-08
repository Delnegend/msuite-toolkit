# Công cụ thêm người dùng vào nhóm

## Mô tả
Công cụ này phân giải danh sách email người dùng thành user ID rồi thêm những
người dùng đó vào một nhóm đích thông qua bulk add endpoint. Người dùng được
thêm theo **lô 10**, gửi đồng thời, kèm thanh tiến trình.

> Cả `group_id` và `emails` đều **bắt buộc** trong `config.toml`.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal có trong danh sách ứng dụng.
- Điền `config.toml` ở thư mục gốc (xem hướng dẫn bên dưới).

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./add-users-to-group.exe
```

Dùng `-config` để chỉ file config khác.

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: lấy từ local storage của Admin Portal, khóa `admin_portal_access_token`.
- `admin_user_id`: User ID của admin trong Admin Portal.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `group_id`: **bắt buộc** — nhóm mà người dùng sẽ được thêm vào.
- `emails`: **bắt buộc** — danh sách email người dùng cần thêm. Mỗi email được phân giải thành user ID trước khi thêm. Ví dụ: `emails = ["a@example.com", "b@example.com"]`.
- `dry_run`: đặt `true` để xem trước (phân giải email và ghi `to-be-added-users.csv`, không thêm thật); đặt `false` để thêm thật (ghi `added-users.csv`).
- `worker_count`: cấu hình số luồng cho việc phân giải email và các lô thêm (mặc định: `100`).

Báo cáo (CSV)
- `added-users.csv` / `to-be-added-users.csv` — người dùng đã (hoặc sẽ) được thêm, kèm user ID đã phân giải.
- `unresolved-emails.csv` — các email không khớp được với người dùng nào (chỉ ghi khi có email phân giải thất bại).

## Ví dụ
Chạy từ thư mục gốc của dự án:

```
./add-users-to-group.exe -config ./config.toml
```
