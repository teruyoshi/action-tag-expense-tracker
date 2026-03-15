import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { TagSelect } from "./TagSelect";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

vi.mock("../api/client", () => ({
  api: {
    getTags: vi.fn(),
    createTag: vi.fn(),
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
  const renderTagSelect = () => render(<MemoryRouter><TagSelect /></MemoryRouter>);

  it("タグ一覧を表示する", async () => {
    renderTagSelect();
    await waitFor(() => {
      expect(screen.getByText("通勤")).toBeInTheDocument();
      expect(screen.getByText("外食")).toBeInTheDocument();
    });
  });

  it("タグがない場合でも新規作成ボタンを表示する", async () => {
    mockApi.getTags.mockResolvedValue([]);
    renderTagSelect();
    await waitFor(() => {
      expect(screen.getByText("+ 新しい支出理由")).toBeInTheDocument();
    });
  });

  it("タグをクリックすると/expense/newに遷移する", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getByText("通勤"));
    expect(mockNavigate).toHaveBeenCalledWith("/expense/new", {
      state: { tag: { id: 1, name: "通勤" } },
    });
  });

  it("戻るボタンで/に遷移する", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(mockNavigate).toHaveBeenCalledWith("/");
  });

  it("新しい支出理由ボタンで入力欄を表示する", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("+ 新しい支出理由"));

    await user.click(screen.getByText("+ 新しい支出理由"));
    expect(screen.getByPlaceholderText("新しい支出理由")).toBeInTheDocument();
    expect(screen.getByText("追加")).toBeInTheDocument();
    expect(screen.getByText("取消")).toBeInTheDocument();
  });

  it("新しい支出理由を作成して遷移する", async () => {
    const newTag = { id: 3, name: "買い物" };
    mockApi.createTag.mockResolvedValue(newTag);
    mockApi.getTags.mockResolvedValueOnce([
      { id: 1, name: "通勤" },
      { id: 2, name: "外食" },
    ]).mockResolvedValueOnce([
      { id: 1, name: "通勤" },
      { id: 2, name: "外食" },
      newTag,
    ]);

    renderTagSelect();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("+ 新しい支出理由"));

    await user.click(screen.getByText("+ 新しい支出理由"));
    await user.type(screen.getByPlaceholderText("新しい支出理由"), "買い物");
    await user.click(screen.getByText("追加"));

    await waitFor(() => {
      expect(mockApi.createTag).toHaveBeenCalledWith("買い物");
      expect(mockNavigate).toHaveBeenCalledWith("/expense/new", {
        state: { tag: newTag },
      });
    });
  });

  it("取消ボタンで入力欄を閉じる", async () => {
    renderTagSelect();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("+ 新しい支出理由"));

    await user.click(screen.getByText("+ 新しい支出理由"));
    await user.click(screen.getByText("取消"));
    expect(screen.getByText("+ 新しい支出理由")).toBeInTheDocument();
    expect(screen.queryByPlaceholderText("新しい支出理由")).not.toBeInTheDocument();
  });
});
