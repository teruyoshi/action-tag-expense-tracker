---
name: change-verifier
description: >
  変更後のコードを検証し、lint・型チェック・テストの問題を検出する。
  「検証して」「チェックして」「変更が壊れてないか確認して」「マージ前の確認」
  などの場面で使う。実装やテスト生成の後に必ずこのスキルを実行すること。
---

# Change Verifier

変更後のコードを **差分ベースで検証する** スキル。
フルチェックはCI/PRで行うため、ここでは変更に関係する部分のみ高速に検証する。

## 検証項目

| チェック | Backend | Frontend |
|---|---|---|
| Lint | `make lint` | `npm run lint` |
| 型チェック | Go compiler | `npm run typecheck` |
| テスト | `make test` | `npm run test` |
| E2E | `make e2e` | `make e2e` |

## 手順

### 1. 変更ファイルの特定

```bash
git diff --name-only
git diff --cached --name-only
```

変更されたファイルから、Backend / Frontend / 両方のどれに影響があるかを判断する。

### 2. Backend 検証

Backend ファイル（`.go`）に変更がある場合：

```bash
make lint
make test
```

エラーがあれば内容を確認し、修正案を提示する。

### 3. Frontend 検証

Frontend ファイル（`.ts`, `.tsx`）に変更がある場合：

```bash
npm run lint
npm run typecheck
```

エラーがあれば内容を確認し、修正案を提示する。

### 4. テスト実行

変更に関連するテストを実行する：

```bash
make test    # Backend テスト
npm run test # Frontend テスト
```

### 5. E2E 検証（必要な場合のみ）

以下の条件に該当する場合にE2Eテストを実行する：
- APIエンドポイントの変更
- UIの変更
- ユーザーフローに影響する変更

```bash
make e2e
```

### 6. 破壊的変更の確認

以下をチェックする：
- 型やインターフェースの変更が他のファイルに影響していないか
- 削除されたものが他で使われていないか
- マイグレーションの互換性

## 出力フォーマット

```
## 検証結果

| チェック | 結果 | 詳細 |
|---|---|---|
| Backend Lint | OK / NG | [詳細] |
| Backend Test | OK / NG | [詳細] |
| Frontend Lint | OK / NG | [詳細] |
| Frontend TypeCheck | OK / NG | [詳細] |
| E2E | OK / SKIP | [詳細] |

## 問題点

1. [問題の説明と修正案]

## マージ可否

[OK: 全チェック通過 / NG: 修正が必要]
```

## 注意

- 検証はコマンドを実際に実行して結果を確認する
- エラーメッセージを正確に読み、的確な修正案を出す
- 判断が難しい警告は人間に報告する
