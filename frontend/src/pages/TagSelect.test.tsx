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

  it("タグがない場合はメッセージを表示する", async () => {
    mockApi.getTags.mockResolvedValue([]);
    renderTagSelect();
    await waitFor(() => {
      expect(screen.getByText("タグがありません。タグ管理から作成してください。")).toBeInTheDocument();
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
});
