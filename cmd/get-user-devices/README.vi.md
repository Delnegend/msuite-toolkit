# Công cụ trích xuất dữ liệu thiết bị người dùng M-Suite

## Mô tả
Công cụ này trích xuất danh sách tất cả các thiết bị liên kết với từng người dùng trong hệ thống. File CSV đầu ra bao gồm User ID, Email, Device ID, Tên thiết bị, Loại thiết bị và thời điểm sử dụng lần cuối (Last Used) cho mỗi thiết bị.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:
```
./get-user-devices.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `user_devices.csv`)
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm admin đang đăng nhập và sao chép `User ID` từ Basic info.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `organizational_unit_id`: tùy chọn ID OU để giới hạn kết quả người dùng theo một đơn vị tổ chức cụ thể (để trống để bỏ qua).

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./get-user-devices.exe
```

Dùng `-config` để chỉ file config khác và `-output` để đổi tên file đầu ra. File đầu ra CSV phân tách bằng dấu `|` và bao gồm các cột sau: `UserID`, `UserEmail`, `DeviceID`, `DeviceName`, `DeviceType`, `LastUsed`.
