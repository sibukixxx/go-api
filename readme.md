# Docker-based Development Guide for Go API Server

## 目次
1. 開発環境のセットアップ
2. Docker開発ワークフロー
3. ベストプラクティス
4. 一般的なコマンドとタスク
5. トラブルシューティング
6. CI/CDパイプライン

## 1. 開発環境のセットアップ

### 必要なツール
- Docker
- Docker Compose
- Git
- お好みのコードエディタ（VSCode推奨）

### 初期セットアップ
1. リポジトリをクローン：
   ```
   git clone https://github.com/your-org/go-api-server.git
   cd go-api-server
   ```

2. 開発用Docker Composeファイルを作成：
   ```yaml
   version: '3'
   services:
     app:
       build:
         context: .
         dockerfile: Dockerfile.dev
       volumes:
         - .:/app
       ports:
         - "8080:8080"
       environment:
         - GO_ENV=development
   ```

3. 開発用Dockerfileを作成（Dockerfile.dev）：
   ```dockerfile
   FROM golang:1.22
   
   WORKDIR /app
   
   COPY go.mod go.sum ./
   RUN go mod download
   
   COPY . .
   
   RUN go install github.com/cosmtrek/air@latest
   
   CMD ["air"]
   ```

## 2. Docker開発ワークフロー

1. 開発サーバーの起動：
   ```
   docker-compose up --build
   ```

2. コードの変更：
    - ローカルでコードを編集
    - 変更は自動的に反映される（airを使用）

3. テストの実行：
   ```
   docker-compose exec app go test ./...
   ```

4. 新しい依存関係の追加：
    - go.modファイルを更新
    - `docker-compose down && docker-compose up --build`を実行

5. データベースマイグレーション（必要な場合）：
   ```
   docker-compose exec app go run migrations/migrate.go
   ```

## 3. ベストプラクティス

- **バージョン管理**：Dockerfileとdocker-compose.ymlファイルをバージョン管理に含める
- **環境変数**：機密情報は.envファイルに保存し、.gitignoreに追加
- **レイヤーキャッシング**：Dockerfileを最適化し、ビルド時間を短縮
- **マルチステージビルド**：本番用Dockerfileではマルチステージビルドを使用
- **軽量ベースイメージ**：本番環境ではalpineベースイメージを使用
- **ヘルスチェック**：Dockerfileにヘルスチェックを追加

## 4. 一般的なコマンドとタスク

- アプリケーションのビルドと起動：
  ```
  docker-compose up --build
  ```

- コンテナ内でコマンドを実行：
  ```
  docker-compose exec app go run main.go
  ```

- ログの確認：
  ```
  docker-compose logs -f app
  ```

- コンテナの停止と削除：
  ```
  docker-compose down
  ```

## 5. トラブルシューティング

- **ポートの競合**：8080ポートが既に使用されている場合、docker-compose.ymlのポートマッピングを変更
- **ホットリロードが機能しない**：airの設定を確認し、必要に応じて.air.tomlファイルを調整
- **依存関係の問題**：go.sumファイルを削除し、`go mod tidy`を実行後、再ビルド

## 6. CI/CDパイプライン

- GitHub Actionsを使用（.github/workflows/main.ymlを参照）
- プルリクエスト時に自動テストとリンターを実行
- mainブランチへのマージ時に自動デプロイ

詳細は`github-actions-workflow`アーティファクトを参照してください。

## 結論

このガイドに従うことで、チームメンバー全員が一貫した開発環境で作業し、効率的にコラボレーションできます。質問や提案がある場合は、チームリーダーに連絡してください。