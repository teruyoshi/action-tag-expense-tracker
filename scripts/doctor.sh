#!/usr/bin/env bash
set -euo pipefail

# doctor.sh
# プロジェクトの整合性を診断する。
# チェック項目:
#   1. repo_map.yaml の整合性
#   2. Makefile の整合性
#   3. Docker の整合性
#   4. ディレクトリ構造
#   5. 依存関係

cd "$(git rev-parse --show-toplevel)"

ERRORS=0
WARNINGS=0

pass() { echo "  ✅ $1"; }
warn() { echo "  ⚠️  $1"; WARNINGS=$((WARNINGS + 1)); }
fail() { echo "  ❌ $1"; ERRORS=$((ERRORS + 1)); }

# ---------- 1. repo_map.yaml ----------

echo "=== 1. repo_map.yaml 整合性 ==="

if [ ! -f "repo_map.yaml" ]; then
  fail "repo_map.yaml が存在しません"
else
  pass "repo_map.yaml が存在します"

  # repo_map に記載されたエントリーポイントが実在するか
  if grep -q "cmd/server/main.go" repo_map.yaml; then
    if [ -f "backend/cmd/server/main.go" ]; then
      pass "Backend エントリーポイントが実在します"
    else
      fail "repo_map に記載の backend/cmd/server/main.go が存在しません"
    fi
  fi

  # handlers/ に記載されたファイルが実在するか
  for f in $(grep -oE 'handlers/[a-z_]+\.go' repo_map.yaml 2>/dev/null || true); do
    if [ -f "backend/$f" ]; then
      pass "$f が実在します"
    else
      fail "repo_map に記載の $f が存在しません"
    fi
  done

  # 実際の handlers/ にあるが repo_map に記載されていないファイル
  for f in backend/handlers/*.go; do
    base=$(basename "$f")
    if ! grep -q "$base" repo_map.yaml 2>/dev/null; then
      warn "$(echo "$f" | sed 's|backend/||') が repo_map.yaml に記載されていません"
    fi
  done

  # repositories/ も同様にチェック
  for f in backend/repositories/*.go; do
    base=$(basename "$f")
    if ! grep -q "$base" repo_map.yaml 2>/dev/null; then
      warn "$(echo "$f" | sed 's|backend/||') が repo_map.yaml に記載されていません"
    fi
  done
fi
echo ""

# ---------- 2. Makefile 整合性 ----------

echo "=== 2. Makefile 整合性 ==="

if [ ! -f "Makefile" ]; then
  fail "Makefile が存在しません"
else
  pass "Makefile が存在します"

  # 必須コマンドの存在チェック
  REQUIRED_TARGETS="dev stop build fmt fmt-frontend fmt-check fmt-check-frontend lint lint-frontend typecheck test test-frontend test-diff e2e quick-check check verify security-check security-check-full"
  for target in $REQUIRED_TARGETS; do
    if grep -qE "^${target}:" Makefile; then
      pass "make $target が定義されています"
    else
      fail "make $target が未定義です"
    fi
  done

  # CLAUDE.md に記載のコマンドが Makefile に存在するか
  if [ -f "CLAUDE.md" ]; then
    for target in $(grep -oP 'make \K[a-z][a-z-]*[a-z]' CLAUDE.md | sort -u); do
      if grep -qE "^${target}:" Makefile; then
        pass "CLAUDE.md 記載の make $target が Makefile に存在します"
      else
        warn "CLAUDE.md 記載の make $target が Makefile に存在しません"
      fi
    done
  fi
fi
echo ""

# ---------- 3. Docker 整合性 ----------

echo "=== 3. Docker 整合性 ==="

if [ ! -f "docker-compose.yml" ]; then
  fail "docker-compose.yml が存在しません"
else
  pass "docker-compose.yml が存在します"

  # Dockerfile の存在チェック
  if [ -f "backend/Dockerfile" ]; then
    pass "backend/Dockerfile が存在します"
  else
    fail "backend/Dockerfile が存在しません"
  fi

  if [ -f "frontend/Dockerfile" ]; then
    pass "frontend/Dockerfile が存在します"
  else
    fail "frontend/Dockerfile が存在しません"
  fi

  if [ -f "e2e/Dockerfile" ]; then
    pass "e2e/Dockerfile が存在します"
  else
    fail "e2e/Dockerfile が存在しません"
  fi

  # docker compose が利用可能か
  if docker compose version > /dev/null 2>&1; then
    pass "docker compose が利用可能です"

    # コンテナが起動しているか
    RUNNING_OUTPUT=$(docker compose ps --status running --format '{{.Name}}' 2>/dev/null || true)
    RUNNING_COUNT=0
    if [ -n "$RUNNING_OUTPUT" ]; then
      RUNNING_COUNT=$(echo "$RUNNING_OUTPUT" | wc -l | tr -d ' ')
    fi
    if [ "$RUNNING_COUNT" -gt 0 ]; then
      pass "Docker コンテナが起動中です ($RUNNING_COUNT サービス)"
    else
      warn "Docker コンテナが起動していません（make dev で起動してください）"
    fi
  else
    warn "docker compose が利用できません"
  fi
fi
echo ""

# ---------- 4. ディレクトリ構造 ----------

echo "=== 4. ディレクトリ構造 ==="

REQUIRED_DIRS="backend backend/cmd/server backend/handlers backend/repositories backend/models backend/migrations frontend frontend/src frontend/src/pages frontend/src/api e2e e2e/tests scripts .claude/skills .github/workflows"
for dir in $REQUIRED_DIRS; do
  if [ -d "$dir" ]; then
    pass "$dir/ が存在します"
  else
    fail "$dir/ が存在しません"
  fi
done
echo ""

# ---------- 5. 依存関係 ----------

echo "=== 5. 依存関係 ==="

# Backend: go.mod
if [ -f "backend/go.mod" ]; then
  pass "backend/go.mod が存在します"
else
  fail "backend/go.mod が存在しません"
fi

# Frontend: package.json
if [ -f "frontend/package.json" ]; then
  pass "frontend/package.json が存在します"

  # node_modules の存在（ホストまたはコンテナ）
  if [ -d "frontend/node_modules" ]; then
    pass "frontend/node_modules が存在します（ホスト）"
  elif docker compose ps --status running --format '{{.Name}}' 2>/dev/null | grep -q frontend; then
    pass "frontend/node_modules はコンテナ内に存在します"
  else
    warn "frontend/node_modules がありません（make dev でコンテナを起動してください）"
  fi

  # テストツールの導入状況
  if grep -q '"vitest"' frontend/package.json 2>/dev/null; then
    pass "Vitest が導入されています"
  else
    warn "Vitest が未導入です（make test-frontend は動作しません）"
  fi

  if grep -q '"playwright"' frontend/package.json 2>/dev/null || [ -d "e2e" ]; then
    pass "E2E テスト環境があります"
  else
    warn "E2E テスト環境が未構築です"
  fi
else
  fail "frontend/package.json が存在しません"
fi
echo ""

# ---------- 結果サマリー ----------

echo "=============================="
echo "  診断結果"
echo "=============================="
echo "  エラー:   $ERRORS"
echo "  警告:     $WARNINGS"
echo ""

if [ $ERRORS -gt 0 ]; then
  echo "❌ $ERRORS 件のエラーがあります。修正してください。"
  exit 1
elif [ $WARNINGS -gt 0 ]; then
  echo "⚠️  $WARNINGS 件の警告があります。確認を推奨します。"
  exit 0
else
  echo "✅ すべてのチェックに通過しました。"
  exit 0
fi
