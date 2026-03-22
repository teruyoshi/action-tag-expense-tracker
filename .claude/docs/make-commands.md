# Makeコマンド一覧

## 開発・運用

```
make dev               # Docker Compose起動（db + backend + frontend）
make stop              # Docker Compose停止
make build             # Docker Composeビルド
make clean             # Docker環境の完全クリーンアップ（コンテナ・ボリューム・イメージ削除）
make ci-up             # CI用サービス起動（VITE_WATCH=false）
```

## フォーマット

```
make fmt               # Backend gofmt
make fmt-frontend      # Frontend Prettier
make fmt-check         # Backend フォーマットチェック（変更なし）
make fmt-check-frontend # Frontend フォーマットチェック（変更なし）
```

## Lint・型チェック

```
make lint              # Backend go vet
make lint-frontend     # Frontend ESLint
make typecheck         # Frontend tsc --noEmit
```

## テスト

```
make test              # Backend go test
make test-frontend     # Frontend Vitest
make e2e               # Playwright E2Eテスト
make test-diff         # 変更ファイルの影響テストのみ実行
```

## 検証（複合コマンド）

```
make quick-check       # fmt-check + fmt-check-frontend + lint + lint-frontend + typecheck（編集直後）
make check             # quick-check + test + test-frontend（実装完了）
make verify            # check + e2e（最終確認）
```

## セキュリティ

```
make security-check      # 日常用（govulncheck + npm audit + gosec + gitleaks）
make security-check-full # フル（security-check + trivy）
make vuln-backend        # Backend govulncheck
make vuln-frontend       # Frontend npm audit
make sec-backend         # Backend gosec
make sec-secrets         # gitleaks シークレット検出
make sec-fs              # trivy リポジトリスキャン
```

## マイグレーション

```
make migrate-up        # DBマイグレーション適用
make migrate-down      # DBマイグレーション1つ戻す
make migrate-create    # 新規マイグレーションファイル作成
```

## 診断

```
make doctor            # プロジェクト整合性チェック
```
