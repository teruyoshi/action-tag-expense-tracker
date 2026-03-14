import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { TagSelect } from "./TagSelect";

vi.mock("../api/client", () => ({
  api: {
    getTags: vi.fn(),
  },
}));

import { api } from "../api/client";
const mockApi = vi.mocked(api);

beforeEach(() => {
  vi.clearAllMocks();
  mockApi.getTags.mockResolvedValue([
    { id: 1, name: "通勤" },
    { id: 2, name: "外食" },
  ]);
});

describe("TagSelect", () => {
  const onSelect = vi.fn();
  const onBack = vi.fn();
  const renderTagSelect = () => render(<TagSelect onSelect={onSelect} onBack={onBack} />);

  it("タグ一覧を表示する", async () => {
    renderTagSelect();
    await waitFor(() => {
      expect(screen.getByText("通勤")).toBeInTheDocument();
      expect(screen.getByText("外食")).toBeInTheDocument();
    });
  });

  it("タグがない場合はメッセージを表示する", async () => {
    mockApi.getTags.mockResolvedValue([]);
    renderTagSelect();
    await waitFor(() => {
      expect(screen.getByText("タグがありません。タグ管理から作成してください。")).toBeInTheDocument();
    });
  });

  it("タグをクリックするとonSelectが呼ばれる", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getByText("通勤"));
    expect(onSelect).toHaveBeenCalledWith({ id: 1, name: "通勤" });
  });

  it("戻るボタンでonBackが呼ばれる", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(onBack).toHaveBeenCalled();
  });
});
