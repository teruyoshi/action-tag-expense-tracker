import { test, expect } from "@playwright/test";

test.describe("Balance更新フロー", () => {
  test("所持金を設定して反映される", async ({ page }) => {
    await page.goto("/");

    // 所持金セクションが表示される
    await expect(page.getByText("所持金")).toBeVisible();

    // 設定ボタンをクリック
    await page.getByRole("button", { name: "設定" }).click();

    // モーダルが表示される
    await expect(page.getByText("所持金を設定")).toBeVisible();

    // 金額を入力して保存
    const input = page.getByPlaceholder("金額を入力");
    await input.fill("250000");
    await page.getByRole("button", { name: "保存" }).click();

    // モーダルが閉じるのを待つ
    await expect(page.getByText("所持金を設定")).not.toBeVisible({
      timeout: 10000,
    });

    // 金額が反映される
    await expect(page.getByText("¥250,000")).toBeVisible();
  });

  test("キャンセルでモーダルが閉じる", async ({ page }) => {
    await page.goto("/");

    await page.getByRole("button", { name: "設定" }).click();
    await expect(page.getByText("所持金を設定")).toBeVisible();

    await page.getByRole("button", { name: "キャンセル" }).click();
    await expect(page.getByText("所持金を設定")).not.toBeVisible();
  });
});
