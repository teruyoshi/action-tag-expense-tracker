---
name: impact-analysis
description: >
  コード変更の影響範囲を分析し、壊れる可能性のある箇所を特定する。
  「影響範囲を調べて」「この変更で何が壊れる？」「破壊的変更はある？」
  「変更前にリスクを確認したい」などの場面で使う。
  実装に着手する前に、必ずこのスキルを使って影響範囲を確認すること。
---

# Impact Analysis

実装前に **変更が何を壊す可能性があるか** を分析するスキル。
破壊的変更を未然に防ぐ。

## なぜ必要か

- 変更は予想外の場所に影響する
- テストが通っても、カバーされていない箇所が壊れる可能性がある
- 事前に把握すれば、対策を含めた実装ができる

## 手順

### 1. 変更内容の把握

以下を明確にする：
- 何を変更するか（ファイル、関数、型など）
- どう変更するか（追加、修正、削除）

### 2. repo_map.yaml の影響分析ルール確認

`repo_map.yaml` の `impact_analysis_rules` セクションを確認：

```yaml
backend_changes:
  check: [backend/services, backend/repositories, backend/handlers, backend/tests]
database_changes:
  check: [backend/models, backend/migrations, backend/repositories]
api_changes:
  check: [api/openapi.yaml, frontend/src/api, backend/handlers]
frontend_changes:
  check: [frontend/src/components, frontend/src/pages, e2e/playwright]
```

### 3. 依存関係の追跡

対象ファイルから以下を調査する：

**Backend (Go)**
- `import` パス → このパッケージを使っているファイル
- 構造体・インターフェースの利用箇所
- 関数の呼び出し元

**Frontend (TypeScript)**
- `import` 文 → このモジュールを使っているファイル
- 型定義の利用箇所
- コンポーネントの使用箇所

**API**
- OpenAPI定義の変更 → backend handler + frontend api client

**Database**
- モデル変更 → repository → service → handler

### 4. リスク評価

各影響箇所について：
- **高リスク**: 型変更、インターフェース変更、DB スキーマ変更
- **中リスク**: ロジック変更、新規依存追加
- **低リスク**: 追加のみ（既存コードに触れない）

### 5. 必要なテストの特定

影響を受ける箇所に対して、以下を確認：
- 既存テストがカバーしているか
- 追加テストが必要か
- E2Eテストへの影響

## 出力フォーマット

```
## 変更対象

- [ファイル名]: [変更内容]

## 影響範囲

| 影響を受けるファイル | 理由 | リスク |
|---|---|---|
| backend/services/xxx.go | 型変更の影響 | 高 |
| frontend/src/api/xxx.ts | API変更の影響 | 中 |

## リスクのある箇所

1. [具体的なリスク説明]

## 必要なテスト

- [ ] backend/tests/xxx_test.go の更新
- [ ] e2e テストの追加
```

## 調査ツール

- `Grep`: import、関数呼び出し、型参照の検索
- `Glob`: 関連ファイルの探索
- `Read`: 依存関係の詳細確認

## 注意

- 推測ではなく、実際のコードを読んで分析する
- 影響範囲の見落としは事故につながる。広めに調査する
- 判断が難しい場合は「要確認」として人間に委ねる
