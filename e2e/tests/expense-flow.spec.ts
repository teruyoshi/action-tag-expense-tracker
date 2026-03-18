import { test, expect } from "@playwright/test";

test.describe("支出入力フロー", () => {
  test("タグ作成→支出入力→保存→Home", async ({ page }) => {
    const tagName = `テスト理由${Date.now()}`;

    // Home → 支出入力
    await page.goto("/");
    await page.getByRole("button", { name: "支出入力" }).click();
    await expect(page.getByText("支出理由を選択")).toBeVisible();

    // タグ作成
    await page.getByRole("button", { name: "+ 新しい支出理由" }).click();
    await page.getByPlaceholder("新しい支出理由").fill(tagName);
    await page.getByRole("button", { name: "追加" }).click();

    // 支出入力画面に自動遷移を待つ（URLで確認）
    await page.waitForURL("**/expense/new", { timeout: 15000 });

    // 金額入力
    await page.getByPlaceholder("金額").fill("1500");
    await page.getByPlaceholder("項目名").fill("テスト項目");

    // 保存
    await page.getByRole("button", { name: "保存" }).click();

    // Homeに戻る
    await page.waitForURL("/", { timeout: 10000 });
    await expect(page.getByRole("heading", { name: "家計簿" })).toBeVisible();
  });

  test("既存タグ選択→支出入力→保存", async ({ page }) => {
    const tagName = `既存タグ${Date.now()}`;

    // まずタグを作成（支出理由管理から）
    await page.goto("/");
    await page.getByRole("button", { name: "支出理由管理" }).click();
    await expect(
      page.getByRole("heading", { name: "支出理由管理" })
    ).toBeVisible();

    // タグを追加
    await page.getByPlaceholder("支出理由名").fill(tagName);
    await page.getByRole("button", { name: "追加" }).click();

    // タグが一覧に表示されるのを確認
    await expect(page.getByText(tagName)).toBeVisible({ timeout: 10000 });

    // Homeに戻る
    await page.getByRole("button", { name: "← 戻る" }).click();
    await expect(page.getByRole("heading", { name: "家計簿" })).toBeVisible();

    // 支出入力でタグを選択
    await page.getByRole("button", { name: "支出入力" }).click();
    await expect(page.getByRole("button", { name: tagName })).toBeVisible({
      timeout: 10000,
    });
    await page.getByRole("button", { name: tagName }).click();

    // 支出入力画面
    await page.waitForURL("**/expense/new", { timeout: 10000 });
    await page.getByPlaceholder("金額").fill("2000");
    await page.getByRole("button", { name: "保存" }).click();

    // Homeに戻る
    await page.waitForURL("/", { timeout: 10000 });
    await expect(page.getByRole("heading", { name: "家計簿" })).toBeVisible();
  });
});
