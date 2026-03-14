import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TagDetails } from "./TagDetails";

vi.mock("../api/client", () => ({
  api: {
    getTagExpenseDetails: vi.fn(),
  },
}));

import { api } from "../api/client";
const mockApi = vi.mocked(api);

beforeEach(() => {
  vi.clearAllMocks();
  mockApi.getTagExpenseDetails.mockResolvedValue([
    { date: "2026-03-01", item: "電車賃", amount: 500 },
    { date: "2026-03-05", item: "", amount: 300 },
  ]);
});

describe("TagDetails", () => {
  const defaultProps = {
    tagId: 1,
    tagName: "通勤",
    year: 2026,
    month: 3,
    onBack: vi.fn(),
  };
  const renderTagDetails = (props = defaultProps) => render(<TagDetails {...props} />);

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

  it("戻るボタンでonBackが呼ばれる", async () => {
    renderTagDetails();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(defaultProps.onBack).toHaveBeenCalled();
  });

  it("正しいパラメータでAPIを呼ぶ", () => {
    renderTagDetails();
    expect(mockApi.getTagExpenseDetails).toHaveBeenCalledWith(2026, 3, 1);
  });
});
