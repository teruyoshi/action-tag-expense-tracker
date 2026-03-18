---
name: slice-implementation-plan
description: >
  Sliceの具体的な実装計画を作成する。
  「実装計画を立てて」「このSliceをどう実装する？」「実装の手順を教えて」などの場面で使う。
  impact-analysis の後、実際のコード生成に入る前にこのスキルを使うこと。
---

# Slice Implementation Plan

1つのSliceに対して **具体的な実装計画** を作成するスキル。
影響範囲分析の結果を踏まえ、安全な実装手順を定める。

## なぜ必要か

- 計画なしに実装すると、手戻りが発生する
- 変更すべきファイルと順序を事前に決めることで、効率的に実装できる
- テスト計画を含めることで、検証漏れを防ぐ

## 前提

- `feature-to-slices` でSliceが定義済み
- `codebase-locator` で実装場所が特定済み
- `impact-analysis` で影響範囲が確認済み

## 手順

### 1. Slice の確認

対象Sliceの以下を確認する：
- 目的（何を実現するか）
- スコープ（どこまでやるか）
- 前提Slice（依存する先行Slice）

### 2. 変更ファイルの列挙

具体的に変更するファイルを一覧化する：
- 新規作成するファイル
- 修正するファイル
- 各ファイルでの変更内容（関数追加、型変更など）

### 3. 実装順序の決定

依存関係を考慮し、以下の順序で実装する：

1. **DB変更**（マイグレーション、モデル）
2. **Repository**（データアクセス）
3. **Service**（必要な場合のみ — 複数repositoryをまたぐ処理、副作用、トランザクション等）
4. **Handler**（APIエンドポイント）
5. **Frontend API Client**（API呼び出し）
6. **Frontend UI**（ページ、コンポーネント）
7. **テスト**（各層のテスト）

デフォルト構造は handler → repository。Service は必要な場合のみ導入する。
変更がない層はスキップする。

### 4. テスト計画

- どのテストを追加するか
- どのテストを修正するか
- E2Eテストが必要か

### 5. Service導入判定

```
## Service導入判定

- 導入: YES / NO
- 理由:
- 対象:
```

### 6. 検証コマンド

```
make quick-check      # 編集直後: fmt-check + fmt-check-frontend + lint + lint-frontend + typecheck
make check            # 実装完了: quick-check + test + test-frontend
make security-check   # 実装完了: 脆弱性 + SAST + シークレット検出
make verify           # E2E変更がある場合: check + e2e
```

## 出力フォーマット

```
## Slice: [名前]

### 目的
[このSliceで実現すること]

### 変更ファイル

| 順序 | ファイル | 操作 | 変更内容 |
|---|---|---|---|
| 1 | backend/migrations/xxx.sql | 新規 | テーブル追加 |
| 2 | backend/models/xxx.go | 新規 | モデル定義 |
| 3 | backend/repositories/xxx.go | 新規 | CRUD操作 |

### Service導入判定

- 導入: YES / NO
- 理由:
- 対象:

### テスト計画

- [ ] backend/handlers/xxx_test.go: [テスト内容]
- [ ] backend/repositories/xxx_test.go: [テスト内容]
- [ ] e2e/tests/xxx.spec.ts: [テスト内容]

### 検証コマンド

make quick-check → make check → make security-check → make verify
```

## 注意

- アーキテクチャの決定はしない。既存パターンに従う
- 実装計画は提案であり、人間が承認してから実装に進む
- スコープが大きすぎる場合は、Sliceのさらなる分割を提案する
