# CLAUDE.md

## プロジェクト概要

フルスタックWebアプリケーション（モノレポ構成）。
Backend は Go、Frontend は React + TypeScript。

## 技術スタック

| レイヤー | 技術 |
|---|---|
| Backend | Go / chi / GORM |
| Frontend | React / TypeScript / Vite |
| DB | MySQL |
| テスト | go test / Vitest / Playwright |
| Lint | go vet / ESLint |
| 環境 | Docker / Docker Compose |
| タスクランナー | Make |
| VCS | Git（trunk-based development） |

## ディレクトリ構造

```
repo/
├── backend/
│   ├── cmd/server/main.go    # エントリーポイント
│   ├── handlers/             # HTTPハンドラ
│   ├── repositories/         # DBアクセス
│   ├── models/               # DBモデル
│   ├── migrations/           # DBマイグレーション
│   └── Dockerfile
├── frontend/src/
│   ├── pages/                # ページコンポーネント
│   ├── components/           # 共通コンポーネント
│   └── api/                  # APIクライアント
├── e2e/                      # E2Eテスト（Playwright）
├── scripts/                  # Make から呼ばれるスクリプト
├── skills/                   # AI開発スキル
├── repo_map.yaml             # AI用コードベース地図
├── Makefile
├── docker-compose.yml
└── CLAUDE.md
```

## コマンド

AIは **Make 経由のみ** でプロジェクトを操作する。

```bash
# --- 開発 ---
make dev              # 開発環境起動（Docker + Frontend devサーバー）
make dev-frontend     # Frontend devサーバーのみ
make stop             # Docker Compose 停止

# --- AI検証（3段階） ---
make quick-check      # 編集直後: lint + lint-frontend + typecheck
make check            # 実装完了: quick-check + test + test-frontend
make verify           # 機能完成: check + e2e

# --- 個別コマンド ---
make lint             # Backend lint (go vet)
make lint-frontend    # Frontend lint (ESLint)
make typecheck        # Frontend 型チェック (tsc --noEmit)
make test             # Backend テスト (go test)
make test-frontend    # Frontend テスト (Vitest)
make test-diff        # 変更ファイルの影響テストのみ実行
make e2e              # E2E テスト (Playwright)

# --- 診断 ---
make doctor           # プロジェクト整合性チェック
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
4. 自動レビュー後（レビューレポートの確認）

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
| `refactor-suggester` | 小さなリファクタリング提案 | レビュー後 |
| `dev-workflow-orchestrator` | 全体フローの統括 | 機能開発全体 |
| `diagnose-project-consistency` | プロジェクト構造の整合性診断 | セットアップ後・構造変更後 |

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
- Backend: `go fmt` + `go vet` に従う
- Frontend: ESLint に従う
- テスト: Backend はテーブルドリブンテスト、Frontend は Vitest

## repo_map.yaml

プロジェクトルートの `repo_map.yaml` はAIがコードベースを探索する際の地図として使う。
新しいタスクに着手する際は、最初にこのファイルを読むこと。

## 検証チェックリスト

コード変更後、段階的に検証する：

- [ ] `make quick-check` — 編集直後（lint + typecheck）
- [ ] `make check` — 実装完了時（lint + typecheck + test）
- [ ] `make verify` — マージ前の最終確認（check + e2e）

AI開発ループ: `編集 → quick-check → 編集 → quick-check → ... → check → verify`
