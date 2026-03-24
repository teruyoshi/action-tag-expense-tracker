---
name: feature-to-slices
description: >
  Featureを実装可能なSlice（小さな作業単位）に分解する。
  ユーザーが「機能を追加したい」「この機能を実装したい」「Feature を分解して」などと言ったときに使う。
  新しい機能開発の最初のステップとして、このスキルを積極的に提案すること。
---

# Feature to Slices

Featureは大きすぎてそのまま実装すると危険。
このスキルは Feature を **Slice（実装単位）** に分解する。

## なぜ分解するのか

- 大きな変更は壊れやすい
- レビューが困難になる
- mainブランチが不安定になる

小さく分割すれば、各Sliceが独立して検証可能になる。

## Sliceのサイズ基準

- 30〜200行の変更
- 1〜3ファイルの変更
- 30分〜2時間の作業量

この範囲を超える場合はさらに分割する。

## 手順

### 1. Feature の理解

ユーザーからFeatureの説明を受け取り、以下を確認する：
- 何を実現するのか
- どのレイヤー（backend / frontend / DB / API）に関係するか
- 既存機能との関連

### 2. repo_map.yaml の確認

プロジェクトルートに `repo_map.yaml` がある場合は読み込み、以下を把握する：
- レイヤー構成（`backend.structure`）
- ルーティングと画面遷移（`frontend.routing`、`frontend.flow`）
- ドメインモデル（`domain.entities`）

### 3. Slice への分解

Featureを以下の観点で分解する：

- **データ層**: DBテーブル追加、マイグレーション
- **ビジネスロジック層**: repository の追加・変更（必要な場合のみ service）
- **API層**: handler の追加・変更
- **フロントエンド層**: ページ、コンポーネント、API クライアント
- **テスト層**: 各層のテスト追加
- **E2E層**: Playwright テスト

### 4. 依存関係の整理

Slice間の依存関係を明確にする。
例：DB変更 → repository → handler → frontend の順序。
（service が必要な場合のみ repository → service → handler）

### 5. 実装順序の提案

依存関係に基づき、安全な実装順序を提案する。
原則として **データ層から UI層へ** 向かう順序にする。

## Service導入判定

各Sliceで service 層が必要かを判定する。デフォルト構造は handler → repository。

以下の場合のみ service を導入する：
- 複数 repository をまたぐ処理
- 副作用を伴う処理
- トランザクションが必要な処理
- 明確なユースケース単位の処理

## 出力フォーマット

```
## Slice一覧

1. [Slice名] - [概要] (変更見込: 約XX行, X files)
2. ...

## 実装レイヤー

- handlerのみ
- handler + repository
- handler + service（例外）

## Service導入判定

- 導入: YES / NO
- 理由:
- 対象:

## 依存関係

Slice 1 → Slice 2 → Slice 3
Slice 2 → Slice 4

## 実装順序

1. Slice 1 (依存なし)
2. Slice 2 (Slice 1完了後)
3. ...
```

## 技術スタック参考

- Backend: Go / chi / GORM / golang-migrate
- Frontend: React / TypeScript / Vite
- DB: MySQL
- E2E: Playwright
- コマンド: Make 経由（`make dev`, `make test`, `make lint`, `make e2e`）

## 注意

- アーキテクチャの決定はしない。分解案を提示し、人間が最終判断する
- 既存のパターンに従う。新しいフレームワークやライブラリを提案しない
- 不明点があれば分解前に質問する
