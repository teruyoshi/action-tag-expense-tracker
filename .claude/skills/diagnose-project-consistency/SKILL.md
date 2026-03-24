---
name: diagnose-project-consistency
description: >
  リポジトリ構造・開発コマンド・テスト・ドキュメント・生成物の整合性を診断し、
  AI開発（AIDD）に適した構造かを評価する。
  ユーザーが「プロジェクトの状態を確認して」「構造を診断して」「整合性チェックして」
  「AIDD対応できてる？」「repo_mapと合ってる？」などと言ったときに使う。
  新しいプロジェクトのセットアップ後や、大きな構造変更の後に積極的に提案すること。
---

# Project Consistency Diagnostic

リポジトリの構造・コマンド・テスト・ドキュメント・生成物の整合性を診断し、
AI駆動開発（AIDD）に適した状態かを評価するスキル。

## 診断の2段階構成

診断は **doctor.sh（機械的チェック）** と **AIチェック（判断が必要なもの）** の2段階で行う。

### Step 1: `make doctor` を実行

以下のカテゴリは `scripts/doctor.sh` が自動チェックする。まず実行して結果を確認する。

| # | カテゴリ | doctor.sh のセクション |
|---|---------|----------------------|
| 1 | repo_map.yaml 整合性 | セクション1: キー存在・APIルート→ハンドラ・エンティティ→モデル |
| 2 | Makefile 整合性 | セクション2: ターゲット定義・CLAUDE.md との一致 |
| 3 | Docker 整合性 | セクション3: docker-compose.yml・Dockerfile・コンテナ状態 |
| 4 | ディレクトリ構造 | セクション4: 必須ディレクトリの存在 |
| 5 | 依存関係 | セクション5: go.mod・package.json・テストツール |
| 6 | テスト整合性 | セクション6: Backend _test.go ペア・Frontend hooks/api テスト |
| 7 | Migration 整合性 | セクション7: up/down ペア・連番チェック |
| 8 | CI/CD 整合性 | セクション8: Make 経由・必須ステップ存在 |
| 9 | Lint/Formatter 設定 | セクション9: ESLint・Prettier・go vet 設定 |
| 10 | スキル一覧整合性 | セクション10: CLAUDE.md ↔ .claude/skills/ ディレクトリ |

### Step 2: AIによる追加チェック

doctor.sh では判断できない以下の4カテゴリをAIがチェックする。

---

#### AI-1: フロントエンドルーティング整合性

repo_map の `frontend.routing` に記載されたルートに対応するページコンポーネントが存在するかを検証する。

1. `repo_map.yaml` の `frontend.routing` を読む
2. 各ルートに対応するページファイルが `frontend/src/pages/` に存在するか Glob で確認する
3. 存在しない場合 → `PHANTOM_ROUTE`

---

#### AI-2: 手打ちコマンド検知

`docker compose`、`go run`、`npm run`、`playwright test` などのコマンドが
ドキュメントで直接実行を案内されていないか検証する。

1. `CLAUDE.md`、`README.md`、`docs/`、`.claude/docs/` 内のドキュメントを Read で読む
2. 以下のコマンドパターンを Grep で検出する：
   - `docker compose`（Makefile 外での使用）
   - `go run` / `go build`
   - `npm run`（package.json の scripts 定義は除外）
   - `playwright test` / `npx`
3. Makefile や scripts 内での使用は OK
4. ドキュメントで「直接実行してください」と案内されている場合 → `MANUAL_COMMAND_DETECTED`

---

#### AI-3: 生成物重複検知

同じ概念を表すファイルが複数存在していないかを検証する。
AIDD では AI が似たファイルを複数生成してしまうことがよくある。

1. Glob でプロジェクトルート直下および `docs/`、`.claude/docs/` 配下のファイルを列挙する
2. 以下のキーワードグループで重複候補を検出する：
   - `map`, `structure`, `tree`, `layout`
   - `guide`, `manual`, `howto`, `instructions`
   - `config`, `settings`, `options`
   - `spec`, `schema`, `definition`
3. 候補が見つかった場合、ファイル内容の先頭50行を比較する
4. 内容が類似している場合 → `DUPLICATE_ARTIFACT`

---

#### AI-4: スキル内の壊れた参照

スキルファイル内で参照されている他のスキル名・ファイルパスが実際に存在するかを検証する。
（スキル一覧の整合性は doctor.sh セクション10でカバー済み。ここではファイル内容の参照チェックのみ。）

1. `.claude/skills/` 配下の全 `SKILL.md` を Read で読む
2. 各スキル内で参照されている他のスキル名・ファイルパスを抽出する
3. 参照先が実際に存在するか Glob で確認する
4. 存在しない場合 → `BROKEN_REFERENCE`

---

## 診断レポート出力フォーマット

全チェックが完了したら、以下のフォーマットでレポートを出力する。

```
# PROJECT CONSISTENCY REPORT

プロジェクト: [プロジェクト名]
診断日時: [日時]
総合スコア: [平均スコア]/10

## スコアサマリー

| # | カテゴリ | スコア | 検出数 | 状態 |
|---|---------|--------|--------|------|
| doctor.sh | 機械的チェック（10項目） | — | N件 | OK/WARN/NG |
| AI-1 | ルーティング整合性 | 9/10 | 1件 | OK |
| AI-2 | 手打ちコマンド検知 | 10/10 | 0件 | OK |
| AI-3 | 生成物重複検知 | 10/10 | 0件 | OK |
| AI-4 | スキル壊れた参照 | 7/10 | 2件 | WARN |

## doctor.sh 結果

[make doctor の出力をそのまま貼る]

## AI チェック詳細

### [各カテゴリの詳細]

検出された問題:
- [ファイルパス]: [問題の種類] - [説明]

修正案:
- [具体的な修正手順]

## 優先対応事項

1. [最もスコアが低いカテゴリの修正案]
2. ...
```

## 状態の判定

- 8〜10/10: `OK` — 問題なし
- 5〜7/10: `WARN` — 対応推奨
- 1〜4/10: `NG` — 要対応
- 該当なし: `SKIP` — 対象外

## 注意事項

- このスキルは診断のみを行う。修正の実行は人間の承認を得てから行う
- スコアは目安であり、プロジェクトの特性に応じて重み付けを変えてよい
- 診断結果は必ず人間に報告し、修正の優先順位は人間が判断する
