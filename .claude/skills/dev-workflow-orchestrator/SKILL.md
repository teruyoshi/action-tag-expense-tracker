---
name: dev-workflow-orchestrator
description: >
  開発ワークフロー全体を統括し、各スキルを正しい順序で実行する。
  「機能を開発したい」「この機能を作って」「開発を始めたい」「ワークフローを開始して」
  「指示書に従って実装して」「この計画で進めて」「このSliceを実装して」「〜を導入して」
  などの場面で使う。Feature開発の全体フローを管理するオーケストレーターとして、
  常にこのスキルを通じて開発を進めること。
  指示書や実装計画が添付されている場合も例外にしない。
  除外: repo_map更新、ドキュメント修正、設定変更など非機能タスク。
---

# Dev Workflow Orchestrator

開発ワークフロー全体を統括するスキル。
各フェーズのスキルを **正しい順序で** 実行し、人間の承認を挟みながら進める。

## ワークフロー全体像

```
Feature
  ↓
feature-to-slices      ← Feature を Slice に分解
  ↓
[人間が Slice を承認]
  ↓
┌─────────────────────── Slice ループ ───────────────────────┐
│                                                            │
│  codebase-locator    ← 実装場所の特定                       │
│    ↓                                                       │
│  impact-analysis     ← 影響範囲の分析                       │
│    ↓                                                       │
│  [人間がリスクを確認]                                        │
│    ↓                                                       │
│  slice-implementation-plan  ← 実装計画の作成                 │
│    ↓                                                       │
│  [人間が計画を承認]                                          │
│    ↓                                                       │
│  safe-code-generator  ← コード生成                          │
│    ↓                                                       │
│  test-generator       ← テスト生成                          │
│    ↓                                                       │
│  change-verifier      ← 検証（lint, test, typecheck, security） │
│    ↓                                                       │
│  self-reviewer        ← 構造レビュー（verify通過後、必須）     │
│    ↓                                                       │
│  [人間がレビュー結果を確認]                                    │
│    ↓                                                       │
│  refactor-suggester   ← リファクタ提案                       │
│    ↓                                                       │
│  Merge → 次の Slice へ                                      │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

## 使い方

### Phase 1: Feature 分解

ユーザーから Feature の説明を受け取ったら：

1. `feature-to-slices` を実行
2. Slice一覧、依存関係、実装順序を提示
3. **ユーザーの承認を待つ**

```
→ Sliceの分割でよいですか？修正があれば教えてください。
```

### Phase 2: Slice 実装ループ

承認されたSliceリストの各Sliceに対して、以下を順番に実行する：

#### Step 1: 探索

`codebase-locator` を実行し、実装場所と既存パターンを特定。

#### Step 2: 影響分析

`impact-analysis` を実行し、リスクを把握。

```
→ 影響範囲を確認しました。リスクの高い箇所があります。進めてよいですか？
```

#### Step 3: 実装計画

`slice-implementation-plan` を実行し、具体的な計画を提示。

```
→ この実装計画で進めてよいですか？
```

#### Step 4: 実装

ユーザーの承認後、`safe-code-generator` でコードを生成する。
このステップではプロダクションコードのみを生成する。テストは次のステップで生成する。

#### Step 5: テスト生成

`test-generator` を実行し、Step 4 で生成したコードに対するテストを生成する。
テスト生成は **省略しない**。実装後に必ず実行すること。

#### Step 6: 検証

`change-verifier` を実行し、全チェックを通す。

検証フローは以下の順序で行う：

```
make quick-check → make check → make security-check → make verify
```

- `security-check` はスキップ禁止。PR前には必ず実行すること
- main相当では `make security-check-full` を考慮する
- 問題があれば修正し、再検証する

#### Step 7: 構造レビュー

verify が通ったら、`self-reviewer` を実行する。**スキップ禁止。**

self-reviewer は構造・責務・リスクの観点でレビューを行う。
技術的検証（lint / test / security）は Step 6 の change-verifier が担当済みのため、ここでは行わない。

レビュー結果を以下の形式で提示する：

```
## Self Review結果
（self-reviewerの出力をそのまま添付）
```

- **APPROVE** の場合 → レポートを提示し、ユーザーの確認を待ってから次へ進む
- **REQUEST CHANGES** の場合 → 問題点と改善提案を提示し、人間の判断を待つ。修正が承認された場合は Step 4（実装）に戻り、修正 → Step 6（検証）→ Step 7（レビュー）をやり直す

#### Step 8: リファクタ提案

self-reviewer が APPROVE を出した後、**コミット・マージに進む前に** `refactor-suggester` を実行する。
**スキップ禁止。** self-reviewer → refactor-suggester の順序を必ず守ること。

リファクタリング候補を提示し、ユーザーが実行を承認した場合は修正 → Step 6（検証）をやり直す。
ユーザーが不要と判断した場合はそのままマージに進む。

#### Step 9: マージ → 次のSlice

refactor-suggester 完了後、マージしてよいか確認し、次のSliceへ進む。

## 人間の承認ポイント

以下のタイミングで必ずユーザーの確認を取る：

1. **Slice分解後** — Sliceの粒度と順序
2. **影響分析後** — リスクの確認
3. **実装計画後** — 計画の承認
4. **構造レビュー後** — self-reviewer のレビュー結果の確認

AIは **これらの承認なしに次のフェーズに進まない**。

## 進捗管理

各Sliceの状態を追跡する：

```
## 進捗

| Slice | 状態 | 備考 |
|---|---|---|
| DBテーブル追加 | 完了 | |
| 通知API追加 | 実装中 | 影響分析完了 |
| 通知送信ロジック | 未着手 | |
```

## 重要原則

1. **変更は小さく保つ** — 各Sliceのサイズ基準を守る
2. **mainは壊さない** — 検証を必ず通す
3. **AIは判断しない** — 設計判断は人間に委ねる
4. **承認なしに進まない** — 各フェーズで人間の確認を取る

## コマンド参照

```bash
make dev              # 開発環境起動（Docker Compose）
make lint             # Backend lint
make test             # Backend テスト
make lint-frontend    # Frontend lint
make typecheck        # Frontend 型チェック
make test-frontend    # Frontend テスト
make e2e              # E2E テスト
make quick-check      # 編集直後の検証
make check            # 実装完了時の検証
make security-check   # 実装完了時のセキュリティ検証（スキップ禁止）
make verify           # マージ前の最終検証
```
