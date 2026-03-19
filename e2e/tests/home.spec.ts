import { test, expect } from "@playwright/test";

test.describe("Home page", () => {
  test("トップページが表示される", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("heading", { name: "家計簿" })).toBeVisible();
  });

  test("所持金セクションが表示される", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByText("所持金")).toBeVisible();
  });

  test("月支出合計セクションが表示される", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByText("月支出合計")).toBeVisible();
  });

  test("支出入力ボタンから画面遷移できる", async ({ page }) => {
    await page.goto("/", { waitUntil: "networkidle" });

    await page.getByRole("button", { name: "支出入力" }).click();

    await expect(page.getByText("支出理由を選択")).toBeVisible();
  });

  test("支出理由管理ボタンから画面遷移できる", async ({ page }) => {
    await page.goto("/", { waitUntil: "networkidle" });

    await page.getByRole("button", { name: "支出理由管理" }).click();

    await expect(page.getByRole("heading", { name: "支出理由管理" })).toBeVisible();
  });
});
