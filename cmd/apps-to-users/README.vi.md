# Công cụ trích xuất ánh xạ Ứng dụng và Người dùng M-Suite

## Mô tả
Công cụ này trích xuất thông tin về việc người dùng nào có quyền truy cập vào những ứng dụng nào. Công cụ tạo ra hai tệp CSV:
- `ONE-to-MANY`: Ánh xạ mỗi ứng dụng tới danh sách người dùng có quyền truy cập.
- `ONE-to-ONE`: Ánh xạ trực tiếp từng người dùng tới mỗi ứng dụng họ có quyền truy cập.

## Các bước nhanh
- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:
```
./apps-to-users.exe
```

## Tham số
- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: tên hậu tố file CSV đầu ra (mặc định: `apps_to_users.csv`). Hai file sẽ được tạo: `ONE-to-MANY_<output>` và `ONE-to-ONE_<output>`.
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`
- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm user admin đang đăng nhập, click vào kết quả, sao chép `User ID` trong Basic info.
- `admin_portal_address`: thay bằng địa chỉ Admin Portal hiện tại (host:port), ví dụ `10.0.0.1:9443`.

## Ghi chú khi chạy
- Sau khi điền `config.toml`, chạy `./apps-to-users.exe`. Các file đầu ra mặc định sẽ xuất hiện cùng thư mục với công cụ.
- Công cụ tạo ra hai file:
  - `ONE-to-MANY_apps_to_users.csv`: Ánh xạ mỗi ứng dụng tới danh sách User ID phân tách bằng dấu phẩy.
  - `ONE-to-ONE_apps_to_users.csv`: Ánh xạ mỗi ứng dụng tới một User ID trên mỗi dòng.
- Dùng `-c` để chỉ file config khác và `-o` để đặt tên file đầu ra khác.

## Ví dụ đầu ra
- `ONE-to-MANY_apps_to_users.csv`:
```
App|Users
App1|user1,user2,user3
App2|user4
```
- `ONE-to-ONE_apps_to_users.csv`:
```
App|User
App1|user1
App1|user2
App1|user3
App2|user4
```
