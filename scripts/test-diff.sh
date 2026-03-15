#!/usr/bin/env bash
set -euo pipefail

# test-diff.sh
# git diff から変更されたファイルを特定し、影響範囲のテストのみ実行する。
# Backend: 変更された .go ファイルのパッケージのテストを実行
# Frontend: 変更された .ts/.tsx ファイルに対応するテストを実行（未導入時はスキップ）

cd "$(git rev-parse --show-toplevel)"

CHANGED_FILES=$(git diff --name-only HEAD 2>/dev/null || true)
STAGED_FILES=$(git diff --cached --name-only 2>/dev/null || true)
ALL_CHANGED=$(echo -e "${CHANGED_FILES}\n${STAGED_FILES}" | sort -u | grep -v '^$' || true)

if [ -z "$ALL_CHANGED" ]; then
  echo "変更ファイルなし。テストをスキップします。"
  exit 0
fi

echo "=== 変更ファイル ==="
echo "$ALL_CHANGED"
echo ""

# ---------- Backend ----------

GO_CHANGED=$(echo "$ALL_CHANGED" | grep '\.go$' || true)
BACKEND_RESULT=0

if [ -n "$GO_CHANGED" ]; then
  # 変更された .go ファイルからパッケージパスを抽出
  GO_PACKAGES=$(echo "$GO_CHANGED" | sed 's|backend/|./|' | xargs -I{} dirname {} | sort -u)

  echo "=== Backend: 影響パッケージ ==="
  echo "$GO_PACKAGES"
  echo ""

  echo "=== Backend: テスト実行 ==="
  for pkg in $GO_PACKAGES; do
    echo "--- $pkg ---"
    docker compose exec -T backend go test "$pkg" || BACKEND_RESULT=1
  done
  echo ""
else
  echo "=== Backend: 変更なし（スキップ）==="
  echo ""
fi

# ---------- Frontend ----------

FRONTEND_CHANGED=$(echo "$ALL_CHANGED" | grep -E '^frontend/.*\.(ts|tsx)$' || true)
FRONTEND_RESULT=0

if [ -n "$FRONTEND_CHANGED" ]; then
  # Vitest が導入されているかチェック
  if grep -q '"vitest"' frontend/package.json 2>/dev/null; then
    # 変更ファイルに対応するテストファイルを探す
    TEST_FILES=""
    for f in $FRONTEND_CHANGED; do
      base=$(echo "$f" | sed 's|frontend/||')
      dir=$(dirname "$base")
      name=$(basename "$base" | sed 's/\.\(ts\|tsx\)$//')
      # .test.ts, .test.tsx, .spec.ts, .spec.tsx を探す
      for ext in test.ts test.tsx spec.ts spec.tsx; do
        candidate="frontend/${dir}/${name}.${ext}"
        if [ -f "$candidate" ]; then
          TEST_FILES="$TEST_FILES $candidate"
        fi
      done
      # 変更ファイル自体がテストファイルの場合
      if echo "$f" | grep -qE '\.(test|spec)\.(ts|tsx)$'; then
        TEST_FILES="$TEST_FILES $f"
      fi
    done

    TEST_FILES=$(echo "$TEST_FILES" | tr ' ' '\n' | sort -u | grep -v '^$' || true)

    if [ -n "$TEST_FILES" ]; then
      echo "=== Frontend: テスト実行 ==="
      echo "$TEST_FILES"
      # frontend/ プレフィックスを除去してコンテナ内パスに変換
      CONTAINER_TEST_FILES=$(echo "$TEST_FILES" | sed 's|frontend/|/app/|g')
      docker compose exec -T frontend npx vitest run $CONTAINER_TEST_FILES || FRONTEND_RESULT=1
    else
      echo "=== Frontend: 対応するテストファイルなし（スキップ）==="
    fi
  else
    echo "=== Frontend: Vitest 未導入（スキップ）==="
  fi
  echo ""
else
  echo "=== Frontend: 変更なし（スキップ）==="
  echo ""
fi

# ---------- 結果 ----------

if [ $BACKEND_RESULT -ne 0 ] || [ $FRONTEND_RESULT -ne 0 ]; then
  echo "❌ 一部テストが失敗しました"
  exit 1
fi

echo "✅ 影響範囲のテストすべて通過"
