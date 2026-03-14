import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { TagManage } from "./TagManage";

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

vi.mock("../api/client", () => ({
  api: {
    getTags: vi.fn(),
    createTag: vi.fn(),
    updateTag: vi.fn(),
    deleteTag: vi.fn(),
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
  mockApi.createTag.mockResolvedValue({ id: 3, name: "" });
  mockApi.updateTag.mockResolvedValue({ id: 1, name: "" });
  mockApi.deleteTag.mockResolvedValue(undefined);
});

describe("TagManage", () => {
  const renderTagManage = () => render(<MemoryRouter><TagManage /></MemoryRouter>);

  it("タグ一覧を表示する", async () => {
    renderTagManage();
    await waitFor(() => {
      expect(screen.getByText("通勤")).toBeInTheDocument();
      expect(screen.getByText("外食")).toBeInTheDocument();
    });
  });

  it("タグがない場合はメッセージを表示する", async () => {
    mockApi.getTags.mockResolvedValue([]);
    renderTagManage();
    await waitFor(() => {
      expect(screen.getByText("タグがありません")).toBeInTheDocument();
    });
  });

  it("新規タグを追加できる", async () => {
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.type(screen.getByPlaceholderText("タグ名"), "買い物");
    await user.click(screen.getByText("追加"));

    expect(mockApi.createTag).toHaveBeenCalledWith("買い物");
  });

  it("空のタグ名では追加しない", async () => {
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getByText("追加"));
    expect(mockApi.createTag).not.toHaveBeenCalled();
  });

  it("タグを編集できる", async () => {
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    const editButtons = screen.getAllByText("編集");
    await user.click(editButtons[0]);

    const editInput = screen.getByDisplayValue("通勤");
    await user.clear(editInput);
    await user.type(editInput, "電車通勤");
    await user.click(screen.getByText("保存"));

    expect(mockApi.updateTag).toHaveBeenCalledWith(1, "電車通勤");
  });

  it("編集をキャンセルできる", async () => {
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getAllByText("編集")[0]);
    expect(screen.getByDisplayValue("通勤")).toBeInTheDocument();

    await user.click(screen.getByText("取消"));
    expect(screen.queryByDisplayValue("通勤")).not.toBeInTheDocument();
    expect(screen.getByText("通勤")).toBeInTheDocument();
  });

  it("タグを削除できる（confirm OK）", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(true);
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getAllByText("削除")[0]);
    expect(mockApi.deleteTag).toHaveBeenCalledWith(1);
  });

  it("タグ削除をキャンセルできる（confirm Cancel）", async () => {
    vi.spyOn(window, "confirm").mockReturnValue(false);
    renderTagManage();
    const user = userEvent.setup();
    await waitFor(() => screen.getByText("通勤"));

    await user.click(screen.getAllByText("削除")[0]);
    expect(mockApi.deleteTag).not.toHaveBeenCalled();
  });

  it("戻るボタンで/に遷移する", async () => {
    renderTagManage();
    const user = userEvent.setup();
    await user.click(screen.getByText("← 戻る"));
    expect(mockNavigate).toHaveBeenCalledWith("/");
  });
});
