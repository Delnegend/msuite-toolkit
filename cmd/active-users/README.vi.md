# M-Suite Active Users CLI

## Mô tả
Công cụ này lấy tất cả người dùng từ Admin Portal và ghi ra một file JSON chứa những người dùng "active" theo cấu hình.

Lọc "active":
- `last_login_threshold_in_month`: loại bỏ người dùng không đăng nhập trong X tháng gần đây (0 = vô hiệu hóa)
- `organizational_unit_id`: tùy chọn giới hạn kết quả theo một OU cụ thể (để trống để bỏ qua)

## Các bước nhanh
- `cmd/active-users/config.tom` đã có trong repo và sẽ được sử dụng mặc định.
- Chạy công cụ từ thư mục này:
```
./active-users.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file JSON đầu ra (mặc định: `active_users.json`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: lấy từ local storage của Admin Portal (khóa `admin_portal_access_token`).
- `admin_user_id`: User ID của admin đang đăng nhập trong Admin Portal.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `last_login_threshold_in_month`: số tháng để tính active (0 = không lọc theo last login).
- `organizational_unit_id`: ID OU nếu muốn lọc theo OU (để trống nếu không dùng).

## Ví dụ
Repo đã chứa `cmd/active-users/config.tom`. `just build` sẽ đóng gói cấu hình
riêng cho lệnh này vào `dist/`; nếu không có, `config.toml` gốc của repo sẽ được dùng.
Để chạy cục bộ:

```
go run ./cmd/active-users -config ./cmd/active-users/config.toml -output actives.json
```
