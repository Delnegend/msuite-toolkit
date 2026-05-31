# Công cụ khóa hàng loạt user M-Suite

## Mô tả
Công cụ này lấy toàn bộ người dùng từ Admin Portal, chọn các tài khoản có `last_login_time` cũ hơn ngưỡng cấu hình, rồi khóa những tài khoản đó. Công cụ tạo báo cáo CSV cho biết trạng thái thành công hoặc thất bại của từng lần khóa.

## Các bước nhanh

- Đảm bảo M-Suite đang mở, đã bật, và ứng dụng Admin Portal xuất hiện trong danh sách ứng dụng.
- Một file `config.toml` đã có trong thư mục này và sẽ được dùng mặc định.
- Chạy file thực thi:

```
./set-inactive-users.exe
```

## Tham số

- `-config`: đường dẫn tới file config (mặc định: `./config.toml`)
- `-output`: đường dẫn file CSV đầu ra (mặc định: `inactive_users.csv`). File này là nhật ký kết quả: CSV phân tách bằng `|` với header `UserID|Result`, giá trị `Result` là `OK` hoặc thông báo lỗi.
- `-h` hoặc `-help`: hiển thị trợ giúp

## Cách điền `config.toml`

- `bearer_token`: mở Admin Portal trong trình duyệt, DevTools -> Application (hoặc Storage) -> Local Storage -> chọn origin của Admin Portal -> tìm khóa `admin_portal_access_token` và sao chép giá trị.
- `admin_user_id`: trong Admin Portal vào Identity > Users, Groups & Unit > Users, tìm user admin đang đăng nhập, click vào kết quả, sao chép `User ID` trong Basic info.
- `admin_portal_address`: thay bằng địa chỉ Admin Portal hiện tại (host:port), ví dụ `10.0.0.1:9443`.
- `last_login_threshold_in_month`: người dùng có `last_login_time` cũ hơn số tháng này sẽ bị xem là inactive.
- `dry_run`: khi `true`, công cụ chỉ liệt kê những người dùng sẽ bị khóa và không khóa tài khoản nào.
- `include_users_with_unknown_last_login`: khi `true`, người dùng không có `last_login_time` hợp lệ sẽ được đưa vào danh sách sẽ bị inactive. Giá trị mặc định là `false`.

## Cảnh báo

Nếu `last_login_threshold_in_month` nhỏ hơn `3` và `dry_run` là `false`, công cụ sẽ yêu cầu hai xác nhận theo thứ tự:

1. Xác nhận ngưỡng nguy hiểm.
2. Xác nhận đã đọc danh sách người dùng sẽ bị inactive.

Chuỗi xác nhận chính xác là:

```
I WANT TO INACTIVE USERS WITH LAST LOGIN OLDER THAN X MONTHS
I HAVE READ THE WOULD-BE-INACTIVE USERS LIST AND WANT TO PROCEED
```

`X` sẽ được thay bằng giá trị ngưỡng đã cấu hình.

## Ghi chú khi chạy

Sau khi điền `config.toml`, chạy:

```
./set-inactive-users.exe
```

Dùng `-config` để chỉ file config khác và `-output` để đổi tên file đầu ra. Công cụ sẽ lấy toàn bộ người dùng, chọn các tài khoản cũ hơn ngưỡng cấu hình, và khóa chúng trừ khi `dry_run` được bật. Khi `dry_run` là `false`, lời nhắc xác nhận ngưỡng nguy hiểm sẽ chạy trước nếu ngưỡng nhỏ hơn 3 tháng, sau đó lời nhắc xác nhận đã đọc danh sách sẽ chạy ngay sau. File đầu ra có dạng:

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
