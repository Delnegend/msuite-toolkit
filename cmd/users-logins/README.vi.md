# Công cụ trích xuất dữ liệu MFA và nhật ký đăng nhập M-Suite

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:
```
./users-logins.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `users_logins.csv`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm user admin đang đăng nhập, click vào kết quả, sao chép `User ID` trong Basic info.
- `admin_portal_address`: thay bằng địa chỉ Admin Portal hiện tại (host:port), ví dụ `10.0.0.1:9443`.

## Ghi chú khi chạy
- Sau khi điền `config.toml`, chạy `./users-logins.exe`. File đầu ra mặc định sẽ xuất hiện cùng thư mục với công cụ.
- Dùng `-c` để chỉ file config khác và `-o` để đặt tên file đầu ra khác.

## Ví dụ đầu ra
```
User ID|Email|MFA|Failed logins
12345|user@example.com|{"type": "totp", "enabled": true}|[{"timestamp": "2023-10-01T10:00:00Z", "ip": "1.2.3.4"}]
```
