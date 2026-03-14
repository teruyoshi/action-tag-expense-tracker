# CLAUDE.md

## プロジェクト概要

フルスタックWebアプリケーション（モノレポ構成）。
Backend は Go、Frontend は React + TypeScript。

## 技術スタック

| レイヤー | 技術 |
|---|---|
| Backend | Go / chi / GORM / golang-migrate |
| Frontend | React / TypeScript / Vite |
| DB | MySQL |
| API契約 | OpenAPI (`api/openapi.yaml`) |
| テスト | go test / Vitest / Playwright |
| Lint | golangci-lint / ESLint / Prettier |
| 環境 | Docker / Docker Compose |
| タスクランナー | Make |
| VCS | Git（trunk-based development） |

## ディレクトリ構造

```
repo/
├── backend/
│   ├── cmd/server/main.go    # エントリーポイント
│   ├── handlers/             # HTTPハンドラ
│   ├── services/             # ビジネスロジック
│   ├── repositories/         # DBアクセス
│   ├── models/               # DBモデル
│   ├── migrations/           # マイグレーション
│   └── tests/                # Backendテスト
├── frontend/src/
│   ├── pages/                # ページコンポーネント
│   ├── components/           # 再利用UIコンポーネント
│   ├── hooks/                # React hooks
│   ├── api/                  # APIクライアント
│   └── utils/                # ユーティリティ
├── api/openapi.yaml          # API仕様
├── e2e/playwright/           # E2Eテスト
├── docker/                   # コンテナ定義
├── skills/                   # AI開発スキル
├── repo_map.yaml             # AI用コードベース地図
├── Makefile
├── docker-compose.yml
└── CLAUDE.md
```

## コマンド

AIは **Make 経由のみ** でプロジェクトを操作する。

```bash
make dev           # 開発環境起動
make lint          # Backend lint (golangci-lint)
make test          # Backend テスト (go test)
make e2e           # E2E テスト (Playwright)
```

```bash
npm run dev        # Frontend 開発サーバー
npm run build      # Frontend ビルド
npm run lint       # Frontend lint (ESLint)
npm run typecheck  # Frontend 型チェック (tsc --noEmit)
npm run test       # Frontend テスト (Vitest)
```

## 開発ワークフロー

このプロジェクトは **安定継続型開発フロー** を採用する。

### 基本原則

1. 変更は小さく保つ（30〜200行 / 1〜3ファイル）
2. main ブランチを常に動く状態に保つ
3. 実装前に影響範囲を確認する
4. AIは分析と生成に使う。設計判断は人間が行う
5. 人間の承認なしに次のフェーズに進まない

### フロー

```
Feature → Slice分解 → [承認] → Slice実装ループ → Merge → 次のSlice
```

各Sliceの実装ループ：

```
探索 → 影響分析 → [承認] → 実装計画 → [承認] → 実装 → テスト → 検証 → [レビュー]
```

### 人間の承認が必要なタイミング

1. Slice分解後（粒度と順序の確認）
2. 影響分析後（リスクの確認）
3. 実装計画後（計画の承認）
4. 実装完了後（コードレビュー）

## AIスキル

`skills/` ディレクトリに開発ワークフローの各フェーズを自動化するスキルがある。

| スキル | 役割 | タイミング |
|---|---|---|
| `feature-to-slices` | FeatureをSliceに分解 | 機能開発の最初 |
| `codebase-locator` | 実装場所の特定 | Slice実装開始時 |
| `impact-analysis` | 影響範囲の分析 | 実装前 |
| `slice-implementation-plan` | 具体的な実装計画の作成 | 影響分析後 |
| `safe-code-generator` | 既存パターンに沿ったコード生成 | 計画承認後 |
| `test-generator` | テストの生成 | 実装後 |
| `change-verifier` | lint・型チェック・テストの実行 | テスト生成後 |
| `refactor-suggester` | 小さなリファクタリング提案 | レビュー後（任意） |
| `dev-workflow-orchestrator` | 全体フローの統括 | 機能開発全体 |

機能開発を始めるときは `dev-workflow-orchestrator` を起点にする。

## AIの役割と制約

### AIがやること

- コードベースの探索と分析
- 影響範囲の特定
- 既存パターンに沿ったコード生成
- テスト生成
- lint / typecheck / test の実行と結果報告
- 小さなリファクタリング提案

### AIがやらないこと

- アーキテクチャ設計の決定
- ドメインモデルの設計
- モジュール境界の変更
- 大規模リファクタリングの実行
- 新しいフレームワークやライブラリの導入
- 人間の承認なしでのフェーズ進行

## コード規約

- 既存コードのパターンに従う。新しいパターンを導入しない
- 関数は小さく保つ
- Backend: `go fmt` + golangci-lint に従う
- Frontend: ESLint + Prettier に従う
- テスト: Backend はテーブルドリブンテスト、Frontend は Vitest

## repo_map.yaml

プロジェクトルートの `repo_map.yaml` はAIがコードベースを探索する際の地図として使う。
新しいタスクに着手する際は、最初にこのファイルを読むこと。

## 検証チェックリスト

コード変更後、マージ前に必ず以下を通す：

- [ ] `make lint` — Backend lint
- [ ] `make test` — Backend テスト
- [ ] `npm run lint` — Frontend lint
- [ ] `npm run typecheck` — Frontend 型チェック
- [ ] `make e2e` — E2Eテスト（UI/API変更時）
