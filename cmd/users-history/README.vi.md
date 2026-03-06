# Công cụ trích xuất Lịch sử Người dùng M-Suite

## Mô tả
Công cụ này trích xuất lịch sử toàn diện cho từng người dùng, bao gồm:
- Trạng thái và cấu hình MFA.
- Danh sách tất cả các thiết bị liên kết với người dùng.
- Địa chỉ IP cuối cùng được biết cho mỗi thiết bị.
- Lịch sử các lần đăng nhập thất bại.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:
```
./users-history.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `users_history.csv`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm user admin đang đăng nhập, click vào kết quả, sao chép `User ID` trong Basic info.
- `admin_portal_address`: thay bằng địa chỉ Admin Portal hiện tại (host:port), ví dụ `10.0.0.1:9443`.

## Ghi chú khi chạy
- Sau khi điền `config.toml`, chạy `./users-mfa-devices-ip.exe`. File đầu ra mặc định sẽ xuất hiện cùng thư mục với công cụ.
- File đầu ra CSV phân tách bằng dấu `|` và bao gồm các cột sau: `UserID`, `UserEmail`, `MFA`, `Device`.
- Các cột `MFA` và `Device` chứa dữ liệu dạng JSON cho các phương thức xác thực đa yếu tố và thông tin thiết bị (bao gồm cả địa chỉ IP cuối cùng được ghi nhận).
- Dùng `-c` để chỉ file config khác và `-o` để đặt tên file đầu ra khác.

## Ví dụ đầu ra
```
UserID|UserEmail|MFA|Device
12345|user@example.com|{"type": "totp", "enabled": true}|[{"id": "dev-1", "name": "iPhone 13", "ip": "1.2.3.4"}]
```
