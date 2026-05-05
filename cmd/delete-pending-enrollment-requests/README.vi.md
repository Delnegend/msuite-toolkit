# Công cụ xóa các enrollment request đang chờ

## Mô tả
Công cụ này lấy toàn bộ enrollment request từ Admin Portal, giữ lại các request có trạng thái `pending`, rồi xóa chúng bằng bulk delete endpoint.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
 - Điền `config.toml` theo hướng dẫn bên dưới.

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./delete-pending-enrollment-requests.exe
```

Dùng `-config` để chỉ file config khác.

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: lấy từ local storage của Admin Portal, khóa `admin_portal_access_token`.
- `admin_user_id`: User ID của admin đang đăng nhập trong Admin Portal.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `worker_count`: cấu hình số luồng cho các request phân trang (mặc định: `100`).

## Ví dụ
Hãy dùng `config.toml` trong thư mục này. Để chạy cục bộ:

```
./delete-pending-enrollment-requests.exe -config ./config.toml
```