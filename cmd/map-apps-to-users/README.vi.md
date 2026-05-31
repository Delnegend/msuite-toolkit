# Công cụ trích xuất ánh xạ Ứng dụng và Người dùng M-Suite

## Mô tả
Công cụ này trích xuất thông tin về việc người dùng nào có quyền truy cập vào những ứng dụng nào. Công cụ tạo ra hai tệp CSV:
- `ONE-to-MANY`: Ánh xạ mỗi ứng dụng tới danh sách người dùng có quyền truy cập.
- `ONE-to-ONE`: Ánh xạ trực tiếp từng người dùng tới mỗi ứng dụng họ có quyền truy cập.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./map-apps-to-users.exe
```

Dùng `-config` để chỉ file config khác và `-output` để đổi tên file đầu ra.

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: tên hậu tố file CSV đầu ra (mặc định: `apps_to_users.csv`). Hai file sẽ được tạo: `ONE-to-MANY_<output>` và `ONE-to-ONE_<output>`.
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm admin đang đăng nhập và sao chép `User ID`.
- `admin_portal_address`: địa chỉ Admin Portal (ví dụ `10.0.0.1:9443`).
- `organizational_unit_id`: tùy chọn ID OU để giới hạn kết quả người dùng theo một đơn vị tổ chức cụ thể (để trống để bỏ qua).
- `filter_by.destination_host`: chuỗi phân tách bằng dấu phẩy chứa các host đích cần lọc. Dùng IP đơn lẻ hoặc khoảng IPv4 như `10.0.0.2-10`.
- `filter_by.destination_port`: chuỗi phân tách bằng dấu phẩy chứa các cổng đích cần lọc. Dùng port đơn lẻ hoặc khoảng như `443-5000`.
- Nếu `filter_by` rỗng, tất cả ứng dụng sẽ được giữ lại.

## Ví dụ cấu hình
```toml
bearer_token = "CHANGE_ME"
admin_user_id = "CHANGE_ME"
admin_portal_address = "10.0.0.1:9443"
worker_count = 16
organizational_unit_id = ""

filter_by = { destination_host = "10.0.0.1, 10.0.0.2-10", destination_port = "442, 443-5000" }
```

Ví dụ này chỉ giữ các ứng dụng có host đích là `10.0.0.1` hoặc một địa chỉ IPv4 trong khoảng từ `10.0.0.2` đến `10.0.0.10`, và cổng là `442` hoặc trong khoảng từ `443` đến `5000`.

## Ghi chú khi chạy
- Sau khi điền `config.toml`, chạy `./map-apps-to-users.exe`. Các file đầu ra mặc định sẽ xuất hiện cùng thư mục với file thực thi.
- Công cụ tạo ra hai file:
  - `ONE-to-MANY_apps_to_users.csv`: Ánh xạ mỗi ứng dụng tới danh sách User ID phân tách bằng dấu phẩy.
  - `ONE-to-ONE_apps_to_users.csv`: Ánh xạ mỗi ứng dụng tới một User ID trên mỗi dòng.

Dùng `-c` để chỉ file config khác và `-o` để đặt tên file đầu ra khác.
