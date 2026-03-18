---
name: self-reviewer
description: >
  実装後のコードを構造・責務・リスクの観点でレビューする。
  verify通過後に必ず実行する。コードの書き換えは行わない。
---

# Self Reviewer

実装後のコードをレビューし、構造・責務・リスクを検出するスキル。

「実装者とは別の視点」で評価する。
「動くコード」ではなく「壊れない構造」を保証するためのもの。

## 実行タイミング

```
make verify 完了後（必須・スキップ禁止）
```

verify を通過していない状態では実行しない。

## 入力

- 変更ファイル一覧（`git diff --name-only`）
- 差分コード（`git diff`）
- impact-analysis 結果（同一 Slice 内で実施済みのもの）
- slice-implementation-plan（同一 Slice 内で作成済みのもの）

## レビュー観点（必須）

### 構造

- handler → repository の原則が守られているか
- service は必要最小限か（ラッパーになっていないか）

### 責務

- handler にビジネスロジックが残っていないか
- service がユースケース単位か
- repository が純粋なDBアクセスか

### 過剰設計

- 不要な抽象化がないか
- interface の乱用がないか

### リスク

- トランザクション不足
- データ不整合の可能性

### テスト

- 分岐が網羅されているか
- モックが適切か

## 出力形式

```md
# レビュー結果

## 1. 概要
変更内容の要約

## 2. 良い点
-

## 3. 問題点（重要度順）

### High
-

### Medium
-

### Low
-

## 4. 設計レビュー

### 責務分離
- handler:
- service:
- repository:

### 不適切なロジック配置
-

### 依存関係の問題
-

## 5. services導入レビュー（該当時）

- 導入は妥当か: YES / NO
- 過剰設計になっていないか
- serviceがラッパーになっていないか

## 6. テストレビュー

- カバレッジは十分か
- 不足ケース

## 7. リスク分析

- 将来の破綻ポイント
- 不整合の可能性

## 8. 改善提案

-

## 9. マージ判断

- APPROVE / REQUEST CHANGES
- 理由:
```

## REQUEST CHANGES の場合

self-reviewer が `REQUEST CHANGES` を出した場合：

1. 問題点と改善提案を人間に提示する
2. 人間の判断を待つ
3. 修正が承認された場合、実装に戻り修正 → 検証 → 再レビュー

## 禁止事項

- コードを書き換えない
- 新しい設計を導入しない
- 検証（lint / test）は行わない（change-verifier の責務）

## 成功条件

- 問題点が具体的なファイル・行で特定されている
- 優先度が明確（High / Medium / Low）
- 修正指針が提示されている

## change-verifier との役割分離

| 観点 | change-verifier | self-reviewer |
|---|---|---|
| lint / test / security | 担当 | 担当しない |
| 設計・責務 | 担当しない | 担当 |
| コード修正 | 修正案を提示 | 提案のみ |
| 実行タイミング | 実装・テスト生成後 | verify 通過後 |
