# Golang の最新バージョンを使用
FROM golang:latest AS builder

WORKDIR /app

# Air をインストール（ホットリロード用）
RUN go install github.com/air-verse/air@latest

# 依存関係をコピーしてインストール
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのコードをコピー
COPY . .

# 最小のランタイム環境
FROM golang:latest
WORKDIR /app
COPY --from=builder /go/bin/air /usr/local/bin/air
COPY . .

CMD ["air"]
