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
createdb botemail
```

### 2. Cấu hình

```bash
cp .env.example .env
# Sửa .env với DB_URL của bạn
```

### 3. Thêm email account

```bash
# Sửa .env với DB_URL, ví dụ:
# DB_URL=postgresql://user:password@localhost:5432/botemail

# Chạy lệnh thêm account:
go run ./cmd/main.go add \
  --email="your-email@gmail.com" \
  --password="your-app-password" \
  --reply-subject="Auto-Reply" \
  --reply-body="Chào bạn, chúng tôi đã nhận được email."
```

### 4. Chạy service

```bash
go run ./cmd/main.go
```

## Quản lý Accounts

### Thêm account mới:
```bash
go run ./cmd/main.go add --email="email2@gmail.com" --password="app-pass"
```

### Xem danh sách accounts:
```bash
psql $DB_URL -c "SELECT id, email, reply_subject, active FROM accounts;"
```

### Tắt/Bật account:
```bash
# Tắt
psql $DB_URL -c "UPDATE accounts SET active = false WHERE id = 1;"

# Bật
psql $DB_URL -c "UPDATE accounts SET active = true WHERE id = 1;"
```

### Xóa account:
```bash
psql $DB_URL -c "DELETE FROM accounts WHERE id = 1;"
```

### Kiểm tra emails đã xử lý:
```bash
psql $DB_URL -c "SELECT * FROM processed_emails;"
```

## Docker

```bash
# Build
docker build -t botemail .

# Run
docker run -e DB_URL="postgresql://user:pass@host:5432/botemail" botemail
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
