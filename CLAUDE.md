# CLAUDE.md

## プロジェクト概要

行動タグ型の家計簿Webアプリ（モノレポ構成）。
Backend: Go / Frontend: React + TypeScript

詳細な構造・ドメイン・APIは `repo_map.yaml` を参照すること。

## 基本ルール

* AIは **必ず repo_map.yaml を最初に読む**
* AIは **Makeコマンド経由のみで操作する**
* 変更は小さく保つ（30〜200行 / 1〜3ファイル）
* 既存パターンに従い、新しい設計を導入しない
* main ブランチを常に動く状態に保つ
* 人間の承認なしに次のフェーズに進まない

## 開発フロー（最重要）

Feature → Slice分解 → [承認] → Slice実装ループ → Merge

### Slice実装ループ

探索 → 影響分析 → [承認]
→ 実装計画 → [承認]
→ 実装 → テスト → 検証
→ 構造レビュー → [人間レビュー]
→ リファクタ提案

## 承認ポイント

以下は必ず人間の承認を得る：

1. Slice分解後
2. 影響分析後
3. 実装計画後
4. 構造レビュー後

## AIスキルの使用

機能開発は必ず `dev-workflow-orchestrator` から開始。直接コード生成は禁止。

| スキル | 役割 | タイミング |
|---|---|---|
| `dev-workflow-orchestrator` | 全体フローの統括 | 機能開発全体 |
| `feature-to-slices` | FeatureをSliceに分解 | 機能開発の最初 |
| `codebase-locator` | 実装場所の特定 | Slice実装開始時 |
| `impact-analysis` | 影響範囲の分析 | 実装前 |
| `slice-implementation-plan` | 具体的な実装計画の作成 | 影響分析後 |
| `safe-code-generator` | 既存パターンに沿ったコード生成 | 計画承認後 |
| `test-generator` | テストの生成 | 実装後 |
| `change-verifier` | lint・型チェック・テストの実行 | テスト生成後 |
| `self-reviewer` | 構造・責務・リスクのレビュー | verify通過後 |
| `refactor-suggester` | 小さなリファクタリング提案 | 構造レビュー後 |
| `diagnose-project-consistency` | プロジェクト構造の整合性診断 | セットアップ後・構造変更後 |

## AIの責務

### やること

* コード探索・分析
* 影響範囲の特定
* 既存パターンに沿った実装
* テスト生成
* 検証（lint / typecheck / test）
* セルフレビュー
* 小規模リファクタ提案

### やらないこと

* アーキテクチャ設計
* ドメイン設計
* モジュール境界の変更
* 大規模変更
* 新技術導入
* 人間の承認なしでのフェーズ進行

## Makeコマンド

Makeコマンドの詳細は `.claude/docs/make-commands.md` を参照すること。

### 検証フロー

編集 → quick-check → … → check → security-check → verify → self-review

```
make quick-check       # fmt-check + fmt-check-frontend + lint + lint-frontend + typecheck（編集直後）
make check             # quick-check + test + test-frontend（実装完了）
make verify            # check + e2e（最終確認）
```

### セキュリティチェックの運用ルール

* 実装完了時: `make security-check`
* PR作成前 / mainマージ後: `make security-check-full`
* gosec 誤検知の場合: `// #nosec` コメントで抑制（理由を併記）
* gitleaks でシークレット検出時: 即座にローテーションし、git履歴からの除去を検討

## コーディング規約

* 既存コードのパターンに従う
* Backend: go fmt / go vet
* Frontend: ESLint / Prettier
* テスト:
  * Backend: テーブルドリブン
  * Frontend: Vitest

## 補足

* 詳細な構造・API・ドメインは `repo_map.yaml` を参照
