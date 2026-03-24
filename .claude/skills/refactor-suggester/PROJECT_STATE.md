# 確認済み状態（2026-03-22時点）

直近の分析で確認済みの事項。再分析時の参考にする。

## Backend
- ヘルパー関数（`writeJSON`, `writeError`, `parseDate`, `parseYearMonth`）は `handlers/helpers.go` に集約済み。重複なし
- 全handlerのエラーレスポンスは `writeError` でJSON統一済み
- Repository層は各ドメインごとに1ファイル。インターフェースは `interfaces.go` に集約済み
- services層は `BalanceService` のみ。他のhandlerはrepository直接呼び出し（意図的な設計）
- Summary系Repository（`TagMonthTotals`, `TagMonthTotalsWithDiff`, `TagExpenseDetails`）はデータなし時に `[]` を返す（`make([]T, 0)` で初期化済み）

## Frontend
- カスタムhooks導入済み: `useBalance`（所持金）、`useTags`（タグCRUD）、`useExpenses`（支出作成・更新）
- `useExpenses` は状態を持たないhook（apiの薄いラッパー）。hookとしての形式統一を優先し現状維持
- `Home.tsx` は `useBalance` hookと `api` を両方importしている（月サマリー取得が直接呼び出し）。`useSummary` hook導入は変更が大きいため見送り中
- 各ページの「戻る」ボタンは似ているが遷移先が異なるため共通化のメリットは薄い
- 共通コンポーネント（`MonthNav`, `BalanceCard`, `TagSummaryCard`）は `components/` に分離済み
- StrictMode は常時有効（CI無効化ロジックは削除済み）
- `Home.tsx` の `getTagTotalsWithDiff` に `?? []` フォールバックあり（Backendも `[]` を返すようになったが防御コードとして残存）

## repo_map.yaml
- `repo_map.yaml` は構造情報のみ保持。リファクタリング時にはレイヤー構成・APIルート・ドメイン定義が実コードと一致しているか確認すること
