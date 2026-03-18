---
name: change-verifier
description: >
  変更後のコードを検証し、lint・型チェック・テストの問題を検出する。
  「検証して」「チェックして」「変更が壊れてないか確認して」「マージ前の確認」
  などの場面で使う。実装やテスト生成の後に必ずこのスキルを実行すること。
---

# Change Verifier

コード変更後に **壊れていないか高速に検証する** スキル。
検証を3段階に分け、状況に応じて最適なレベルを選択する。

## 役割の範囲

このスキルは **技術的検証のみ** を担当する：

- lint / format チェック
- 型チェック
- テスト実行
- セキュリティチェック

**設計レビュー・構造判断は行わない。** それらは `self-reviewer` の責務。

## 絶対ルール

**Make コマンド以外を直接実行してはいけない。**

```
# 禁止
go test ./...
npm run lint
npx tsc --noEmit
npx playwright test

# 必ず Make 経由
make quick-check
make check
make verify
```

理由: Docker / npm / go / playwright の差異を Make が吸収している。
直接実行すると環境差異でCI結果と一致しなくなる。

## 4段階の検証フロー

### レベル1: quick-check（編集直後）

```bash
make quick-check
```

| 実行内容 | コマンド |
|---|---|
| Backend フォーマットチェック | `make fmt-check` |
| Frontend フォーマットチェック | `make fmt-check-frontend` |
| Backend lint | `make lint` |
| Frontend lint | `make lint-frontend` |
| Frontend 型チェック | `make typecheck` |

**いつ使うか:** ファイルを編集した直後。最も頻繁に使う。
**目安:** 5〜15秒

AIの開発ループでは、編集のたびにこれを回す。
テストまで毎回実行すると速度が落ちるため、lint と型チェックのみで素早くフィードバックを得る。

### レベル2: check（実装完了時）

```bash
make check
```

| 実行内容 | コマンド |
|---|---|
| quick-check 全項目 | fmt-check + fmt-check-frontend + lint + lint-frontend + typecheck |
| Backend テスト | `make test` |
| Frontend テスト | `make test-frontend` |

**いつ使うか:** ある程度まとまった編集が終わったとき。コミット前。
**目安:** 10〜60秒

quick-check を内包しているため、別途 quick-check を実行する必要はない。

### レベル3: security-check（実装完了時）

```bash
make security-check
```

| 実行内容 | コマンド |
|---|---|
| Backend 脆弱性チェック | `make vuln-backend` |
| Frontend 脆弱性チェック | `make vuln-frontend` |
| Backend 静的セキュリティ解析 | `make sec-backend` |
| シークレット漏洩検出 | `make sec-secrets` |

**いつ使うか:** 実装完了時、check と合わせて実行する。PR前は必須。
**目安:** 10〜30秒

**security-check はスキップ禁止。** PR前には必ず実行すること。

### レベル4: verify（機能完成時）

```bash
make verify
```

| 実行内容 | コマンド |
|---|---|
| check 全項目 | fmt-check + fmt-check-frontend + lint + lint-frontend + typecheck + test + test-frontend |
| E2E テスト | `make e2e` |

**いつ使うか:** 機能が完成し、マージ前の最終確認。
**目安:** 30秒〜数分

API エンドポイントの変更、UI の変更、ユーザーフローに影響する変更があった場合は必ず実行する。

## 差分テスト（速度重視）

通常の `make test` が全テストを実行するのに対し、差分テストは変更ファイルに関係するテストのみ実行する。

```bash
make test-diff
```

**いつ使うか:** quick-check は通ったが、テストも確認したい場合。全テストを回すほどではない場合。
フルテスト（`make check`）の3〜5倍速い。

## AI開発ループ

```
編集 → quick-check → 編集 → quick-check → ... → check → security-check → verify
```

このループが AIDD の速度の核心。具体的には：

1. **編集中:** `make quick-check` を繰り返す
2. **実装完了:** `make check`（または `make test-diff` で影響テストのみ）
3. **実装完了:** `make security-check`（脆弱性 + SAST + シークレット検出）
4. **機能完成:** `make verify`（E2E含む完全検証）

## 手順

### 1. 変更ファイルの特定

```bash
git diff --name-only
git diff --cached --name-only
```

変更されたファイルから Backend / Frontend / 両方のどれに影響があるかを判断する。

### 2. 検証レベルの判断

| 状況 | レベル |
|---|---|
| ファイルを1〜2個編集した直後 | `make quick-check` |
| 差分に対するテストだけ確認したい | `make test-diff` |
| まとまった実装が完了した | `make check` |
| 実装完了、PR前 | `make security-check` |
| 機能完成、マージ前 | `make verify` |
| プロジェクト構造に不安がある | `make doctor` |

### 3. 実行と結果確認

選択したレベルのコマンドを実行し、結果を確認する。
エラーがあれば内容を読み、修正案を提示する。

### 4. 破壊的変更の確認

以下をチェックする：
- 型やインターフェースの変更が他のファイルに影響していないか
- 削除されたものが他で使われていないか
- マイグレーションの互換性

## 出力フォーマット

```
## 検証結果

- quick-check: OK / NG
- check: OK / NG
- security-check: OK / NG
- verify: OK / NG

### 詳細

| チェック | 結果 | 詳細 |
|---|---|---|
| Backend Format | OK / NG | [詳細] |
| Frontend Format | OK / NG | [詳細] |
| Backend Lint | OK / NG | [詳細] |
| Frontend Lint | OK / NG | [詳細] |
| Frontend TypeCheck | OK / NG | [詳細] |
| Backend Test | OK / NG / SKIP | [詳細] |
| Frontend Test | OK / NG / SKIP | [詳細] |
| Security Check | OK / NG / SKIP | [詳細] |
| E2E | OK / NG / SKIP | [詳細] |

## 問題点

1. [問題の説明と修正案]

## 判定

[OK: 全チェック通過 / NG: 修正が必要]
```

## 注意

- 検証はコマンドを実際に実行して結果を確認する
- エラーメッセージを正確に読み、的確な修正案を出す
- 判断が難しい警告は人間に報告する
- 速度を優先し、不必要に高いレベルの検証を行わない
