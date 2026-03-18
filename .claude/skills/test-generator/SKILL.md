---
name: test-generator
description: >
  実装に対するテストを生成する。ユニットテスト、エッジケース、E2Eテストのカバー。
  「テストを書いて」「テストを追加して」「テストが足りない」「エッジケースをカバーして」
  などの場面で使う。実装後にテストが不足している場合も積極的に提案すること。
---

# Test Generator

実装に対して **適切なテストを生成する** スキル。

## テスト戦略

プロジェクトの3層テスト：

| 層 | ツール | 対象 | 実行 |
|---|---|---|---|
| Backend Unit | `go test` | handlers, repositories（services がある場合は services も） | `make test` |
| Frontend Unit | Vitest | components, hooks, utils | `make test-frontend` |
| E2E | Playwright | ユーザーフロー全体 | `make e2e` |

## 手順

### 1. 変更差分の特定

まず `git diff` で今回の変更差分を確認し、テスト対象を絞り込む：

```bash
git diff --name-only HEAD   # 変更ファイル一覧
git diff HEAD               # 変更内容の詳細
```

- 変更・追加されたファイルを特定する
- 変更された関数・コンポーネントをリストアップする
- **差分のあるコードに対してテストを書く**。無関係なコードのテストは書かない

### 2. テスト対象の確認

差分のある関数・コンポーネントについて、以下を把握する：
- 関数・メソッドの入出力
- エッジケース（空値、境界値、エラー）
- 外部依存（DB、API）

### 3. 既存テストパターンの確認

同じ層の既存テストを読み、パターンを把握する：
- テストファイルの配置場所
- テスト関数の命名
- セットアップ・ティアダウンの方法
- モック・スタブの使い方

### 4. テスト生成

#### Backend (Go)

```go
// テーブルドリブンテストを基本とする
func TestXxx(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        // 正常系
        // エッジケース
        // エラーケース
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実行
        })
    }
}
```

#### Frontend (Vitest)

```typescript
// 既存テストのパターンに従う
describe('ComponentName', () => {
  it('正常系の動作', () => {
    // テスト
  });

  it('エッジケース', () => {
    // テスト
  });
});
```

#### E2E (Playwright)

```typescript
// ユーザーの操作フローをテスト
test('ユーザーが○○できる', async ({ page }) => {
  // ページ遷移
  // 操作
  // 検証
});
```

### 5. カバレッジ確認

以下のケースがカバーされているか確認：

- **正常系**: 期待通りの入力での動作
- **境界値**: 0、空文字、最大値など
- **エラー系**: 不正入力、DB エラー、API エラー
- **権限**: 認証・認可が関わる場合

## 出力フォーマット

```
## 生成テスト

| テストファイル | テスト内容 | 種別 |
|---|---|---|
| backend/handlers/xxx_test.go | [内容] | unit |
| backend/repositories/xxx_test.go | [内容] | unit |
| e2e/tests/xxx.spec.ts | [内容] | e2e |

## エッジケース

- [ケース1]: [テスト方法]
- [ケース2]: [テスト方法]
```

## 検証

テスト生成後、必ず実行して通ることを確認する：

```bash
make test          # Backend
make test-frontend # Frontend
make e2e           # E2E（変更がある場合）
```

## テスト配置ルール

テストファイルは対象コードと同じディレクトリに置く（co-located）：

```
backend/handlers/xxx.go      → backend/handlers/xxx_test.go
backend/repositories/xxx.go  → backend/repositories/xxx_test.go
backend/services/xxx.go      → backend/services/xxx_test.go（存在する場合のみ）
```

service がない場合は handler テスト + repository テストのみ。
service がある場合は service 単体テストも追加する。

## 注意

- 既存テストのパターンに従う
- 過度なモックは避ける。プロジェクトのテスト方針に合わせる
- テストが通らない場合は、テストコード側を修正する（実装側を勝手に変更しない）
