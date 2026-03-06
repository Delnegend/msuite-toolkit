# Công cụ khóa hàng loạt user M-Suite

## Mô tả
Công cụ này thực hiện khóa hàng loạt tài khoản người dùng dựa trên danh sách User ID được cung cấp. Nó xử lý từng người dùng và tạo báo cáo CSV cho biết trạng thái thành công hoặc thất bại của việc khóa đối với mỗi tài khoản.

## Các bước nhanh

- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Điền `config.toml` theo hướng dẫn bên dưới.
- Chuẩn bị file input: tạo file tên `users.txt` mỗi dòng chứa một `user id` (không có header). Ví dụ:

```
12345
67890
```

- Mở terminal trong thư mục này (nhấp chuột phải vào thư mục này và chọn "Open in Terminal") và chạy:

```
./inactive-users.exe -input users.txt
```

## Tham số

- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `inactive_users.csv`). File này là nhật ký kết quả: CSV phân tách bằng `|` với header `UserID|Result`, giá trị `Result` là `OK` hoặc thông báo lỗi.
- `-input`: đường dẫn tới file input (bắt buộc). Mỗi dòng của file input chứa một `user id` cần khóa.
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`

- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm user admin đang đăng nhập, click vào kết quả, sao chép `User ID` trong Basic info.
- `admin_portal_address`: thay bằng địa chỉ Admin Portal hiện tại (host:port), ví dụ `10.0.0.1:9443`.

## Ghi chú khi chạy

- Sau khi điền `config.toml`, chạy `./inactive-users.exe -input users.txt`. Công cụ sẽ gọi Admin Portal để khóa từng user và ghi kết quả vào file CSV đầu ra (xem `-output`).
- File đầu ra có dạng:

```
UserID|Result
12345|OK
67890|request failed: unexpected status code: 500
```

## Ví dụ đầu ra
```
UserID|Result
12345|OK
67890|request failed: unexpected status code: 500
```
