# Công cụ trích xuất Chính sách Cấp phát (Provision Policies) M-Suite

## Mô tả
Công cụ này trích xuất tất cả các chính sách cấp phát từ hệ thống M-Suite. File CSV đầu ra bao gồm ID chính sách, Tên chính sách, Thời gian tạo và cấu hình đầy đủ của chính sách ở định dạng JSON.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:
```
./get-provision-policies.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `provision_policies.csv`)
- `-h` hoặc `-help`: hiển thị trợ giúp
