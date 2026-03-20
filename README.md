# BotEmail - Auto Reply Service

Service tự động reply email cho Gmail.

## Tính năng

- Kết nối mailbox qua IMAP
- Phát hiện email chưa đọc
- Gửi reply tự động
- Tránh duplicate reply (dùng SQLite)
- Log kết quả xử lý

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
│       └── sqlite.go     # Lưu email đã xử lý
├── data/                 # Folder chứa database
├── .env                  # Config (email, password)
└── Dockerfile
```

## Cách chạy

### 1. Cài đặt

```bash
# Tạo thư mục data
mkdir -p data

# Tạo file .env
touch .env
```

### 2. Cấu hình

Sửa file `.env`:

```env
EMAIL_USER=your_email@gmail.com
EMAIL_PASS=your_app_password
```

**Lấy App Password:**
1. Bật 2-Step Verification: https://myaccount.google.com/security
2. Tạo App Password: https://myaccount.google.com/apppasswords

### 3. Chạy

```bash
go run ./cmd/main.go
```

### 4. Tùy chỉnh reply

```env
REPLY_SUBJECT=Auto-Reply
REPLY_BODY=Hello World
```

## Docker

```bash
# Build
docker build -t botemail .

# Run
docker run -e EMAIL_USER=xxx -e EMAIL_PASS=xxx botemail
```

## Database

File `data/emails.db` lưu ID các email đã reply để tránh duplicate.

Kiểm tra dữ liệu:
```bash
sqlite3 data/emails.db "SELECT * FROM processed_emails;"
```

## Log output

```
🚀 Service started. Monitoring mailbox...
📩 Sending auto-reply to: sender@gmail.com
✅ Success: Reply sent to sender@gmail.com
```
