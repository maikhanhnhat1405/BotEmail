# BotEmail - Auto Reply Service

Service tự động reply email cho nhiều tài khoản.

## Tính năng

- Hỗ trợ nhiều tài khoản email
- Mỗi email có cấu hình riêng (IMAP, SMTP, nội dung reply)
- Phát hiện email chưa đọc
- Gửi reply tự động
- Tránh duplicate reply (dùng PostgreSQL)
- Bật/tắt account dễ dàng

## Cấu trúc project

```
BotEmail/
├── cmd/main.go           # Entry point
├── internal/
│   ├── config/           # Đọc config từ .env
│   ├── email/
│   │   ├── imap.go       # Kết nối IMAP
│   │   └── smtp.go       # Gửi email qua SMTP
│   ├── service/
│   │   └── worker.go     # Logic xử lý chính
│   └── store/
│       └── postgres.go    # Lưu accounts & email đã xử lý
├── .env
└── Dockerfile
```

## Database Schema

### Bảng accounts
Lưu thông tin các tài khoản email:

| Column | Mô tả |
|--------|--------|
| id | ID tự động |
| email | Địa chỉ email |
| password | App Password |
| imap_host | Server IMAP |
| imap_port | Port IMAP |
| smtp_host | Server SMTP |
| smtp_port | Port SMTP |
| reply_subject | Tiêu đề reply |
| reply_body | Nội dung reply |
| active | Bật/tắt (true/false) |

### Bảng processed_emails
Lưu ID email đã reply để tránh duplicate.

## Cách chạy

### 1. Cài đặt Database

```bash
# Tạo database PostgreSQL
createdb botemail -U nhat
```

### 2. Cấu hình

Tạo file `.env`:
```env
DB_URL=postgresql://nhat:123456@localhost:5432/botemail?sslmode=disable
```

### 3. Chạy

```bash
go run ./cmd/main.go
```

Lần đầu chạy sẽ tự động tạo account mặc định.

## Quản lý Accounts

### Thêm account qua SQL:

```sql
-- Thêm account mới
INSERT INTO accounts (email, password, imap_host, imap_port, smtp_host, smtp_port, reply_subject, reply_body)
VALUES ('email2@gmail.com', 'app_password', 'imap.gmail.com', '993', 'smtp.gmail.com', '587', 'Auto-Reply', 'Chào bạn!');

-- Xem danh sách accounts
SELECT * FROM accounts;

-- Tắt account
UPDATE accounts SET active = false WHERE id = 1;

-- Bật account
UPDATE accounts SET active = true WHERE id = 1;

-- Xóa account
DELETE FROM accounts WHERE id = 1;
```

### Kiểm tra emails đã xử lý:
```sql
SELECT * FROM processed_emails;
```

## Docker

```bash
# Build
docker build -t botemail .

# Run
docker run -e DB_URL=postgresql://nhat:123456@host:5432/botemail botemail
```

## Log output

```
🚀 Service started. Monitoring mailbox...
📋 Processing 2 account(s)...
🔄 Checking inbox: email1@gmail.com
📬 Found 3 unread email(s) for: email1@gmail.com
📩 Sending auto-reply to: sender@example.com
✅ Success: Reply sent to sender@example.com
🔄 Checking inbox: email2@gmail.com
📭 No unread emails for: email2@gmail.com
```
