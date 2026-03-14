---
name: codebase-locator
description: >
  実装すべき場所をコードベースから特定する。
  「どこに実装すればいい？」「このファイルはどこ？」「実装場所を教えて」
  「既存のパターンを確認したい」などの場面で使う。
  新しいコードを書く前に、まず既存コードの構造を理解するために積極的に使うこと。
---

# Codebase Locator

実装前に **どこに何を置くべきか** を特定するスキル。
探索時間を大幅に削減し、既存パターンとの一貫性を保つ。

## なぜ必要か

- 間違った場所にコードを置くとアーキテクチャが崩れる
- 既存パターンを無視すると保守性が下がる
- 探索に時間をかけすぎると開発効率が落ちる

## 手順

### 1. repo_map.yaml の確認

プロジェクトルートの `repo_map.yaml` を最初に読む。
特に以下のセクションを確認：
- `directory_structure`: モジュール構成
- `quick_lookup`: タスク別の参照先

### 2. タスクの分類

変更の種類を特定する：

| タスク種別 | 主な対象ディレクトリ |
|---|---|
| APIエンドポイント追加 | `api/openapi.yaml`, `backend/handlers`, `backend/services`, `backend/repositories` |
| DBテーブル追加 | `backend/models`, `backend/migrations` |
| フロントエンドページ追加 | `frontend/src/pages`, `frontend/src/components`, `frontend/src/api` |
| ビジネスロジック変更 | `backend/services`, `backend/repositories` |
| テスト追加 | `backend/tests`, `e2e/playwright` |

### 3. 既存パターンの調査

対象ディレクトリ内の既存ファイルを読み、以下を把握する：
- ファイル命名規則
- 関数・構造体の命名パターン
- import の構成
- エラーハンドリングのパターン
- テストの書き方

### 4. 関連ファイルの特定

変更対象だけでなく、影響を受ける可能性のあるファイルも列挙する。

## 出力フォーマット

```
## 実装対象

| ファイル | 目的 | 操作 |
|---|---|---|
| backend/handlers/xxx.go | HTTPハンドラ | 新規作成 |
| backend/services/xxx.go | ビジネスロジック | 新規作成 |

## 既存パターン

- ハンドラの例: backend/handlers/existing.go
- サービスの例: backend/services/existing.go

## 関連ファイル

- backend/services/other.go (依存関係あり)
```

## 探索ツール

以下のツールを使って調査する：
- `Glob`: ファイルパターン検索
- `Grep`: コード内容検索
- `Read`: ファイル読み込み

`repo_map.yaml` の `quick_lookup` セクションに該当するタスクがあれば、そこに記載されたパスから探索を開始する。

## 注意

- 新しいディレクトリ構造を提案しない。既存構造に従う
- 推測ではなく、実際にファイルを読んで確認する
- パターンが見つからない場合は、人間に確認する
