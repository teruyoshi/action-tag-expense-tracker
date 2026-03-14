import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ExpenseInput } from "./ExpenseInput";

vi.mock("../api/client", () => ({
  api: {
    createEvent: vi.fn(),
    createExpense: vi.fn(),
  },
}));

import { api } from "../api/client";
const mockApi = vi.mocked(api);

beforeEach(() => {
  vi.clearAllMocks();
  mockApi.createEvent.mockResolvedValue({ id: 1, date: "2026-03-14", action_tag_id: 1 });
  mockApi.createExpense.mockResolvedValue({ id: 1, event_id: 1, item: "", amount: 500 });
});

describe("ExpenseInput", () => {
  const tag = { id: 1, name: "通勤" };
  const onBack = vi.fn();
  const onSaved = vi.fn();
  const renderExpenseInput = () => render(<ExpenseInput tag={tag} onBack={onBack} onSaved={onSaved} />);

  it("タグ名と日付が表示される", () => {
    renderExpenseInput();
    expect(screen.getByDisplayValue("通勤")).toBeInTheDocument();
    expect(screen.getByDisplayValue(new Date().toISOString().split("T")[0])).toBeInTheDocument();
  });

  it("項目追加ボタンで行が増える", async () => {
    renderExpenseInput();
    const user = userEvent.setup();

    expect(screen.getAllByPlaceholderText("金額")).toHaveLength(1);
    await user.click(screen.getByText("+ 項目追加"));
    expect(screen.getAllByPlaceholderText("金額")).toHaveLength(2);
  });

  it("金額を入力して保存できる", async () => {
    renderExpenseInput();
    const user = userEvent.setup();

    await user.type(screen.getByPlaceholderText("項目名"), "電車賃");
    await user.type(screen.getByPlaceholderText("金額"), "500");
    await user.click(screen.getByText("保存"));

    expect(mockApi.createEvent).toHaveBeenCalledWith(
      expect.any(String),
      1
    );
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, "電車賃", 500);
    expect(onSaved).toHaveBeenCalled();
  });

  it("金額が0以下の行は送信されない", async () => {
    renderExpenseInput();
    const user = userEvent.setup();

    // 金額を入力しない状態で保存
    await user.click(screen.getByText("保存"));
    expect(mockApi.createEvent).not.toHaveBeenCalled();
  });

  it("戻るボタンでonBackが呼ばれる", async () => {
    renderExpenseInput();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(onBack).toHaveBeenCalled();
  });

  it("複数行の支出を保存できる", async () => {
    renderExpenseInput();
    const user = userEvent.setup();

    // 1行目を入力
    const itemInputs = screen.getAllByPlaceholderText("項目名");
    const amountInputs = screen.getAllByPlaceholderText("金額");
    await user.type(itemInputs[0], "電車賃");
    await user.type(amountInputs[0], "500");

    // 2行目を追加して入力
    await user.click(screen.getByText("+ 項目追加"));
    const newAmountInputs = screen.getAllByPlaceholderText("金額");
    await user.type(screen.getAllByPlaceholderText("項目名")[1], "バス代");
    await user.type(newAmountInputs[1], "300");

    await user.click(screen.getByText("保存"));

    expect(mockApi.createExpense).toHaveBeenCalledTimes(2);
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, "電車賃", 500);
    expect(mockApi.createExpense).toHaveBeenCalledWith(1, "バス代", 300);
  });
});
