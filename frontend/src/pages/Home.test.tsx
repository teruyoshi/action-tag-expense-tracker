import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { Home } from "./Home";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

vi.mock("../api/client", () => ({
  api: {
    getMonthTotal: vi.fn(),
    getTagTotals: vi.fn(),
    getBalance: vi.fn(),
    updateBalance: vi.fn(),
  },
}));

import { api } from "../api/client";
const mockApi = vi.mocked(api);

beforeEach(() => {
  vi.clearAllMocks();
  mockApi.getMonthTotal.mockResolvedValue({ total: 15000 });
  mockApi.getTagTotals.mockResolvedValue([
    { tag_id: 1, tag: "通勤", total: 5000 },
    { tag_id: 2, tag: "外食", total: 10000 },
  ]);
  mockApi.getBalance.mockResolvedValue({ id: 1, amount: 100000, updated_at: "" });
});

describe("Home", () => {
  const renderHome = () => render(<MemoryRouter><Home /></MemoryRouter>);

  it("月支出合計を表示する", async () => {
    renderHome();
    await waitFor(() => {
      expect(screen.getByText("¥15,000")).toBeInTheDocument();
    });
  });

  it("所持金を表示する", async () => {
    renderHome();
    await waitFor(() => {
      expect(screen.getByText("¥100,000")).toBeInTheDocument();
    });
  });

  it("タグ別合計を表示する", async () => {
    renderHome();
    await waitFor(() => {
      expect(screen.getByText("通勤")).toBeInTheDocument();
      expect(screen.getByText("¥5,000")).toBeInTheDocument();
      expect(screen.getByText("外食")).toBeInTheDocument();
      expect(screen.getByText("¥10,000")).toBeInTheDocument();
    });
  });

  it("タグデータがない場合は「データなし」を表示する", async () => {
    mockApi.getTagTotals.mockResolvedValue([]);
    renderHome();
    await waitFor(() => {
      expect(screen.getByText("データなし")).toBeInTheDocument();
    });
  });

  it("月ナビゲーションで前月に移動する", async () => {
    renderHome();
    const user = userEvent.setup();
    const now = new Date();
    const curYear = now.getFullYear();
    const curMonth = now.getMonth() + 1;
    await waitFor(() => screen.getByText(`${curYear}年${curMonth}月`));

    await user.click(screen.getByText("<"));

    let expectedYear = curYear;
    let expectedMonth = curMonth - 1;
    if (expectedMonth === 0) {
      expectedMonth = 12;
      expectedYear--;
    }
    expect(screen.getByText(`${expectedYear}年${expectedMonth}月`)).toBeInTheDocument();
  });

  it("支出入力ボタンで/tags/selectに遷移する", async () => {
    renderHome();
    const user = userEvent.setup();
    await user.click(screen.getByText("支出入力"));
    expect(mockNavigate).toHaveBeenCalledWith("/tags/select");
  });

  it("タグ管理ボタンで/tags/manageに遷移する", async () => {
    renderHome();
    const user = userEvent.setup();
    await user.click(screen.getByText("タグ管理"));
    expect(mockNavigate).toHaveBeenCalledWith("/tags/manage");
  });

  it("タグクリックでタグ詳細に遷移する", async () => {
    renderHome();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getByText("通勤"));

    expect(mockNavigate).toHaveBeenCalledWith(
      expect.stringContaining("/tags/1/details")
    );
  });

  it("所持金モーダルで金額を更新できる", async () => {
    mockApi.updateBalance.mockResolvedValue({ id: 1, amount: 50000, updated_at: "" });
    renderHome();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("¥100,000"));

    await user.click(screen.getByText("設定"));
    expect(screen.getByText("所持金を設定")).toBeInTheDocument();

    const input = screen.getByPlaceholderText("金額を入力");
    await user.clear(input);
    await user.type(input, "50000");
    await user.click(screen.getByText("保存"));

    expect(mockApi.updateBalance).toHaveBeenCalledWith(50000);
    await waitFor(() => {
      expect(screen.getByText("¥50,000")).toBeInTheDocument();
    });
  });
});
