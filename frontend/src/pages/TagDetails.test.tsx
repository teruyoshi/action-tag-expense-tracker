import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { TagDetails } from "./TagDetails";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useNavigate: () => mockNavigate,
    useParams: () => ({ tagId: "1" }),
    useSearchParams: () => [new URLSearchParams("year=2026&month=3&name=通勤")],
  };
});

vi.mock("../api/client", () => ({
  api: {
    getTagExpenseDetails: vi.fn(),
    updateExpense: vi.fn(),
  },
}));

import { api } from "../api/client";
const mockApi = vi.mocked(api);

beforeEach(() => {
  vi.clearAllMocks();
  mockApi.getTagExpenseDetails.mockResolvedValue([
    { id: 1, date: "2026-03-01", item: "電車賃", amount: 500 },
    { id: 2, date: "2026-03-05", item: "", amount: 300 },
  ]);
});

describe("TagDetails", () => {
  const renderTagDetails = () => render(<MemoryRouter><TagDetails /></MemoryRouter>);

  it("タグ名と年月を表示する", () => {
    renderTagDetails();
    expect(screen.getByText("通勤")).toBeInTheDocument();
    expect(screen.getByText("2026年3月")).toBeInTheDocument();
  });

  it("明細を表示する", async () => {
    renderTagDetails();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
      expect(screen.getByText("¥500")).toBeInTheDocument();
      expect(screen.getByText("—")).toBeInTheDocument();
      expect(screen.getByText("¥300")).toBeInTheDocument();
    });
  });

  it("合計を表示する", async () => {
    renderTagDetails();
    await waitFor(() => {
      expect(screen.getByText("¥800")).toBeInTheDocument();
    });
  });

  it("明細がない場合は「データなし」を表示する", async () => {
    mockApi.getTagExpenseDetails.mockResolvedValue([]);
    renderTagDetails();
    await waitFor(() => {
      expect(screen.getByText("データなし")).toBeInTheDocument();
    });
  });

  it("戻るボタンで/に遷移する", async () => {
    renderTagDetails();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(mockNavigate).toHaveBeenCalledWith("/");
  });

  it("正しいパラメータでAPIを呼ぶ", () => {
    renderTagDetails();
    expect(mockApi.getTagExpenseDetails).toHaveBeenCalledWith(2026, 3, 1);
  });

  it("行をクリックすると編集モードになる", async () => {
    renderTagDetails();
    const user = userEvent.setup();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
    });
    await user.click(screen.getByText("電車賃"));
    expect(screen.getByDisplayValue("電車賃")).toBeInTheDocument();
    expect(screen.getByDisplayValue("500")).toBeInTheDocument();
    expect(screen.getByText("保存")).toBeInTheDocument();
    expect(screen.getByText("取消")).toBeInTheDocument();
  });

  it("取消ボタンで編集モードを終了する", async () => {
    renderTagDetails();
    const user = userEvent.setup();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
    });
    await user.click(screen.getByText("電車賃"));
    await user.click(screen.getByText("取消"));
    expect(screen.queryByDisplayValue("電車賃")).not.toBeInTheDocument();
    expect(screen.getByText("電車賃")).toBeInTheDocument();
  });

  it("保存ボタンで更新APIを呼び明細を更新する", async () => {
    mockApi.updateExpense.mockResolvedValue({ id: 1, event_id: 1, item: "バス代", amount: 600 });
    renderTagDetails();
    const user = userEvent.setup();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
    });
    await user.click(screen.getByText("電車賃"));

    const itemInput = screen.getByDisplayValue("電車賃");
    await user.clear(itemInput);
    await user.type(itemInput, "バス代");

    const amountInput = screen.getByDisplayValue("500");
    await user.clear(amountInput);
    await user.type(amountInput, "600");

    await user.click(screen.getByText("保存"));

    expect(mockApi.updateExpense).toHaveBeenCalledWith(1, "バス代", 600);
    await waitFor(() => {
      expect(screen.getByText("バス代")).toBeInTheDocument();
      expect(screen.getByText("¥600")).toBeInTheDocument();
    });
  });

  it("金額が0以下の場合はalertを表示して更新しない", async () => {
    const alertSpy = vi.spyOn(window, "alert").mockImplementation(() => {});
    renderTagDetails();
    const user = userEvent.setup();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
    });
    await user.click(screen.getByText("電車賃"));

    const amountInput = screen.getByDisplayValue("500");
    await user.clear(amountInput);
    await user.type(amountInput, "0");

    await user.click(screen.getByText("保存"));

    expect(alertSpy).toHaveBeenCalledWith("金額は1以上を入力してください");
    expect(mockApi.updateExpense).not.toHaveBeenCalled();
    alertSpy.mockRestore();
  });

  it("更新API失敗時にalertを表示する", async () => {
    const alertSpy = vi.spyOn(window, "alert").mockImplementation(() => {});
    mockApi.updateExpense.mockRejectedValue(new Error("server error"));
    renderTagDetails();
    const user = userEvent.setup();
    await waitFor(() => {
      expect(screen.getByText("電車賃")).toBeInTheDocument();
    });
    await user.click(screen.getByText("電車賃"));
    await user.click(screen.getByText("保存"));

    await waitFor(() => {
      expect(alertSpy).toHaveBeenCalledWith("更新に失敗しました");
    });
    // 編集モードが維持され、元の値がinputに残っていること
    expect(screen.getByDisplayValue("500")).toBeInTheDocument();
    alertSpy.mockRestore();
  });
});
