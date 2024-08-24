# ビルドステージ
FROM golang:1.23 AS builder

WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドステージから実行可能ファイルをコピー
COPY --from=builder /app/main .

# アプリケーションを実行
CMD ["./main"]